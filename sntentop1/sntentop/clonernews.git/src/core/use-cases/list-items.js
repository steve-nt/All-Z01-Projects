// This use case turns a feed type plus pagination inputs into a sorted page of Hacker News items.
// The public API is the createListItemsUseCase factory, which returns an async executor that always resolves to a Result object.
// The implementation stays pure by depending only on the injected API contract and by avoiding DOM or fetch access.

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/** @typedef {import('../entities/item.js').HnItem} HnItem */
/** @typedef {import('../interfaces/api-interface.js').ApiInterface} ApiInterface */
/** @typedef {import('../interfaces/api-interface.js').FeedType} FeedType */

/**
 * @typedef {Object} ListItemsInput
 * @property {FeedType} type
 * @property {number=} page
 * @property {number=} limit
 */

/**
 * @typedef {Object} ListItemsOutput
 * @property {HnItem[]} items
 * @property {boolean} hasMore
 */

/**
 * @typedef {Object} ListItemsUseCaseDependencies
 * @property {ApiInterface} api
 */

/**
 * @callback ListItemsUseCase
 * @param {ListItemsInput} input
 * @returns {Promise<Result<ListItemsOutput>>}
 */

// Ticket TA-2 sets a default page size of twenty items for deterministic pagination windows.
export const DEFAULT_LIST_ITEMS_LIMIT = 20;
// The six-request cap aligns with API-discipline rules and protects upstream endpoints from bursts.
export const ITEM_FETCH_CONCURRENCY_CAP = 6;

// The supported feed types mirror the API contract so invalid requests fail before any adapter call.
const FEED_TYPES = Object.freeze(['top', 'new', 'ask', 'show', 'job', 'poll']);
const POLL_SOURCE_FEED_TYPES = Object.freeze(['top', 'new', 'ask', 'show']);
const FALLBACK_POLL_ITEM_IDS = Object.freeze([
  47543948, 47385325, 47282783, 47210537, 47203860, 47043735, 46886457, 46866688,
  46809755, 46694649,
]);

// Result helpers keep the control flow readable without repeating object literals.
const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

// Pagination and feed validation stay strict so the use case remains deterministic.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;
const isFeedType = (value) => FEED_TYPES.includes(value);

// Unknown thrown values are collapsed to user-safe strings so callers always receive a Result error.
const toErrorMessage = (value) => {
  if (value instanceof Error && value.message.length > 0) {
    return value.message;
  }

  if (typeof value === 'string' && value.length > 0) {
    return value;
  }

  return 'Unknown error.';
};

// Newer items should sort first, with IDs acting as a deterministic tie-break when timestamps match.
const compareNewestFirst = (left, right) => {
  const timeDiff = (right.time ?? 0) - (left.time ?? 0);

  if (timeDiff !== 0) {
    return timeDiff;
  }

  return right.id - left.id;
};

/**
 * @param {ListItemsUseCaseDependencies} dependencies
 * @returns {ListItemsUseCase}
 */
