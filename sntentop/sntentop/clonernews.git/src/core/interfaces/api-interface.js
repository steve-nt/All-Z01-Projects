// The API interface defines the transport contract that infra adapters must satisfy without exposing fetch details.
// It keeps core use cases decoupled from the HN transport layer and from any browser or DOM APIs.
// The exported method-name list is used as a narrow contract check in adapters and tests.
/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/**
 * @typedef {'top' | 'new' | 'ask' | 'show' | 'job'} FeedType
 */

/**
 * @callback GetItem
 * @param {number} id
 * @returns {Promise<Result<import('../entities/item.js').HnItem | null>>}
 */

/**
 * @callback GetFeedIds
 * @param {FeedType} type
 * @returns {Promise<Result<number[]>>}
 */

/**
 * @callback GetUpdates
 * @returns {Promise<Result<number[]>>}
 */

/**
 * @typedef {Object} ApiInterface
 * @property {GetItem} getItem
 * @property {GetFeedIds} getFeedIds
 * @property {GetUpdates} getUpdates
 */

// Method-name metadata is exported so adapter conformance checks can stay declarative.
export const API_METHOD_NAMES = Object.freeze(['getItem', 'getFeedIds', 'getUpdates']);
