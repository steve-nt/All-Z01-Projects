import { Temporal } from '@js-temporal/polyfill';

export const ITEM_ROUTE_PATTERN = /#\/item\/(\d+)/;

export const FEED_TAB_SELECTORS = Object.freeze({
  stories: Object.freeze([
    '[data-testid="feed-tab-stories"]',
    '[data-testid="feed-tab-story"]',
    '[data-testid="feed-tab-top"]',
    '[data-testid="tab-stories"]',
  ]),
  jobs: Object.freeze([
    '[data-testid="feed-tab-jobs"]',
    '[data-testid="feed-tab-job"]',
    '[data-testid="tab-jobs"]',
  ]),
  polls: Object.freeze([
    '[data-testid="feed-tab-polls"]',
    '[data-testid="feed-tab-poll"]',
    '[data-testid="tab-polls"]',
  ]),
  ask: Object.freeze(['[data-testid="feed-tab-ask"]', '[data-testid="tab-ask"]']),
  show: Object.freeze(['[data-testid="feed-tab-show"]', '[data-testid="tab-show"]']),
});

export const POST_LINK_SELECTORS = Object.freeze({
  stories: Object.freeze([
    '[data-testid="story-link"]',
    '[data-testid="feed-item-story-link"]',
    '[data-testid="feed-story-link"]',
    '[data-testid="story-item"] [data-testid="post-link"]',
    '[data-testid="feed-item"][data-item-type="story"] [data-testid="post-link"]',
    '[data-testid="feed-item"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="story"]',
    'main a[href^="#/item/"]',
  ]),
  jobs: Object.freeze([
    '[data-testid="job-link"]',
    '[data-testid="feed-item-job-link"]',
    '[data-testid="job-item"] [data-testid="post-link"]',
    '[data-testid="feed-item"][data-item-type="job"] [data-testid="post-link"]',
    '[data-testid="feed-item"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="job"]',
    'main a[href^="#/item/"]',
  ]),
  polls: Object.freeze([
    '[data-testid="poll-link"]',
    '[data-testid="feed-item-poll-link"]',
    '[data-testid="poll-item"] [data-testid="post-link"]',
    '[data-testid="feed-item"][data-item-type="poll"] [data-testid="post-link"]',
    '[data-testid="feed-item"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="poll"]',
    'main a[href^="#/item/"]',
  ]),
});

export const STRICT_POST_LINK_SELECTORS = Object.freeze({
  stories: Object.freeze([
    '[data-testid="story-link"]',
    '[data-testid="feed-item-story-link"]',
    '[data-testid="feed-story-link"]',
    '[data-testid="feed-item"][data-item-type="story"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="story"]',
  ]),
  jobs: Object.freeze([
    '[data-testid="job-link"]',
    '[data-testid="feed-item-job-link"]',
    '[data-testid="feed-item"][data-item-type="job"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="job"]',
  ]),
  polls: Object.freeze([
    '[data-testid="poll-link"]',
    '[data-testid="feed-item-poll-link"]',
    '[data-testid="feed-item"][data-item-type="poll"] [data-testid="post-link"]',
    'a[data-testid="post-link"][data-item-type="poll"]',
  ]),
});

export const FEED_ITEM_SELECTORS = Object.freeze([
  '[data-testid="feed-item"]',
  '[data-testid="story-item"]',
  '[data-testid="post-item"]',
  '[data-testid="post-row"]',
]);

export const FEED_ITEM_TIME_SELECTORS = Object.freeze([
  '[data-testid="post-time"]',
  '[data-testid="feed-item-time"]',
  '[data-testid="item-time"]',
  '[data-testid="time"]',
  'time',
]);

export const FEED_ITEM_TIME_ATTRIBUTES = Object.freeze([
  'data-time',
  'data-timestamp',
  'data-unix-time',
  'data-created-at',
  'datetime',
]);

export const LOAD_MORE_SELECTORS = Object.freeze([
  '[data-testid="load-more"]',
  '[data-testid="feed-load-more"]',
  '[data-testid="load-more-button"]',
]);

export const LOAD_MORE_SENTINEL_SELECTORS = Object.freeze([
  '[data-testid="load-more-sentinel"]',
  '[data-testid="feed-sentinel"]',
  '[data-testid="pagination-sentinel"]',
]);

export const COMMENT_ITEM_SELECTORS = Object.freeze([
  '[data-testid="comment-item"]',
  '[data-testid="comment-node"]',
  '[data-testid="comment"]',
]);

export const COMMENT_TIME_SELECTORS = Object.freeze([
  '[data-testid="comment-time"]',
  '[data-testid="time"]',
  'time',
]);

export const COMMENT_TIME_ATTRIBUTES = Object.freeze([
  'data-time',
  'data-timestamp',
  'data-unix-time',
  'data-created-at',
  'datetime',
]);

export const COMMENT_PARENT_SELECTORS = Object.freeze([
  '[data-testid="comment-parent-post-id"]',
  '[data-testid="comment-parent-id"]',
  '[data-testid="comment-parent"]',
]);

export const COMMENT_PARENT_ATTRIBUTES = Object.freeze([
  'data-parent-post-id',
  'data-parent-id',
  'data-parent',
]);

export const LIVE_BANNER_SELECTORS = Object.freeze([
  '[data-testid="live-banner"]',
  '[data-testid="live-updates-banner"]',
  '[data-testid="updates-banner"]',
]);

export const LIVE_BANNER_ACTION_SELECTORS = Object.freeze([
  '[data-testid="live-banner-refresh"]',
  '[data-testid="live-banner-action"]',
  '[data-testid="live-banner-button"]',
]);

const asSelectorList = (selectors) => (Array.isArray(selectors) ? selectors : [selectors]);

