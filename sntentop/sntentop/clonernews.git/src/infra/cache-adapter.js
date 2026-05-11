/*
 * Purpose: Provide an in-memory TTL cache adapter for HN API reads.
 * Public API: createCacheAdapter({ ttlMs }) returns get, set, has, invalidate, and clear methods.
 * Implementation notes: The adapter stores only fresh entries and returns Result-shaped outcomes.
 */
import { Temporal } from '@js-temporal/polyfill';

const DEFAULT_TTL_MS = 60_000;

const ok = (data) => ({ ok: true, data });

const getNowEpochMilliseconds = () => Temporal.Now.instant().epochMilliseconds;

export const createCacheAdapter = ({ ttlMs = DEFAULT_TTL_MS } = {}) => {
  const store = new Map();

  const set = async (key, value) => {
    const expiresAt = getNowEpochMilliseconds() + ttlMs;
    // Keep the write isolated so callers never need to manage expiry metadata themselves.
    store.set(key, { value, expiresAt });

    return ok(undefined);
  };

  const get = async (key) => {
    // Filter expired entries before exposing cached data to callers.
    const entry = getFreshEntry(key);

    return ok(entry?.value);
  };

  const has = async (key) => {
    // Reuse the same freshness check so presence and retrieval stay consistent.
    const entry = getFreshEntry(key);

    return ok(entry !== undefined);
  };

  const invalidate = async (key) => {
    // Remove a single key so tests and cache misses can force a clean read path.
    store.delete(key);

    return ok(undefined);
  };

  const clear = async () => {
    // Clear the cache wholesale when the adapter needs a fresh baseline.
    store.clear();

    return ok(undefined);
  };

  const getFreshEntry = (key) => {
    // Read once so the freshness decision is based on the same snapshot.
    const entry = store.get(key);

    if (entry === undefined) {
      return undefined;
    }

    // Expired entries are removed eagerly so stale data never survives a lookup.
    if (entry.expiresAt < getNowEpochMilliseconds()) {
      store.delete(key);
      return undefined;
    }

    return entry;
  };

  return {
    get,
    set,
    has,
    invalidate,
    clear,
  };
};

export default createCacheAdapter;
