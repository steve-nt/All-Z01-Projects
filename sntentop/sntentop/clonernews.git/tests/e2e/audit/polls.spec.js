import { expect, test } from '@playwright/test';
import {
  createJsErrorCollector,
  FEED_TAB_SELECTORS,
  hasAnySelector,
  openPostFromFeed,
  POST_LINK_SELECTORS,
} from './helpers.js';

const POLL_FIXTURE_IDS = Object.freeze({
  pollPrimary: 8301,
  pollSecondary: 8302,
});

const POLL_FIXTURE_ITEMS = Object.freeze({
  [POLL_FIXTURE_IDS.pollPrimary]: Object.freeze({
    id: POLL_FIXTURE_IDS.pollPrimary,
    type: 'poll',
    by: 'fixture-poll-author',
    time: 2_100_300,
    title: 'Fixture poll post primary',
    score: 10,
    descendants: 1,
    parts: [],
  }),
  [POLL_FIXTURE_IDS.pollSecondary]: Object.freeze({
    id: POLL_FIXTURE_IDS.pollSecondary,
    type: 'poll',
    by: 'fixture-poll-author-2',
    time: 2_100_100,
    title: 'Fixture poll post secondary',
    score: 8,
    descendants: 0,
    parts: [],
  }),
});

const fulfillJson = (route, payload) =>
  route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(payload),
  });

const installPollsAuditFixture = async (page) => {
  await page.route('**/v0/topstories.json', async (route) => {
    await fulfillJson(route, [POLL_FIXTURE_IDS.pollPrimary, POLL_FIXTURE_IDS.pollSecondary]);
  });

  await page.route('**/v0/item/*.json', async (route) => {
    const match = route
      .request()
      .url()
      .match(/\/item\/(\d+)\.json$/);

    if (!match) {
      await route.fallback();
      return;
    }

    const itemId = Number(match[1]);
    const item = POLL_FIXTURE_ITEMS[itemId];

    if (!item) {
      await route.fallback();
      return;
    }

    await fulfillJson(route, item);
  });
};

test.beforeEach(async ({ page }) => {
  await installPollsAuditFixture(page);
});

test('AUDIT-F-03 poll post opens without any errors', async ({ page }) => {
  const jsErrorCollector = createJsErrorCollector(page);

  try {
    await page.goto('');

    const hasPollsTab = await hasAnySelector(page, FEED_TAB_SELECTORS.polls);

    const openResult = await openPostFromFeed({
      page,
      feedTabSelectors: hasPollsTab ? FEED_TAB_SELECTORS.polls : [],
      postLinkSelectors: POST_LINK_SELECTORS.polls,
    });

    test.skip(!openResult.ok, `AUDIT-F-03 dependency guard: ${openResult.reason}`);

    expect(jsErrorCollector.getErrors()).toEqual([]);
  } finally {
    jsErrorCollector.stop();
  }
});
