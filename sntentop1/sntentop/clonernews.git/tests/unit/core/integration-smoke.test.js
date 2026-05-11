// These smoke tests verify real Track A wiring across adapter and use-case boundaries with fetch mocked at the network edge.
// Public coverage focuses on TA-5 scenarios plus a TA-3 integration path to prove comment-tree and parent-link behavior.
// Constraints: keep tests deterministic, non-DOM, and strictly Result-contract based.
import { afterEach, describe, expect, it, vi } from 'vitest';

import { createGetItemUseCase } from '../../../src/core/use-cases/get-item.js';
import { createListItemsUseCase } from '../../../src/core/use-cases/list-items.js';
import { createPollUpdatesUseCase } from '../../../src/core/use-cases/poll-updates.js';
import { createHnApiAdapter } from '../../../src/infra/hn-api-adapter.js';

// Keep endpoint matching centralized so queue fixtures stay concise across scenarios.
const HN_API_PATH_PREFIX = '/v0';

// Response helper mirrors the minimal fetch Response surface the adapter inspects.
const createJsonResponse = (payload, options = {}) => {
  // Default status and content-type emulate healthy HN API responses.
  const status = options.status ?? 200;
  const contentType = options.contentType ?? 'application/json; charset=utf-8';

  return {
    // ok mirrors fetch semantics so adapter status handling is exercised authentically.
    ok: status >= 200 && status < 300,
    status,
    headers: {
      // Header lookup stays case-insensitive to match real Response headers behavior.
      get(name) {
        return name.toLowerCase() === 'content-type' ? contentType : null;
      },
    },
    // json() returns queued payloads asynchronously to preserve adapter async flow.
    json: vi.fn(async () => payload),
  };
};

// Story fixture includes only fields required by entity validation and list sorting assertions.
const createStoryItem = ({ id, time }) => ({
  id,
  type: 'story',
  time,
  title: `Story ${id}`,
  by: `user-${id}`,
});

// Comment fixture keeps IDs, timestamps, and parent links deterministic for tree assertions.
const createCommentItem = ({ id, time, kids = [] }) => ({
  id,
  type: 'comment',
  time,
  text: `Comment ${id}`,
  by: `user-${id}`,
  kids,
});

// URL helper ensures fixture path keys exactly match adapter endpoint construction.
const itemPath = (id) => `${HN_API_PATH_PREFIX}/item/${id}.json`;
const topStoriesPath = `${HN_API_PATH_PREFIX}/topstories.json`;
const updatesPath = `${HN_API_PATH_PREFIX}/updates.json`;

// Queue-backed fetch mock models ordered endpoint calls and fails loudly on unexpected requests.
const createQueuedFetchMock = (queuedPayloadsByPath) => {
  // Clone each path queue so each test consumes payloads without mutating fixture literals.
  const queueByPath = new Map(
    Object.entries(queuedPayloadsByPath).map(([path, payloads]) => [path, [...payloads]]),
  );

  return vi.fn(async (input) => {
    // Parse request path so assertions are decoupled from protocol/host differences.
    const requestPath = new URL(String(input)).pathname;
    const queue = queueByPath.get(requestPath);

    // Throwing here protects tests from silently passing when adapter routing regresses.
    if (!queue || queue.length === 0) {
      throw new Error(`Unexpected fetch request for ${requestPath}.`);
    }

    // Consume payloads FIFO to represent repeated calls against the same endpoint.
    const nextPayload = queue.shift();
    return createJsonResponse(nextPayload);
  });
};

// Path extraction keeps endpoint assertions readable without repeating URL parsing logic.
const getRequestedPaths = (fetchFn) =>
  fetchFn.mock.calls.map(([input]) => new URL(String(input)).pathname);

afterEach(() => {
  // Reset mocks after each scenario so cross-test call history cannot mask regressions.
  vi.restoreAllMocks();
});

