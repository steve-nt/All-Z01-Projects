/**
 * Purpose: Define the TC-5 controller contract for polling-driven banner visibility and refresh resets.
 * Public API: createLiveBannerController({ pollUpdates, view, onRefresh }) -> { mount, refresh, unmount }.
 * Constraints: Tests stay pure, mock poll-updates and view dependencies directly, and avoid DOM or app-shell wiring.
 */

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { createLiveBannerController } from '../../../src/features/live-banner/live-banner.controller.js';

// Result helper keeps async refresh and poll fixtures aligned with the repo-wide explicit Result contract.
const ok = (data) => ({ ok: true, data });

// Hidden state is shared across scenarios so all non-visible assertions target the same public shape.
const createHiddenBannerState = () =>
  Object.freeze({
    isVisible: false,
    updateCount: 0,
  });

// Visible state assertions pin count and visibility while allowing elapsed labels to evolve by timer ticks.
const expectVisibleBannerState = (renderSpy, updateCount) => {
  expect(renderSpy).toHaveBeenLastCalledWith(
    expect.objectContaining({
      isVisible: true,
      updateCount,
      elapsedLabel: expect.any(String),
    }),
  );
};

// Summary assertions pin the copy contract shown after refresh acknowledges pending updates.
const expectSummaryBannerState = (renderSpy, messageText) => {
  expect(renderSpy).toHaveBeenLastCalledWith(
    expect.objectContaining({
      isVisible: true,
      messageText,
      showRefresh: false,
      showClear: true,
    }),
  );
};

// Poll payload helper mirrors the TA-4 shape so controller tests stay compatible with the core contract.
const createPollPayload = ({
  newIds = [],
  currentIds = [],
  previousIds = [],
  polledAtMs = 10_000,
  isFirstPoll = false,
} = {}) =>
  Object.freeze({
    newIds: Object.freeze([...newIds]),
    currentIds: Object.freeze([...currentIds]),
    previousIds: Object.freeze([...previousIds]),
    polledAtMs,
    isFirstPoll,
  });

// Harness centralizes the poll-updates subscription channel so tests can emit payloads deterministically.
const createPollUpdatesHarness = () => {
  let subscriber = null;
  const unsubscribe = vi.fn(() => {
    subscriber = null;
  });

  return {
    pollUpdates: {
      // poll is mocked so timer-driven tests can verify cadence without any core implementation dependency.
      poll: vi.fn(async () => ok(createPollPayload())),
      // subscribe captures the active subscriber and returns an unsubscribe function for cleanup assertions.
      subscribe: vi.fn((callback) => {
        subscriber = callback;
        return unsubscribe;
      }),
    },
    // emit lets each test push synthetic TA-4 payloads through the subscribed controller callback.
    emit(payload) {
      subscriber?.(payload);
    },
    unsubscribe,
  };
};

// View harness keeps controller tests focused on state translation rather than DOM rendering details.
const createViewHarness = () => ({
  render: vi.fn(),
});

beforeEach(() => {
  // Fake timers make the five-second polling contract deterministic and fast to assert.
  vi.useFakeTimers();
});

afterEach(() => {
  // Timer restoration prevents one controller test from leaking scheduled work into the next one.
  vi.useRealTimers();
  // Spy restoration keeps call counts isolated across all test cases.
  vi.restoreAllMocks();
});

