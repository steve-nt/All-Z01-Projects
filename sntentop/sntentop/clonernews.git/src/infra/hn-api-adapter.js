// This adapter isolates Hacker News transport concerns so core code only sees validated Result objects.
// The public API is createHnApiAdapter, which hides fetch, cache, timeout, and payload validation details.
// The module stays infra-only so core use cases never depend on transport concerns or DOM APIs.
import { isHnItemType } from '../core/entities/item.js';
import { createCacheAdapter } from './cache-adapter.js';

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/** @typedef {import('../core/interfaces/api-interface.js').ApiInterface} ApiInterface */
/** @typedef {import('../core/interfaces/storage-interface.js').StorageInterface} StorageInterface */

export const HN_API_BASE_URL = 'https://hacker-news.firebaseio.com/v0';
export const DEFAULT_FETCH_TIMEOUT_MS = 8000;

// Keep endpoint selection explicit so invalid feed types fail fast instead of building URLs dynamically.
const FEED_ENDPOINT_BY_TYPE = Object.freeze({
  top: '/topstories.json',
  new: '/newstories.json',
  ask: '/askstories.json',
  show: '/showstories.json',
  job: '/jobstories.json',
});

// Updates live on a dedicated endpoint so callers do not need endpoint-string duplication.
const UPDATES_ENDPOINT = '/updates.json';

// Pin request behavior so the adapter stays predictable across browsers and avoids credential leakage.
const FETCH_OPTIONS = Object.freeze({
  method: 'GET',
  mode: 'cors',
  credentials: 'omit',
  cache: 'no-store',
  redirect: 'follow',
  headers: Object.freeze({
    Accept: 'application/json',
  }),
});

const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

// Narrow guard helpers keep the payload validators readable and deterministic.
const isObjectRecord = (value) =>
  value !== null && typeof value === 'object' && !Array.isArray(value);
// Positive IDs are the only valid identifiers for HN items and feed entries.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;
// Timestamps and descendants can legitimately be zero, so this guard stays separate.
const isNonNegativeInteger = (value) => Number.isInteger(value) && value >= 0;
// Scores may be negative in some datasets, so the validator only checks integer-ness here.
const isInteger = (value) => Number.isInteger(value);
// HN nested arrays always contain item IDs, so each entry must pass the positive-ID guard.
const isIntegerArray = (value) =>
  Array.isArray(value) && value.every((entry) => isPositiveInteger(entry));
// Response checks stay case-insensitive because servers may vary the charset suffix.
const isJsonContentType = (value) => value.toLowerCase().includes('application/json');

// Normalize the caller-provided base URL so request assembly never has to special-case slashes.
const normalizeBaseUrl = (value) => {
  if (typeof value !== 'string') {
    return HN_API_BASE_URL;
  }

  const trimmedValue = value.trim();

  if (trimmedValue.length === 0) {
    return HN_API_BASE_URL;
  }

  return trimmedValue.endsWith('/') ? trimmedValue.slice(0, -1) : trimmedValue;
};

// Transport policy is HTTPS-only, with localhost HTTP allowed for local development and tests.
const validateBaseUrl = (value) => {
  const normalizedBaseUrl = normalizeBaseUrl(value);

  let parsedBaseUrl;

  try {
    parsedBaseUrl = new URL(normalizedBaseUrl);
  } catch {
    return err('Base URL must be a valid absolute URL.');
  }

  const isLocalhost =
    parsedBaseUrl.hostname === 'localhost' ||
    parsedBaseUrl.hostname === '127.0.0.1' ||
    parsedBaseUrl.hostname === '::1' ||
    parsedBaseUrl.hostname === '[::1]';
  const isHttps = parsedBaseUrl.protocol === 'https:';
  const isLocalHttp = isLocalhost && parsedBaseUrl.protocol === 'http:';

  if (!isHttps && !isLocalHttp) {
    return err('Base URL must use HTTPS (HTTP is allowed only for localhost testing).');
  }

  return ok(normalizedBaseUrl);
};

// Preserve useful failure detail while still reducing arbitrary thrown values to strings.
const toErrorMessage = (value) => {
  if (value instanceof Error) {
    return value.message;
  }

  if (typeof value === 'string' && value.length > 0) {
    return value;
  }

  return 'Unknown error.';
};

// Abort is the one expected exceptional fetch path, so it gets a dedicated predicate.
const isAbortError = (value) => value instanceof Error && value.name === 'AbortError';

// Array validation is shared by feed and update endpoints, so the label keeps error messages specific.
const validateIdArray = (value, label) => {
  if (!Array.isArray(value)) {
    return err(`${label} must be an array.`);
  }

  if (!value.every((entry) => isPositiveInteger(entry))) {
    return err(`${label} must contain positive integer IDs only.`);
  }

  return ok(value);
};

