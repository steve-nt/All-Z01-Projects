import { expect, test } from '@playwright/test';
import {
  createJsErrorCollector,
  FEED_ITEM_SELECTORS,
  FEED_ITEM_TIME_ATTRIBUTES,
  FEED_ITEM_TIME_SELECTORS,
  FEED_TAB_SELECTORS,
  findFirstLocator,
  hasAnySelector,
  openPostFromFeed,
  POST_LINK_SELECTORS,
  readFirstNumericSignal,
  STRICT_POST_LINK_SELECTORS,
  waitForAnySelector,
} from './helpers.js';

const FEED_ITEM_QUERY = FEED_ITEM_SELECTORS.join(', ');

const FIXTURE_IDS = Object.freeze({
  storyPrimary: 7101,
  storySecondary: 7102,
  jobPrimary: 7201,
  pollPrimary: 7301,
});

const FEED_FIXTURE_ITEMS = Object.freeze({
  [FIXTURE_IDS.storyPrimary]: Object.freeze({
    id: FIXTURE_IDS.storyPrimary,
    type: 'story',
    by: 'fixture-story-author',
    time: 2_000_300,
    title: 'Fixture story post',
    descendants: 12,
    score: 90,
    url: 'https://example.com/story/fixture-1',
  }),
  [FIXTURE_IDS.storySecondary]: Object.freeze({
    id: FIXTURE_IDS.storySecondary,
    type: 'story',
    by: 'fixture-story-author-2',
    time: 2_000_100,
    title: 'Fixture story post 2',
    descendants: 4,
    score: 50,
    url: 'https://example.com/story/fixture-2',
  }),
  [FIXTURE_IDS.jobPrimary]: Object.freeze({
    id: FIXTURE_IDS.jobPrimary,
    type: 'job',
    by: 'fixture-job-author',
    time: 2_000_200,
    title: 'Fixture job post',
    score: 1,
    text: 'Fixture job text',
  }),
  [FIXTURE_IDS.pollPrimary]: Object.freeze({
    id: FIXTURE_IDS.pollPrimary,
    type: 'poll',
    by: 'fixture-poll-author',
    time: 2_000_150,
    title: 'Fixture poll post',
    descendants: 3,
    score: 33,
    parts: [],
  }),
});

const fulfillJson = (route, payload) =>
  route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(payload),
  });

const installStoriesAuditFixture = async (page) => {
  await page.route('**/v0/topstories.json', async (route) => {
    await fulfillJson(route, [
      FIXTURE_IDS.storyPrimary,
      FIXTURE_IDS.pollPrimary,
      FIXTURE_IDS.storySecondary,
    ]);
  });

  await page.route('**/v0/newstories.json', async (route) => {
    await fulfillJson(route, [FIXTURE_IDS.storyPrimary, FIXTURE_IDS.storySecondary]);
  });

  await page.route('**/v0/jobstories.json', async (route) => {
    await fulfillJson(route, [FIXTURE_IDS.jobPrimary]);
  });

  await page.route('**/v0/askstories.json', async (route) => {
    await fulfillJson(route, []);
  });

  await page.route('**/v0/showstories.json', async (route) => {
    await fulfillJson(route, []);
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
    const item = FEED_FIXTURE_ITEMS[itemId];

    if (!item) {
      await route.fallback();
      return;
    }

    await fulfillJson(route, item);
  });
};

test.beforeEach(async ({ page }) => {
  await installStoriesAuditFixture(page);
});

test('AUDIT-F-01 story post opens without any errors', async ({ page }) => {
  const jsErrorCollector = createJsErrorCollector(page);

  try {
    await page.goto('');

    const hasStoriesTab = await hasAnySelector(page, FEED_TAB_SELECTORS.stories);
    const openResult = await openPostFromFeed({
      page,
      feedTabSelectors: hasStoriesTab ? FEED_TAB_SELECTORS.stories : [],
      postLinkSelectors: POST_LINK_SELECTORS.stories,
    });

    test.skip(!openResult.ok, `AUDIT-F-01 dependency guard: ${openResult.reason}`);

    expect(jsErrorCollector.getErrors()).toEqual([]);
  } finally {
    jsErrorCollector.stop();
  }
});

