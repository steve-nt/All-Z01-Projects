import { Temporal } from '@js-temporal/polyfill';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { createCacheAdapter } from '../../../src/infra/cache-adapter.js';

describe('cache adapter', () => {
  let currentEpochMilliseconds = 1_000;

  beforeEach(() => {
    currentEpochMilliseconds = 1_000;

    vi.spyOn(Temporal.Now, 'instant').mockImplementation(() =>
      Temporal.Instant.fromEpochMilliseconds(currentEpochMilliseconds),
    );
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('stores and returns a cached value', async () => {
    const cache = createCacheAdapter();
    const item = { id: 42, title: 'Cached story' };

    await cache.set(42, item);

    await expect(cache.get(42)).resolves.toEqual({
      ok: true,
      data: item,
    });
  });

  it('returns undefined for missing keys', async () => {
    const cache = createCacheAdapter();

    await expect(cache.get(999)).resolves.toEqual({
      ok: true,
      data: undefined,
    });
  });

  it('reports whether a fresh entry exists', async () => {
    const cache = createCacheAdapter();

    await expect(cache.has('story')).resolves.toEqual({
      ok: true,
      data: false,
    });

    await cache.set('story', { id: 1 });

    await expect(cache.has('story')).resolves.toEqual({
      ok: true,
      data: true,
    });
  });

  it('invalidates a single cached entry', async () => {
    const cache = createCacheAdapter();

    await cache.set('story', { id: 1 });
    await expect(cache.invalidate('story')).resolves.toEqual({
      ok: true,
      data: undefined,
    });
    await expect(cache.get('story')).resolves.toEqual({
      ok: true,
      data: undefined,
    });
  });

  it('clears all cached entries', async () => {
    const cache = createCacheAdapter();

    await cache.set('story', { id: 1 });
    await cache.set('job', { id: 2 });

    await expect(cache.clear()).resolves.toEqual({
      ok: true,
      data: undefined,
    });
    await expect(cache.get('story')).resolves.toEqual({
      ok: true,
      data: undefined,
    });
    await expect(cache.get('job')).resolves.toEqual({
      ok: true,
      data: undefined,
    });
  });

  it('expires entries after the configured ttl on get', async () => {
    const cache = createCacheAdapter({ ttlMs: 10 });

    await cache.set('story', { id: 1 });
    currentEpochMilliseconds += 11;

    await expect(cache.get('story')).resolves.toEqual({
      ok: true,
      data: undefined,
    });
  });

  it('treats expired entries as missing in has', async () => {
    const cache = createCacheAdapter({ ttlMs: 10 });

    await cache.set('story', { id: 1 });
    currentEpochMilliseconds += 11;

    await expect(cache.has('story')).resolves.toEqual({
      ok: true,
      data: false,
    });
  });
});
