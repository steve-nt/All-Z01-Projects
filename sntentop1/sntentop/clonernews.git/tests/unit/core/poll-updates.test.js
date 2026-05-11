/**
 * Purpose: Verify TA-4 polling semantics, throttling, and signal emissions with deterministic clocks.
 * Public API: Exercises createPollUpdatesUseCase via poll() and subscribe() only.
 * Constraints: Tests remain pure unit tests with mocked adapters and no network or DOM dependencies.
 */

import { afterEach, describe, expect, it, vi } from 'vitest';

import { createPollUpdatesUseCase } from '../../../src/core/use-cases/poll-updates.js';

// Helper keeps Result fixtures explicit and readable inside test scenarios.
const ok = (data) => ({ ok: true, data });
// Helper mirrors adapter error Results so failure paths stay realistic.
const err = (error) => ({ ok: false, error });

const createManualClock = (initialMs = 0) => {
  // Mutable clock state allows deterministic throttle boundary testing.
  let nowMs = initialMs;

  return {
    // now() is injected into the use case so tests can drive time progression precisely.
    now: () => nowMs,
    // advance() simulates elapsed time between polls without real timers.
    advance: (deltaMs) => {
      nowMs += deltaMs;
      return nowMs;
    },
  };
};

const createQueuedUpdatesApi = (results) => {
  // Cursor tracks sequential responses to model repeated polling cycles.
  let cursor = 0;

  return {
    getUpdates: vi.fn(async () => {
      // Clamp to the last result to keep behavior stable if poll is called extra times.
      const safeIndex = Math.min(cursor, Math.max(results.length - 1, 0));
      cursor += 1;
      return results[safeIndex];
    }),
  };
};

const createSignalHarness = () => {
  // value mirrors the latest signal payload to emulate minimal signal semantics.
  let value = null;
  // Set is used so duplicate subscriptions do not emit twice.
  const subscribers = new Set();
  // emissions captures every set() payload for assertion.
  const emissions = [];

  return {
    createSignal: vi.fn((initialValue) => {
      // Seed the signal state from the factory-provided initial payload.
      value = initialValue;

      return {
        set(nextValue) {
          // Support function updaters and direct payload assignment.
          const resolvedValue = typeof nextValue === 'function' ? nextValue(value) : nextValue;
          value = resolvedValue;
          emissions.push(resolvedValue);

          // Broadcast each resolved payload to all current subscribers.
          for (const subscriber of [...subscribers]) {
            subscriber(resolvedValue);
          }

          return value;
        },
        subscribe(callback) {
          subscribers.add(callback);
          return () => {
            subscribers.delete(callback);
          };
        },
      };
    }),
    emissions,
  };
};

const createHarness = ({ apiResults, startMs = 10_000, minIntervalMs = 5_000 } = {}) => {
  // Manual clock lets each test assert exact boundary behavior at 4,999ms and 5,000ms.
  const clock = createManualClock(startMs);
  // Queue-backed API provides deterministic responses for each poll call.
  const api = createQueuedUpdatesApi(apiResults);
  // Signal harness allows direct inspection of emitted payloads.
  const signalHarness = createSignalHarness();

  const useCase = createPollUpdatesUseCase({
    // API adapter dependency is mocked by the queue helper above.
    api,
    // nowMs keeps polledAtMs deterministic across all test runs.
    nowMs: clock.now,
    // minIntervalMs drives throttle behavior under test.
    minIntervalMs,
    // createSignal lets tests capture payload emissions without external state libs.
    createSignal: signalHarness.createSignal,
  });

  return { useCase, clock, api, signalHarness };
};

afterEach(() => {
  // Reset spies so call counts never leak across tests.
  vi.restoreAllMocks();
});