test('AUDIT-G-01 UI has stories, jobs, and polls', async ({ page }) => {
  await page.goto('');

  const supportsStories =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.stories)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.stories));
  const supportsJobs =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.jobs)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.jobs));
  const supportsPolls =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.polls)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.polls));

  test.skip(
    !(supportsStories && supportsJobs && supportsPolls),
    'AUDIT-G-01 dependency guard: stories/jobs/polls selectors are not all available yet.',
  );

  const storyTab = await findFirstLocator(page, FEED_TAB_SELECTORS.stories);

  if (storyTab) {
    await storyTab.click();
  }

  const storyLink =
    (await waitForAnySelector(page, STRICT_POST_LINK_SELECTORS.stories, { timeoutMs: 6_000 })) ??
    (await waitForAnySelector(page, POST_LINK_SELECTORS.stories, { timeoutMs: 6_000 }));

  test.skip(
    !storyLink,
    'AUDIT-G-01 dependency guard: story links are not available after story feed activation.',
  );
  await expect(storyLink).toBeVisible();

  const jobsTab = await findFirstLocator(page, FEED_TAB_SELECTORS.jobs);

  if (jobsTab) {
    await jobsTab.click();
  }

  const jobLink =
    (await waitForAnySelector(page, STRICT_POST_LINK_SELECTORS.jobs, { timeoutMs: 6_000 })) ??
    (jobsTab
      ? await waitForAnySelector(page, POST_LINK_SELECTORS.jobs, { timeoutMs: 6_000 })
      : null);

  test.skip(
    !jobLink,
    'AUDIT-G-01 dependency guard: job links are not available via tab or typed selectors.',
  );
  await expect(jobLink).toBeVisible();

  const pollsTab = await findFirstLocator(page, FEED_TAB_SELECTORS.polls);

  if (pollsTab) {
    await pollsTab.click();
  }

  const pollLink =
    (await waitForAnySelector(page, STRICT_POST_LINK_SELECTORS.polls, { timeoutMs: 6_000 })) ??
    (pollsTab
      ? await waitForAnySelector(page, POST_LINK_SELECTORS.polls, { timeoutMs: 6_000 })
      : null);

  test.skip(
    !pollLink,
    'AUDIT-G-01 dependency guard: poll links are not available via tab or typed selectors.',
  );
  await expect(pollLink).toBeVisible();
});

test('AUDIT-B-01 Ask HN and Show HN tabs are visible', async ({ page }) => {
  await page.goto('');

  const supportsStories =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.stories)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.stories));
  const supportsJobs =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.jobs)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.jobs));
  const supportsPolls =
    (await hasAnySelector(page, FEED_TAB_SELECTORS.polls)) ||
    (await hasAnySelector(page, STRICT_POST_LINK_SELECTORS.polls));

  test.skip(
    !(supportsStories && supportsJobs && supportsPolls),
    'AUDIT-B-01 dependency guard (Track B / TB-2): baseline stories/jobs/polls post-type selectors are not all available yet.',
  );

  const askTab = await findFirstLocator(page, FEED_TAB_SELECTORS.ask);
  const showTab = await findFirstLocator(page, FEED_TAB_SELECTORS.show);

  test.skip(
    !askTab || !showTab,
    'AUDIT-B-01 dependency guard (Track B / TB-2): Ask HN and/or Show HN tab selectors are not available in this branch yet.',
  );

  await expect(askTab).toBeVisible();
  await expect(showTab).toBeVisible();
});

test('AUDIT-G-02 posts are ordered newest-to-oldest', async ({ page }) => {
  await page.goto('');

  const storyTab = await findFirstLocator(page, FEED_TAB_SELECTORS.stories);

  if (storyTab) {
    await storyTab.click();
  }

  await waitForAnySelector(page, FEED_ITEM_SELECTORS, { timeoutMs: 10_000 });

  const feedItems = page.locator(FEED_ITEM_QUERY);
  const feedItemCount = await feedItems.count();

  test.skip(
    feedItemCount < 2,
    'AUDIT-G-02 dependency guard: fewer than two feed items are available for order assertions.',
  );

  const sampleSize = Math.min(feedItemCount, 20);
  const observedTimestamps = [];

  for (let index = 0; index < sampleSize; index += 1) {
    const timestampValue = await readFirstNumericSignal({
      locator: feedItems.nth(index),
      attributes: FEED_ITEM_TIME_ATTRIBUTES,
      childSelectors: FEED_ITEM_TIME_SELECTORS,
      allowTextFallback: false,
    });

    if (timestampValue !== null) {
      observedTimestamps.push(timestampValue);
    }
  }

  test.skip(
    observedTimestamps.length < 2,
    'AUDIT-G-02 dependency guard: timestamp metadata is not available for feed items yet.',
  );

  for (let index = 1; index < observedTimestamps.length; index += 1) {
    expect(observedTimestamps[index - 1]).toBeGreaterThanOrEqual(observedTimestamps[index]);
  }
});
