/*
 * Purpose: Coordinate live-update polling signals with the live-banner view without owning DOM creation.
 * Public API: createLiveBannerController({ pollUpdates, view, onRefresh?, pollIntervalMs? }) -> { mount(), refresh(), unmount() }.
 * Constraints: This module subscribes to the live-update polling contract, renders simple banner state, and keeps timer cleanup local to the controller.
 */

import { MIN_POLL_INTERVAL_MS } from '../../core/use-cases/poll-updates.js';
import { Temporal } from '@js-temporal/polyfill';

/**
 * @template T
 * @typedef {{ ok: true, data: T } | { ok: false, error: string }} Result
 */

/**
 * @typedef {Object} PollUpdatesPayload
 * @property {readonly number[]} newIds
 * @property {readonly number[]} currentIds
 * @property {readonly number[]} previousIds
 * @property {number} polledAtMs
 * @property {boolean} isFirstPoll
 */

/**
 * @typedef {Object} PollUpdatesContract
 * @property {() => Promise<Result<PollUpdatesPayload>>} poll
 * @property {(callback: (payload: PollUpdatesPayload) => void) => () => void} subscribe
 */

/**
 * @typedef {Object} LiveBannerViewContract
 * @property {(state: { isVisible: boolean, updateCount: number, elapsedLabel?: string }) => void} render
 */

/**
 * @typedef {Object} LiveBannerControllerDependencies
 * @property {PollUpdatesContract} pollUpdates
 * @property {LiveBannerViewContract} view
 * @property {(() => Promise<Result<unknown> | void>)=} onRefresh
 * @property {(() => Promise<Result<unknown> | void>)=} onClear
 * @property {number=} pollIntervalMs
 * @property {readonly number[]=} initialPendingUpdateIds
 * @property {number | null=} initialPendingSinceMs
 * @property {() => number=} nowMs
 */

const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });
const POLL_INTERVAL_SAFETY_MARGIN_MS = 100;

// Hidden state is shared so every reset path renders the same stable contract.
const HIDDEN_BANNER_STATE = Object.freeze({
  isVisible: false,
  updateCount: 0,
});

// Positive integer counts are the only values that should surface a visible banner.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Result success checks stay local so refresh can support both explicit Results and void callbacks.
const isOkResult = (value) =>
  value !== null && typeof value === 'object' && 'ok' in value && value.ok === true;

// Summary payload parsing uses integer normalization so malformed callback data cannot break rendering.
const toNonNegativeInteger = (value, fallback = 0) =>
  Number.isInteger(value) && value >= 0 ? value : fallback;

// Unexpected thrown values are normalized so refresh callers always receive a Result on failure.
const toErrorMessage = (value) => {
  if (value instanceof Error && value.message.length > 0) {
    return value.message;
  }

  if (typeof value === 'string' && value.length > 0) {
    return value;
  }

  return 'Unknown error.';
};

const toNowMs = () => Temporal.Now.instant().epochMilliseconds;

const toElapsedLabel = (elapsedMs) => {
  const elapsedSeconds = Math.max(1, Math.floor(elapsedMs / 1000));

  if (elapsedSeconds < 60) {
    return `${elapsedSeconds} second${elapsedSeconds === 1 ? '' : 's'}`;
  }

  const elapsedMinutes = Math.floor(elapsedSeconds / 60);

  if (elapsedMinutes < 60) {
    return `${elapsedMinutes} minute${elapsedMinutes === 1 ? '' : 's'}`;
  }

  const elapsedHours = Math.floor(elapsedMinutes / 60);
  return `${elapsedHours} hour${elapsedHours === 1 ? '' : 's'}`;
};

// Visible-state creation stays centralized so count-driven rendering never duplicates shape logic.
const createVisibleBannerState = (updateCount, elapsedLabel) =>
  Object.freeze({
    isVisible: true,
    updateCount,
    elapsedLabel,
  });

// Summary copy is centralized so post-refresh messaging stays consistent across controller call sites.
const toRefreshSummaryMessage = ({ movedCount, newCount, updatedCount }) =>
  `Moved ${movedCount} refreshed post${movedCount === 1 ? '' : 's'} to top (${newCount} new, ${updatedCount} updated).`;

