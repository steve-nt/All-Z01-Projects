/**
 * Purpose: Poll the HN updates endpoint and emit deterministic diff payloads for consumers.
 * Public API: createPollUpdatesUseCase(deps) -> { poll, subscribe } where poll returns a Result payload.
 * Constraints: Core-only module with zero DOM/fetch access; all effects are delegated via injected adapters.
 */

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/** @typedef {import('../interfaces/api-interface.js').ApiInterface} ApiInterface */

/**
 * @typedef {Object} PollUpdatesPayload
 * @property {readonly number[]} newIds
 * @property {readonly number[]} currentIds
 * @property {readonly number[]} previousIds
 * @property {number} polledAtMs
 * @property {boolean} isFirstPoll
 */

/**
 * @typedef {Object} PollUpdatesSignal
 * @property {(nextValue: PollUpdatesPayload | ((currentValue: PollUpdatesPayload) => PollUpdatesPayload)) => PollUpdatesPayload} set
 * @property {(callback: (value: PollUpdatesPayload) => void) => () => void} subscribe
 */

/**
 * @typedef {Object} PollUpdatesUseCaseDependencies
 * @property {ApiInterface} api
 * @property {number=} minIntervalMs
 * @property {() => number=} nowMs
 * @property {(initialValue: PollUpdatesPayload) => PollUpdatesSignal=} createSignal
 */

/**
 * @typedef {Object} PollUpdatesUseCase
 * @property {() => Promise<Result<PollUpdatesPayload>>} poll
 * @property {(callback: (value: PollUpdatesPayload) => void) => () => void} subscribe
 */

export const MIN_POLL_INTERVAL_MS = 5000;

const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

const toErrorMessage = (value) => {
  if (value instanceof Error && value.message.length > 0) {
    return value.message;
  }

  if (typeof value === 'string' && value.length > 0) {
    return value;
  }

  return 'Unknown error.';
};

const createFallbackSignal = (initialValue) => {
  let value = initialValue;
  // A Set prevents duplicate subscriptions while keeping deletion O(1).
  const subscribers = new Set();

  return {
    set(nextValue) {
      // Support both direct value writes and updater callbacks to match signal contracts.
      const resolvedValue = typeof nextValue === 'function' ? nextValue(value) : nextValue;
      value = resolvedValue;

      for (const subscriber of [...subscribers]) {
        subscriber(value);
      }

      return value;
    },
    subscribe(callback) {
      // Non-function subscribers are ignored to preserve a safe no-op contract.
      if (typeof callback !== 'function') {
        return () => {};
      }

      subscribers.add(callback);

      return () => subscribers.delete(callback);
    },
  };
};

const defaultNowMs = () => {
  // Temporal is the primary source because it yields stable epoch milliseconds across runtimes.
  if (typeof Temporal !== 'undefined' && typeof Temporal.Now?.instant === 'function') {
    return Temporal.Now.instant().epochMilliseconds;
  }

  // When Temporal is unavailable, reconstruct epoch milliseconds from performance clocks.
  if (typeof performance !== 'undefined' && typeof performance.now === 'function') {
    const timeOrigin =
      typeof performance.timeOrigin === 'number' && Number.isFinite(performance.timeOrigin)
        ? performance.timeOrigin
        : 0;

    return timeOrigin + performance.now();
  }

  // Keep fallback deterministic for tests; callers can inject a stronger clock via nowMs.
  return 0;
};

const createPayload = ({ newIds, currentIds, previousIds, polledAtMs, isFirstPoll }) =>
  Object.freeze({
    newIds: Object.freeze([...newIds]),
    currentIds: Object.freeze([...currentIds]),
    previousIds: Object.freeze([...previousIds]),
    polledAtMs,
    isFirstPoll,
  });

const INITIAL_PAYLOAD = createPayload({
  newIds: [],
  currentIds: [],
  previousIds: [],
  polledAtMs: -1,
  isFirstPoll: true,
});

/**
 * @param {PollUpdatesUseCaseDependencies} dependencies
 * @returns {PollUpdatesUseCase}
 */
export const createPollUpdatesUseCase = ({
  api,
  minIntervalMs = MIN_POLL_INTERVAL_MS,
  nowMs = defaultNowMs,
  createSignal = createFallbackSignal,
} = {}) => {
  // The floor prevents accidental violation of the hard 5-second policy.
  const pollIntervalMs = Math.max(minIntervalMs, MIN_POLL_INTERVAL_MS);

  const updatesSignal =
    typeof createSignal === 'function'
      ? createSignal(INITIAL_PAYLOAD)
      : createFallbackSignal(INITIAL_PAYLOAD);

  let lastPollAtMs = null;
  // previousIds stores last successful payload to support diff computation.
  let previousIds = [];
  // previousIdSet keeps diff checks linear over current IDs.
  let previousIdSet = new Set();
  // hasPolled avoids inferring "first poll" from array emptiness, which fails on empty responses.
  let hasPolled = false;

  const poll = async () => {
    if (!api || typeof api.getUpdates !== 'function') {
      return err('Poll updates use case requires an API adapter with a getUpdates method.');
    }

    if (!isPositiveInteger(pollIntervalMs)) {
      return err('Poll interval must be a positive integer in milliseconds.');
    }

    if (typeof nowMs !== 'function') {
      return err('Poll updates use case requires a nowMs clock function.');
    }

    if (!updatesSignal || typeof updatesSignal.set !== 'function') {
      return err('Poll updates use case requires a signal with a set method.');
    }

    const currentTimeMs = nowMs();

    if (!Number.isFinite(currentTimeMs)) {
      return err('Clock function must return a finite timestamp in milliseconds.');
    }

    if (lastPollAtMs !== null && currentTimeMs - lastPollAtMs < pollIntervalMs) {
      return err(`Polling is throttled. Wait at least ${pollIntervalMs} ms between polls.`);
    }

    lastPollAtMs = currentTimeMs;

    try {
      const updatesResult = await api.getUpdates();

      if (!updatesResult.ok) {
        return err(updatesResult.error);
      }

      if (!Array.isArray(updatesResult.data)) {
        return err('Updates result payload is invalid. Expected an array of item IDs.');
      }

      if (!updatesResult.data.every((id) => isPositiveInteger(id))) {
        return err('Updates result payload contains invalid item IDs.');
      }

      const currentIds = [...updatesResult.data];
      // First-poll semantics depend on successful history, not previous payload cardinality.
      const isFirstPoll = !hasPolled;
      const newIds = currentIds.filter((id) => !previousIdSet.has(id));

      const payload = createPayload({
        newIds,
        currentIds,
        previousIds,
        polledAtMs: currentTimeMs,
        isFirstPoll,
      });

      updatesSignal.set(payload);

      // Commit state only after successful payload publication so retries stay consistent.
      previousIds = [...currentIds];
      previousIdSet = new Set(currentIds);
      hasPolled = true;

      return ok(payload);
    } catch (unexpectedError) {
      return err(`Failed to poll updates: ${toErrorMessage(unexpectedError)}`);
    }
  };

  const subscribe = (callback) => {
    if (!updatesSignal || typeof updatesSignal.subscribe !== 'function') {
      return () => {};
    }

    return updatesSignal.subscribe(callback);
  };

  return Object.freeze({
    poll,
    subscribe,
  });
};