export const createListItemsUseCase = ({ api } = {}) => {
  /**
   * @param {ListItemsInput} input
   * @returns {Promise<Result<ListItemsOutput>>}
   */
  const execute = async (input) => {
    // The use case depends on the API contract, so fail fast if wiring is incomplete.
    if (!api || typeof api.getFeedIds !== 'function' || typeof api.getItem !== 'function') {
      return err(
        'List items use case requires an API adapter with getFeedIds and getItem methods.',
      );
    }

    // Null and primitive inputs are rejected before destructuring so the executor always returns a Result.
    if (input === null || typeof input !== 'object' || Array.isArray(input)) {
      return err('List items use case expects an input object.');
    }

    const { type, page = 1, limit = DEFAULT_LIST_ITEMS_LIMIT } = input;

    // Invalid feed types are rejected up front so no downstream fetch work is wasted.
    if (!isFeedType(type)) {
      return err('Invalid feed type. Expected one of: top, new, ask, show, job, poll.');
    }

    // Page must stay positive so offset math never produces negative slices.
    if (!isPositiveInteger(page)) {
      return err('Invalid page. Page must be a positive integer.');
    }

    // Limit must also stay positive so the use case never returns empty windows by accident.
    if (!isPositiveInteger(limit)) {
      return err('Invalid limit. Limit must be a positive integer.');
    }

    try {
      // Polls do not have a dedicated endpoint, so they are derived from the top feed.
      const sourceFeedType = type === 'poll' ? 'top' : type;
      // Feed IDs define the canonical order window, so pagination starts from the ID list itself.
      const feedIdsResult = await api.getFeedIds(sourceFeedType);

      if (!feedIdsResult.ok) {
        return err(feedIdsResult.error);
      }

      // Defensive validation keeps malformed adapter payloads from breaking pagination math.
      if (!Array.isArray(feedIdsResult.data)) {
        return err('Feed ID result payload is invalid.');
      }

      const feedIds = feedIdsResult.data;

      if (type === 'poll') {
        const offset = (page - 1) * limit;
        const requiredWindowEnd = offset + limit;
        const requiredCountForHasMore = requiredWindowEnd + 1;
        let collectedPollItems = [];

        // Track processed IDs so feed overlap never triggers duplicate item fetches.
        const visitedItemIds = new Set();

        // Poll collection stays batched to preserve the six-request concurrency cap.
        const collectPollItemsFromIds = async (candidateIds) => {
          for (
            let index = 0;
            index < candidateIds.length && collectedPollItems.length < requiredCountForHasMore;
            index += ITEM_FETCH_CONCURRENCY_CAP
          ) {
            const batchIds = candidateIds
              .slice(index, index + ITEM_FETCH_CONCURRENCY_CAP)
              .filter((id) => !visitedItemIds.has(id));

            if (batchIds.length === 0) {
              continue;
            }

            for (const id of batchIds) {
              visitedItemIds.add(id);
            }

            const batchResults = await Promise.all(batchIds.map((id) => api.getItem(id)));

            const batchPollItemsResult = batchResults.reduce(
              (accumulator, itemResult, batchIndex) => {
                if (!accumulator.ok) {
                  return accumulator;
                }

                const itemId = batchIds[batchIndex];

                if (!itemResult.ok) {
                  return err(`Failed to fetch item ${itemId}: ${itemResult.error}`);
                }

                if (!itemResult.data || itemResult.data.type !== 'poll') {
                  return accumulator;
                }

                return ok([...accumulator.data, itemResult.data]);
              },
              ok([]),
            );

            if (!batchPollItemsResult.ok) {
              return batchPollItemsResult;
            }

            collectedPollItems = [...collectedPollItems, ...batchPollItemsResult.data];
          }

          return ok(collectedPollItems);
        };

        // Poll pages start from top stories, then expand to related feeds if needed.
        const sourceFeedTypes = POLL_SOURCE_FEED_TYPES;

        for (const sourceType of sourceFeedTypes) {
          if (collectedPollItems.length >= requiredCountForHasMore) {
            break;
          }

          const sourceIdsResult =
            sourceType === 'top'
              ? ok(feedIds)
              : await (async () => {
                  const fetchedSourceIdsResult = await api.getFeedIds(sourceType);

                  if (!fetchedSourceIdsResult.ok) {
                    return fetchedSourceIdsResult;
                  }

                  if (!Array.isArray(fetchedSourceIdsResult.data)) {
                    return err('Feed ID result payload is invalid.');
                  }

                  return ok(fetchedSourceIdsResult.data);
                })();

          if (!sourceIdsResult.ok) {
            return sourceIdsResult;
          }

          const sourceCollectionResult = await collectPollItemsFromIds(sourceIdsResult.data);

          if (!sourceCollectionResult.ok) {
            return sourceCollectionResult;
          }
        }

        // Deterministic fallback IDs keep the Polls tab useful when active feeds temporarily contain zero polls.
        if (collectedPollItems.length === 0) {
          const fallbackCollectionResult = await collectPollItemsFromIds(FALLBACK_POLL_ITEM_IDS);

          if (!fallbackCollectionResult.ok) {
            return fallbackCollectionResult;
          }
        }

        const pageItems = collectedPollItems
          .slice(offset, requiredWindowEnd)
          .toSorted(compareNewestFirst);
        const hasMore = collectedPollItems.length > requiredWindowEnd;

        return ok({ items: pageItems, hasMore });
      }

      // Offset math is derived from the requested page so the same page always maps to the same slice.
      const offset = (page - 1) * limit;
      const pageIds = feedIds.slice(offset, offset + limit);
      // hasMore is computed from the source list size, not from the fetched item count.
      const hasMore = offset + limit < feedIds.length;

      // Empty slices can return immediately without touching the item endpoint.
      if (pageIds.length === 0) {
        return ok({ items: [], hasMore });
      }

      let collectedItems = [];

      // Chunking keeps the item endpoint concurrency bounded to avoid request bursts.
      for (let index = 0; index < pageIds.length; index += ITEM_FETCH_CONCURRENCY_CAP) {
        const batchIds = pageIds.slice(index, index + ITEM_FETCH_CONCURRENCY_CAP);
        // Promise.all stays safe here because each batch is intentionally capped at six requests.
        const batchResults = await Promise.all(batchIds.map((id) => api.getItem(id)));

        // Batch reduction preserves the original page order while skipping null items.
        const batchItemsResult = batchResults.reduce((accumulator, itemResult, batchIndex) => {
          if (!accumulator.ok) {
            return accumulator;
          }

          const itemId = batchIds[batchIndex];

          if (!itemResult.ok) {
            return err(`Failed to fetch item ${itemId}: ${itemResult.error}`);
          }

          if (!itemResult.data) {
            return accumulator;
          }

          // Copy-on-write accumulation keeps reducer state immutable across all batch iterations.
          return ok([...accumulator.data, itemResult.data]);
        }, ok([]));

        if (!batchItemsResult.ok) {
          return batchItemsResult;
        }

        // The collected list accumulates all successful batch results before the final sort.
        collectedItems = [...collectedItems, ...batchItemsResult.data];
      }

      // Top feed preserves raw API ranking order; other feeds remain newest-first.
      const shouldSortNewestFirst = type !== 'top';
      const orderedItems = shouldSortNewestFirst
        ? collectedItems.toSorted(compareNewestFirst)
        : collectedItems;

      return ok({ items: orderedItems, hasMore });
    } catch (unexpectedError) {
      // Unexpected exceptions are normalized into Result errors so callers never have to catch here.
      return err(`Failed to list items: ${toErrorMessage(unexpectedError)}`);
    }
  };

  return execute;
};
