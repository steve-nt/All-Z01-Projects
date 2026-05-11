// These tests lock down the TA-2 feed pagination contract so Track A can evolve without breaking core feed behavior.
// Public coverage focuses on paging, ordering, concurrency, and input validation at the use-case boundary.
// The suite uses mocked API adapters so Track A stays isolated from transport and cache concerns.
import { afterEach, describe, expect, it, vi } from 'vitest';

const USE_CASE_MODULE_PATH = '../../../src/core/use-cases/list-items.js';

// Result helpers keep the test fixtures aligned with the use case contract.
const ok = (data) => ({ ok: true, data });
// Error fixture helper keeps failed Result objects consistent across negative-path cases.
const err = (error) => ({ ok: false, error });

// Sequential IDs make pagination slices and ordering assertions easy to read.
const createSequentialIds = (count, start = 1) =>
  Array.from({ length: count }, (_, index) => start + index);

// Each fixture only needs the fields the use case actually sorts or forwards.
const createStoryItem = ({ id, time }) => ({
  id,
  type: 'story',
  time,
  title: `Story ${id}`,
});

// Poll fixtures cover the derived poll feed path that scans top-story IDs.
const createPollItem = ({ id, time }) => ({
  id,
  type: 'poll',
  time,
  title: `Poll ${id}`,
});

// Microtask flushing lets the concurrency test observe the batch boundary before releasing promises.
const flushMicrotasks = async (turns = 4) => {
  for (let index = 0; index < turns; index += 1) {
    await Promise.resolve();
  }
};

// Dynamic import keeps the test aligned with whatever export shape the module currently provides.
const loadListItemsFactory = async () => {
  let moduleNamespace;

  try {
    moduleNamespace = await import(USE_CASE_MODULE_PATH);
  } catch (loadError) {
    throw new Error(
      `Could not load ${USE_CASE_MODULE_PATH}. Implement TA-2 use case before running these tests.`,
      { cause: loadError },
    );
  }

  const factory =
    moduleNamespace.createListItemsUseCase ??
    moduleNamespace.createListItems ??
    moduleNamespace.default;

  if (typeof factory !== 'function') {
    throw new Error(
      'Expected list-items module to export a factory function as createListItemsUseCase, createListItems, or default.',
    );
  }

  return factory;
};

// The factory accepts either api or apiAdapter in a few historical shapes, so the test supports both.
const createListItemsInvoker = async (api) => {
  const factory = await loadListItemsFactory();
  const useCase = factory({ api, apiAdapter: api });

  if (typeof useCase === 'function') {
    return useCase;
  }

  if (useCase && typeof useCase.execute === 'function') {
    return useCase.execute.bind(useCase);
  }

  throw new Error('Expected list-items factory to return a function or an object with execute().');
};

// Helper to keep the scenario setup short and focused on the asserted behavior.
const runListItems = async ({ api, type = 'top', page, limit }) => {
  const listItems = await createListItemsInvoker(api);
  // Request construction stays explicit so tests can intentionally omit page/limit defaults.
  const request = { type };

  if (page !== undefined) {
    request.page = page;
  }

  if (limit !== undefined) {
    request.limit = limit;
  }

  return listItems(request);
};

afterEach(() => {
  // Each test needs a clean mock slate so concurrency assertions stay deterministic.
  vi.restoreAllMocks();
});

