// Signal tests prove the tiny reactive primitive still behaves predictably under dependency changes.
// Public API under test: createSignal, createEffect, and createComputed semantics.
// Constraints: tests stay synchronous and deterministic so dependency tracking regressions are easy to isolate.
import { describe, expect, it } from 'vitest';

import { createComputed, createEffect, createSignal } from '../../../src/shared/signals.js';

describe('signals', () => {
  it('updates a signal value through get and set', () => {
    // Basic read/write coverage protects the simplest state flow from regressions.
    const count = createSignal(1);

    expect(count.get()).toBe(1);

    count.set(3);

    expect(count.get()).toBe(3);
  });

  it('re-runs effects when dependencies change', () => {
    // Effects must re-run only when their source signal changes so the dependency tracker stays correct.
    const count = createSignal(1);
    const values = [];

    createEffect(() => {
      values.push(count.get());
    });

    count.set(2);

    expect(values).toEqual([1, 2]);
  });

  it('derives computed values from signals', () => {
    // Computed values should stay in sync automatically without manual subscriptions.
    const count = createSignal(2);
    const doubled = createComputed(() => count.get() * 2);

    expect(doubled.get()).toBe(4);

    count.set(5);

    expect(doubled.get()).toBe(10);
  });
});
