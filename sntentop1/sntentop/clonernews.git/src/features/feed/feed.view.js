/*
 * Purpose: Render feed UI with tab navigation, story cards, skeleton loaders, and infinite scroll.
 * Public API: createFeedView({ container, controller }) mounts and manages DOM updates.
 * Implementation notes: IntersectionObserver tracks sentinel element for page loads; uses DocumentFragment
 * for batch DOM insertions; observer cleanup prevents memory leaks; loading flag prevents duplicate requests.
 */

import './feed.css';
import { appendChildren, clearElement, createElement } from '../../shared/dom-helpers.js';
import { formatRelativeTime } from '../../shared/time-format.js';

/**
 * @typedef {Object} FeedViewDependencies
 * @property {HTMLElement} container
 * @property {Object} controller
 */

/**
 * @typedef {Object} FeedView
 * @property {() => void} mount
 * @property {() => void} unmount
 * @property {(tabType: string) => Promise<void>} selectTab
 * @property {(updatedItemIds?: readonly number[], options?: { isMarkerReplay?: boolean }) => Promise<{ updatedItemIds: readonly number[] }>} applyRefreshMarkers
 * @property {() => void} clearRefreshMarkers
 */

const FEED_TYPES = Object.freeze([
  { id: 'top', label: 'Top' },
  { id: 'new', label: 'New' },
  { id: 'job', label: 'Jobs' },
  { id: 'poll', label: 'Polls' },
  { id: 'ask', label: 'Ask HN' },
  { id: 'show', label: 'Show HN' },
]);

const FEED_TYPE_LABELS = Object.freeze(
  Object.fromEntries(FEED_TYPES.map((feedType) => [feedType.id, feedType.label])),
);

const TAB_TEST_IDS_BY_FEED_TYPE = Object.freeze({
  top: 'feed-tab-top',
  new: 'feed-tab-new',
  job: 'feed-tab-jobs',
  poll: 'feed-tab-polls',
  ask: 'feed-tab-ask',
  show: 'feed-tab-show',
});

// Protocol whitelist for safe external links.
const SAFE_EXTERNAL_PROTOCOLS = new Set(['http:', 'https:']);

/**
 * Normalize an external href and reject non-http(s) protocols.
 * @param {unknown} rawUrl
 * @returns {string | null}
 */
const toSafeExternalHref = (rawUrl) => {
  if (typeof rawUrl !== 'string') {
    return null;
  }

  try {
    const parsedUrl = new URL(rawUrl);
    return SAFE_EXTERNAL_PROTOCOLS.has(parsedUrl.protocol) ? parsedUrl.href : null;
  } catch {
    return null;
  }
};

const toFeedItemType = (rawType) =>
  typeof rawType === 'string' && rawType.length > 0 ? rawType : 'story';

const toDetailHref = (itemId) => `#/item/${itemId}`;

const toTabTestId = (feedTypeId) => TAB_TEST_IDS_BY_FEED_TYPE[feedTypeId] ?? `tab-${feedTypeId}`;

const toFeedLabel = (feedTypeId) => FEED_TYPE_LABELS[feedTypeId] ?? 'Feed';
const toValidItemId = (value) => (Number.isInteger(value) && value > 0 ? value : null);

const toRefreshMarkerByItemId = (items, updatedItemIds) => {
  // Marker rendering is only meaningful when refreshed IDs are available for the current feed view.
  if (updatedItemIds.size === 0) {
    return new Map();
  }

  const markerByItemId = new Map();

  for (const item of items) {
    const itemId = toValidItemId(item?.id);

    if (itemId === null) {
      continue;
    }

    // Only refreshed IDs receive a badge marker.
    const markerKind = updatedItemIds.has(itemId) ? 'updated' : undefined;

    if (markerKind === undefined) {
      continue;
    }

    markerByItemId.set(itemId, markerKind);
  }

  return markerByItemId;
};

/**
 * Create a skeleton card placeholder for loading state.
 * @returns {HTMLElement}
 */
const createSkeletonCard = () => {
  return createElement('div', {
    className: 'feed-card feed-card--skeleton',
    attributes: { 'data-testid': 'skeleton-card' },
  });
};