describe('live-banner controller', () => {
  it('subscribes to poll-updates and starts a five-second polling timer on mount', async () => {
    // Mount should establish the signal subscription and the periodic poll loop without showing stale UI.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
    });

    controller.mount();

    expect(pollUpdatesHarness.pollUpdates.subscribe).toHaveBeenCalledTimes(1);
    expect(view.render).toHaveBeenCalledWith(createHiddenBannerState());

    await vi.advanceTimersByTimeAsync(5_099);

    expect(pollUpdatesHarness.pollUpdates.poll).not.toHaveBeenCalled();

    await vi.advanceTimersByTimeAsync(1);

    expect(pollUpdatesHarness.pollUpdates.poll).toHaveBeenCalledTimes(1);

    controller.unmount();
  });

  it('ignores the first poll payload even when it contains update ids', () => {
    // The first successful poll seeds baseline state and should never announce historical updates to the user.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
    });

    controller.mount();
    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [701, 702],
        currentIds: [701, 702],
        isFirstPoll: true,
      }),
    );

    expect(view.render).toHaveBeenLastCalledWith(createHiddenBannerState());

    controller.unmount();
  });

  it('shows the banner only when a non-first poll includes diff ids', () => {
    // Empty diffs should keep the banner hidden so the UI only interrupts when there is actionable freshness.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
    });

    controller.mount();
    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [],
        currentIds: [810, 811],
        previousIds: [810, 811],
        isFirstPoll: false,
      }),
    );

    expect(view.render).toHaveBeenLastCalledWith(createHiddenBannerState());

    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [812, 813, 814],
        currentIds: [810, 811, 812, 813, 814],
        previousIds: [810, 811],
        isFirstPoll: false,
      }),
    );

    expectVisibleBannerState(view.render, 3);

    controller.unmount();
  });

  it('keeps the banner visible after later zero-diff polls until refresh clears pending updates', () => {
    // Once updates are announced, later polls without diffs should not silently dismiss the pending refresh action.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
    });

    controller.mount();
    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [820, 821],
        currentIds: [810, 811, 820, 821],
        previousIds: [810, 811],
        isFirstPoll: false,
      }),
    );

    expectVisibleBannerState(view.render, 2);

    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [],
        currentIds: [810, 811, 820, 821],
        previousIds: [810, 811, 820, 821],
        isFirstPoll: false,
      }),
    );

    expectVisibleBannerState(view.render, 2);

    controller.unmount();
  });

  it('renders a refresh summary message after refresh completes', async () => {
    // Refresh should acknowledge pending updates with explicit moved/new/updated counts.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const onRefresh = vi.fn(async () =>
      ok({
        movedCount: 2,
        newCount: 0,
        updatedCount: 2,
      }),
    );
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh,
    });

    controller.mount();
    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [901, 902],
        currentIds: [901, 902],
        previousIds: [900],
        isFirstPoll: false,
      }),
    );

    expectVisibleBannerState(view.render, 2);

    await controller.refresh();

    expect(onRefresh).toHaveBeenCalledTimes(1);
    expectSummaryBannerState(view.render, 'Moved 2 refreshed posts to top (0 new, 2 updated).');

    controller.unmount();
  });

  it('resets the banner after clear completes without requiring refresh', async () => {
    // Clear should dismiss pending updates while leaving feed data untouched.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const onClear = vi.fn(async () => ok(undefined));
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
      onClear,
    });

    controller.mount();
    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [931, 932],
        currentIds: [931, 932],
        previousIds: [930],
        isFirstPoll: false,
      }),
    );

    expectVisibleBannerState(view.render, 2);

    await controller.clear();

    expect(onClear).toHaveBeenCalledTimes(1);
    expect(view.render).toHaveBeenLastCalledWith(createHiddenBannerState());

    controller.unmount();
  });

  it('cleans up the poll subscription and timer on unmount', async () => {
    // Unmount must stop future polling work so the feature can be mounted and removed without leaks.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
    });

    controller.mount();
    controller.unmount();

    expect(pollUpdatesHarness.unsubscribe).toHaveBeenCalledTimes(1);

    await vi.advanceTimersByTimeAsync(15_000);

    expect(pollUpdatesHarness.pollUpdates.poll).not.toHaveBeenCalled();

    pollUpdatesHarness.emit(
      createPollPayload({
        newIds: [999],
        currentIds: [999],
        isFirstPoll: false,
      }),
    );

    expect(view.render).toHaveBeenCalledTimes(1);
  });

  it('hydrates a previously pending banner state on mount', () => {
    // Route remounts should keep the pending counter visible until refresh acknowledges it.
    const pollUpdatesHarness = createPollUpdatesHarness();
    const view = createViewHarness();
    const controller = createLiveBannerController({
      pollUpdates: pollUpdatesHarness.pollUpdates,
      view,
      onRefresh: vi.fn(async () => ok(undefined)),
      initialPendingUpdateIds: [1001, 1002, 1003],
      initialPendingSinceMs: 9_000,
      nowMs: () => 11_000,
    });

    controller.mount();

    expectVisibleBannerState(view.render, 3);

    controller.unmount();
  });
});
