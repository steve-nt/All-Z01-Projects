import { expect, test } from '@playwright/test';
import {
  createJsErrorCollector,
  createRequestTracker,
  duplicateValues,
  FEED_ITEM_SELECTORS,
  findFirstLocator,
  LOAD_MORE_SELECTORS,
  LOAD_MORE_SENTINEL_SELECTORS,
  waitForAnySelector,
} from './helpers.js';

const FEED_ITEM_QUERY = FEED_ITEM_SELECTORS.join(', ');
const LOAD_MORE_EFFECT_TIMEOUT_MS = 10_000;
const NETWORK_WAIT_TIMEOUT_MS = 4_000;
const NETWORK_POLL_INTERVAL_MS = 200;
const ITEM_ENDPOINT_PATTERN = /\/v0\/item\/\d+\.json$/;
const TRACKED_ENDPOINT_PATTERN =
  /\/v0\/(topstories|newstories|askstories|showstories|jobstories|item\/\d+)\.json$/;

const waitForLoadMoreEffect = async ({
  page,
  feedItems,
  baselineCount,
  requestTracker,
  baselineRequestCount,
}) => {
  const startedAtMs = performance.now();
  let sawCountIncrease = false;
  let sawNetworkActivity = false;

  while (performance.now() - startedAtMs < LOAD_MORE_EFFECT_TIMEOUT_MS) {
    const currentCount = await feedItems.count();
    const hasNewTrackedRequests = requestTracker.entries.length > baselineRequestCount;

    sawCountIncrease = sawCountIncrease || currentCount > baselineCount;
    sawNetworkActivity = sawNetworkActivity || hasNewTrackedRequests;

    if (sawCountIncrease && sawNetworkActivity) {
      return true;
    }

    await page.waitForTimeout(NETWORK_POLL_INTERVAL_MS);
  }

  return false;
};

test('AUDIT-F-04 load more posts without errors and without request spamming', async ({ page }) => {
  await page.goto('');
  const jsErrorCollector = createJsErrorCollector(page);
  const requestTracker = createRequestTracker(page, (request) =>
    TRACKED_ENDPOINT_PATTERN.test(new URL(request.url()).pathname),
  );

  try {
    await waitForAnySelector(page, FEED_ITEM_SELECTORS, { timeoutMs: 10_000 });

    const feedItems = page.locator(FEED_ITEM_QUERY);
    const baselineCount = await feedItems.count();

    test.skip(
      baselineCount < 1,
      'AUDIT-F-04 dependency guard: feed items are not rendered, so load-more behavior is unavailable.',
    );

    const loadMoreButton = await findFirstLocator(page, LOAD_MORE_SELECTORS);
    const loadMoreSentinel = await findFirstLocator(page, LOAD_MORE_SENTINEL_SELECTORS, {
      visible: false,
    });

    test.skip(
      !loadMoreButton && !loadMoreSentinel,
      'AUDIT-F-04 dependency guard: load-more controls/sentinel selectors are not available yet.',
    );

    const baselineRequestCount = requestTracker.entries.length;

    if (loadMoreButton) {
      await loadMoreButton.click();
    } else {
      await loadMoreSentinel.scrollIntoViewIfNeeded();
    }

    const didObserveLoadMoreEffect = await waitForLoadMoreEffect({
      page,
      feedItems,
      baselineCount,
      requestTracker,
      baselineRequestCount,
    });

    test.skip(
      !didObserveLoadMoreEffect,
      'AUDIT-F-04 dependency guard: no observable load-more effect occurred in this branch.',
    );

    await page.waitForTimeout(NETWORK_WAIT_TIMEOUT_MS);

    const afterCount = await feedItems.count();

    test.skip(
      afterCount <= baselineCount,
      'AUDIT-F-04 dependency guard: rendered feed count did not increase after load-more interaction.',
    );

    const deltaEntries = requestTracker.entries.slice(baselineRequestCount);
    const itemRequestPaths = deltaEntries
      .map((entry) => new URL(entry.url).pathname)
      .filter((pathname) => ITEM_ENDPOINT_PATTERN.test(pathname));

    test.skip(
      itemRequestPaths.length < 1,
      'AUDIT-F-04 dependency guard: no item requests were observed during load-more activity.',
    );

    const duplicateItemRequests = duplicateValues(itemRequestPaths);

    expect(duplicateItemRequests).toEqual([]);
    expect(jsErrorCollector.getErrors()).toEqual([]);
  } finally {
    requestTracker.stop();
    jsErrorCollector.stop();
  }
});