/**
 * Create an animated loading panel that sits above skeleton placeholders.
 * @param {string} feedTypeId
 * @returns {HTMLElement}
 */
const createLoadingScreen = (feedTypeId) => {
  const feedLabel = toFeedLabel(feedTypeId);

  return createElement('section', {
    className: 'feed__loading-screen',
    attributes: {
      'data-testid': 'feed-loading-screen',
      'aria-live': 'polite',
    },
    children: [
      createElement('p', {
        className: 'feed__loading-eyebrow',
        text: `${feedLabel} stream selected`,
      }),
      createElement('h3', {
        className: 'feed__loading-title',
        text: `Loading ${feedLabel} posts`,
      }),
      createElement('p', {
        className: 'feed__loading-copy',
        text: 'Curating fresh stories for this feed type.',
      }),
      createElement('div', {
        className: 'feed__loading-orbit',
        attributes: {
          'aria-hidden': 'true',
        },
        children: [
          createElement('span', {
            className: 'feed__loading-orbit-ring',
          }),
          createElement('span', {
            className: 'feed__loading-orbit-dot',
          }),
        ],
      }),
    ],
  });
};

/**
 * Create a compact loading card that appears while the next page is being fetched.
 * @returns {HTMLElement}
 */
const createLoadMoreCard = () => {
  return createElement('article', {
    className: 'feed-card feed-card--loading-more',
    attributes: {
      'data-testid': 'feed-loading-more-card',
      'aria-live': 'polite',
    },
    children: [
      createElement('p', {
        className: 'feed-card__loading-more-copy',
        text: 'Loading more posts...',
      }),
    ],
  });
};

/**
 * Create a story/job card from an HN item.
 * @param {Object} item - HN item with id, title, score, by, time, descendants
 * @returns {HTMLElement}
 */
const createStoryCard = (item, refreshMarkerKind, shouldAnimateRefreshHighlight = false) => {
  const timeAgo = formatRelativeTime(item.time);
  const commentCount = item.descendants ?? 0;
  const score = item.score ?? 0;
  const author = item.by ?? 'Anonymous';
  const feedItemType = toFeedItemType(item.type);
  const safeExternalHref = toSafeExternalHref(item.url);

  // Build the card structure with safe DOM methods.
  const card = createElement('article', {
    className: 'feed-card',
    attributes: {
      'data-testid': `feed-card-${item.id}`,
      'data-item-id': String(item.id),
    },
  });

  if (shouldAnimateRefreshHighlight && refreshMarkerKind === 'updated') {
    card.classList.add('feed-card--refresh-highlight');
  }

  if (refreshMarkerKind === 'updated') {
    const badgeText = '~';
    const badgeLabel = 'Refreshed post since last refresh';

    // Right-side status badge makes refresh impact scannable without changing card hierarchy.
    card.append(
      createElement('span', {
        className: `feed-card__refresh-marker feed-card__refresh-marker--${refreshMarkerKind}`,
        text: badgeText,
        attributes: {
          'data-testid': `feed-card-refresh-marker-${item.id}`,
          'aria-label': badgeLabel,
          title: badgeLabel,
        },
      }),
    );
  }

  const feedItem = createElement('div', {
    className: 'feed-card__content',
    attributes: {
      'data-testid': 'feed-item',
      'data-item-id': String(item.id),
      'data-item-type': feedItemType,
      'data-time': Number.isInteger(item.time) ? String(item.time) : '',
    },
  });

  // The detail link is the primary CTA for audit flows and app routing.
  const detailLink = createElement('a', {
    text: item.title,
    attributes: {
      href: toDetailHref(item.id),
      'data-testid': 'post-link',
      'data-item-type': feedItemType,
    },
  });

  const titleEl = createElement('h3', {
    className: 'feed-card__title',
    children: [detailLink],
  });

  let sourceLink = null;

  if (safeExternalHref) {
    sourceLink = createElement('a', {
      className: 'feed-card__source-link',
      text: 'source',
      attributes: {
        href: safeExternalHref,
        rel: 'noopener noreferrer',
        target: '_blank',
        'data-testid': 'card-title-link',
      },
    });
  }

  // Metadata row: score, author, time, comment count.
  const metadata = createElement('div', {
    className: 'feed-card__meta',
  });

  const scoreEl = createElement('span', {
    className: 'feed-card__score',
    text: `${score} points`,
    attributes: { 'data-testid': 'card-score' },
  });

  const authorEl = createElement('span', {
    className: 'feed-card__author',
    text: `by ${author}`,
    attributes: { 'data-testid': 'card-author' },
  });

  const timeEl = createElement('time', {
    className: 'feed-card__time',
    text: timeAgo,
    attributes: {
      'data-testid': 'post-time',
      'data-time': Number.isInteger(item.time) ? String(item.time) : '',
    },
  });

  const commentsEl = createElement('span', {
    className: 'feed-card__comments',
    text: `${commentCount} comment${commentCount !== 1 ? 's' : ''}`,
    attributes: { 'data-testid': 'card-comments' },
  });

  appendChildren(metadata, [scoreEl, authorEl, timeEl, commentsEl]);

  if (sourceLink) {
    appendChildren(feedItem, [titleEl, sourceLink, metadata]);
  } else {
    appendChildren(feedItem, [titleEl, metadata]);
  }

  appendChildren(card, [feedItem]);

  return card;
};

