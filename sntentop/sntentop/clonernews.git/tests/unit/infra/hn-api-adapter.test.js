// Adapter tests focus on Result shapes, timeout handling, cache behavior, and transport-policy enforcement.
// The suite exercises the createHnApiAdapter public API via injected fetch mocks only.
// Constraints: tests stay network-free and validate only observable adapter contracts.
import { afterEach, describe, expect, it, vi } from 'vitest';

import { createHnApiAdapter } from '../../../src/infra/hn-api-adapter.js';

const createAbortError = () => {
  // DOMException is preferred when available because fetch aborts usually surface that exact error shape.
  if (typeof DOMException === 'function') {
    return new DOMException('The operation was aborted.', 'AbortError');
  }

  const error = new Error('The operation was aborted.');
  error.name = 'AbortError';
  return error;
};

const createJsonResponse = (payload, options = {}) => {
  // The mock response mirrors the subset of the Fetch API this adapter actually inspects.
  const status = options.status ?? 200;
  const contentType = options.contentType ?? 'application/json; charset=utf-8';

  return {
    ok: status >= 200 && status < 300,
    status,
    headers: {
      get(name) {
        return name.toLowerCase() === 'content-type' ? contentType : null;
      },
    },
    json: vi.fn(async () => payload),
  };
};

const expectOkResult = (result) => {
  // These helpers keep the assertions focused on the Result contract instead of test boilerplate.
  expect(result).toHaveProperty('ok', true);
  expect(result).toHaveProperty('data');
};

const expectErrorResult = (result, messagePattern) => {
  // Error assertions validate both the Result shape and meaningful operator-facing message content.
  expect(result).toHaveProperty('ok', false);
  expect(result).toHaveProperty('error');
  expect(String(result.error)).toMatch(messagePattern);
};

afterEach(() => {
  vi.restoreAllMocks();
  vi.useRealTimers();
});

