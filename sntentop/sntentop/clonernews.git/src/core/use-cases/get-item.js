// This use case resolves a root Hacker News item and its comment tree into a deterministic Result payload.
// The public API is createGetItemUseCase, which returns an async executor for stories, jobs, and polls.
// The implementation stays pure by depending only on the injected API contract and by avoiding DOM access.

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/** @typedef {import('../entities/item.js').HnItem} HnItem */
/** @typedef {import('../interfaces/api-interface.js').ApiInterface} ApiInterface */

/**
 * @typedef {Omit<HnItem, 'parts'> & { parts?: HnItem[] }} ItemWithResolvedParts
 */

/**
 * @typedef {Object} TreeNode
 * @property {HnItem} item
 * @property {TreeNode[]} comments
 */

/**
 * @typedef {Object} GetItemOutput
 * @property {ItemWithResolvedParts} item
 * @property {TreeNode[]} comments
 */

/**
 * @typedef {Object} GetItemUseCaseDependencies
 * @property {ApiInterface} api
 * @property {number=} maxDepth
 */

/**
 * @callback GetItemUseCase
 * @param {number} id
 * @returns {Promise<Result<GetItemOutput>>}
 */

// Depth constants are exported so tests and future callers share one authoritative recursion policy.
export const MAX_COMMENT_TREE_DEPTH = 5;
export const DEFAULT_GET_ITEM_MAX_DEPTH = MAX_COMMENT_TREE_DEPTH;
export const GET_ITEM_FETCH_CONCURRENCY_CAP = 6;

// Small Result constructors keep success/error return shapes consistent across all branches.
const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

// Positive-integer validation protects recursion, pagination, and adapter calls from invalid IDs.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Unknown thrown values are normalized to stable strings so callers always receive predictable errors.
const toErrorMessage = (value) => {
  if (value instanceof Error && value.message.length > 0) {
    return value.message;
  }

  if (typeof value === 'string' && value.length > 0) {
    return value;
  }

  return 'Unknown error.';
};

// Newest-first ordering is stable by timestamp, then by ID for deterministic ties.
const compareNewestFirstItems = (left, right) => {
  const timeDiff = (right.time ?? 0) - (left.time ?? 0);

  if (timeDiff !== 0) {
    return timeDiff;
  }

  return right.id - left.id;
};

// Node-level comparator delegates to item comparator so ordering logic stays single-sourced.
const compareNewestFirstNodes = (left, right) => compareNewestFirstItems(left.item, right.item);

/**
 * @param {number} maxConcurrent
 * @returns {(task: () => Promise<unknown>) => Promise<unknown>}
 */
const createConcurrencyLimiter = (maxConcurrent) => {
  // Active count tracks in-flight work so the limiter enforces a strict upper bound.
  let activeCount = 0;
  // Queue entries are resolver callbacks waiting for the next available execution slot.
  const waitingResolvers = [];

  const runNext = () => {
    // No-op when capacity is full or no queued work exists.
    if (activeCount >= maxConcurrent || waitingResolvers.length === 0) {
      return;
    }

    // Reserve a slot before waking the next queued waiter to avoid double-claim races.
    activeCount += 1;
    const resolveNext = waitingResolvers.shift();
    resolveNext();
  };

  const acquire = async () => {
    // Fast path takes an available slot immediately.
    if (activeCount < maxConcurrent) {
      activeCount += 1;
      return;
    }

    // Otherwise queue until release() wakes this waiter.
    await new Promise((resolve) => {
      waitingResolvers.push(resolve);
    });
  };

  const release = () => {
    // Releasing a slot immediately triggers the next queued task, if any.
    activeCount -= 1;
    runNext();
  };

  return async (task) => {
    await acquire();

    try {
      return await task();
    } finally {
      release();
    }
  };
};

/**
 * @template T
 * @param {number[]} ids
 * @param {(id: number) => Promise<Result<T | null>>} resolveById
 * @returns {Promise<Result<T[]>>}
 */
const collectResolvedItemsInBatches = async (ids, resolveById) => {
  if (!Array.isArray(ids) || ids.length === 0) {
    return ok([]);
  }

  let collectedItems = [];

  // Chunking keeps all comment and poll-option retrieval aligned with the global max-six request rule.
  for (let index = 0; index < ids.length; index += GET_ITEM_FETCH_CONCURRENCY_CAP) {
    const batchIds = ids.slice(index, index + GET_ITEM_FETCH_CONCURRENCY_CAP);
    const batchResults = await Promise.all(batchIds.map((id) => resolveById(id)));

    // Reducer composes the batch into one Result so callers can handle success/failure uniformly.
    const batchCollectedItemsResult = batchResults.reduce((accumulator, itemResult) => {
      if (!accumulator.ok) {
        return accumulator;
      }

      // First error short-circuits the batch so callers keep the original failure cause.
      if (!itemResult.ok) {
        return itemResult;
      }

      // Null entries represent missing/deleted items and are intentionally skipped.
      if (itemResult.data === null) {
        return accumulator;
      }

      // Immutable accumulation avoids hidden mutations between reduce iterations.
      return ok([...accumulator.data, itemResult.data]);
    }, ok([]));

    if (!batchCollectedItemsResult.ok) {
      return batchCollectedItemsResult;
    }

    // Batch results are appended immutably so ordering stays deterministic across chunks.
    collectedItems = [...collectedItems, ...batchCollectedItemsResult.data];
  }

  return ok(collectedItems);
};

/**
 * @param {(id: number) => Promise<Result<HnItem | null>>} getItemById
 * @param {number} id
 * @param {string} itemLabel
 * @returns {Promise<Result<HnItem | null>>}
 */