/**
 * Create the fallback tab navigation bar when shell-level tabs are unavailable.
 * @returns {{ nav: HTMLElement, tabButtons: Array<{ feedTypeId: string, button: HTMLButtonElement }> }}
 */
const createTabNavigation = () => {
  const nav = createElement('nav', {
    className: 'feed__tabs',
    attributes: { 'data-testid': 'feed-tabs' },
  });

  const tabButtons = FEED_TYPES.map((feedType) => {
    const button = createElement('button', {
      className: 'feed__tab-button',
      text: feedType.label,
      attributes: {
        'data-testid': toTabTestId(feedType.id),
        'data-feed-type': feedType.id,
        type: 'button',
      },
    });

    return {
      feedTypeId: feedType.id,
      button,
    };
  });

  appendChildren(
    nav,
    tabButtons.map((entry) => entry.button),
  );

  return {
    nav,
    tabButtons,
  };
};

/**
 * Create the items list container.
 * @returns {HTMLElement}
 */
const createItemsList = () => {
  return createElement('div', {
    className: 'feed__list',
    attributes: { 'data-testid': 'feed-list' },
  });
};

/**
 * Create hidden sentinel element for IntersectionObserver.
 * @returns {HTMLElement}
 */
const createSentinel = () => {
  return createElement('div', {
    className: 'feed__sentinel',
    attributes: {
      'data-testid': 'feed-sentinel',
      'aria-hidden': 'true',
    },
  });
};

/**
 * Create floating action button that jumps to the top of the page.
 * @returns {HTMLButtonElement}
 */
const createScrollTopButton = () => {
  return createElement('button', {
    className: 'feed__scroll-top-button',
    text: 'Top',
    attributes: {
      type: 'button',
      'data-testid': 'feed-scroll-top-button',
      'aria-label': 'Back to top',
      title: 'Back to top',
    },
  });
};

/**
 * Render loading skeleton cards.
 * @param {HTMLElement} listContainer
 * @returns {void}
 */
const renderSkeletons = (listContainer, feedTypeId = 'top') => {
  clearElement(listContainer);

  const loadingScreen = createLoadingScreen(feedTypeId);
  const skeletons = Array.from({ length: 10 }, createSkeletonCard);

  appendChildren(listContainer, [loadingScreen, ...skeletons]);
};

/**
 * Render story cards from items array.
 * @param {HTMLElement} listContainer
 * @param {Array<Object>} items
 * @returns {void}
 */
const renderStoryCards = (
  listContainer,
  items,
  refreshMarkerByItemId = new Map(),
  shouldAnimateRefreshHighlight = false,
) => {
  clearElement(listContainer);
  const cards = items.map((item) => {
    const itemId = toValidItemId(item?.id);
    return createStoryCard(
      item,
      itemId === null ? undefined : refreshMarkerByItemId.get(itemId),
      shouldAnimateRefreshHighlight,
    );
  });
  appendChildren(listContainer, cards);
};

/**
 * Append story cards to existing list (infinite scroll).
 * @param {HTMLElement} listContainer
 * @param {Array<Object>} items
 * @returns {void}
 */