export const createJsErrorCollector = (page) => {
  const errors = [];

  const onConsoleMessage = (message) => {
    if (message.type() === 'error') {
      errors.push(`console: ${message.text()}`);
    }
  };

  const onPageError = (error) => {
    errors.push(`pageerror: ${error.message}`);
  };

  page.on('console', onConsoleMessage);
  page.on('pageerror', onPageError);

  return {
    getErrors: () => [...errors],
    stop: () => {
      page.off('console', onConsoleMessage);
      page.off('pageerror', onPageError);
    },
  };
};

export const createRequestTracker = (page, predicate = () => true) => {
  const entries = [];

  const onRequest = (request) => {
    if (!predicate(request)) {
      return;
    }

    entries.push({
      url: request.url(),
      method: request.method(),
      startedAtMs: performance.now(),
    });
  };

  page.on('request', onRequest);

  return {
    entries,
    stop: () => {
      page.off('request', onRequest);
    },
  };
};

export const findFirstLocator = async (page, selectors, { visible = true } = {}) => {
  for (const selector of asSelectorList(selectors)) {
    const locator = page.locator(selector).first();
    const count = await locator.count();

    if (count < 1) {
      continue;
    }

    if (visible && !(await locator.isVisible())) {
      continue;
    }

    return locator;
  }

  return null;
};

export const waitForAnySelector = async (
  page,
  selectors,
  { visible = true, timeoutMs = 8_000, pollIntervalMs = 200 } = {},
) => {
  const startedAtMs = performance.now();
  const selectorList = asSelectorList(selectors);

  while (performance.now() - startedAtMs < timeoutMs) {
    const locator = await findFirstLocator(page, selectorList, { visible });

    if (locator) {
      return locator;
    }

    await page.waitForTimeout(pollIntervalMs);
  }

  return null;
};

export const hasAnySelector = async (page, selectors, options) => {
  const locator = await findFirstLocator(page, selectors, options);

  return locator !== null;
};

export const waitForItemRoute = async (page, timeoutMs = 10_000) => {
  try {
    await page.waitForURL(ITEM_ROUTE_PATTERN, { timeout: timeoutMs });

    return true;
  } catch {
    return false;
  }
};

export const openPostFromFeed = async ({ page, feedTabSelectors = [], postLinkSelectors }) => {
  if (feedTabSelectors.length > 0) {
    const feedTab = await waitForAnySelector(page, feedTabSelectors, {
      timeoutMs: 8_000,
    });

    if (!feedTab) {
      return {
        ok: false,
        reason: 'Required feed tab selectors are not available in this branch yet.',
      };
    }

    await feedTab.click();
  }

  const postLink = await waitForAnySelector(page, postLinkSelectors, {
    timeoutMs: 10_000,
  });

  if (!postLink) {
    return {
      ok: false,
      reason: 'Required post link selectors are not available in this branch yet.',
    };
  }

  await postLink.click();

  const hasItemRoute = await waitForItemRoute(page);

  if (!hasItemRoute) {
    return {
      ok: false,
      reason: 'The detail route #/item/:id is not available in this branch yet.',
    };
  }

  return { ok: true };
};

export const extractFirstInteger = (value) => {
  if (typeof value !== 'string') {
    return null;
  }

  const trimmedValue = value.trim();

  if (trimmedValue === '') {
    return null;
  }

  if (/^-?\d+$/.test(trimmedValue)) {
    const fullInteger = Number(trimmedValue);

    return Number.isFinite(fullInteger) ? fullInteger : null;
  }

  try {
    return Number(Temporal.Instant.from(trimmedValue).epochSeconds);
  } catch {
    // Ignore non-Temporal values and continue with best-effort integer extraction.
  }

  const match = trimmedValue.match(/-?\d+/);

  if (!match) {
    return null;
  }

  const parsed = Number(match[0]);

  return Number.isFinite(parsed) ? parsed : null;
};

export const readFirstNumericSignal = async ({
  locator,
  attributes,
  childSelectors,
  allowTextFallback = true,
}) => {
  for (const attributeName of attributes) {
    const attributeValue = await locator.getAttribute(attributeName);
    const parsedAttributeValue = extractFirstInteger(attributeValue);

    if (parsedAttributeValue !== null) {
      return parsedAttributeValue;
    }
  }

  for (const childSelector of childSelectors) {
    const childLocator = locator.locator(childSelector).first();

    if ((await childLocator.count()) < 1) {
      continue;
    }

    for (const attributeName of attributes) {
      const childAttributeValue = await childLocator.getAttribute(attributeName);
      const parsedChildAttributeValue = extractFirstInteger(childAttributeValue);

      if (parsedChildAttributeValue !== null) {
        return parsedChildAttributeValue;
      }
    }

    if (allowTextFallback) {
      const childText = await childLocator.textContent();
      const parsedChildText = extractFirstInteger(childText);

      if (parsedChildText !== null) {
        return parsedChildText;
      }
    }
  }

  if (!allowTextFallback) {
    return null;
  }

  const locatorText = await locator.textContent();

  return extractFirstInteger(locatorText);
};

export const duplicateValues = (values) => {
  const seenValues = new Set();
  const duplicateValueSet = new Set();

  for (const value of values) {
    if (seenValues.has(value)) {
      duplicateValueSet.add(value);
    }

    seenValues.add(value);
  }

  return [...duplicateValueSet].toSorted();
};

export const intervalSeriesMs = (entries) =>
  entries.slice(1).map((entry, index) => entry.startedAtMs - entries[index].startedAtMs);

export const parseItemIdFromUrl = (url) => {
  const match = ITEM_ROUTE_PATTERN.exec(url);

  if (!match) {
    return null;
  }

  const parsedId = Number(match[1]);

  return Number.isInteger(parsedId) ? parsedId : null;
};