describe('poll-updates use case', () => {
  it('returns all update IDs on the first successful poll', async () => {
    const harness = createHarness({ apiResults: [ok([901, 902, 903])] });

    const result = await harness.useCase.poll();

    expect(harness.api.getUpdates).toHaveBeenCalledTimes(1);
    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [901, 902, 903],
          currentIds: [901, 902, 903],
          previousIds: [],
          polledAtMs: 10_000,
          isFirstPoll: true,
        }),
      }),
    );
  });

  it('returns only newly observed IDs on the second poll', async () => {
    const harness = createHarness({ apiResults: [ok([11, 12, 13]), ok([12, 13, 14, 15])] });

    const firstResult = await harness.useCase.poll();
    // Advance exactly one interval so the second call is not throttled.
    harness.clock.advance(5_000);
    const secondResult = await harness.useCase.poll();

    expect(firstResult).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [11, 12, 13],
          isFirstPoll: true,
        }),
      }),
    );
    expect(secondResult).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [14, 15],
          currentIds: [12, 13, 14, 15],
          previousIds: [11, 12, 13],
          polledAtMs: 15_000,
          isFirstPoll: false,
        }),
      }),
    );
  });

  it('marks only the first successful poll as first even when the first payload is empty', async () => {
    // This case reproduces the TA-4 regression where first-poll state was inferred from previousIds length.
    const harness = createHarness({ apiResults: [ok([]), ok([55])] });

    const firstResult = await harness.useCase.poll();
    harness.clock.advance(5_000);
    const secondResult = await harness.useCase.poll();

    expect(firstResult).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [],
          isFirstPoll: true,
        }),
      }),
    );
    expect(secondResult).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [55],
          previousIds: [],
          isFirstPoll: false,
        }),
      }),
    );
  });

  it('enforces the minimum 5-second throttle and returns an error Result inside the window', async () => {
    const harness = createHarness({ apiResults: [ok([71, 72]), ok([71, 72, 73])] });

    const firstResult = await harness.useCase.poll();
    // Move to just before the boundary to verify throttling.
    harness.clock.advance(4_999);
    const throttledResult = await harness.useCase.poll();
    // Move one additional millisecond to the exact boundary.
    harness.clock.advance(1);
    const boundaryResult = await harness.useCase.poll();

    expect(firstResult).toEqual(expect.objectContaining({ ok: true }));
    expect(throttledResult).toEqual(
      expect.objectContaining({
        ok: false,
        error: expect.stringMatching(/throttled/i),
      }),
    );
    expect(boundaryResult).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          newIds: [73],
          isFirstPoll: false,
        }),
      }),
    );
    // Network should run only for first and boundary calls.
    expect(harness.api.getUpdates).toHaveBeenCalledTimes(2);
  });

  it('emits payloads through subscribe and injected signal channels', async () => {
    const harness = createHarness({ apiResults: [ok([201, 202]), ok([201, 202, 203])] });
    const observedBySubscription = [];
    // Subscribe through the public API to verify downstream live-banner compatibility.
    const unsubscribe = harness.useCase.subscribe((payload) => {
      observedBySubscription.push(payload);
    });

    try {
      await harness.useCase.poll();
      harness.clock.advance(5_000);
      await harness.useCase.poll();
    } finally {
      unsubscribe();
    }

    // Signal harness receives payloads directly from updatesSignal.set.
    expect(harness.signalHarness.emissions).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ newIds: [201, 202], isFirstPoll: true }),
        expect.objectContaining({ newIds: [203], isFirstPoll: false }),
      ]),
    );
    // Public subscribe channel should observe the same semantic sequence.
    expect(observedBySubscription).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ newIds: [201, 202], isFirstPoll: true }),
        expect.objectContaining({ newIds: [203], isFirstPoll: false }),
      ]),
    );
  });

  it('propagates getUpdates failures as Result errors', async () => {
    const harness = createHarness({ apiResults: [err('updates-fetch-failed')] });

    const result = await harness.useCase.poll();

    expect(harness.api.getUpdates).toHaveBeenCalledTimes(1);
    expect(result).toEqual(
      expect.objectContaining({
        ok: false,
        error: expect.stringMatching(/updates-fetch-failed/i),
      }),
    );
  });

  it('does not require window or document to run polling logic', async () => {
    const harness = createHarness({ apiResults: [ok([301])] });
    // Capture original globals so we can restore test environment after assertion.
    const hadWindow = Object.hasOwn(globalThis, 'window');
    const hadDocument = Object.hasOwn(globalThis, 'document');
    const originalWindow = globalThis.window;
    const originalDocument = globalThis.document;

    // Remove DOM globals to prove core logic has no browser API coupling.
    Reflect.deleteProperty(globalThis, 'window');
    Reflect.deleteProperty(globalThis, 'document');

    try {
      const result = await harness.useCase.poll();
      expect(result).toEqual(
        expect.objectContaining({
          ok: true,
          data: expect.objectContaining({ newIds: [301] }),
        }),
      );
    } finally {
      // Restore globals exactly to prevent side effects in other tests.
      if (hadWindow) {
        globalThis.window = originalWindow;
      }

      if (hadDocument) {
        globalThis.document = originalDocument;
      }
    }
  });
});