// Refresh summary state keeps clear visible while hiding refresh to avoid accidental duplicate reloads.
const createRefreshSummaryBannerState = ({ movedCount, newCount, updatedCount }) =>
  Object.freeze({
    isVisible: true,
    updateCount: Math.max(1, movedCount),
    messageText: toRefreshSummaryMessage({ movedCount, newCount, updatedCount }),
    showRefresh: false,
    showClear: true,
  });

// Refresh callbacks can optionally return count details; missing data falls back to pending update count.
const toRefreshSummaryCounts = (refreshResult, fallbackMovedCount) => {
  if (!isOkResult(refreshResult)) {
    return null;
  }

  const payload = refreshResult.data;

  if (payload === undefined || payload === null || typeof payload !== 'object') {
    return {
      movedCount: toNonNegativeInteger(fallbackMovedCount),
      newCount: toNonNegativeInteger(fallbackMovedCount),
      updatedCount: 0,
    };
  }

  const movedCount = toNonNegativeInteger(payload.movedCount, toNonNegativeInteger(fallbackMovedCount));
  const updatedCount = toNonNegativeInteger(payload.updatedCount, 0);
  const normalizedUpdatedCount = Math.min(updatedCount, movedCount);
  const derivedNewCount = Math.max(0, movedCount - normalizedUpdatedCount);
  const newCount = toNonNegativeInteger(payload.newCount, derivedNewCount);
  const normalizedNewCount = Math.min(newCount, movedCount);

  // Count totals are clamped so the summary sentence always remains internally consistent.
  const normalizedMovedCount = Math.max(movedCount, normalizedNewCount + normalizedUpdatedCount);

  return {
    movedCount: normalizedMovedCount,
    newCount: normalizedNewCount,
    updatedCount: normalizedUpdatedCount,
  };
};

/**
 * @param {LiveBannerControllerDependencies} dependencies
 * @returns {{ mount(): void, refresh(): Promise<Result<void>>, clear(): Promise<Result<void>>, unmount(): void }}
 */