const fetchItemById = async (getItemById, id, itemLabel) => {
  const itemResult = await getItemById(id);

  if (!itemResult.ok) {
    return err(`Failed to fetch ${itemLabel} ${id}: ${itemResult.error}`);
  }

  return ok(itemResult.data);
};

/**
 * @param {GetItemUseCaseDependencies=} dependencies
 * @returns {GetItemUseCase}
 */
export const createGetItemUseCase = ({ api, maxDepth = DEFAULT_GET_ITEM_MAX_DEPTH } = {}) => {
  // Wiring guard keeps misconfigured dependency injection from reaching runtime traversal paths.
  if (!api || typeof api.getItem !== 'function') {
    return async () => err('Get item use case requires an API adapter with a getItem method.');
  }

  // Invalid maxDepth is rejected immediately so recursion semantics remain explicit and predictable.
  if (!isPositiveInteger(maxDepth)) {
    return async () => err('Invalid max depth. maxDepth must be a positive integer.');
  }

  const depthLimit = Math.min(maxDepth, MAX_COMMENT_TREE_DEPTH);
  // A single limiter instance keeps total in-flight getItem calls bounded across all recursion branches.
  const runWithFetchLimit = createConcurrencyLimiter(GET_ITEM_FETCH_CONCURRENCY_CAP);
  // Adapter calls are wrapped so every fetch in this use case goes through one global concurrency gate.
  const getItemById = async (id) =>
    /** @type {Promise<Result<HnItem | null>>} */ (runWithFetchLimit(() => api.getItem(id)));

  /**
   * @param {number[]} commentIds
   * @param {number} parentId
   * @param {number} depth
   * @param {Set<number>} ancestorIds
   * @returns {Promise<Result<TreeNode[]>>}
   */
  const buildCommentForest = async (commentIds, parentId, depth, ancestorIds) => {
    // Guard clauses keep recursion bounded and avoid extra adapter work on empty branches.
    if (!Array.isArray(commentIds) || commentIds.length === 0 || depth > depthLimit) {
      return ok([]);
    }

    // Batch resolution reuses shared error/null handling so each recursion level stays consistent.
    const commentNodesResult = await collectResolvedItemsInBatches(
      commentIds,
      async (commentId) => {
        // Ancestor cycle checks prevent self-referential kids graphs from re-walking the same path.
        if (ancestorIds.has(commentId)) {
          return ok(null);
        }

        const commentResult = await fetchItemById(getItemById, commentId, 'comment');

        if (!commentResult.ok) {
          return commentResult;
        }

        if (commentResult.data === null) {
          // Null comments are skipped so deleted nodes do not pollute the resulting tree.
          return ok(null);
        }

        // Parent-link propagation ensures each comment carries a direct reference for audit checks.
        const commentWithParent = {
          ...commentResult.data,
          parent: parentId,
        };

        if (depth >= depthLimit || !Array.isArray(commentWithParent.kids)) {
          // Depth/leaf exit keeps traversal bounded while still returning the current comment node.
          return ok({
            item: commentWithParent,
            comments: [],
          });
        }

        const nestedCommentsResult = await buildCommentForest(
          commentWithParent.kids,
          commentWithParent.id,
          depth + 1,
          new Set([...ancestorIds, commentWithParent.id]),
        );

        if (!nestedCommentsResult.ok) {
          return nestedCommentsResult;
        }

        return ok({
          item: commentWithParent,
          comments: nestedCommentsResult.data,
        });
      },
    );

    if (!commentNodesResult.ok) {
      return commentNodesResult;
    }

    // Sorting once at the end keeps each depth level deterministic without extra work during traversal.
    return ok(commentNodesResult.data.toSorted(compareNewestFirstNodes));
  };

  /**
   * @param {number} id
   * @returns {Promise<Result<GetItemOutput>>}
   */
  const execute = async (id) => {
    if (!isPositiveInteger(id)) {
      return err('Invalid item ID. Item ID must be a positive integer.');
    }

    try {
      const rootItemResult = await fetchItemById(getItemById, id, 'item');

      if (!rootItemResult.ok) {
        return rootItemResult;
      }

      if (rootItemResult.data === null) {
        return err(`Item ${id} was not found.`);
      }

      // Root-item caching avoids repeated reads and keeps following branch logic concise.
      const rootItem = rootItemResult.data;

      // Poll options are resolved only when parts exist so non-poll items avoid unnecessary adapter calls.
      const resolvedPartsResult = Array.isArray(rootItem.parts)
        ? await collectResolvedItemsInBatches(rootItem.parts, (partId) =>
            fetchItemById(getItemById, partId, 'poll option'),
          )
        : ok([]);

      if (!resolvedPartsResult.ok) {
        return err(resolvedPartsResult.error);
      }

      const commentsResult = await buildCommentForest(
        rootItem.kids ?? [],
        rootItem.id,
        1,
        new Set([rootItem.id]),
      );

      if (!commentsResult.ok) {
        return commentsResult;
      }

      // The returned root payload keeps parts canonicalized as resolved poll-option items when applicable.
      const itemWithResolvedParts = Array.isArray(rootItem.parts)
        ? {
            ...rootItem,
            parts: resolvedPartsResult.data,
          }
        : {
            ...rootItem,
          };

      return ok({
        item: itemWithResolvedParts,
        comments: commentsResult.data,
      });
    } catch (unexpectedError) {
      // Throwable adapter/runtime faults are wrapped so callers never rely on exceptions for control flow.
      return err(`Failed to get item ${id}: ${toErrorMessage(unexpectedError)}`);
    }
  };

  return execute;
};
