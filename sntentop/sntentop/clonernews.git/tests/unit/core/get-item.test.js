// These tests pin the TA-3 item-tree contract so Track A can refactor recursion without breaking behavior.
// Public coverage focuses on comment-tree assembly, poll-part resolution, and Result error propagation.
// The fixture helpers normalize multiple result shapes emitted by the use case under test.
import { afterEach, describe, expect, it, vi } from 'vitest';

// Module path is centralized so dynamic imports stay consistent across all test helpers.
const USE_CASE_MODULE_PATH = '../../../src/core/use-cases/get-item.js';

// Result helpers mirror the production contract to keep assertions explicit and uniform.
const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

// Tree helpers accept multiple node shapes so assertions can stay stable through internal refactors.
const asArray = (value) => (Array.isArray(value) ? value : []);

// Structured cloning keeps fixture reuse safe by preventing accidental in-test mutation.
const cloneValue = (value) => {
  if (typeof structuredClone === 'function') {
    return structuredClone(value);
  }

  return JSON.parse(JSON.stringify(value));
};

// Node payload access stays strict so tests fail if the TreeNode contract shape regresses.
const getNodePayload = (node) => {
  if (!node || typeof node !== 'object' || !('item' in node)) {
    return null;
  }

  return node.item;
};

// Children access stays strict so nested-tree contract drift is caught by this suite.
const getNodeChildren = (node) => {
  if (!node || typeof node !== 'object' || !Array.isArray(node.comments)) {
    return [];
  }

  return node.comments;
};

// ID extraction centralizes node-shape handling so tree traversals stay concise.
const getNodeId = (node) => {
  const payload = getNodePayload(node);
  return payload?.id;
};

// Depth-first lookup mirrors how callers consume nested comments by ID.
const findNodeById = (nodes, targetId) => {
  for (const node of nodes) {
    if (getNodeId(node) === targetId) {
      return node;
    }

    const nestedMatch = findNodeById(getNodeChildren(node), targetId);

    if (nestedMatch) {
      return nestedMatch;
    }
  }

  return null;
};

// Flattening IDs provides a compact way to assert depth-limit behavior in the final tree.
const collectNodeIds = (nodes) => {
  const collected = [];

  for (const node of nodes) {
    const nodeId = getNodeId(node);

    if (nodeId !== undefined) {
      collected.push(nodeId);
    }

    const childIds = collectNodeIds(getNodeChildren(node));

    for (const childId of childIds) {
      collected.push(childId);
    }
  }

  return collected;
};

// Poll-option extraction is strict so parts must be attached in the canonical output location.
const extractAttachedPollPartIds = (resultData) => {
  const parts = resultData?.item?.parts;

  if (!Array.isArray(parts) || parts.length === 0) {
    return [];
  }

  return parts.map((part) => part.id);
};

// Microtask flushing lets concurrency tests observe batching boundaries between Promise.all chunks.
const flushMicrotasks = async (turns = 4) => {
  for (let index = 0; index < turns; index += 1) {
    await Promise.resolve();
  }
};

// Deterministic API fixtures keep test behavior stable while still exercising real use-case flows.
const createDeterministicApi = (itemsById) => ({
  getItem: vi.fn(async (id) => {
    if (!itemsById.has(id)) {
      return err(`missing-item-${id}`);
    }

    return ok(cloneValue(itemsById.get(id)));
  }),
});

// Dynamic import keeps this suite resilient to export-name changes while still enforcing a factory contract.
const loadGetItemFactory = async () => {
  let moduleNamespace;

  try {
    moduleNamespace = await import(USE_CASE_MODULE_PATH);
  } catch (loadError) {
    throw new Error(
      `Could not load ${USE_CASE_MODULE_PATH}. Implement TA-3 use case before running these tests.`,
      { cause: loadError },
    );
  }

  const factory =
    moduleNamespace.createGetItemUseCase ??
    moduleNamespace.createGetItem ??
    moduleNamespace.default;

  if (typeof factory !== 'function') {
    throw new Error(
      'Expected get-item module to export a factory function as createGetItemUseCase, createGetItem, or default.',
    );
  }

  return factory;
};

// Factory normalization keeps tests compatible with both function-style and object-style use-case exports.
const createGetItemInvoker = async (api, options = {}) => {
  const factory = await loadGetItemFactory();
  const useCase = factory({ api, apiAdapter: api, ...options });

  if (typeof useCase === 'function') {
    return useCase;
  }

  if (useCase && typeof useCase.execute === 'function') {
    return useCase.execute.bind(useCase);
  }

  throw new Error('Expected get-item factory to return a function or an object with execute().');
};