export const createLiveBannerController = ({
  pollUpdates,
  view,
  onRefresh = async () => ok(undefined),
  onClear = async () => ok(undefined),
  pollIntervalMs = MIN_POLL_INTERVAL_MS,
  initialPendingUpdateIds = [],
  initialPendingSinceMs = null,
  nowMs = toNowMs,
} = {}) => {
  // Mount tracking keeps repeated mount calls from creating duplicate subscriptions or timers.
  let isMounted = false;
  // Pending update ids stay latched so the banner remains visible until the user explicitly refreshes.
  let pendingUpdateIds = new Set(
    initialPendingUpdateIds.filter((value) => Number.isInteger(value) && value > 0),
  );
  // The first-seen timestamp is used for "in the last X" elapsed messaging.
  let pendingSinceMs = Number.isInteger(initialPendingSinceMs) ? initialPendingSinceMs : null;
  // The unsubscribe handle is stored so unmount can detach from the live-update stream cleanly.
  let unsubscribe = null;
  // The interval identifier is stored so timer cleanup stays deterministic on unmount.
  let pollTimerId = null;
  // A second timer keeps elapsed-time copy moving from seconds to minutes while updates are pending.
  let elapsedTimerId = null;

  // Rendering stays behind one helper so every state transition shares the same defensive gate.
  const render = (state) => {
    if (typeof view?.render !== 'function') {
      return;
    }

    view.render(state);
  };

  const renderCurrentState = () => {
    if (!isPositiveInteger(pendingUpdateIds.size)) {
      render(HIDDEN_BANNER_STATE);
      return;
    }

    if (!isPositiveInteger(pendingSinceMs)) {
      pendingSinceMs = nowMs();
    }

    const elapsedMs = Math.max(0, nowMs() - pendingSinceMs);
    render(createVisibleBannerState(pendingUpdateIds.size, toElapsedLabel(elapsedMs)));
  };

  // Signal payload translation keeps update-diff semantics out of the view layer.
  const handlePollUpdates = (payload) => {
    // First polls seed baseline history and should never surface historical updates in the banner.
    if (payload?.isFirstPoll === true) {
      renderCurrentState();
      return;
    }

    // Only array payloads contribute new unseen ids to the latched pending-update state.
    const newIds = Array.isArray(payload?.newIds) ? payload.newIds : [];

    // Newly discovered ids are merged into the pending set so later empty polls do not dismiss the banner.
    for (const id of newIds) {
      if (!isPositiveInteger(id)) {
        continue;
      }

      if (!isPositiveInteger(pendingUpdateIds.size)) {
        pendingSinceMs = Number.isInteger(payload?.polledAtMs) ? payload.polledAtMs : nowMs();
      }

      pendingUpdateIds.add(id);
    }

    renderCurrentState();
  };

  const runPollTick = () => {
    if (typeof pollUpdates?.poll !== 'function') {
      return;
    }

    try {
      const pollPromise = pollUpdates.poll();

      // Rejected poll promises are swallowed here because the banner should not crash the app shell.
      if (pollPromise && typeof pollPromise.catch === 'function') {
        void pollPromise.catch(() => {});
      }
    } catch {
      // Synchronous poll failures are ignored because retry happens on the next interval tick.
    }
  };

  return Object.freeze({
    // Mount renders the initial hidden state, subscribes to update notifications, and starts periodic polling.
    mount() {
      if (isMounted) {
        return;
      }

      isMounted = true;
      renderCurrentState();

      if (typeof pollUpdates?.subscribe === 'function') {
        unsubscribe = pollUpdates.subscribe(handlePollUpdates);
      }

      // A small safety margin avoids scheduler jitter causing effective intervals to dip below 5 seconds.
      const normalizedIntervalMs = Math.max(
        pollIntervalMs,
        MIN_POLL_INTERVAL_MS + POLL_INTERVAL_SAFETY_MARGIN_MS,
      );

      if (isPositiveInteger(normalizedIntervalMs)) {
        pollTimerId = setInterval(runPollTick, normalizedIntervalMs);
      }

      elapsedTimerId = setInterval(() => {
        if (isPositiveInteger(pendingUpdateIds.size)) {
          renderCurrentState();
        }
      }, 1_000);
    },

    // Refresh delegates actual feed reload logic upward and clears the banner only after success.
    async refresh() {
      if (typeof onRefresh !== 'function') {
        return err('Live banner controller requires an onRefresh callback.');
      }

      try {
        const pendingCountBeforeRefresh = pendingUpdateIds.size;
        const refreshResult = await onRefresh();

        // Void refresh callbacks are treated as success so shell integrations can stay lightweight.
        if (refreshResult === undefined || isOkResult(refreshResult)) {
          // Successful refresh acknowledges every pending update currently surfaced by the banner.
          pendingUpdateIds = new Set();
          pendingSinceMs = null;

          const refreshSummaryCounts = toRefreshSummaryCounts(
            refreshResult,
            pendingCountBeforeRefresh,
          );

          if (refreshSummaryCounts !== null && refreshSummaryCounts.movedCount > 0) {
            render(createRefreshSummaryBannerState(refreshSummaryCounts));
          } else {
            render(HIDDEN_BANNER_STATE);
          }

          return ok(undefined);
        }

        // Failed Results are passed through unchanged so callers can decide how to surface them.
        return refreshResult;
      } catch (unexpectedError) {
        return err(`Live banner refresh failed: ${toErrorMessage(unexpectedError)}`);
      }
    },

    // Clear dismisses pending updates without triggering a feed reload.
    async clear() {
      if (typeof onClear !== 'function') {
        return err('Live banner controller requires an onClear callback.');
      }

      try {
        const clearResult = await onClear();

        if (clearResult === undefined || isOkResult(clearResult)) {
          pendingUpdateIds = new Set();
          pendingSinceMs = null;
          render(HIDDEN_BANNER_STATE);
          return ok(undefined);
        }

        return clearResult;
      } catch (unexpectedError) {
        return err(`Live banner clear failed: ${toErrorMessage(unexpectedError)}`);
      }
    },

    // Unmount detaches all external resources but intentionally avoids an extra render side effect.
    unmount() {
      if (!isMounted) {
        return;
      }

      isMounted = false;
      // Unmount clears latched update state so remounts start from a clean baseline.
      pendingUpdateIds = new Set();

      if (pollTimerId !== null) {
        clearInterval(pollTimerId);
        pollTimerId = null;
      }

      if (elapsedTimerId !== null) {
        clearInterval(elapsedTimerId);
        elapsedTimerId = null;
      }

      unsubscribe?.();
      unsubscribe = null;
    },
  });
};
