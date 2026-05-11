/*
 * Purpose: Manage feed view state, coordinate tab switching, and handle infinite scroll pagination.
 * Public API: createFeedController({ listItems }) returns { switchTab, mount, unmount, loadMorePage }.
 * Implementation notes: Tracks active feed type, manages loading flag to prevent duplicate requests,
 * handles page pagination with stale-response tracking, integrates IntersectionObserver via loadMorePage.
 */

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/**
 * @typedef {Object} FeedControllerDependencies
 * @property {(input: {type: string, page?: number, limit?: number}) => Promise<Result<{items: any[], hasMore: boolean}>>} listItems
 */

/**
 * @typedef {Object} FeedController
 * @property {(tabType: string, onStateChange: Function) => Promise<void>} switchTab
 * @property {(onStateChange: Function) => void} mount
 * @property {() => void} unmount
 * @property {(onStateChange: Function) => Promise<void>} loadMorePage
 */

const FEED_TYPES = Object.freeze(['top', 'new', 'job', 'ask', 'show', 'poll']);
const DEFAULT_PAGE_SIZE = 20;

/**
 * @param {FeedControllerDependencies} dependencies
 * @returns {FeedController}
 */
export const createFeedController = ({ listItems } = {}) => {
  // Track current feed type, pagination state, and abort controller.
  let currentFeedType = 'top';
  let currentAbortController = null;
  let currentRequestId = 0;

  // Pagination tracking: current page and whether more items exist.
  let currentPage = 1;
  let hasMorePages = true;

  // Loading flag prevents duplicate requests during in-flight fetch.
  let isLoading = false;

  // Validate dependencies on wiring to fail fast.
  if (typeof listItems !== 'function') {
    throw new Error('Feed controller requires a listItems function.');
  }

  /**
   * Internal fetch helper: increments request ID, sets loading flag, calls listItems use-case.
   * @param {string} feedType
   * @param {number} page
   * @param {Function} onStateChange
   * @param {boolean} isAppending
   * @returns {Promise<void>}
   */
  const performFetch = async (feedType, page, onStateChange, isAppending) => {
    // Cancel any previous in-flight request.
    if (currentAbortController) {
      currentAbortController.abort();
    }

    // Increment request ID to ignore stale responses from older requests.
    currentRequestId += 1;
    const requestId = currentRequestId;

    // Set loading flag to prevent IntersectionObserver from triggering duplicate fetches.
    isLoading = true;
    currentAbortController = new AbortController();

    // Only show loading state for initial page (not for infinite scroll appends).
    if (!isAppending) {
      onStateChange([], true, undefined, hasMorePages, false);
    }

    try {
      const normalizedResult = await listItems({
        type: feedType,
        page,
        limit: DEFAULT_PAGE_SIZE,
      });

      // Ignore stale responses if a newer request has already started.
      if (requestId !== currentRequestId) {
        isLoading = false;
        return;
      }

      if (!normalizedResult.ok) {
        // On error, keep existing items and show error message (don't clear list).
        onStateChange(undefined, false, normalizedResult.error, hasMorePages, isAppending);
        isLoading = false;
        return;
      }

      // Update pagination state from the API response.
      hasMorePages = normalizedResult.data.hasMore;

      // Call onStateChange with items, loading state, error, and hasMore flag.
      onStateChange(normalizedResult.data.items, false, undefined, hasMorePages, isAppending);
      isLoading = false;
    } catch (error) {
      // Ignore errors from aborted or stale requests.
      if (requestId !== currentRequestId) {
        isLoading = false;
        return;
      }

      // Normalize error message for safe display.
      const errorMsg = error instanceof Error ? error.message : 'Failed to load feed.';
      onStateChange(undefined, false, errorMsg, hasMorePages, isAppending);
      isLoading = false;
    }
  };

  /**
   * Fetch initial page of the current feed type (page reset to 1).
   * @param {Function} onStateChange
   * @returns {Promise<void>}
   */
  const fetchInitialPage = async (onStateChange) => {
    currentPage = 1;
    hasMorePages = true;
    await performFetch(currentFeedType, currentPage, onStateChange, false);
  };

  /**
   * Fetch the next page and append items to the existing list.
   * Called by IntersectionObserver when sentinel becomes visible.
   * @param {Function} onStateChange
   * @returns {Promise<void>}
   */
  const loadMorePage = async (onStateChange) => {
    // Prevent overlapping requests (loading flag checked by caller too).
    if (isLoading || !hasMorePages) {
      return;
    }

    // Increment page for next batch of items.
    currentPage += 1;
    await performFetch(currentFeedType, currentPage, onStateChange, true);
  };

  /**
   * Switch to a new feed tab: reset pagination and fetch first page.
   * @param {string} feedType
   * @param {Function} onStateChange
   * @returns {Promise<void>}
   */
  const switchTab = async (feedType, onStateChange) => {
    if (!FEED_TYPES.includes(feedType)) {
      return;
    }

    currentFeedType = feedType;
    await fetchInitialPage(onStateChange);
  };

  /**
   * Mount the controller and initialize the default feed on first page.
   * @param {Function} onStateChange
   * @returns {void}
   */
  const mount = (onStateChange) => {
    void fetchInitialPage(onStateChange);
  };

  /**
   * Unmount and cancel any pending requests; clean up abort controller.
   * @returns {void}
   */
  const unmount = () => {
    if (currentAbortController) {
      currentAbortController.abort();
    }

    isLoading = false;
  };

  return {
    switchTab,
    mount,
    unmount,
    loadMorePage,
  };
};

export default createFeedController;