const appendStoryCards = (
  listContainer,
  items,
  refreshMarkerByItemId = new Map(),
  shouldAnimateRefreshHighlight = false,
) => {
  // Use DocumentFragment for batch insert to avoid reflow thrashing.
  const fragment = document.createDocumentFragment();
  const cards = items.map((item) => {
    const itemId = toValidItemId(item?.id);
    return createStoryCard(
      item,
      itemId === null ? undefined : refreshMarkerByItemId.get(itemId),
      shouldAnimateRefreshHighlight,
    );
  });
  appendChildren(fragment, cards);
  listContainer.append(fragment);
};

/**
 * Show a loading card at the bottom of an already-rendered feed.
 * @param {HTMLElement} listContainer
 * @returns {void}
 */
const renderLoadMoreCard = (listContainer) => {
  const existingLoadMoreCard = listContainer.querySelector('[data-testid="feed-loading-more-card"]');

  if (existingLoadMoreCard !== null) {
    return;
  }

  listContainer.append(createLoadMoreCard());
};

/**
 * Remove the bottom loading card once pagination request settles.
 * @param {HTMLElement} listContainer
 * @returns {void}
 */
const clearLoadMoreCard = (listContainer) => {
  const existingLoadMoreCard = listContainer.querySelector('[data-testid="feed-loading-more-card"]');

  if (existingLoadMoreCard !== null) {
    existingLoadMoreCard.remove();
  }
};

/**
 * Show an end-of-feed card when pagination is exhausted.
 * @param {HTMLElement} listContainer
 * @returns {void}
 */
const renderEndOfFeedCard = (listContainer) => {
  const existingEndCard = listContainer.querySelector('[data-testid="feed-end-of-feed-card"]');

  if (existingEndCard !== null) {
    return;
  }

  listContainer.append(
    createElement('article', {
      className: 'feed-card feed-card--end-of-feed',
      attributes: {
        'data-testid': 'feed-end-of-feed-card',
        'aria-live': 'polite',
      },
      children: [
        createElement('p', {
          className: 'feed-card__end-of-feed-copy',
          text: 'No more posts to load.',
        }),
      ],
    }),
  );
};

/**
 * Remove the end-of-feed card whenever loading resumes or feed context changes.
 * @param {HTMLElement} listContainer
 * @returns {void}
 */
const clearEndOfFeedCard = (listContainer) => {
  const existingEndCard = listContainer.querySelector('[data-testid="feed-end-of-feed-card"]');

  if (existingEndCard !== null) {
    existingEndCard.remove();
  }
};

/**
 * Render error message.
 * @param {HTMLElement} listContainer
 * @param {string} errorMessage
 * @returns {void}
 */
const renderError = (listContainer, errorMessage) => {
  clearElement(listContainer);
  const errorEl = createElement('div', {
    className: 'feed__error',
    attributes: { 'data-testid': 'feed-error' },
    children: [
      createElement('p', {
        text: `Error loading feed: ${errorMessage}`,
      }),
      createElement('p', {
        className: 'feed__error-hint',
        text: 'Try clicking a tab to retry.',
      }),
    ],
  });
  listContainer.append(errorEl);
};

/**
 * Create the feed view with infinite scroll support.
 * @param {FeedViewDependencies} dependencies
 * @returns {FeedView}
 */