// Sanitizing payloads here keeps malformed API data from reaching the cache or callers.
const validateItemPayload = (value) => {
  if (!isObjectRecord(value)) {
    return err('Item payload must be an object.');
  }

  const item = /** @type {Record<string, unknown>} */ (value);

  // Validate the API shape before it can enter the cache or escape to callers.
  if (!isPositiveInteger(item.id)) {
    return err('Item payload field "id" must be a positive integer.');
  }

  if (!isHnItemType(item.type)) {
    return err('Item payload field "type" is invalid.');
  }

  if ('by' in item && typeof item.by !== 'string') {
    return err('Item payload field "by" must be a string when provided.');
  }

  if ('time' in item && !isNonNegativeInteger(item.time)) {
    return err('Item payload field "time" must be a non-negative integer when provided.');
  }

  if ('text' in item && typeof item.text !== 'string') {
    return err('Item payload field "text" must be a string when provided.');
  }

  if ('kids' in item && !isIntegerArray(item.kids)) {
    return err('Item payload field "kids" must be an array of positive integer IDs when provided.');
  }

  if ('url' in item && typeof item.url !== 'string') {
    return err('Item payload field "url" must be a string when provided.');
  }

  if ('score' in item && !isInteger(item.score)) {
    return err('Item payload field "score" must be an integer when provided.');
  }

  if ('title' in item && typeof item.title !== 'string') {
    return err('Item payload field "title" must be a string when provided.');
  }

  if ('parts' in item && !isIntegerArray(item.parts)) {
    return err(
      'Item payload field "parts" must be an array of positive integer IDs when provided.',
    );
  }

  if ('descendants' in item && !isNonNegativeInteger(item.descendants)) {
    return err('Item payload field "descendants" must be a non-negative integer when provided.');
  }

  if ('parent' in item && !isPositiveInteger(item.parent)) {
    return err('Item payload field "parent" must be a positive integer when provided.');
  }

  const sanitizedItem = {
    id: item.id,
    type: item.type,
  };

  if (typeof item.by === 'string') {
    sanitizedItem.by = item.by;
  }

  if (typeof item.time === 'number') {
    sanitizedItem.time = item.time;
  }

  if (typeof item.text === 'string') {
    sanitizedItem.text = item.text;
  }

  if (Array.isArray(item.kids)) {
    sanitizedItem.kids = item.kids;
  }

  if (typeof item.url === 'string') {
    sanitizedItem.url = item.url;
  }

  if (typeof item.score === 'number') {
    sanitizedItem.score = item.score;
  }

  if (typeof item.title === 'string') {
    sanitizedItem.title = item.title;
  }

  if (Array.isArray(item.parts)) {
    sanitizedItem.parts = item.parts;
  }

  if (typeof item.descendants === 'number') {
    sanitizedItem.descendants = item.descendants;
  }

  if (typeof item.parent === 'number') {
    sanitizedItem.parent = item.parent;
  }

  return ok(sanitizedItem);
};

// Updates may arrive as either a raw array or an object wrapper, so both shapes are accepted.
const validateUpdatesPayload = (value) => {
  if (Array.isArray(value)) {
    return validateIdArray(value, 'Updates payload');
  }

  if (!isObjectRecord(value)) {
    return err('Updates payload must be an object with an "items" array.');
  }

  if (!('items' in value)) {
    return err('Updates payload is missing the "items" field.');
  }

  return validateIdArray(value.items, 'Updates payload field "items"');
};

// Consolidate timeout, response, and JSON parsing rules so every endpoint behaves the same way.
const readJsonResponse = async ({ fetchFn, url, endpointLabel, timeoutMs, allowNotFound }) => {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => {
    controller.abort();
  }, timeoutMs);

  try {
    const response = await fetchFn(url, { ...FETCH_OPTIONS, signal: controller.signal });

    if (allowNotFound && response.status === 404) {
      return ok(null);
    }

    if (!response.ok) {
      return err(`Request to ${endpointLabel} failed with status ${response.status}.`);
    }

    const contentType = response.headers.get('content-type') ?? '';

    if (!isJsonContentType(contentType)) {
      return err(`Request to ${endpointLabel} returned unsupported content type "${contentType}".`);
    }

    try {
      return ok(await response.json());
    } catch (parseError) {
      return err(
        `Request to ${endpointLabel} returned malformed JSON: ${toErrorMessage(parseError)}`,
      );
    }
  } catch (requestError) {
    if (isAbortError(requestError)) {
      return err(`Request to ${endpointLabel} timed out after ${timeoutMs}ms.`);
    }

    return err(`Request to ${endpointLabel} failed: ${toErrorMessage(requestError)}`);
  } finally {
    // Clear the timeout no matter how the request settles so timers never leak across calls.
    clearTimeout(timeoutId);
  }
};

// Cache reads are wrapped because cache implementations are allowed to fail without breaking API reads.
const safeCacheGet = async (cache, key) => {
  try {
    const cacheResult = await cache.get(key);

    if (!cacheResult.ok) {
      return ok(undefined);
    }

    return ok(cacheResult.data);
  } catch {
    // Cache read faults are intentionally ignored so transport reads remain available.
    return ok(undefined);
  }
};

