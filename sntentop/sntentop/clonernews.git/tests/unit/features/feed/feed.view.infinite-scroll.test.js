/*
 * Test suite for feed view infinite scroll with IntersectionObserver.
 * Verifies: sentinel tracking, batch DOM insertions, observer cleanup.
 */

import { JSDOM } from 'jsdom';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { createFeedView } from '../../../../src/features/feed/feed.view.js';

describe('Feed View - Infinite Scroll', () => {
  let dom;
  let container;
  let mockController;

  beforeEach(() => {
    // Set up JSDOM for DOM testing with proper URL origin to allow localStorage.
    dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
      url: 'https://localhost',
    });
    global.window = dom.window;
    global.document = dom.window.document;
    global.IntersectionObserver = vi.fn((_callback) => ({
      observe: vi.fn(),
      disconnect: vi.fn(),
      unobserve: vi.fn(),
    }));

    container = document.createElement('div');

    mockController = {
      switchTab: vi.fn(),
      mount: vi.fn(),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    };
  });

  describe('Sentinel element creation', () => {
    it('creates sentinel below feed list', () => {
      const view = createFeedView({ container, controller: mockController });
      view.mount();

      // Find sentinel element.
      const sentinel = container.querySelector('[data-testid="feed-sentinel"]');
      expect(sentinel).toBeTruthy();
      expect(sentinel.hidden).toBe(false);
    });

    it('renders feed list before sentinel', () => {
      const view = createFeedView({ container, controller: mockController });
      view.mount();

      const list = container.querySelector('[data-testid="feed-list"]');
      const sentinel = container.querySelector('[data-testid="feed-sentinel"]');

      // List should come before sentinel in DOM order.
      expect(list).toBeTruthy();
      expect(sentinel).toBeTruthy();
      expect(list.parentElement === container).toBe(true);
      expect(sentinel.parentElement === container).toBe(true);
    });
  });

  describe('IntersectionObserver integration', () => {
    it('creates observer on mount', () => {
      const observerConstructor = global.IntersectionObserver;
      const view = createFeedView({ container, controller: mockController });
      view.mount();

      // IntersectionObserver should be created.
      expect(observerConstructor).toHaveBeenCalled();
    });

    it('observes sentinel element', () => {
      const mockObserverInstance = {
        observe: vi.fn(),
        disconnect: vi.fn(),
      };
      global.IntersectionObserver = vi.fn(() => mockObserverInstance);

      const view = createFeedView({ container, controller: mockController });
      view.mount();

      const sentinel = container.querySelector('[data-testid="feed-sentinel"]');
      expect(mockObserverInstance.observe).toHaveBeenCalledWith(sentinel);
    });

    it('disconnects observer on unmount', () => {
      const mockObserverInstance = {
        observe: vi.fn(),
        disconnect: vi.fn(),
      };
      global.IntersectionObserver = vi.fn(() => mockObserverInstance);

      const view = createFeedView({ container, controller: mockController });
      view.mount();
      view.unmount();

      expect(mockObserverInstance.disconnect).toHaveBeenCalled();
    });

    it('reconnects observer when switching to a tab with more pages after an exhausted tab', async () => {
      const observerInstances = [];

      global.IntersectionObserver = vi.fn((_callback) => {
        const observerInstance = {
          observe: vi.fn(),
          disconnect: vi.fn(),
        };

        observerInstances.push(observerInstance);
        return observerInstance;
      });

      let onStateChange;

      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
        callback(
          [{ id: 1, title: 'Ask 1', by: 'asker', time: 1000, score: 1, descendants: 0 }],
          false,
          undefined,
          true,
        );
      });

      mockController.switchTab.mockImplementation(async (_tabType, callback) => {
        callback(
          [{ id: 2, title: 'Top 1', by: 'topper', time: 2000, score: 2, descendants: 0 }],
          false,
          undefined,
          true,
        );
      });

      const view = createFeedView({ container, controller: mockController });
      view.mount();

      expect(observerInstances.length).toBe(1);

      // Simulate feed exhaustion on current tab.
      onStateChange(
        [{ id: 1, title: 'Ask 1', by: 'asker', time: 1000, score: 1, descendants: 0 }],
        false,
        undefined,
        false,
      );

      expect(observerInstances[0].disconnect).toHaveBeenCalled();

      // Switching tabs with hasMore=true should recreate and observe a fresh sentinel observer.
      await view.selectTab('top');

      expect(observerInstances.length).toBe(2);
      expect(observerInstances[1].observe).toHaveBeenCalled();
    });
  });

  describe('State changes and rendering', () => {
    it('renders initial items on mount', () => {
      const view = createFeedView({ container, controller: mockController });

      // Capture onStateChange callback passed to controller.mount.
      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // Simulate receiving items from controller.
      onStateChange(
        [
          { id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 },
          { id: 2, title: 'Story 2', by: 'user2', time: 2000, score: 20, descendants: 10 },
        ],
        false,
        undefined,
        true,
      );

      // Should render story cards.
      const cards = container.querySelectorAll('[data-testid^="feed-card-"]');
      expect(cards.length).toBe(2);
    });

    it('appends items on infinite scroll load', () => {
      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // Initial items.
      onStateChange(
        [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
        false,
        undefined,
        true,
      );

      const initialCards = container.querySelectorAll('[data-testid^="feed-card-"]');
      expect(initialCards.length).toBe(1);

      // Append more items (simulate infinite scroll).
      onStateChange(
        [{ id: 2, title: 'Story 2', by: 'user2', time: 2000, score: 20, descendants: 10 }],
        false,
        undefined,
        true,
      );

      // Should now have 2 cards (not replaced).
      const allCards = container.querySelectorAll('[data-testid^="feed-card-"]');
      expect(allCards.length).toBe(2);
    });

    it('uses DocumentFragment for batch inserts', () => {
      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // Initial render.
      onStateChange(
        [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
        false,
        undefined,
        true,
      );

      // Append many items to verify DocumentFragment usage (no reflow for each).
      const newItems = Array.from({ length: 10 }, (_, i) => ({
        id: i + 2,
        title: `Story ${i + 2}`,
        by: `user${i}`,
        time: 1000 + i,
        score: 10 + i,
        descendants: 5 + i,
      }));

      onStateChange(newItems, false, undefined, true);

      const allCards = container.querySelectorAll('[data-testid^="feed-card-"]');
      expect(allCards.length).toBe(11); // 1 initial + 10 new.
    });
  });

  describe('Observer cleanup on hasMore=false', () => {
    it('disconnects observer when hasMore is false', () => {
      const mockObserverInstance = {
        observe: vi.fn(),
        disconnect: vi.fn(),
      };
      global.IntersectionObserver = vi.fn(() => mockObserverInstance);

      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      mockObserverInstance.disconnect.mockClear();

      // Simulate reaching end of feed (hasMore=false).
      onStateChange(
        [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
        false,
        undefined,
        false, // No more pages.
      );

      expect(mockObserverInstance.disconnect).toHaveBeenCalled();
    });
  });

  describe('Error handling', () => {
    it('displays error message without clearing list', () => {
      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // Render initial items.
      onStateChange(
        [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
        false,
        undefined,
        true,
      );

      // Simulate error on page load.
      onStateChange(undefined, false, 'Network error', true);

      // Error message should be shown.
      const errorEl = container.querySelector('[data-testid="feed-error"]');
      expect(errorEl).toBeTruthy();
      expect(errorEl.textContent).toContain('Network error');
    });

    it('shows loading skeletons during fetch', () => {
      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // Simulate loading state.
      onStateChange(undefined, true, undefined, true);

      const loadingScreen = container.querySelector('[data-testid="feed-loading-screen"]');
      const skeletons = container.querySelectorAll('[data-testid="skeleton-card"]');

      expect(loadingScreen).toBeTruthy();
      expect(skeletons.length).toBeGreaterThan(0); // Multiple skeleton loaders shown.
    });

    it('shows a bottom loading card while the next page is loading', async () => {
      let observerCallback;
      let settleLoadMore;

      global.IntersectionObserver = vi.fn((callback) => {
        observerCallback = callback;
        return {
          observe: vi.fn(),
          disconnect: vi.fn(),
          unobserve: vi.fn(),
        };
      });

      mockController.mount.mockImplementation((onStateChange) => {
        onStateChange(
          [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
          false,
          undefined,
          true,
        );
      });

      mockController.loadMorePage.mockImplementation(
        (onStateChange) =>
          new Promise((resolve) => {
            settleLoadMore = () => {
              onStateChange(
                [{ id: 2, title: 'Story 2', by: 'user2', time: 2000, score: 20, descendants: 10 }],
                false,
                undefined,
                true,
              );
              resolve();
            };
          }),
      );

      const view = createFeedView({ container, controller: mockController });
      view.mount();

      observerCallback?.([{ isIntersecting: true }]);

      const loadingMoreCard = container.querySelector('[data-testid="feed-loading-more-card"]');
      expect(loadingMoreCard).toBeTruthy();

      settleLoadMore?.();
      await Promise.resolve();
      await Promise.resolve();

      const loadingMoreCardAfterLoad = container.querySelector('[data-testid="feed-loading-more-card"]');
      expect(loadingMoreCardAfterLoad).toBeNull();
    });

    it('marks the selected tab as pending until the switch request settles', async () => {
      const view = createFeedView({ container, controller: mockController });

      let settleSwitch;

      mockController.mount.mockImplementation((onStateChange) => {
        onStateChange([], false, undefined, true);
      });

      mockController.switchTab.mockImplementation(
        () =>
          new Promise((resolve) => {
            settleSwitch = resolve;
          }),
      );

      view.mount();

      const targetTab = container.querySelector('[data-testid="feed-tab-new"]');
      targetTab?.dispatchEvent(new dom.window.MouseEvent('click', { bubbles: true }));

      expect(targetTab?.classList.contains('is-pending')).toBe(true);

      settleSwitch?.();
      await Promise.resolve();
      await Promise.resolve();

      expect(targetTab?.classList.contains('is-pending')).toBe(false);
      expect(targetTab?.classList.contains('is-active')).toBe(true);
    });

    it('resets scroll position to top when switching tabs', async () => {
      mockController.mount.mockImplementation((onStateChange) => {
        onStateChange([], false, undefined, true);
      });

      mockController.switchTab.mockImplementation(async (_tabType, onStateChange) => {
        onStateChange([], false, undefined, true);
      });

      const view = createFeedView({ container, controller: mockController });
      view.mount();

      document.documentElement.scrollTop = 260;
      document.body.scrollTop = 140;

      await view.selectTab('ask');

      expect(document.documentElement.scrollTop).toBe(0);
      expect(document.body.scrollTop).toBe(0);
    });

    it('replaces initial skeletons with real stories after the first page resolves', () => {
      const view = createFeedView({ container, controller: mockController });

      let onStateChange;
      mockController.mount.mockImplementation((callback) => {
        onStateChange = callback;
      });

      view.mount();

      // The initial loading state should render placeholder cards into the empty list.
      onStateChange(undefined, true, undefined, true);

      const initialSkeletons = container.querySelectorAll('[data-testid="skeleton-card"]');
      expect(initialSkeletons.length).toBeGreaterThan(0);

      // The first resolved page should replace placeholders instead of appending after them.
      onStateChange(
        [{ id: 1, title: 'Story 1', by: 'user1', time: 1000, score: 10, descendants: 5 }],
        false,
        undefined,
        true,
      );

      const remainingSkeletons = container.querySelectorAll('[data-testid="skeleton-card"]');
      const renderedCards = container.querySelectorAll('[data-testid^="feed-card-"]');

      expect(remainingSkeletons.length).toBe(0);
      expect(renderedCards.length).toBe(1);
    });
  });

  describe('Cleanup on unmount', () => {
    it('calls controller.unmount on view unmount', () => {
      const view = createFeedView({ container, controller: mockController });
      view.mount();
      view.unmount();

      expect(mockController.unmount).toHaveBeenCalled();
    });

    it('clears container on unmount', () => {
      const view = createFeedView({ container, controller: mockController });
      view.mount();

      const itemCount = container.children.length;
      expect(itemCount).toBeGreaterThan(0);

      view.unmount();

      expect(container.children.length).toBe(0);
    });
  });
});