export const createFeedView = ({ container, controller }) => {
  if (!container) {
    throw new Error('Feed view requires a container element.');
  }

  if (!controller) {
    throw new Error('Feed view requires a controller.');
  }

  let listContainer;
  let sentinel;
  let intersectionObserver = null;
  let cleanupTabListeners = () => {};
  let cleanupScrollTopButtonListener = () => {};
  let mountedTabButtons = [];
  let activeTabType = 'top';
  let pendingTabType = null;
  let loadingTabType = 'top';
  let hasMorePages = true;
  let isLoadingMore = false;
  let shouldAnimateRefreshHighlights = false;
  let loadedPostsCount = 0;
  let currentRenderedItemIds = new Set();
  let refreshUpdatedItemIds = new Set();

  // Counter updates are scoped to this mounted view so switching routes or refresh resets naturally.
  const setLoadedPostsCount = (count) => {
    loadedPostsCount = Math.max(0, count);

    const loadedCountElement = document.querySelector('[data-testid="feed-loaded-count"]');

    if (!loadedCountElement) {
      return;
    }

    loadedCountElement.textContent = `Loaded: ${loadedPostsCount}`;
    loadedCountElement.setAttribute('data-count', String(loadedPostsCount));
  };

  // Keep tab switches predictable by always resetting viewport position before loading a new feed.
  const scrollViewportToTop = () => {
    if (typeof window === 'undefined') {
      return;
    }

    const isJsDomEnvironment =
      typeof window.navigator !== 'undefined' &&
      typeof window.navigator.userAgent === 'string' &&
      window.navigator.userAgent.toLowerCase().includes('jsdom');

    if (!isJsDomEnvironment && typeof window.scrollTo === 'function') {
      try {
        window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
      } catch {
        // Browsers that reject object-form scroll options still get the scrollTop fallback below.
      }
    }

    if (document.documentElement) {
      document.documentElement.scrollTop = 0;
    }

    if (document.body) {
      document.body.scrollTop = 0;
    }
  };

  // Reset helper keeps every non-refresh navigation path free of stale marker state.
  const resetRefreshVisualState = () => {
    refreshUpdatedItemIds = new Set();
    shouldAnimateRefreshHighlights = false;
  };

  // Badge rendering updates existing cards in place so refresh clicks do not trigger data refetch.
  const syncRenderedRefreshMarkers = () => {
    const renderedCards = [
      ...listContainer.querySelectorAll('[data-testid^="feed-card-"][data-item-id]'),
    ];

    for (const card of renderedCards) {
      card.classList.remove('feed-card--refresh-highlight');
      card.querySelectorAll('.feed-card__refresh-marker').forEach((marker) => {
        marker.remove();
      });

      const itemId = Number(card.getAttribute('data-item-id'));

      if (!Number.isInteger(itemId) || !refreshUpdatedItemIds.has(itemId)) {
        continue;
      }

      if (shouldAnimateRefreshHighlights) {
        card.classList.add('feed-card--refresh-highlight');
      }

      card.append(
        createElement('span', {
          className: 'feed-card__refresh-marker feed-card__refresh-marker--updated',
          text: '~',
          attributes: {
            'data-testid': `feed-card-refresh-marker-${itemId}`,
            'aria-label': 'Refreshed post since last refresh',
            title: 'Refreshed post since last refresh',
          },
        }),
      );
    }

    shouldAnimateRefreshHighlights = false;
  };

  // Recreate and observe the sentinel whenever pagination becomes available again.
  const ensureObserverConnected = () => {
    if (!sentinel || !hasMorePages) {
      return;
    }

    if (!intersectionObserver) {
      intersectionObserver = new IntersectionObserver(handleSentinelIntersect, {
        root: null,
        rootMargin: '100px',
        threshold: 0.01,
      });
    }

    intersectionObserver.observe(sentinel);
  };

  // Disconnect observer when the current feed exhausts available pages.
  const disconnectObserver = () => {
    if (!intersectionObserver) {
      return;
    }

    intersectionObserver.disconnect();
    intersectionObserver = null;
  };

  const syncTabVisualState = () => {
    for (const tabEntry of mountedTabButtons) {
      const isActiveTab = tabEntry.feedTypeId === activeTabType;
      const isPendingTab = tabEntry.feedTypeId === pendingTabType;

      tabEntry.button.classList.toggle('is-active', isActiveTab);
      tabEntry.button.classList.toggle('is-pending', isPendingTab);
      tabEntry.button.setAttribute('aria-pressed', String(isActiveTab));
      tabEntry.button.setAttribute('aria-busy', String(isPendingTab));
    }
  };

  /**
   * Handle state changes from controller (items, loading, error, hasMore).
   * For initial page: replaces all items. For appended page: adds to list.
   * @param {Array<Object> | undefined} items
   * @param {boolean} isLoading
   * @param {string | undefined} error
   * @param {boolean} hasMore
  * @param {boolean | undefined=} isAppending
   * @returns {void}
   */
  const onStateChange = (items, isLoading, error, hasMore, isAppending) => {
    hasMorePages = hasMore;

    // Any settled state should clear pagination loading affordance.
    if (!isLoading) {
      isLoadingMore = false;
      clearLoadMoreCard(listContainer);
    }

    // Loading state: show skeletons for initial load.
    if (isLoading) {
      clearEndOfFeedCard(listContainer);
      renderSkeletons(listContainer, loadingTabType);
      setLoadedPostsCount(0);
      ensureObserverConnected();
      return;
    }

    // Error state: show error message (don't clear list if appending).
    if (error) {
      clearEndOfFeedCard(listContainer);
      renderError(listContainer, error);
      return;
    }

    // Items present: render or append them.
    if (Array.isArray(items)) {
      const shouldRenderRefreshMarkers = refreshUpdatedItemIds.size > 0;
      const itemsForRender = items;
      const itemIds = new Set(
        itemsForRender
          .map((item) => toValidItemId(item?.id))
          .filter((itemId) => itemId !== null),
      );

      // Only previously rendered story cards should trigger append behavior; skeletons and errors should be replaced.
      const hasRenderedStories =
        listContainer.querySelector('[data-testid^="feed-card-"][data-item-id]') !== null;
      const shouldAppend = isAppending === true || (isAppending === undefined && hasRenderedStories);
      const refreshMarkerByItemId = shouldRenderRefreshMarkers
        ? toRefreshMarkerByItemId(itemsForRender, refreshUpdatedItemIds)
        : new Map();

      clearEndOfFeedCard(listContainer);

      // For initial page, items are already filtered by controller; replace list.
      if (!shouldAppend || !hasRenderedStories) {
        renderStoryCards(
          listContainer,
          itemsForRender,
          refreshMarkerByItemId,
          shouldAnimateRefreshHighlights,
        );
        setLoadedPostsCount(itemsForRender.length);
        currentRenderedItemIds = itemIds;
      } else {
        // For appended pages, add to existing list via DocumentFragment.
        appendStoryCards(
          listContainer,
          itemsForRender,
          refreshMarkerByItemId,
          shouldAnimateRefreshHighlights,
        );
        setLoadedPostsCount(loadedPostsCount + itemsForRender.length);
        currentRenderedItemIds = currentRenderedItemIds.union(itemIds);
      }

      shouldAnimateRefreshHighlights = false;
    }

    // Update observer based on hasMore: disconnect if no more pages to load.
    if (!hasMorePages) {
      if (listContainer.querySelector('[data-testid^="feed-card-"][data-item-id]') !== null) {
        renderEndOfFeedCard(listContainer);
      }
      disconnectObserver();
      return;
    }

    clearEndOfFeedCard(listContainer);
    ensureObserverConnected();
  };

  /**
   * Handle tab click via controller.switchTab.
   * @param {string} tabType
   * @returns {Promise<void>}
   */
  const handleTabClick = async (
    tabType,
  ) => {
    scrollViewportToTop();
    loadingTabType = tabType;
    pendingTabType = tabType;

    resetRefreshVisualState();

    setLoadedPostsCount(0);
    syncTabVisualState();

    try {
      await controller.switchTab(tabType, onStateChange);
      activeTabType = tabType;
    } finally {
      pendingTabType = null;
      syncTabVisualState();
    }
  };

  // Refresh markers apply to already-rendered cards and persist for later renders in the same tab.
  const applyRefreshMarkers = async (updatedItemIds = [], { isMarkerReplay = false } = {}) => {
    refreshUpdatedItemIds = new Set(
      updatedItemIds.map((itemId) => toValidItemId(itemId)).filter((itemId) => itemId !== null),
    );
    shouldAnimateRefreshHighlights = !isMarkerReplay;
    syncRenderedRefreshMarkers();

    return {
      updatedItemIds: [...refreshUpdatedItemIds],
    };
  };

  // Clear removes every marker from the active feed view.
  const clearRefreshMarkers = () => {
    resetRefreshVisualState();
    syncRenderedRefreshMarkers();
  };

  /**
   * Attach click handlers to mounted tab buttons and mark the default selection.
   * @param {Array<{ feedTypeId: string, button: HTMLElement }>} tabButtons
   * @returns {void}
   */
  const bindTabButtons = (tabButtons) => {
    const listenerCleanups = tabButtons.map((tabEntry) => {
      const handleClick = () => {
        void handleTabClick(tabEntry.feedTypeId);
      };

      tabEntry.button.addEventListener('click', handleClick);

      return () => {
        tabEntry.button.removeEventListener('click', handleClick);
      };
    });

    cleanupTabListeners = () => {
      for (const cleanup of listenerCleanups) {
        cleanup();
      }
    };

    mountedTabButtons = tabButtons;

    activeTabType = 'top';
    pendingTabType = null;
    loadingTabType = 'top';
    currentRenderedItemIds = new Set();
    resetRefreshVisualState();
    syncTabVisualState();
  };

  /**
   * Intersection callback: triggered when sentinel becomes visible.
   * Calls controller.loadMorePage to fetch next page.
   * @param {IntersectionObserverEntry[]} entries
   * @returns {void}
   */
  const handleSentinelIntersect = (entries) => {
    for (const entry of entries) {
      // Trigger load only when sentinel enters viewport.
      if (!entry.isIntersecting || isLoadingMore || !hasMorePages) {
        continue;
      }

      const hasRenderedStories =
        listContainer.querySelector('[data-testid^="feed-card-"][data-item-id]') !== null;

      if (!hasRenderedStories) {
        continue;
      }

      isLoadingMore = true;
      renderLoadMoreCard(listContainer);

      void Promise.resolve(controller.loadMorePage(onStateChange)).finally(() => {
        isLoadingMore = false;
        clearLoadMoreCard(listContainer);
      });
    }
  };

  /**
   * Mount the feed view into the container.
   * @returns {void}
   */
  const mount = () => {
    clearElement(container);
    currentRenderedItemIds = new Set();
    resetRefreshVisualState();
    setLoadedPostsCount(0);

    // Prefer shell-level sticky tabs when available, with in-view fallback for isolated mounts.
    const shellTabButtons = FEED_TYPES.map((feedType) => ({
      feedTypeId: feedType.id,
      button: document.querySelector(`[data-testid="${toTabTestId(feedType.id)}"]`),
    })).filter((tabEntry) => tabEntry.button !== null);

    const hasCompleteShellTabs = shellTabButtons.length === FEED_TYPES.length;

    listContainer = createItemsList();
    sentinel = createSentinel();

    const scrollTopButton = createScrollTopButton();
    const handleScrollTopClick = () => {
      scrollViewportToTop();
    };
    scrollTopButton.addEventListener('click', handleScrollTopClick);
    cleanupScrollTopButtonListener = () => {
      scrollTopButton.removeEventListener('click', handleScrollTopClick);
    };

    if (hasCompleteShellTabs) {
      appendChildren(container, [listContainer, sentinel, scrollTopButton]);
      bindTabButtons(shellTabButtons);
    } else {
      const fallbackTabs = createTabNavigation();
      appendChildren(container, [
        fallbackTabs.nav,
        listContainer,
        sentinel,
        scrollTopButton,
      ]);
      bindTabButtons(fallbackTabs.tabButtons);
    }

    // Create IntersectionObserver to watch sentinel for page loads.
    ensureObserverConnected();

    // Initialize with default feed on first page.
    controller.mount(onStateChange);
  };

  /**
   * Unmount and clean up: disconnect observer, abort pending requests.
   * @returns {void}
   */
  const unmount = () => {
    cleanupTabListeners();
    cleanupTabListeners = () => {};
    cleanupScrollTopButtonListener();
    cleanupScrollTopButtonListener = () => {};
    mountedTabButtons = [];
    currentRenderedItemIds = new Set();
    resetRefreshVisualState();

    // Disconnect observer to prevent memory leaks.
    if (intersectionObserver) {
      intersectionObserver.disconnect();
      intersectionObserver = null;
    }

    // Abort any pending requests via controller.
    controller.unmount();
    setLoadedPostsCount(0);
    clearElement(container);
  };

  return {
    mount,
    unmount,
    selectTab: handleTabClick,
    applyRefreshMarkers,
    clearRefreshMarkers,
  };
};

export default createFeedView;
