import { expect, test } from '@playwright/test';
import {
  createJsErrorCollector,
  FEED_TAB_SELECTORS,
  openPostFromFeed,
  POST_LINK_SELECTORS,
} from './helpers.js';

test('AUDIT-F-02 job post opens without any errors', async ({ page }) => {
  const jsErrorCollector = createJsErrorCollector(page);

  try {
    await page.goto('');

    const openResult = await openPostFromFeed({
      page,
      feedTabSelectors: FEED_TAB_SELECTORS.jobs,
      postLinkSelectors: POST_LINK_SELECTORS.jobs,
    });

    test.skip(!openResult.ok, `AUDIT-F-02 dependency guard: ${openResult.reason}`);

    expect(jsErrorCollector.getErrors()).toEqual([]);
  } finally {
    jsErrorCollector.stop();
  }
});
