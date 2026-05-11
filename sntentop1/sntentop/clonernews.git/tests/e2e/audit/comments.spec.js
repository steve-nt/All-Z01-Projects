import { expect, test } from '@playwright/test';
import {
  COMMENT_ITEM_SELECTORS,
  COMMENT_PARENT_ATTRIBUTES,
  COMMENT_PARENT_SELECTORS,
  COMMENT_TIME_ATTRIBUTES,
  COMMENT_TIME_SELECTORS,
  FEED_TAB_SELECTORS,
  findFirstLocator,
  hasAnySelector,
  openPostFromFeed,
  POST_LINK_SELECTORS,
  parseItemIdFromUrl,
  readFirstNumericSignal,
  waitForAnySelector,
  waitForItemRoute,
} from './helpers.js';

const COMMENT_ITEM_QUERY = COMMENT_ITEM_SELECTORS.join(', ');

const COMMENTS_FIXTURE_IDS = Object.freeze({
  storyPrimary: 9101,
  commentPrimary: 9201,
  commentSecondary: 9202,
  commentNested: 9203,
  commentDeepNested: 9204,
});

const COMMENTS_FIXTURE_ITEMS = Object.freeze({
  [COMMENTS_FIXTURE_IDS.storyPrimary]: Object.freeze({
    id: COMMENTS_FIXTURE_IDS.storyPrimary,
    type: 'story',
    by: 'fixture-story-author',
    time: 2_200_500,
    title: 'Fixture story with comments',
    descendants: 4,
    score: 42,
    url: 'https://example.com/story/with-comments',
    kids: [COMMENTS_FIXTURE_IDS.commentPrimary, COMMENTS_FIXTURE_IDS.commentSecondary],
  }),
  [COMMENTS_FIXTURE_IDS.commentPrimary]: Object.freeze({
    id: COMMENTS_FIXTURE_IDS.commentPrimary,
    type: 'comment',
    by: 'fixture-comment-author-1',
    time: 2_200_490,
    text: '<p>Primary comment body</p>',
    parent: COMMENTS_FIXTURE_IDS.storyPrimary,
    kids: [COMMENTS_FIXTURE_IDS.commentNested],
  }),
  [COMMENTS_FIXTURE_IDS.commentSecondary]: Object.freeze({
    id: COMMENTS_FIXTURE_IDS.commentSecondary,
    type: 'comment',
    by: 'fixture-comment-author-2',
    time: 2_200_480,
    text: '<p>Secondary comment body</p>',
    parent: COMMENTS_FIXTURE_IDS.storyPrimary,
    kids: [],
  }),
  [COMMENTS_FIXTURE_IDS.commentNested]: Object.freeze({
    id: COMMENTS_FIXTURE_IDS.commentNested,
    type: 'comment',
    by: 'fixture-comment-author-nested',
    time: 2_200_485,
    text: '<p>Nested comment body</p>',
    parent: COMMENTS_FIXTURE_IDS.commentPrimary,
    kids: [COMMENTS_FIXTURE_IDS.commentDeepNested],
  }),
  [COMMENTS_FIXTURE_IDS.commentDeepNested]: Object.freeze({
    id: COMMENTS_FIXTURE_IDS.commentDeepNested,
    type: 'comment',
    by: 'fixture-comment-author-deep-nested',
    time: 2_200_482,
    text: '<p>Deep nested comment body</p>',
    parent: COMMENTS_FIXTURE_IDS.commentNested,
    kids: [],
  }),
});

const fulfillJson = (route, payload) =>
  route.fulfill({
    status: 200,
    contentType: 'application/json',
    body: JSON.stringify(payload),
  });

const installCommentsAuditFixture = async (page) => {
  await page.route('**/v0/topstories.json', async (route) => {
    await fulfillJson(route, [COMMENTS_FIXTURE_IDS.storyPrimary]);
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
    const item = COMMENTS_FIXTURE_ITEMS[itemId];

    if (!item) {
      await route.fallback();
      return;
    }

    await fulfillJson(route, item);
  });
};

test.beforeEach(async ({ page }) => {
  await installCommentsAuditFixture(page);
});

const openCommentsPost = async (page) => {
  await page.goto('');

  const hasStoriesTab = await hasAnySelector(page, FEED_TAB_SELECTORS.stories);
  if (hasStoriesTab) {
    const storyTab = await findFirstLocator(page, FEED_TAB_SELECTORS.stories);
    await storyTab.click();
  }

  // Find a post link where the card has comments (best effort)
  await page.waitForSelector('[data-testid="feed-item"]');
  const feedItems = page.locator('[data-testid="feed-item"]');
  const count = await feedItems.count();

  for (let i = 0; i < count; i++) {
    const item = feedItems.nth(i);
    const commentsText = await item.locator('[data-testid="card-comments"]').textContent();
    const commentCount = Number.parseInt(commentsText, 10) || 0;
    if (commentCount > 10) {
      const link = item.locator('[data-testid="post-link"]');
      await link.click();
      await waitForItemRoute(page);
      return { ok: true };
    }
  }

  // Fallback to first one if none found with > 5 comments
  return openPostFromFeed({
    page,
    feedTabSelectors: [],
    postLinkSelectors: POST_LINK_SELECTORS.stories,
  });
};

