// The storage interface keeps caching abstract so core logic can stay detached from Map or browser storage details.
// It gives adapters a small contract for cache lookups, writes, invalidation, and clearing.
// The method list is intentionally tiny so cache implementations stay easy to audit.
/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/**
 * @callback StorageGet
 * @param {string | number} key
 * @returns {Promise<Result<unknown>>}
 */

/**
 * @callback StorageSet
 * @param {string | number} key
 * @param {unknown} value
 * @returns {Promise<Result<void>>}
 */

/**
 * @typedef {Object} StorageInterface
 * @property {StorageGet} get
 * @property {StorageSet} set
 * @property {(key: string | number) => Promise<Result<boolean>>} has
 * @property {(key: string | number) => Promise<Result<void>>} invalidate
 * @property {() => Promise<Result<void>>} clear
 */

// Method-name metadata makes storage contract verification explicit in tests and adapters.
export const STORAGE_METHOD_NAMES = Object.freeze(['get', 'set', 'has', 'invalidate', 'clear']);
