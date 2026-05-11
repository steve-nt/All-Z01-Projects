// @vitest-environment jsdom

/*
 * Purpose: Verify feed-view card title links enforce safe external URL protocols.
 * Public API: createFeedView mount behavior and rendered data-testid DOM contract.
 * Notes: Tests keep controller mocks minimal so URL trust-boundary behavior is isolated.
 */

import { beforeEach, describe, expect, it, vi } from 'vitest';

import { createFeedView } from '../../../../src/features/feed/feed.view.js';

// Mock IntersectionObserver for mount tests.
beforeEach(() => {
  global.IntersectionObserver = vi.fn(() => ({
    observe: vi.fn(),
    disconnect: vi.fn(),
    unobserve: vi.fn(),
  }));
});

// Test-id lookup is centralized so assertions target the same public selectors used by e2e tests.
const getByTestId = (root, testId) => root.querySelector(`[data-testid="${testId}"]`);

// Controller factory keeps the view test focused on rendering outcomes rather than use-case plumbing.
const createControllerMock = (items) => {
  return {
    // Mount callback simulates an already-resolved feed so the view renders cards synchronously in the test.
    // New signature: onStateChange(items, isLoading, error, hasMore)
    mount: vi.fn((onStateChange) => {
      onStateChange(items, false, undefined, true);
    }),
    // Tab switching returns a resolved promise because the view click handler awaits this method.
    switchTab: vi.fn(async (_tabType, onStateChange) => {
      onStateChange(items, false, undefined, true);
    }),
    // Unmount exists for API completeness and to keep teardown assertions available if needed later.
    unmount: vi.fn(),
  };
};

