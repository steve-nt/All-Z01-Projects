/*
 * Feed controller unit tests.
 * Validates tab switching, feed fetching, and error handling.
 */

import { beforeEach, describe, expect, it, vi } from 'vitest';
import { createFeedController } from '../../../../src/features/feed/feed.controller.js';

describe('Feed Controller', () => {
  let mockListItems;
  let mockContainer;
  let controller;

  beforeEach(() => {
    // Mock the listItems use-case.
    mockListItems = vi.fn();

    // Mock container element.
    mockContainer = {
      querySelector: vi.fn(),
    };

    controller = createFeedController({
      listItems: mockListItems,
      container: mockContainer,
    });
  });

  describe('initialization', () => {
    it('should throw if listItems is not a function', () => {
      expect(() => {
        createFeedController({
          listItems: null,
          container: mockContainer,
        });
      }).toThrow('Feed controller requires a listItems function.');
    });
  });

  describe('switchTab', () => {
    it('should call listItems with correct feed type', async () => {
      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: [], hasMore: false },
      });

      const onStateChange = vi.fn();
      await controller.switchTab('new', onStateChange);

      expect(mockListItems).toHaveBeenCalledWith({
        type: 'new',
        page: 1,
        limit: 20,
      });
    });

    it('should call onStateChange with loading state before fetch', async () => {
      mockListItems.mockImplementation(
        () =>
          new Promise((resolve) => {
            setTimeout(() => resolve({ ok: true, data: { items: [], hasMore: false } }), 10);
          }),
      );

      const onStateChange = vi.fn();

      // Don't await to catch the loading state call.
      controller.switchTab('top', onStateChange);

      expect(onStateChange).toHaveBeenCalledWith([], true, undefined, true, false);
    });

    it('should pass items to onStateChange on successful fetch', async () => {
      const mockItems = [
        { id: 1, title: 'Story 1', score: 100, by: 'user1', time: 1000, descendants: 5 },
      ];

      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: mockItems, hasMore: false },
      });

      const onStateChange = vi.fn();
      await controller.switchTab('top', onStateChange);

      // Should have: loading call, then success call (items, isLoading, error, hasMore)
      expect(onStateChange).toHaveBeenNthCalledWith(2, mockItems, false, undefined, false, false);
    });

    it('should pass error to onStateChange on API failure', async () => {
      mockListItems.mockResolvedValue({
        ok: false,
        error: 'Failed to fetch feed',
      });

      const onStateChange = vi.fn();
      await controller.switchTab('top', onStateChange);

      expect(onStateChange).toHaveBeenNthCalledWith(
        2,
        undefined,
        false,
        'Failed to fetch feed',
        true,
        false,
      );
    });

    it('should reject invalid feed types', async () => {
      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: [], hasMore: false },
      });

      const onStateChange = vi.fn();
      await controller.switchTab('invalid', onStateChange);

      expect(mockListItems).not.toHaveBeenCalled();
    });
  });

  describe('mount', () => {
    it('should fetch default feed on mount', async () => {
      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: [], hasMore: false },
      });

      const onStateChange = vi.fn();
      controller.mount(onStateChange);

      // Allow async to settle.
      await new Promise((resolve) => setTimeout(resolve, 10));

      expect(mockListItems).toHaveBeenCalledWith({
        type: 'top',
        page: 1,
        limit: 20,
      });
    });
  });

  describe('unmount', () => {
    it('should not throw when unmounting', () => {
      expect(() => {
        controller.unmount();
      }).not.toThrow();
    });
  });

  describe('tab switching scenarios', () => {
    it('should clear loading state on success', async () => {
      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: [{ id: 1 }], hasMore: false },
      });

      const onStateChange = vi.fn();
      await controller.switchTab('top', onStateChange);

      // Final call should have isLoading = false.
      const lastCall = onStateChange.mock.calls[onStateChange.mock.calls.length - 1];
      expect(lastCall[1]).toBe(false);
    });

    it('should handle rapid tab switches', async () => {
      mockListItems.mockResolvedValue({
        ok: true,
        data: { items: [], hasMore: false },
      });

      const onStateChange = vi.fn();

      // Switch tabs rapidly.
      await controller.switchTab('top', onStateChange);
      await controller.switchTab('new', onStateChange);
      await controller.switchTab('job', onStateChange);

      // Should have called listItems 3 times (one per switch).
      expect(mockListItems).toHaveBeenCalledTimes(3);
    });

    it('keeps poll feed results even when the first page is empty', async () => {
      mockListItems.mockResolvedValueOnce({
        ok: true,
        data: { items: [], hasMore: false },
      });

      const onStateChange = vi.fn();
      await controller.switchTab('poll', onStateChange);

      expect(mockListItems).toHaveBeenNthCalledWith(1, {
        type: 'poll',
        page: 1,
        limit: 20,
      });
      expect(onStateChange).toHaveBeenNthCalledWith(
        2,
        [],
        false,
        undefined,
        false,
        false,
      );
    });

    it('does not switch to top feed when poll requests resolve slowly', async () => {
      vi.useFakeTimers();

      try {
        mockListItems.mockImplementationOnce(
          () =>
            new Promise((resolve) => {
              setTimeout(() => {
                resolve({
                  ok: true,
                  data: { items: [{ id: 88, title: 'Slow poll result' }], hasMore: true },
                });
              }, 6_000);
            }),
        );

        const onStateChange = vi.fn();
        const switchPromise = controller.switchTab('poll', onStateChange);

        await vi.advanceTimersByTimeAsync(6_000);
        await switchPromise;

        expect(mockListItems).toHaveBeenCalledTimes(1);
        expect(mockListItems).toHaveBeenNthCalledWith(1, {
          type: 'poll',
          page: 1,
          limit: 20,
        });
        expect(mockListItems).not.toHaveBeenCalledWith({
          type: 'top',
          page: 1,
          limit: 20,
        });
        expect(onStateChange).toHaveBeenNthCalledWith(
          2,
          [{ id: 88, title: 'Slow poll result' }],
          false,
          undefined,
          true,
          false,
        );
      } finally {
        vi.useRealTimers();
      }
    });
  });
});