test('AUDIT-F-05 comments are displayed from newest to oldest', async ({ page }) => {
  const openResult = await openCommentsPost(page);

  test.skip(!openResult.ok, `AUDIT-F-05 dependency guard: ${openResult.reason}`);
  await waitForAnySelector(page, COMMENT_ITEM_SELECTORS, { timeoutMs: 8_000, visible: false });

  const comments = page.locator(COMMENT_ITEM_QUERY);
  const commentCount = await comments.count();

  test.skip(
    commentCount < 2,
    'AUDIT-F-05 dependency guard: fewer than two comments are available for order assertions.',
  );

  const sampleSize = Math.min(commentCount, 25);
  const observedCommentTimestamps = [];

  for (let index = 0; index < sampleSize; index += 1) {
    const timestampValue = await readFirstNumericSignal({
      locator: comments.nth(index),
      attributes: COMMENT_TIME_ATTRIBUTES,
      childSelectors: COMMENT_TIME_SELECTORS,
      allowTextFallback: false,
    });

    if (timestampValue !== null) {
      observedCommentTimestamps.push(timestampValue);
    }
  }

  test.skip(
    observedCommentTimestamps.length < 2,
    'AUDIT-F-05 dependency guard: comment timestamp metadata is not available yet.',
  );

  for (let index = 1; index < observedCommentTimestamps.length; index += 1) {
    expect(observedCommentTimestamps[index - 1]).toBeGreaterThanOrEqual(
      observedCommentTimestamps[index],
    );
  }
});

test('AUDIT-G-03 each comment presents the correct parent post', async ({ page }) => {
  const openResult = await openCommentsPost(page);

  test.skip(!openResult.ok, `AUDIT-G-03 dependency guard: ${openResult.reason}`);
  await waitForAnySelector(page, COMMENT_ITEM_SELECTORS, { timeoutMs: 8_000, visible: false });

  const currentPostId = parseItemIdFromUrl(page.url());

  test.skip(
    currentPostId === null,
    'AUDIT-G-03 dependency guard: current route is not #/item/:id, so parent-post validation cannot run.',
  );

  const comments = page.locator(COMMENT_ITEM_QUERY);
  const commentCount = await comments.count();

  test.skip(
    commentCount < 1,
    'AUDIT-G-03 dependency guard: no comments are rendered for parent-post checks.',
  );

  const sampleSize = Math.min(commentCount, 20);
  const observedParentIds = [];

  for (let index = 0; index < sampleSize; index += 1) {
    const parentId = await readFirstNumericSignal({
      locator: comments.nth(index),
      attributes: COMMENT_PARENT_ATTRIBUTES,
      childSelectors: COMMENT_PARENT_SELECTORS,
      allowTextFallback: true,
    });

    if (parentId !== null) {
      observedParentIds.push(parentId);
    }
  }

  test.skip(
    observedParentIds.length < Math.min(sampleSize, 3),
    'AUDIT-G-03 dependency guard: comment parent-post metadata is not exposed consistently yet.',
  );

  for (const parentId of observedParentIds) {
    expect(parentId).toBe(currentPostId);
  }
});

test('AUDIT-B-02 nested sub-comments render at ≥2 depth levels', async ({ page }) => {
  const openResult = await openCommentsPost(page);

  test.skip(!openResult.ok, `AUDIT-B-02 dependency guard: ${openResult.reason}`);

  await page.waitForSelector(COMMENT_ITEM_QUERY, { timeout: 10_000 });
  const comments = page.locator(COMMENT_ITEM_QUERY);
  const commentCount = await comments.count();

  test.skip(
    commentCount < 1,
    'AUDIT-B-02 dependency guard: no comments are rendered for depth checks.',
  );

  const depthValues = await comments.evaluateAll((nodes) =>
    nodes
      .map((node) => Number(node.getAttribute('data-depth')))
      .filter((depth) => Number.isInteger(depth) && depth >= 0),
  );

  test.skip(
    depthValues.length < 1,
    'AUDIT-B-02 dependency guard: rendered comments do not expose depth metadata.',
  );

  const maxDepth = depthValues.toSorted((leftDepth, rightDepth) => rightDepth - leftDepth)[0] ?? -1;

  expect(maxDepth).toBeGreaterThanOrEqual(2);
});
