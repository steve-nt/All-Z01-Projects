/**
 * Small global store for framework apps.
 *
 * setState semantics:
 * - updater function: setState((prev) => nextState)
 * - direct value/object: setState(nextState)
 * - plain object + plain object => shallow merge
 * - otherwise => full replace
 *
 * Immutability expectation:
 * - treat state as immutable in app code
 * - never mutate state in place (objects/arrays)
 * - always pass a new value/object to setState
 */
export function createStore(initialState) {
  let state = initialState;
  /** @type {Set<() => void>} */
  const listeners = new Set();

  const isPlainObject = (value) =>
    value !== null &&
    typeof value === "object" &&
    Object.getPrototypeOf(value) === Object.prototype;

  function getState() {
    return state;
  }

  function setState(next) {
    const prevState = state;
    const resolvedNext = typeof next === "function" ? next(prevState) : next;
    const shouldMerge = isPlainObject(prevState) && isPlainObject(resolvedNext);
    const nextState = shouldMerge
      ? { ...prevState, ...resolvedNext }
      : resolvedNext;
    state = nextState;

    // Notify a snapshot so listener mutations during emit are safe.
    [...listeners].forEach((listener) => {
      listener();
    });
  }

  function subscribe(listener) {
    if (typeof listener !== "function") {
      throw new TypeError("store.subscribe(listener): listener must be a function.");
    }

    listeners.add(listener);

    // Idempotent unsubscribe to support safe mount/unmount cleanup.
    return function unsubscribe() {
      listeners.delete(listener);
    };
  }

  return {
    getState,
    setState,
    subscribe,
  };
}