describe('feed view url safety', () => {
  it('renders a detail link plus safe external source anchor when the item url uses https', () => {
    const container = document.createElement('section');
    const controller = createControllerMock([
      { id: 1, title: 'Safe story', url: 'https://example.com/path', by: 'alice', time: 0 },
    ]);

    const view = createFeedView({ container, controller });
    view.mount();

    const detailLink = getByTestId(container, 'post-link');
    const sourceLink = getByTestId(container, 'card-title-link');

    expect(detailLink?.getAttribute('href')).toBe('#/item/1');
    expect(sourceLink?.getAttribute('href')).toBe('https://example.com/path');

    view.unmount();
  });

  it('hides the external source link when the item url uses an unsafe protocol', () => {
    const container = document.createElement('section');
    const controller = createControllerMock([
      { id: 2, title: 'Unsafe story', url: 'javascript:alert(1)', by: 'bob', time: 0 },
    ]);

    const view = createFeedView({ container, controller });
    view.mount();

    const detailLink = getByTestId(container, 'post-link');
    const sourceLink = getByTestId(container, 'card-title-link');

    expect(detailLink?.getAttribute('href')).toBe('#/item/2');
    expect(sourceLink).toBeNull();

    view.unmount();
  });

  it('replaces existing cards on tab switches instead of appending stale results', async () => {
    const container = document.createElement('section');
    const loadedCount = document.createElement('span');
    loadedCount.setAttribute('data-testid', 'feed-loaded-count');
    loadedCount.textContent = 'Loaded: 0';
    document.body.append(loadedCount);

    const initialItems = [
      { id: 101, title: 'Top story', by: 'alice', time: 1 },
      { id: 102, title: 'Top story 2', by: 'bob', time: 2 },
    ];
    const switchedItems = [{ id: 201, title: 'New story', by: 'carol', time: 3 }];
    const controller = {
      mount: vi.fn((onStateChange) => {
        onStateChange(initialItems, false, undefined, true, false);
      }),
      switchTab: vi.fn(async (_tabType, onStateChange) => {
        onStateChange(switchedItems, false, undefined, true, false);
      }),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    };

    const view = createFeedView({ container, controller });
    view.mount();

    expect(container.querySelector('[data-testid="feed-card-101"]')).not.toBeNull();
    expect(container.querySelector('[data-testid="feed-card-102"]')).not.toBeNull();
    expect(loadedCount.textContent).toBe('Loaded: 2');

    await view.selectTab('new');

    expect(container.querySelector('[data-testid="feed-card-101"]')).toBeNull();
    expect(container.querySelector('[data-testid="feed-card-102"]')).toBeNull();
    expect(container.querySelector('[data-testid="feed-card-201"]')).not.toBeNull();
    expect(loadedCount.textContent).toBe('Loaded: 1');

    view.unmount();
    loadedCount.remove();
  });

  it('increments the loaded counter when appended pages are rendered', () => {
    const container = document.createElement('section');
    const loadedCount = document.createElement('span');
    loadedCount.setAttribute('data-testid', 'feed-loaded-count');
    loadedCount.textContent = 'Loaded: 0';
    document.body.append(loadedCount);

    let stateChangeCallback;
    const controller = {
      mount: vi.fn((onStateChange) => {
        stateChangeCallback = onStateChange;
        onStateChange([{ id: 301, title: 'Page 1', by: 'dave', time: 10 }], false, undefined, true, false);
      }),
      switchTab: vi.fn(async () => {}),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    };

    const view = createFeedView({ container, controller });
    view.mount();

    stateChangeCallback(
      [{ id: 302, title: 'Page 2', by: 'erin', time: 11 }],
      false,
      undefined,
      true,
      true,
    );

    expect(container.querySelector('[data-testid="feed-card-301"]')).not.toBeNull();
    expect(container.querySelector('[data-testid="feed-card-302"]')).not.toBeNull();
    expect(loadedCount.textContent).toBe('Loaded: 2');

    view.unmount();
    loadedCount.remove();
  });

  it('renders right-side refresh markers only for refreshed posts', async () => {
    const container = document.createElement('section');
    const newItems = [
      { id: 501, title: 'Old post', by: 'alice', time: 20 },
      { id: 502, title: 'Will update', by: 'bob', time: 21 },
      { id: 503, title: 'Brand new post', by: 'carol', time: 22 },
    ];
    const controller = {
      mount: vi.fn((onStateChange) => {
        onStateChange([], false, undefined, true, false);
      }),
      switchTab: vi.fn(async (_tabType, onStateChange) => {
        onStateChange(newItems, false, undefined, true, false);
      }),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    };

    const view = createFeedView({ container, controller });
    view.mount();

    await view.selectTab('new');
    await view.applyRefreshMarkers([502, 503]);

    const updatedMarker = getByTestId(container, 'feed-card-refresh-marker-502');
    const alsoUpdatedMarker = getByTestId(container, 'feed-card-refresh-marker-503');
    const untouchedMarker = getByTestId(container, 'feed-card-refresh-marker-501');
    const orderedCardIds = [...container.querySelectorAll('[data-testid^="feed-card-"][data-item-id]')].map((card) =>
      Number(card.getAttribute('data-item-id')),
    );

    expect(updatedMarker?.classList.contains('feed-card__refresh-marker--updated')).toBe(true);
    expect(alsoUpdatedMarker?.classList.contains('feed-card__refresh-marker--updated')).toBe(true);
    expect(untouchedMarker).toBeNull();
    expect(orderedCardIds).toEqual([501, 502, 503]);

    view.unmount();
  });

  it('shows an end-of-feed card when there are no more posts to load', () => {
    const container = document.createElement('section');
    const controller = {
      mount: vi.fn((onStateChange) => {
        onStateChange(
          [{ id: 801, title: 'Last story', by: 'alice', time: 90 }],
          false,
          undefined,
          false,
          false,
        );
      }),
      switchTab: vi.fn(async () => {}),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    };

    const view = createFeedView({ container, controller });
    view.mount();

    const endCard = getByTestId(container, 'feed-end-of-feed-card');

    expect(endCard).not.toBeNull();
    expect(endCard?.textContent).toContain('No more posts to load.');

    view.unmount();
  });

  it('clears applied refresh markers when clearRefreshMarkers is called', async () => {
    const container = document.createElement('section');
    const cards = [
      { id: 1, title: 'Story 1', by: 'alice', time: 1 },
      { id: 2, title: 'Story 2', by: 'bob', time: 2 },
      { id: 3, title: 'Story 3', by: 'carol', time: 3 },
    ];

    const controller = {
      mount: vi.fn((onStateChange) => {
        onStateChange([], false, undefined, true, false);
      }),
      switchTab: vi.fn(async (_tabType, onStateChange) => {
        onStateChange(cards, false, undefined, true, false);
      }),
      loadMorePage: vi.fn(),
      unmount: vi.fn(),
    };

    const view = createFeedView({ container, controller });
    view.mount();

    await view.selectTab('new');
    await view.applyRefreshMarkers([2, 3]);

    expect(getByTestId(container, 'feed-card-refresh-marker-2')).not.toBeNull();
    expect(getByTestId(container, 'feed-card-refresh-marker-3')).not.toBeNull();

    view.clearRefreshMarkers();

    expect(getByTestId(container, 'feed-card-refresh-marker-2')).toBeNull();
    expect(getByTestId(container, 'feed-card-refresh-marker-3')).toBeNull();

    view.unmount();
  });

  it('renders a floating back-to-top button that resets scroll positions', () => {
    const container = document.createElement('section');
    const controller = createControllerMock([
      { id: 601, title: 'Story for top button', by: 'alice', time: 50 },
    ]);

    const view = createFeedView({ container, controller });
    view.mount();

    const topButton = getByTestId(container, 'feed-scroll-top-button');

    document.documentElement.scrollTop = 420;
    document.body.scrollTop = 420;

    topButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(topButton).not.toBeNull();
    expect(document.documentElement.scrollTop).toBe(0);
    expect(document.body.scrollTop).toBe(0);

    view.unmount();
  });
});