describe('list-items use case', () => {
  // The default page size must stay aligned with TA-2 acceptance criteria.
  it('returns the first page successfully with default limit=20', async () => {
    const feedIds = createSequentialIds(25);
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const result = await runListItems({ api, type: 'top', page: 1 });

    expect(api.getFeedIds).toHaveBeenCalledTimes(1);
    expect(api.getFeedIds).toHaveBeenCalledWith('top');
    expect(api.getItem).toHaveBeenCalledTimes(20);

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          hasMore: true,
          items: expect.any(Array),
        }),
      }),
    );

    expect(result.data.items).toHaveLength(20);
    expect(result.data.items.map((item) => item.id)).toEqual(createSequentialIds(20));
  });

  // Omitting page should fall back to page 1 so callers can request the first page with only a type.
  it('defaults page to 1 when page is omitted', async () => {
    const feedIds = createSequentialIds(30);
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const result = await runListItems({ api, type: 'top' });
    const fetchedIds = api.getItem.mock.calls.map(([id]) => id);

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(fetchedIds).toEqual(feedIds.slice(0, 20));
  });

  // Invalid top-level inputs must fail before destructuring so the executor always returns a Result.
  it.each([
    ['null input', null],
    ['primitive input', 42],
    ['array input', []],
  ])('rejects %s before destructuring', async (_, input) => {
    const api = {
      getFeedIds: vi.fn(async () => ok(createSequentialIds(10))),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const listItems = await createListItemsInvoker(api);
    const result = await listItems(input);

    expect(result).toEqual({ ok: false, error: 'List items use case expects an input object.' });
    expect(api.getFeedIds).not.toHaveBeenCalled();
    expect(api.getItem).not.toHaveBeenCalled();
  });

  // Input validation should stop malformed paging and feed values before any adapter work starts.
  it.each([
    [
      'invalid feed type',
      { type: 'archive', page: 1 },
      'Invalid feed type. Expected one of: top, new, ask, show, job, poll.',
    ],
    ['invalid page', { type: 'top', page: 0 }, 'Invalid page. Page must be a positive integer.'],
    [
      'invalid limit',
      { type: 'top', page: 1, limit: 0 },
      'Invalid limit. Limit must be a positive integer.',
    ],
  ])('rejects %s', async (_, input, expectedError) => {
    const api = {
      getFeedIds: vi.fn(async () => ok(createSequentialIds(10))),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const listItems = await createListItemsInvoker(api);
    const result = await listItems(input);

    expect(result).toEqual({ ok: false, error: expectedError });
    expect(api.getFeedIds).not.toHaveBeenCalled();
    expect(api.getItem).not.toHaveBeenCalled();
  });

  // Page math should select the correct window without overfetching outside the slice.
  it('fetches only the correct page slice by page and limit', async () => {
    const feedIds = createSequentialIds(45);
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: 1_000 + id }))),
    };

    const result = await runListItems({ api, type: 'new', page: 2, limit: 10 });
    const fetchedIds = api.getItem.mock.calls.map(([id]) => id);

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(fetchedIds).toEqual(feedIds.slice(10, 20));
    expect(api.getItem).toHaveBeenCalledTimes(10);
  });

  // Poll feed IDs are derived from top stories because the HN API has no dedicated pollstories endpoint.
  it('derives poll pages from top feed IDs and filters only poll items', async () => {
    const feedIds = [1, 2, 3, 4, 5, 6];
    const itemsById = new Map([
      [1, createStoryItem({ id: 1, time: 10 })],
      [2, createPollItem({ id: 2, time: 20 })],
      [3, createStoryItem({ id: 3, time: 30 })],
      [4, createStoryItem({ id: 4, time: 40 })],
      [5, createPollItem({ id: 5, time: 50 })],
      [6, createStoryItem({ id: 6, time: 60 })],
    ]);

    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(itemsById.get(id))),
    };

    const pageOneResult = await runListItems({ api, type: 'poll', page: 1, limit: 1 });
    const pageTwoResult = await runListItems({ api, type: 'poll', page: 2, limit: 1 });

    expect(api.getFeedIds).toHaveBeenNthCalledWith(1, 'top');
    expect(api.getFeedIds).toHaveBeenNthCalledWith(2, 'top');

    expect(pageOneResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageOneResult.data.items).toHaveLength(1);
    expect(pageOneResult.data.items[0].type).toBe('poll');
    expect(pageOneResult.data.hasMore).toBe(true);

    expect(pageTwoResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageTwoResult.data.items).toHaveLength(1);
    expect(pageTwoResult.data.items[0].type).toBe('poll');
    expect(pageTwoResult.data.hasMore).toBe(false);
  });

  // When active feeds have no polls, the deterministic fallback IDs should still populate the poll tab.
  it('uses fallback poll IDs when top/new/ask/show feeds contain no poll items', async () => {
    const feedIds = [11, 12, 13];
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => {
        if (id >= 46_000_000) {
          return ok(createPollItem({ id, time: id }));
        }

        return ok(createStoryItem({ id, time: id }));
      }),
    };

    const result = await runListItems({ api, type: 'poll', page: 1, limit: 3 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.items.length).toBeGreaterThan(0);
    expect(result.data.items.every((item) => item.type === 'poll')).toBe(true);
  });

  // Sorting happens at the end, so the original feed order should never leak into the response.
  it('always orders output newest-to-oldest by item.time', async () => {
    const feedIds = [101, 102, 103, 104];
    const itemsById = new Map([
      [101, createStoryItem({ id: 101, time: 20 })],
      [102, createStoryItem({ id: 102, time: 80 })],
      [103, createStoryItem({ id: 103, time: 40 })],
      [104, createStoryItem({ id: 104, time: 60 })],
    ]);

    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(itemsById.get(id))),
    };

    const result = await runListItems({ api, type: 'ask', page: 1, limit: 4 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.items.map((item) => item.id)).toEqual([102, 104, 103, 101]);
    expect(result.data.items.map((item) => item.time)).toEqual([80, 60, 40, 20]);
  });

  // Top feed must preserve Hacker News API ranking order instead of applying local time sorting.
  it('preserves raw feed order for top items', async () => {
    const feedIds = [900, 901, 902, 903];
    const itemsById = new Map([
      [900, createStoryItem({ id: 900, time: 50 })],
      [901, createStoryItem({ id: 901, time: 300 })],
      [902, createStoryItem({ id: 902, time: 120 })],
      [903, createStoryItem({ id: 903, time: 700 })],
    ]);

    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(itemsById.get(id))),
    };

    const result = await runListItems({ api, type: 'top', page: 1, limit: 4 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.items.map((item) => item.id)).toEqual([900, 901, 902, 903]);
  });

  // Equal timestamps still need a deterministic fallback so repeated runs stay reproducible.
  it('breaks equal timestamps by descending item id', async () => {
    const feedIds = [301, 302, 303];
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: 50 }))),
    };

    const result = await runListItems({ api, type: 'show', page: 1, limit: 3 });

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.items.map((item) => item.id)).toEqual([303, 302, 301]);
  });

  // hasMore must reflect the source ID list, not the number of items successfully fetched.
  it('sets hasMore correctly for both non-terminal and terminal pages', async () => {
    const pageOneWithMoreApi = {
      getFeedIds: vi.fn(async () => ok(createSequentialIds(21))),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };
    const pageOneTerminalApi = {
      getFeedIds: vi.fn(async () => ok(createSequentialIds(20))),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const pageOneWithMoreResult = await runListItems({
      api: pageOneWithMoreApi,
      type: 'show',
      page: 1,
    });
    const pageOneTerminalResult = await runListItems({
      api: pageOneTerminalApi,
      type: 'show',
      page: 1,
    });

    expect(pageOneWithMoreResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageOneWithMoreResult.data.hasMore).toBe(true);

    expect(pageOneTerminalResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageOneTerminalResult.data.hasMore).toBe(false);
  });

  // hasMore should remain accurate on later pages and when the request is past the end of the feed.
  it('sets hasMore correctly for later and out-of-range pages', async () => {
    const feedIds = createSequentialIds(25);
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const pageTwoResult = await runListItems({ api, type: 'top', page: 2, limit: 10 });
    const pageFourResult = await runListItems({ api, type: 'top', page: 4, limit: 10 });

    expect(pageTwoResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageTwoResult.data.hasMore).toBe(true);

    expect(pageFourResult).toEqual(expect.objectContaining({ ok: true }));
    expect(pageFourResult.data.hasMore).toBe(false);
    expect(pageFourResult.data.items).toEqual([]);
  });

  // Concurrency must stay capped at six to match the core batching requirement.
  it('never exceeds 6 concurrent in-flight item fetches', async () => {
    const feedIds = createSequentialIds(12);
    const pendingResolves = [];
    let inFlight = 0;
    let maxInFlight = 0;

    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(
        (id) =>
          new Promise((resolve) => {
            inFlight += 1;
            maxInFlight = Math.max(maxInFlight, inFlight);

            pendingResolves.push(() => {
              inFlight -= 1;
              resolve(ok(createStoryItem({ id, time: id })));
            });
          }),
      ),
    };

    const listItems = await createListItemsInvoker(api);
    const resultPromise = listItems({ type: 'top', page: 1, limit: 12 });

    await flushMicrotasks();
    expect(api.getItem).toHaveBeenCalledTimes(6);
    expect(maxInFlight).toBeLessThanOrEqual(6);

    for (const release of pendingResolves.slice(0, 6)) {
      release();
    }

    await flushMicrotasks();
    expect(api.getItem).toHaveBeenCalledTimes(12);
    expect(maxInFlight).toBeLessThanOrEqual(6);

    for (const release of pendingResolves.slice(6)) {
      release();
    }

    const result = await resultPromise;

    expect(result).toEqual(expect.objectContaining({ ok: true }));
    expect(result.data.items).toHaveLength(12);
    expect(maxInFlight).toBe(6);
  });

  // Feed ID failures should short-circuit before any item requests are attempted.
  it('propagates feed ID fetch errors and skips item fetches', async () => {
    const feedError = err('feed-id-fetch-failed');
    const api = {
      getFeedIds: vi.fn(async () => feedError),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const result = await runListItems({ api, type: 'job', page: 1 });

    expect(api.getFeedIds).toHaveBeenCalledTimes(1);
    expect(api.getItem).not.toHaveBeenCalled();
    expect(result).toEqual(feedError);
  });

  // Item fetch failures should bubble up as a failed Result instead of being swallowed.
  it('propagates errors when one or more item fetches fail', async () => {
    const feedIds = createSequentialIds(8);
    const api = {
      getFeedIds: vi.fn(async () => ok(feedIds)),
      getItem: vi.fn(async (id) => {
        if (id === 3 || id === 7) {
          return err('item-fetch-failed');
        }

        return ok(createStoryItem({ id, time: id }));
      }),
    };

    const result = await runListItems({ api, type: 'top', page: 1, limit: 8 });

    expect(api.getFeedIds).toHaveBeenCalledTimes(1);
    expect(api.getItem).toHaveBeenCalled();
    expect(result).toEqual(expect.objectContaining({ ok: false }));
    expect(String(result.error)).toMatch(/item-fetch-failed/i);
  });

  // Defensive payload validation should fail fast when feed IDs are not returned as an array.
  it('returns an error result for malformed feed-id payloads', async () => {
    const api = {
      getFeedIds: vi.fn(async () => ok({ ids: [1, 2, 3] })),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const result = await runListItems({ api, type: 'top', page: 1 });

    expect(result).toEqual({ ok: false, error: 'Feed ID result payload is invalid.' });
    expect(api.getItem).not.toHaveBeenCalled();
  });

  // Unexpected thrown exceptions should still be wrapped into a stable Result error shape.
  it('wraps unexpected thrown errors in a failed Result', async () => {
    const api = {
      getFeedIds: vi.fn(async () => {
        throw new Error('boom');
      }),
      getItem: vi.fn(async (id) => ok(createStoryItem({ id, time: id }))),
    };

    const result = await runListItems({ api, type: 'top', page: 1 });

    expect(result).toEqual(expect.objectContaining({ ok: false }));
    expect(String(result.error)).toMatch(/Failed to list items: boom/i);
  });
});
