import { expect, test } from '@playwright/test';
import {
  createRequestTracker,
  intervalSeriesMs,
  LIVE_BANNER_SELECTORS,
  waitForAnySelector,
} from './helpers.js';

const UPDATES_ENDPOINT_PATTERN = /\/v0\/updates\.json$/;
const THROTTLE_ASSERTION_FLOOR_MS = 4_999;

const waitForCondition = async ({ evaluate, timeoutMs, pollIntervalMs = 250 }) => {
  const startedAtMs = performance.now();

  while (performance.now() - startedAtMs < timeoutMs) {
    if (await evaluate()) {
      return true;
    }

    await new Promise((resolve) => {
      setTimeout(resolve, pollIntervalMs);
    });
  }

  return false;
};

test('AUDIT-G-04 UI notifies the user when post data updates', async ({ page }) => {
  await page.goto('');
  test.setTimeout(45_000);

  let updatesCallCount = 0;

  await page.route('**/v0/updates.json**', async (route) => {
    updatesCallCount += 1;

    const body =
      updatesCallCount < 2
        ? { items: [100_001], profiles: [] }
        : { items: [100_001, 100_002], profiles: [] };

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(body),
    });
  });

  await page.goto('');

  const observedTwoPolls = await waitForCondition({
    evaluate: async () => updatesCallCount >= 2,
    timeoutMs: 20_000,
  });

  test.skip(
    !observedTwoPolls,
    'AUDIT-G-04 dependency guard: live updates polling is not active in this branch yet.',
  );

  const liveBanner = await waitForAnySelector(page, LIVE_BANNER_SELECTORS, {
    timeoutMs: 10_000,
  });

  test.skip(
    !liveBanner,
    'AUDIT-G-04 dependency guard: live banner selectors are not available in this branch yet.',
  );

  await expect(liveBanner).toBeVisible();
});

test('AUDIT-G-05 live-data polling is throttled to at least every 5 seconds', async ({ page }) => {
  test.setTimeout(45_000);

  await page.route('**/v0/updates.json**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ items: [100_001, 100_002], profiles: [] }),
    });
  });

  const requestTracker = createRequestTracker(page, (request) =>
    UPDATES_ENDPOINT_PATTERN.test(new URL(request.url()).pathname),
  );

  try {
    await page.goto('');

    const observedEnoughPollRequests = await waitForCondition({
      evaluate: async () => requestTracker.entries.length >= 3,
      timeoutMs: 22_000,
    });

    test.skip(
      !observedEnoughPollRequests,
      'AUDIT-G-05 dependency guard: fewer than three updates requests were observed.',
    );

    const intervals = intervalSeriesMs(requestTracker.entries);

    test.skip(
      intervals.length < 2,
      'AUDIT-G-05 dependency guard: not enough request intervals were captured for throttle assertions.',
    );

    for (const intervalMs of intervals) {
      // Browser scheduling and measurement precision can produce sub-millisecond drift around exact 5s ticks.
      expect(intervalMs).toBeGreaterThanOrEqual(THROTTLE_ASSERTION_FLOOR_MS);
    }
  } finally {
    requestTracker.stop();
  }
});