// Scenario runner centralizes factory wiring so individual tests stay focused on behavior assertions.
const runGetItem = async ({ api, id, options }) => {
  const getItem = await createGetItemInvoker(api, options);
  return getItem(id);
};

afterEach(() => {
  // Restoring mocks prevents call-history leakage between deep recursive test scenarios.
  vi.restoreAllMocks();
});

describe('get-item use case', () => {
  // Input validation must fail fast so adapter calls never run for impossible IDs.
  it('returns a failed Result for invalid item ids without calling the adapter', async () => {
    const api = {
      getItem: vi.fn(async () => ok(null)),
    };

    const result = await runGetItem({ api, id: 0 });

    expect(result).toEqual({
      ok: false,
      error: 'Invalid item ID. Item ID must be a positive integer.',
    });
    expect(api.getItem).not.toHaveBeenCalled();
  });

  // Null root payloads from the adapter must map to a stable not-found Result contract.
  it('returns a failed Result when the root item is not found', async () => {
    const api = {
      getItem: vi.fn(async () => ok(null)),
    };

    const result = await runGetItem({ api, id: 4040 });

    expect(result).toEqual({
      ok: false,
      error: 'Item 4040 was not found.',
    });
    expect(api.getItem).toHaveBeenCalledTimes(1);
    expect(api.getItem).toHaveBeenCalledWith(4040);
  });

  // Nested-fetch failures must propagate so comment-tree corruption is never silently ignored.
  it('propagates nested comment fetch failures as a failed Result', async () => {
    const storyWithMissingNestedComment = new Map([
      [
        70,
        {
          id: 70,
          type: 'story',
          time: 7_000,
          title: 'Nested failure root',
          kids: [701],
        },
      ],
      [
        701,
        {
          id: 701,
          type: 'comment',
          time: 701,
          text: 'Parent comment',
          kids: [702],
        },
      ],
      // 702 intentionally missing to force a nested fetch failure.
    ]);

    const api = createDeterministicApi(storyWithMissingNestedComment);
    const result = await runGetItem({ api, id: 70 });

    expect(result).toEqual(expect.objectContaining({ ok: false }));
    expect(String(result.error)).toMatch(/missing-item-702|Failed to fetch comment 702/i);
  });

  // Poll-option failures should surface as Result errors because poll rendering depends on parts integrity.
  it('returns a failed Result when poll option fetching fails', async () => {
    const api = {
      getItem: vi.fn(async (id) => {
        if (id === 80) {
          return ok({
            id: 80,
            type: 'poll',
            time: 8_000,
            title: 'Poll root',
            parts: [801],
            kids: [],
          });
        }

        if (id === 801) {
          return err('poll-option-fetch-failed');
        }

        return err(`unexpected-id-${id}`);
      }),
    };

    const result = await runGetItem({ api, id: 80 });

    expect(result).toEqual(expect.objectContaining({ ok: false }));
    expect(String(result.error)).toMatch(/poll-option-fetch-failed/i);
  });

  // Thrown adapter errors should still resolve as Result failures instead of throwing to callers.
  it('wraps unexpected thrown adapter errors in a failed Result', async () => {
    const api = {
      getItem: vi.fn(async () => {
        throw new Error('boom');
      }),
    };

    const result = await runGetItem({ api, id: 90 });

    expect(result).toEqual(expect.objectContaining({ ok: false }));
    expect(String(result.error)).toMatch(/Failed to get item 90: boom/i);
  });

  // Primary TA-3 behavior validates sorting, nesting, and structure for story comment trees.
  it('returns story comments as a newest-to-oldest tree at each depth level', async () => {
    const storyAndComments = new Map([
      [
        10,
        {
          id: 10,
          type: 'story',
          time: 1_000,
          title: 'Story root',
          kids: [101, 102, 103],
        },
      ],
      [
        101,
        {
          id: 101,
          type: 'comment',
          time: 100,
          text: 'oldest top-level',
          kids: [1011, 1012],
        },
      ],
      [
        102,
        {
          id: 102,
          type: 'comment',
          time: 300,
          text: 'newest top-level',
          kids: [1021, 1022],
        },
      ],
      [
        103,
        {
          id: 103,
          type: 'comment',
          time: 200,
          text: 'middle top-level',
          kids: [],
        },
      ],
      [
        1011,
        {
          id: 1011,
          type: 'comment',
          time: 10,
          text: 'oldest nested under 101',
          kids: [],
        },
      ],
      [
        1012,
        {
          id: 1012,
          type: 'comment',
          time: 40,
          text: 'newest nested under 101',
          kids: [],
        },
      ],
      [
        1021,
        {
          id: 1021,
          type: 'comment',
          time: 50,
          text: 'older nested under 102',
          kids: [],
        },
      ],
      [
        1022,
        {
          id: 1022,
          type: 'comment',
          time: 70,
          text: 'newer nested under 102',
          kids: [],
        },
      ],
    ]);

    const api = createDeterministicApi(storyAndComments);
    const result = await runGetItem({ api, id: 10 });

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          item: expect.objectContaining({ id: 10, type: 'story' }),
          comments: expect.any(Array),
        }),
      }),
    );

    const levelOne = asArray(result.data.comments);
    // Enforce canonical TreeNode shape at each level so contract regressions fail loudly.
    expect(levelOne.every((node) => node?.item && Array.isArray(node.comments))).toBe(true);
    expect(levelOne.map((node) => getNodeId(node))).toEqual([102, 103, 101]);

    const under102 = getNodeChildren(findNodeById(levelOne, 102));
    const under101 = getNodeChildren(findNodeById(levelOne, 101));

    expect(under102.map((node) => getNodeId(node))).toEqual([1022, 1021]);
    expect(under101.map((node) => getNodeId(node))).toEqual([1012, 1011]);
  });

  // Equal-time ordering must be deterministic so repeated renders do not shuffle comment positions.
  it('orders equal-time comments deterministically by descending id', async () => {
    const storyWithEqualTimes = new Map([
      [
        11,
        {
          id: 11,
          type: 'story',
          time: 1_100,
          title: 'Equal time story',
          kids: [111, 112, 113],
        },
      ],
      [111, { id: 111, type: 'comment', time: 90, kids: [] }],
      [112, { id: 112, type: 'comment', time: 90, kids: [] }],
      [113, { id: 113, type: 'comment', time: 90, kids: [] }],
    ]);

    const api = createDeterministicApi(storyWithEqualTimes);
    const result = await runGetItem({ api, id: 11 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.comments.map((node) => node.item.id)).toEqual([113, 112, 111]);
  });

  // Job items should return an empty comment tree because jobs do not have threaded discussions.
  it('returns a job post with no comments', async () => {
    const jobOnlyItems = new Map([
      [
        20,
        {
          id: 20,
          type: 'job',
          time: 2_000,
          title: 'Remote JS role',
          text: 'Now hiring',
        },
      ],
    ]);

    const api = createDeterministicApi(jobOnlyItems);
    const result = await runGetItem({ api, id: 20 });

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          item: expect.objectContaining({ id: 20, type: 'job' }),
          comments: [],
        }),
      }),
    );
    expect(api.getItem).toHaveBeenCalledTimes(1);
    expect(api.getItem).toHaveBeenCalledWith(20);
  });

  // Poll roots must include resolved option items to support downstream poll rendering logic.
  it('returns a poll with fetched poll options attached', async () => {
    const pollAndParts = new Map([
      [
        30,
        {
          id: 30,
          type: 'poll',
          time: 3_000,
          title: 'Favorite runtime?',
          parts: [301, 302],
          kids: [],
        },
      ],
      [
        301,
        {
          id: 301,
          type: 'pollopt',
          time: 3_100,
          text: 'Node.js',
          score: 10,
        },
      ],
      [
        302,
        {
          id: 302,
          type: 'pollopt',
          time: 3_200,
          text: 'Deno',
          score: 8,
        },
      ],
    ]);

    const api = createDeterministicApi(pollAndParts);
    const result = await runGetItem({ api, id: 30 });

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          item: expect.objectContaining({ id: 30, type: 'poll' }),
          comments: expect.any(Array),
        }),
      }),
    );

    expect(api.getItem).toHaveBeenCalledWith(301);
    expect(api.getItem).toHaveBeenCalledWith(302);

    const attachedPartIds = extractAttachedPollPartIds(result.data);
    expect(attachedPartIds).toEqual(expect.arrayContaining([301, 302]));
    expect(attachedPartIds).toHaveLength(2);
  });

  // Concurrency must remain capped globally to satisfy API discipline and avoid burst overload.
  it('never exceeds 6 concurrent in-flight comment fetches', async () => {
    const rootId = 60;
    const topLevelCommentIds = Array.from({ length: 6 }, (_, index) => 601 + index);
    const childCommentIdsByParent = new Map(
      topLevelCommentIds.map((parentId, parentIndex) => [
        parentId,
        Array.from({ length: 6 }, (_, childIndex) => 701 + parentIndex * 6 + childIndex),
      ]),
    );
    const pendingResolves = [];
    let inFlight = 0;
    let maxInFlight = 0;

    const api = {
      getItem: vi.fn((id) => {
        if (id === rootId) {
          return Promise.resolve(
            ok({
              id: rootId,
              type: 'story',
              time: 6_000,
              title: 'Concurrency root',
              kids: topLevelCommentIds,
            }),
          );
        }

        const nestedKids = childCommentIdsByParent.get(id) ?? [];

        return new Promise((resolve) => {
          inFlight += 1;
          maxInFlight = Math.max(maxInFlight, inFlight);

          pendingResolves.push(() => {
            inFlight -= 1;
            resolve(
              ok({
                id,
                type: 'comment',
                time: id,
                kids: nestedKids,
              }),
            );
          });
        });
      }),
    };

    const getItem = await createGetItemInvoker(api);
    const resultPromise = getItem(rootId);

    await flushMicrotasks();
    expect(maxInFlight).toBeLessThanOrEqual(6);

    // Releasing pending requests in waves simulates deep-tree expansion under load.
    let releasedCount = 0;
    let releasedCursor = 0;
    const expectedResolvedRequests = 42;

    for (
      let iteration = 0;
      iteration < 200 && releasedCount < expectedResolvedRequests;
      iteration += 1
    ) {
      await flushMicrotasks();
      for (const release of pendingResolves.slice(releasedCursor)) {
        release();
        releasedCount += 1;
      }
      releasedCursor = pendingResolves.length;
    }

    const result = await resultPromise;

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.comments).toHaveLength(6);
    expect(releasedCount).toBe(expectedResolvedRequests);
    expect(maxInFlight).toBe(6);
  });

  // Default depth cap protects UI responsiveness on extremely deep comment threads.
  it('enforces a strict comment recursion max depth of 5', async () => {
    const deepChainItems = new Map([
      [
        40,
        {
          id: 40,
          type: 'story',
          time: 4_000,
          title: 'Deep thread root',
          kids: [401],
        },
      ],
      [401, { id: 401, type: 'comment', time: 410, kids: [402] }],
      [402, { id: 402, type: 'comment', time: 420, kids: [403] }],
      [403, { id: 403, type: 'comment', time: 430, kids: [404] }],
      [404, { id: 404, type: 'comment', time: 440, kids: [405] }],
      [405, { id: 405, type: 'comment', time: 450, kids: [406] }],
      [406, { id: 406, type: 'comment', time: 460, kids: [] }],
    ]);

    const api = createDeterministicApi(deepChainItems);
    const result = await runGetItem({ api, id: 40 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));

    const calledIds = api.getItem.mock.calls.map(([id]) => id);
    expect(calledIds).not.toContain(406);

    const treeIds = collectNodeIds(asArray(result.data.comments));
    expect(treeIds).toContain(405);
    expect(treeIds).not.toContain(406);
  });

  // Lower maxDepth settings must truncate traversal earlier for caller-controlled performance tradeoffs.
  it('respects a custom lower maxDepth configuration', async () => {
    const shallowDepthItems = new Map([
      [
        41,
        {
          id: 41,
          type: 'story',
          time: 4_100,
          title: 'Shallow depth root',
          kids: [411],
        },
      ],
      [411, { id: 411, type: 'comment', time: 411, kids: [412] }],
      [412, { id: 412, type: 'comment', time: 412, kids: [413] }],
      [413, { id: 413, type: 'comment', time: 413, kids: [] }],
    ]);

    const api = createDeterministicApi(shallowDepthItems);
    const result = await runGetItem({ api, id: 41, options: { maxDepth: 2 } });

    expect(result).toEqual(expect.objectContaining({ ok: true }));

    const calledIds = api.getItem.mock.calls.map(([id]) => id);
    expect(calledIds).not.toContain(413);

    const treeIds = collectNodeIds(asArray(result.data.comments));
    expect(treeIds).toContain(412);
    expect(treeIds).not.toContain(413);
  });

  // Invalid maxDepth configuration should be rejected during use-case construction.
  it('returns a failed Result when maxDepth configuration is invalid', async () => {
    const api = {
      getItem: vi.fn(async () => err('should-not-be-called')),
    };

    const result = await runGetItem({ api, id: 42, options: { maxDepth: 0 } });

    expect(result).toEqual({
      ok: false,
      error: 'Invalid max depth. maxDepth must be a positive integer.',
    });
    expect(api.getItem).not.toHaveBeenCalled();
  });

  // Oversized maxDepth input must still clamp to the global depth ceiling of five.
  it('clamps maxDepth above 5 to the strict global depth cap', async () => {
    const deepChainItems = new Map([
      [
        43,
        {
          id: 43,
          type: 'story',
          time: 4_300,
          title: 'Clamp depth root',
          kids: [431],
        },
      ],
      [431, { id: 431, type: 'comment', time: 431, kids: [432] }],
      [432, { id: 432, type: 'comment', time: 432, kids: [433] }],
      [433, { id: 433, type: 'comment', time: 433, kids: [434] }],
      [434, { id: 434, type: 'comment', time: 434, kids: [435] }],
      [435, { id: 435, type: 'comment', time: 435, kids: [436] }],
      [436, { id: 436, type: 'comment', time: 436, kids: [] }],
    ]);

    const api = createDeterministicApi(deepChainItems);
    const result = await runGetItem({ api, id: 43, options: { maxDepth: 99 } });

    expect(result).toEqual(expect.objectContaining({ ok: true }));

    const calledIds = api.getItem.mock.calls.map(([id]) => id);
    expect(calledIds).not.toContain(436);
  });

  // Dependency wiring guard ensures the use case fails predictably when getItem is missing.
  it('returns a failed Result when API wiring is missing getItem', async () => {
    const result = await runGetItem({ api: {}, id: 44 });

    expect(result).toEqual({
      ok: false,
      error: 'Get item use case requires an API adapter with a getItem method.',
    });
  });

  // Parent references are required for audit checks and nested comment attribution in UI.
  it('attaches the correct parent id to each nested comment level', async () => {
    const nestedComments = new Map([
      [
        50,
        {
          id: 50,
          type: 'story',
          time: 5_000,
          title: 'Parent mapping story',
          kids: [501],
        },
      ],
      [
        501,
        {
          id: 501,
          type: 'comment',
          time: 510,
          text: 'first level',
          kids: [502],
        },
      ],
      [
        502,
        {
          id: 502,
          type: 'comment',
          time: 520,
          text: 'second level',
          kids: [503],
        },
      ],
      [
        503,
        {
          id: 503,
          type: 'comment',
          time: 530,
          text: 'third level',
          kids: [],
        },
      ],
    ]);

    const api = createDeterministicApi(nestedComments);
    const result = await runGetItem({ api, id: 50 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));

    const comments = asArray(result.data.comments);
    const firstLevel = getNodePayload(findNodeById(comments, 501));
    const secondLevel = getNodePayload(findNodeById(comments, 502));
    const thirdLevel = getNodePayload(findNodeById(comments, 503));

    expect(firstLevel).toEqual(expect.objectContaining({ id: 501, parent: 50 }));
    expect(secondLevel).toEqual(expect.objectContaining({ id: 502, parent: 501 }));
    expect(thirdLevel).toEqual(expect.objectContaining({ id: 503, parent: 502 }));
  });

  // Cycle handling must short-circuit repeated ancestors so recursive graphs cannot loop forever.
  it('short-circuits cyclic kids graphs without duplicate traversal', async () => {
    const cyclicComments = new Map([
      [
        95,
        {
          id: 95,
          type: 'story',
          time: 9_500,
          title: 'Cycle root',
          kids: [951],
        },
      ],
      [
        951,
        {
          id: 951,
          type: 'comment',
          time: 951,
          kids: [952],
        },
      ],
      [
        952,
        {
          id: 952,
          type: 'comment',
          time: 952,
          kids: [951],
        },
      ],
    ]);

    const api = createDeterministicApi(cyclicComments);
    const result = await runGetItem({ api, id: 95 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));

    const calledIds = api.getItem.mock.calls.map(([id]) => id);
    expect(calledIds).toEqual([95, 951, 952]);

    const treeIds = collectNodeIds(asArray(result.data.comments));
    expect(treeIds).toEqual([951, 952]);
  });
});
