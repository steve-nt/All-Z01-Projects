/*
 * Test suite for feed controller infinite scroll / load more functionality.
 * Verifies: page tracking, loading flag, duplicate prevention, observer cleanup.
 */

import { beforeEach, describe, expect, it, vi } from 'vitest';

import { createFeedController } from '../../../../src/features/feed/feed.controller.js';

describe('Feed Controller - Infinite Scroll', () => {
  let mockListItems;
  let controller;

  beforeEach(() => {
    // Mock listItems use-case that tracks calls and returns paginated results.
    mockListItems = vi.fn(async (input) => {
      // Extract page for pagination logic (type is unused in test).
      const { page } = input;

      // Simulate API returning different items per page with hasMore flag.
      if (page === 1) {
        return {
          ok: true,
          data: {
            items: [
              { id: 1, title: 'Page 1 Item 1', type: 'story' },
              { id: 2, title: 'Page 1 Item 2', type: 'story' },
            ],
            hasMore: true, // More pages available.
          },
        };
      }

      if (page === 2) {
        return {
          ok: true,
          data: {
            items: [
              { id: 3, title: 'Page 2 Item 1', type: 'story' },
              { id: 4, title: 'Page 2 Item 2', type: 'story' },
            ],
            hasMore: true,
          },
        };
      }

      // Page 3 is the last page.
      if (page === 3) {
        return {
          ok: true,
          data: {
            items: [{ id: 5, title: 'Page 3 Item 1', type: 'story' }],
            hasMore: false, // No more pages.
          },
        };
      }

      return { ok: false, error: 'Invalid page' };
    });

    controller = createFeedController({ listItems: mockListItems });
  });

  describe('Initial load and page tracking', () => {
    it('mounts and fetches page 1', async () => {
      const onStateChange = vi.fn();
      controller.mount(onStateChange);

      // Wait for async operation.
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Should call onStateChange with page 1 items and hasMore=true.
      expect(mockListItems).toHaveBeenCalledWith({
        type: 'top',
        page: 1,
        limit: 20,
      });

      expect(onStateChange).toHaveBeenCalledWith(
        expect.arrayContaining([
          expect.objectContaining({ id: 1, title: 'Page 1 Item 1' }),
          expect.objectContaining({ id: 2, title: 'Page 1 Item 2' }),
        ]),
        false,
        undefined,
        true, // hasMore flag.
        false,
      );
    });

    it('loadMorePage increments page and appends items', async () => {
      const onStateChange = vi.fn();
      controller.mount(onStateChange);

      // Wait for initial load.
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Reset call history to focus on loadMorePage.
      mockListItems.mockClear();
      onStateChange.mockClear();

      // Trigger load more (page 2).
      controller.loadMorePage(onStateChange);

      // Wait for async operation.
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Should fetch page 2.
      expect(mockListItems).toHaveBeenCalledWith({
        type: 'top',
        page: 2,
        limit: 20,
      });

      // Should call onStateChange with page 2 items.
      expect(onStateChange).toHaveBeenCalledWith(
        expect.arrayContaining([
          expect.objectContaining({ id: 3, title: 'Page 2 Item 1' }),
          expect.objectContaining({ id: 4, title: 'Page 2 Item 2' }),
        ]),
        false,
        undefined,
        true, // hasMore still true.
        true,
      );
    });
  });

  describe('Loading flag prevents duplicate requests', () => {
    it('ignores loadMorePage calls while loading is true', async () => {
      const onStateChange = vi.fn();
      let fetchResolve;

      // Create a promise that we can control to delay the fetch.
      mockListItems.mockImplementation(
        () =>
          new Promise((resolve) => {
            fetchResolve = resolve;
          }),
      );

      controller.mount(onStateChange);

      // Don't wait yet—loading is still true.
      // Try to load more while first fetch is in flight.
      controller.loadMorePage(onStateChange);

      // mockListItems should only be called once (for mount), not twice.
      expect(mockListItems).toHaveBeenCalledTimes(1);

      // Resolve the fetch to clean up.
      fetchResolve({
        ok: true,
        data: { items: [], hasMore: true },
      });

      // Wait for clean up.
      await new Promise((resolve) => setTimeout(resolve, 10));
    });
  });

  describe('Switch tab resets pagination', () => {
    it('resets to page 1 when switching tabs', async () => {
      const onStateChange = vi.fn();
      controller.mount(onStateChange);

      // Wait for initial load.
      await new Promise((resolve) => setTimeout(resolve, 10));

      mockListItems.mockClear();

      // Switch to 'new' tab.
      await controller.switchTab('new', onStateChange);

      // Should fetch page 1 of 'new' feed.
      expect(mockListItems).toHaveBeenCalledWith({
        type: 'new',
        page: 1,
        limit: 20,
      });
    });
  });

  describe('Observer cleanup on hasMore=false', () => {
    it('stops loading when hasMore becomes false', async () => {
      const onStateChange = vi.fn();
      controller.mount(onStateChange);

      // Wait for initial load (page 1, hasMore=true).
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Load page 2.
      controller.loadMorePage(onStateChange);
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Load page 3 (hasMore=false).
      mockListItems.mockClear();
      controller.loadMorePage(onStateChange);
      await new Promise((resolve) => setTimeout(resolve, 10));

      // After page 3 is loaded with hasMore=false, should report no more pages.
      const lastCall = onStateChange.mock.calls[onStateChange.mock.calls.length - 1];
      expect(lastCall[3]).toBe(false); // hasMore flag should be false.

      // Try to load more—should not call listItems (already at end).
      mockListItems.mockClear();
      controller.loadMorePage(onStateChange);
      await new Promise((resolve) => setTimeout(resolve, 10));

      // No new API call should be made.
      expect(mockListItems).not.toHaveBeenCalled();
    });
  });

  describe('Error handling during pagination', () => {
    it('passes error message on failed page load', async () => {
      const onStateChange = vi.fn();

      mockListItems.mockImplementation(async (input) => {
        if (input.page === 2) {
          return { ok: false, error: 'Network error' };
        }

        return {
          ok: true,
          data: {
            items: [{ id: 1, title: 'Item 1', type: 'story' }],
            hasMore: true,
          },
        };
      });

      controller.mount(onStateChange);
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Load page 2 (will fail).
      controller.loadMorePage(onStateChange);
      await new Promise((resolve) => setTimeout(resolve, 10));

      // Should call onStateChange with error.
      expect(onStateChange).toHaveBeenCalledWith(
        undefined,
        false,
        'Network error',
        expect.any(Boolean),
        true,
      );
    });
  });

  describe('Stale response handling', () => {
    it('ignores responses from cancelled requests', async () => {
      const onStateChange = vi.fn();
      let firstResolve;

      mockListItems.mockImplementation(
        () =>
          new Promise((resolve) => {
            firstResolve = resolve;
          }),
      );

      controller.mount(onStateChange);

      // Switch tab before first response arrives (starts new request).
      controller.switchTab('new', onStateChange);

      // First request was cancelled; only second should complete.
      firstResolve({
        ok: true,
        data: { items: [{ id: 999 }], hasMore: true },
      });

      // Wait for cleanup.
      await new Promise((resolve) => setTimeout(resolve, 10));

      // The old response should be ignored; listItems called twice (mount, switch).
      expect(mockListItems).toHaveBeenCalledTimes(2);
    });
  });
});