// Cache writes are intentionally best-effort so transport success is never lost to storage issues.
const safeCacheSet = async (cache, key, value) => {
  try {
    // Cache writes are best-effort so a broken cache never blocks fresh API data.
    await cache.set(key, value);
  } catch {
    // Cache write failures should not break successful API reads.
  }
};

const getFeedEndpoint = (type) => {
  if (typeof type !== 'string') {
    return undefined;
  }

  // Direct map lookup keeps allowed feed types centralized in one immutable table.
  return FEED_ENDPOINT_BY_TYPE[type];
};

/**
 * Create the Hacker News API adapter with injectable fetch, timeout, and storage dependencies.
 *
 * The returned object is intentionally tiny so core use-cases can stay isolated from transport concerns.
 * @param {Object} [options]
 * @param {string=} options.baseUrl
 * @param {number=} options.timeoutMs
 * @param {typeof fetch=} options.fetchFn
 * @param {StorageInterface=} options.cache
 * @returns {ApiInterface}
 */
export const createHnApiAdapter = (options = {}) => {
  // Keep the base URL and timeout deterministic even when callers omit configuration.
  const baseUrlResult = validateBaseUrl(options.baseUrl);
  const timeoutMs = isPositiveInteger(options.timeoutMs)
    ? options.timeoutMs
    : DEFAULT_FETCH_TIMEOUT_MS;
  // Allow dependency injection in tests while still defaulting to the browser fetch implementation.
  const fetchFn = options.fetchFn ?? globalThis.fetch?.bind(globalThis);
  // Use the shared cache adapter by default so infra behavior matches production wiring.
  const cache = options.cache ?? createCacheAdapter();

  if (!baseUrlResult.ok) {
    const invalidBaseUrl = async () => err(baseUrlResult.error);

    return Object.freeze({
      getItem: invalidBaseUrl,
      getFeedIds: invalidBaseUrl,
      getUpdates: invalidBaseUrl,
    });
  }

  const baseUrl = baseUrlResult.data;

  if (typeof fetchFn !== 'function') {
    // Surface an explicit Result instead of throwing so callers can handle the environment gap safely.
    const fetchUnavailable = async () => err('Global fetch is not available in this environment.');

    return Object.freeze({
      getItem: fetchUnavailable,
      getFeedIds: fetchUnavailable,
      getUpdates: fetchUnavailable,
    });
  }

  const getItem = async (id) => {
    if (!isPositiveInteger(id)) {
      return err('Item ID must be a positive integer.');
    }

    // Check cache first so repeated item requests do not waste network calls.
    const cachedValueResult = await safeCacheGet(cache, id);

    if (cachedValueResult.ok && cachedValueResult.data !== undefined) {
      // Null entries are cached explicitly so missing items do not keep hitting the network.
      if (cachedValueResult.data === null) {
        return ok(null);
      }

      // Re-validate cached payloads so stale or externally injected data cannot leak through.
      const validatedCachedItem = validateItemPayload(cachedValueResult.data);

      if (validatedCachedItem.ok) {
        return validatedCachedItem;
      }
    }

    const endpoint = `/item/${id}.json`;
    const itemResponse = await readJsonResponse({
      fetchFn,
      url: `${baseUrl}${endpoint}`,
      endpointLabel: endpoint,
      timeoutMs,
      allowNotFound: true,
    });

    if (!itemResponse.ok) {
      return itemResponse;
    }

    if (itemResponse.data === null) {
      // Cache null explicitly so repeated 404 lookups stay cheap.
      await safeCacheSet(cache, id, null);
      return ok(null);
    }

    const validatedItem = validateItemPayload(itemResponse.data);

    if (!validatedItem.ok) {
      return validatedItem;
    }

    await safeCacheSet(cache, id, validatedItem.data);
    return ok(validatedItem.data);
  };

  const getFeedIds = async (type) => {
    // Reject unsupported feed types before any fetch is attempted.
    const endpoint = getFeedEndpoint(type);

    if (endpoint === undefined) {
      return err('Feed type must be one of: top, new, ask, show, job.');
    }

    const feedResponse = await readJsonResponse({
      fetchFn,
      url: `${baseUrl}${endpoint}`,
      endpointLabel: endpoint,
      timeoutMs,
      allowNotFound: false,
    });

    if (!feedResponse.ok) {
      return feedResponse;
    }

    // Feeds must resolve to a strict id array so use-cases can page them safely.
    return validateIdArray(feedResponse.data, `Feed payload for "${type}"`);
  };

  const getUpdates = async () => {
    // Updates use the same transport rules as feeds so throttling and validation stay uniform.
    const updatesResponse = await readJsonResponse({
      fetchFn,
      url: `${baseUrl}${UPDATES_ENDPOINT}`,
      endpointLabel: UPDATES_ENDPOINT,
      timeoutMs,
      allowNotFound: false,
    });

    if (!updatesResponse.ok) {
      return updatesResponse;
    }

    // Updates payloads vary slightly by endpoint shape, so normalize both accepted forms.
    return validateUpdatesPayload(updatesResponse.data);
  };

  return Object.freeze({
    getItem,
    getFeedIds,
    getUpdates,
  });
};

export default createHnApiAdapter;
