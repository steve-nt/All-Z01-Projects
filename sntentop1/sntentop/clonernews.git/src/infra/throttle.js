/*
 * Purpose: Provide timing utilities for rate-limiting repeated function calls.
 * Public API: throttle(fn, ms) and debounce(fn, ms).
 * Implementation notes: The wrappers keep their own timer state and never touch DOM or network APIs.
 */

export const throttle = (fn, ms) => {
  let isThrottled = false;

  return (...args) => {
    // Ignore calls inside the active window so the wrapped function cannot fire too often.
    if (isThrottled) {
      return;
    }

    // Invoke immediately so the first call stays responsive.
    fn(...args);
    isThrottled = true;

    // Release the gate after the cooldown so the next allowed call can proceed.
    setTimeout(() => {
      isThrottled = false;
    }, ms);
  };
};

export const debounce = (fn, ms) => {
  let timeoutId;

  return (...args) => {
    // Reset the timer on every call so only the final burst survives.
    clearTimeout(timeoutId);

    timeoutId = setTimeout(() => {
      // Defer execution until the caller stops invoking the wrapper.
      fn(...args);
    }, ms);
  };
};
