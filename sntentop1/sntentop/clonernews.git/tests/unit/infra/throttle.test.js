import { afterEach, describe, expect, it, vi } from 'vitest';

import { debounce, throttle } from '../../../src/infra/throttle.js';

afterEach(() => {
  vi.useRealTimers();
  vi.restoreAllMocks();
});

describe('throttle utility', () => {
  it('runs immediately on the first call', () => {
    vi.useFakeTimers();

    const fn = vi.fn();
    const throttled = throttle(fn, 100);

    throttled('first');

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('first');
  });

  it('blocks repeated calls during the throttle window', () => {
    vi.useFakeTimers();

    const fn = vi.fn();
    const throttled = throttle(fn, 100);

    throttled('first');
    throttled('second');
    throttled('third');

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('first');
  });

  it('allows a later call after the throttle window expires', () => {
    vi.useFakeTimers();

    const fn = vi.fn();
    const throttled = throttle(fn, 100);

    throttled('first');
    vi.advanceTimersByTime(100);
    throttled('second');

    expect(fn).toHaveBeenCalledTimes(2);
    expect(fn).toHaveBeenNthCalledWith(1, 'first');
    expect(fn).toHaveBeenNthCalledWith(2, 'second');
  });
});

describe('debounce utility', () => {
  it('waits until the quiet period ends before running', () => {
    vi.useFakeTimers();

    const fn = vi.fn();
    const debounced = debounce(fn, 100);

    debounced('first');

    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(100);

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('first');
  });

  it('forgets earlier calls and only runs the last one', () => {
    vi.useFakeTimers();

    const fn = vi.fn();
    const debounced = debounce(fn, 100);

    debounced('first');
    vi.advanceTimersByTime(50);
    debounced('second');
    vi.advanceTimersByTime(50);
    debounced('third');
    vi.advanceTimersByTime(99);

    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(1);

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('third');
  });
});