describe('hn api adapter', () => {
  it('returns ok data for a successful item response', async () => {
    // A happy-path item fetch proves the adapter keeps the API payload intact.
    const item = {
      by: 'dhouston',
      id: 8863,
      score: 111,
      time: 1175714200,
      title: 'My YC app: Dropbox - Throw away your USB drive',
      type: 'story',
      url: 'http://www.getdropbox.com/u/2/screencast.html',
    };

    const fetchFn = vi.fn(async () => createJsonResponse(item));
    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getItem(8863);

    expectOkResult(result);
    expect(result.data).toEqual(item);
    expect(fetchFn).toHaveBeenCalledTimes(1);
    expect(String(fetchFn.mock.calls[0][0])).toMatch(/\/item\/8863\.json/);
    expect(fetchFn.mock.calls[0][1]).toEqual(
      expect.objectContaining({ signal: expect.any(Object) }),
    );
  });

  it('returns an error result when request times out', async () => {
    // Timeout handling matters because every request must be bounded by AbortController.
    vi.useFakeTimers();

    const fetchFn = vi.fn((_, init = {}) => {
      return new Promise((_, reject) => {
        const signal = init.signal;

        if (signal?.aborted) {
          reject(createAbortError());
          return;
        }

        signal?.addEventListener(
          'abort',
          () => {
            reject(createAbortError());
          },
          { once: true },
        );
      });
    });

    const adapter = createHnApiAdapter({ fetchFn });
    const resultPromise = adapter.getItem(12345);

    await vi.advanceTimersByTimeAsync(8100);
    const result = await resultPromise;

    expectErrorResult(result, /timed out|timeout|abort/i);
    expect(fetchFn).toHaveBeenCalledTimes(1);
  });

  it('returns an error result for malformed feed response payloads', async () => {
    // Feed validation keeps invalid payloads from leaking into pagination logic.
    const fetchFn = vi.fn(async () => createJsonResponse({ not: 'an-array' }));
    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getFeedIds('top');

    expectErrorResult(result, /array|payload/i);
    expect(fetchFn).toHaveBeenCalledTimes(1);
    expect(String(fetchFn.mock.calls[0][0])).toMatch(/\/topstories\.json/);
  });

  it('returns an error result for non-ok status responses like 404', async () => {
    // A non-200 status should still return a structured Result instead of throwing.
    const fetchFn = vi.fn(async () => createJsonResponse({ error: 'Not Found' }, { status: 404 }));
    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getFeedIds('job');

    expectErrorResult(result, /404|status/i);
    expect(fetchFn).toHaveBeenCalledTimes(1);
  });

  it('returns an error result when fetch rejects with a network failure', async () => {
    // Network failures are the common transport failure mode, so they need explicit coverage.
    const fetchFn = vi.fn(async () => {
      throw new TypeError('Failed to fetch');
    });

    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getFeedIds('new');

    expectErrorResult(result, /failed to fetch|network/i);
    expect(fetchFn).toHaveBeenCalledTimes(1);
  });

  it('requests /updates.json and normalizes the items wrapper payload', async () => {
    // The updates endpoint can return an object wrapper, so the adapter must normalize it to an id array.
    const fetchFn = vi.fn(async () => createJsonResponse({ items: [301, 302, 303] }));
    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getUpdates();

    expectOkResult(result);
    expect(result.data).toEqual([301, 302, 303]);
    expect(fetchFn).toHaveBeenCalledTimes(1);
    expect(String(fetchFn.mock.calls[0][0])).toMatch(/\/updates\.json$/);
  });

  it('returns an error result when the request is explicitly aborted', async () => {
    // Abort errors should be normalized to the same safe Result flow as timeouts.
    const fetchFn = vi.fn(async () => {
      throw createAbortError();
    });

    const adapter = createHnApiAdapter({ fetchFn });
    const result = await adapter.getUpdates();

    expectErrorResult(result, /timed out|abort/i);
    expect(fetchFn).toHaveBeenCalledTimes(1);
  });

  it('uses cached data for repeated item lookups and skips the second network call', async () => {
    // Cache hits must short-circuit fetches so repeated lookups remain cheap.
    const item = {
      by: 'cache-user',
      id: 1001,
      score: 10,
      time: 1700000000,
      title: 'Cached item',
      type: 'story',
      url: 'https://example.com/cached',
    };

    const fetchFn = vi.fn(async () => createJsonResponse(item));
    const adapter = createHnApiAdapter({ fetchFn });
    const firstResult = await adapter.getItem(1001);
    const secondResult = await adapter.getItem(1001);

    expectOkResult(firstResult);
    expectOkResult(secondResult);
    expect(firstResult.data).toEqual(item);
    expect(secondResult.data).toEqual(item);
    expect(fetchFn).toHaveBeenCalledTimes(1);
  });

  it('rejects insecure non-localhost HTTP base URLs', async () => {
    // Enforcing HTTPS by default prevents accidental plaintext production transport.
    const fetchFn = vi.fn(async () => createJsonResponse([1, 2, 3]));
    const adapter = createHnApiAdapter({ fetchFn, baseUrl: 'http://example.com/v0' });
    const result = await adapter.getFeedIds('top');

    expectErrorResult(result, /https|localhost/i);
    expect(fetchFn).not.toHaveBeenCalled();
  });

  it('allows localhost HTTP base URLs for local testing', async () => {
    // Localhost is the only allowed HTTP exception so local mocks remain easy to run.
    const fetchFn = vi.fn(async () => createJsonResponse([10, 20, 30]));
    const adapter = createHnApiAdapter({ fetchFn, baseUrl: 'http://localhost:4010/v0' });
    const result = await adapter.getFeedIds('top');

    expectOkResult(result);
    expect(result.data).toEqual([10, 20, 30]);
    expect(String(fetchFn.mock.calls[0][0])).toContain('http://localhost:4010/v0/topstories.json');
  });
});