describe('integration smoke: adapter + use-cases', () => {
  it('fresh load fetches feed ids, batch-fetches items, and preserves top feed API order', async () => {
    // Queue one feed-id response plus one payload for each requested story item.
    const fetchFn = createQueuedFetchMock({
      [topStoriesPath]: [[11, 12, 13, 14]],
      [itemPath(11)]: [createStoryItem({ id: 11, time: 100 })],
      [itemPath(12)]: [createStoryItem({ id: 12, time: 500 })],
      [itemPath(13)]: [createStoryItem({ id: 13, time: 300 })],
      [itemPath(14)]: [createStoryItem({ id: 14, time: 500 })],
    });

    // Use the real adapter/use-case composition to validate cross-module behavior.
    const api = createHnApiAdapter({ fetchFn });
    const listItems = createListItemsUseCase({ api });

    // Request a full four-item page to exercise sorting and hasMore semantics.
    const result = await listItems({ type: 'top', page: 1, limit: 4 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.hasMore).toBe(false);
    // Top feed should keep Hacker News API ranking order for the requested page.
    expect(result.data.items.map((item) => item.id)).toEqual([11, 12, 13, 14]);

    const requestedPaths = getRequestedPaths(fetchFn);
    // Feed list must be fetched once for the requested page calculation.
    expect(requestedPaths.filter((path) => path === topStoriesPath)).toHaveLength(1);
    // Exactly four item fetches should occur for the four requested feed IDs.
    expect(
      requestedPaths.filter((path) => path.startsWith(`${HN_API_PATH_PREFIX}/item/`)),
    ).toHaveLength(4);
  });

  it('second load reuses adapter cache so item network fetches are not repeated', async () => {
    // Reusing the same ID list across two loads allows cache-hit behavior to be observed directly.
    const ids = [21, 22, 23];
    const fetchFn = createQueuedFetchMock({
      [topStoriesPath]: [ids, ids],
      [itemPath(21)]: [createStoryItem({ id: 21, time: 400 })],
      [itemPath(22)]: [createStoryItem({ id: 22, time: 300 })],
      [itemPath(23)]: [createStoryItem({ id: 23, time: 200 })],
    });

    const api = createHnApiAdapter({ fetchFn });
    const listItems = createListItemsUseCase({ api });

    // First run populates adapter cache, second run should reuse cached item payloads.
    const firstResult = await listItems({ type: 'top', page: 1, limit: 3 });
    const secondResult = await listItems({ type: 'top', page: 1, limit: 3 });

    expect(firstResult).toEqual(expect.objectContaining({ ok: true }));
    expect(secondResult).toEqual(expect.objectContaining({ ok: true }));
    // Item order should remain stable between cold and warm cache paths.
    expect(secondResult.data.items.map((item) => item.id)).toEqual([21, 22, 23]);

    const requestedPaths = getRequestedPaths(fetchFn);
    // Feed IDs are intentionally fetched on each list invocation.
    expect(requestedPaths.filter((path) => path === topStoriesPath)).toHaveLength(2);
    // Item endpoints should only be fetched once because second run is cache-backed.
    expect(requestedPaths.filter((path) => path === itemPath(21))).toHaveLength(1);
    expect(requestedPaths.filter((path) => path === itemPath(22))).toHaveLength(1);
    expect(requestedPaths.filter((path) => path === itemPath(23))).toHaveLength(1);
  });

  it('get-item integration resolves a sorted comment tree and attaches parent ids', async () => {
    // Story fixture includes top-level kids so get-item traverses both depth and sibling ordering.
    const fetchFn = createQueuedFetchMock({
      [itemPath(500)]: [
        {
          id: 500,
          type: 'story',
          time: 8_000,
          title: 'Story 500',
          by: 'root-user',
          kids: [501, 502],
        },
      ],
      [itemPath(501)]: [createCommentItem({ id: 501, time: 100, kids: [503] })],
      [itemPath(502)]: [createCommentItem({ id: 502, time: 300, kids: [] })],
      [itemPath(503)]: [createCommentItem({ id: 503, time: 250, kids: [] })],
    });

    // Compose the get-item use case with the real adapter to verify TA-3 wiring through network boundary mocks.
    const api = createHnApiAdapter({ fetchFn });
    const getItem = createGetItemUseCase({ api });

    const result = await getItem(500);

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.item.id).toBe(500);

    // Top-level comments should be sorted newest-first by timestamp.
    expect(result.data.comments.map((node) => node.item.id)).toEqual([502, 501]);
    // Parent propagation is required for audit checks and nested attribution.
    expect(result.data.comments[0].item.parent).toBe(500);
    expect(result.data.comments[1].item.parent).toBe(500);

    // Nested comments should resolve under the correct parent comment node.
    expect(result.data.comments[1].comments.map((node) => node.item.id)).toEqual([503]);
    expect(result.data.comments[1].comments[0].item.parent).toBe(501);

    const requestedPaths = getRequestedPaths(fetchFn);
    // Exactly four item endpoints are required: root story, two top-level comments, one nested comment.
    expect(
      requestedPaths.filter((path) => path.startsWith(`${HN_API_PATH_PREFIX}/item/`)),
    ).toHaveLength(4);
  });

  it('computes live-update diffs correctly across two poll cycles', async () => {
    // Two queued update payloads allow deterministic first-poll and diff-poll assertions.
    const fetchFn = createQueuedFetchMock({
      [updatesPath]: [{ items: [301, 302, 303] }, { items: [302, 303, 304, 305] }],
    });

    const api = createHnApiAdapter({ fetchFn });
    // Injected clock keeps throttle behavior deterministic without real timers.
    let now = 1_000_000;

    const pollUpdates = createPollUpdatesUseCase({
      api,
      nowMs: () => now,
      minIntervalMs: 5_000,
    });

    // First poll seeds history, second poll computes newIds against previous IDs.
    const firstResult = await pollUpdates.poll();
    now += 5_000;
    const secondResult = await pollUpdates.poll();

    expect(firstResult).toEqual(expect.objectContaining({ ok: true }));
    expect(firstResult.data.isFirstPoll).toBe(true);
    expect(firstResult.data.newIds).toEqual([301, 302, 303]);
    expect(firstResult.data.currentIds).toEqual([301, 302, 303]);
    expect(firstResult.data.previousIds).toEqual([]);

    expect(secondResult).toEqual(expect.objectContaining({ ok: true }));
    expect(secondResult.data.isFirstPoll).toBe(false);
    expect(secondResult.data.newIds).toEqual([304, 305]);
    expect(secondResult.data.currentIds).toEqual([302, 303, 304, 305]);
    expect(secondResult.data.previousIds).toEqual([301, 302, 303]);

    const requestedPaths = getRequestedPaths(fetchFn);
    // Updates endpoint should be called once per successful poll.
    expect(requestedPaths.filter((path) => path === updatesPath)).toHaveLength(2);
  });
});
