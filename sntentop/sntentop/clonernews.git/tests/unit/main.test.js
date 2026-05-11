// @vitest-environment jsdom

// Main-app tests pin the route wiring between the shared router contract and the feed/detail/live-banner features.
// Public API under test: createApp start/stop behavior for feed and item routes.
// Constraints: tests inject fake router, controller, and view dependencies so no network or real shell code runs.
import { describe, expect, it, vi } from 'vitest';

import { createApp } from '../../src/main.js';

// Route fixtures stay tiny so each test can focus on one transition contract at a time.
const feedRoute = {
  name: 'feed',
  path: '/feed/top',
  params: { feedType: 'top' },
};

const pollFeedRoute = {
  name: 'feed',
  path: '/feed/poll',
  params: { feedType: 'poll' },
};

const itemRoute = {
  name: 'item',
  path: '/item/42',
  params: { id: 42 },
};

// Microtask flushing lets async route hydration settle before DOM assertions run.
const flushMicrotasks = async (turns = 4) => {
  for (let index = 0; index < turns; index += 1) {
    await Promise.resolve();
  }
};

// A tiny fake router keeps the tests focused on composition instead of real hashchange behavior.
const createFakeRouter = (initialRoute = feedRoute) => {
  let routeHandler = null;

  return {
    start: vi.fn(() => initialRoute),
    stop: vi.fn(),
    navigate: vi.fn(),
    subscribe: vi.fn((handler) => {
      routeHandler = handler;
      return () => {
        routeHandler = null;
      };
    }),
    emit(route) {
      routeHandler?.({ detail: route });
    },
  };
};

// Test container creation mirrors the app shell structure without depending on index.html.
const createTestContainer = () => {
  const app = document.createElement('div');
  app.id = 'app';

  const featureContainer = document.createElement('div');
  featureContainer.className = 'app-feature';

  const bootScreen = document.createElement('main');
  bootScreen.className = 'boot-screen';
  bootScreen.dataset.testid = 'boot-screen';
  bootScreen.textContent = 'Boot screen';

  featureContainer.append(bootScreen);
  app.append(featureContainer);
  document.body.replaceChildren(app);

  return featureContainer;
};

// Feed harness keeps app-shell tests deterministic by avoiding the real feed implementation.
const createFeedHarness = () => {
  const mount = vi.fn();
  const unmount = vi.fn();
  const selectTab = vi.fn(async () => {});
  const applyRefreshMarkers = vi.fn(async (updatedItemIds = []) => ({
    updatedItemIds,
  }));
  const clearRefreshMarkers = vi.fn();

  const createFeedViewImpl = vi.fn(({ container }) => {
    const element = document.createElement('section');
    element.dataset.testid = 'feed-view';

    mount.mockImplementation(() => {
      container.replaceChildren(element);
    });

    unmount.mockImplementation(() => {
      element.remove();
    });

    return {
      mount,
      unmount,
      selectTab,
      applyRefreshMarkers,
      clearRefreshMarkers,
    };
  });

  return {
    createFeedViewImpl,
    feedController: {
      switchTab: vi.fn(),
      mount: vi.fn(),
      unmount: vi.fn(),
      loadMorePage: vi.fn(),
    },
    mount,
    applyRefreshMarkers,
    clearRefreshMarkers,
    selectTab,
    unmount,
  };
};

// Banner harness keeps app-shell tests deterministic by avoiding real poll timers and network adapters.
const createLiveBannerHarness = () => {
  const mount = vi.fn();
  const unmount = vi.fn();
  const refresh = vi.fn(async () => ({ ok: true, data: undefined }));
  const destroy = vi.fn();

  return {
    createLiveBannerViewImpl: vi.fn(() => {
      const element = document.createElement('section');
      element.dataset.testid = 'live-banner';
      return {
        element,
        render: vi.fn(),
        destroy,
      };
    }),
    createLiveBannerControllerImpl: vi.fn(() => ({
      mount,
      unmount,
      refresh,
    })),
    pollUpdates: {
      poll: vi.fn(async () => ({ ok: true, data: {} })),
      subscribe: vi.fn(() => () => {}),
    },
    destroy,
    mount,
    refresh,
    unmount,
  };
};

describe('main app route wiring', () => {
  it('mounts the post-detail flow when the initial route is an item', async () => {
    // Initial item routes should render loading first, then the mapped success state once data arrives.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter(itemRoute);
    const load = vi.fn(async () => ({
      ok: true,
      data: {
        item: { id: 42, type: 'story' },
        comments: [],
        viewModel: {
          id: 42,
          type: 'story',
          title: 'Route-wired story',
          author: 'alice',
          time: null,
          score: 5,
          url: null,
          hasUrl: false,
          text: '',
          hasText: false,
        },
      },
    }));
    const render = vi.fn();
    const destroy = vi.fn();
    const createPostDetailViewImpl = vi.fn(() => {
      const element = document.createElement('section');
      element.dataset.testid = 'fake-detail-view';
      return { element, render, destroy };
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: { load },
      createPostDetailViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(load).toHaveBeenCalledWith(42);
    expect(createPostDetailViewImpl).toHaveBeenCalledTimes(1);
    expect(render).toHaveBeenNthCalledWith(1, { status: 'loading' });
    expect(render).toHaveBeenNthCalledWith(
      2,
      expect.objectContaining({
        status: 'success',
        data: expect.objectContaining({
          viewModel: expect.objectContaining({
            title: 'Route-wired story',
          }),
        }),
      }),
    );
    expect(featureContainer.querySelector('[data-testid="fake-detail-view"]')).not.toBeNull();
    expect(feedHarness.mount).not.toHaveBeenCalled();
    expect(liveBannerHarness.mount).not.toHaveBeenCalled();

    app.stop();
  });

  it('mounts the feed view and live banner when the initial route is feed', async () => {
    // Feed routes should compose the banner above the feed content and mount both pieces together.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter(feedRoute);

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.createFeedViewImpl).toHaveBeenCalledTimes(1);
    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.createLiveBannerControllerImpl).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.mount).toHaveBeenCalledTimes(1);
    expect(featureContainer.querySelector('[data-testid="app-feed-route"]')).not.toBeNull();
    expect(featureContainer.querySelector('[data-testid="app-feed-banner"]')).not.toBeNull();
    expect(featureContainer.querySelector('[data-testid="feed-view"]')).not.toBeNull();
    expect(featureContainer.querySelector('[data-testid="live-banner"]')).not.toBeNull();

    app.stop();

    expect(feedHarness.unmount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.unmount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.destroy).toHaveBeenCalledTimes(1);
  });

  it('restores the selected feed tab from the initial route on refresh', async () => {
    // Feed routes with a tab segment should apply that tab immediately after mount.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter(pollFeedRoute);

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenCalledWith('poll');

    app.stop();
  });

  it('applies refresh markers on the current feed when the banner refresh action is clicked', async () => {
    // Refresh should only apply marker symbols and keep the current feed/tab selection unchanged.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const router = createFakeRouter({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });
    const mount = vi.fn();
    const unmount = vi.fn();
    const destroy = vi.fn();

    const createLiveBannerViewImpl = vi.fn(({ onRefresh }) => {
      const element = document.createElement('section');
      const refreshButton = document.createElement('button');
      refreshButton.type = 'button';
      refreshButton.dataset.testid = 'live-banner-refresh';
      refreshButton.addEventListener('click', () => {
        onRefresh();
      });
      element.append(refreshButton);

      return {
        element,
        render: vi.fn(),
        destroy,
      };
    });

    const createLiveBannerControllerImpl = vi.fn(({ onRefresh }) => ({
      mount,
      unmount,
      refresh: vi.fn(async () => onRefresh()),
    }));

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: {
        poll: vi.fn(async () => ({ ok: true, data: {} })),
        subscribe: vi.fn(() => () => {}),
      },
      createLiveBannerViewImpl,
      createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('new');

    const refreshButton = featureContainer.querySelector('[data-testid="live-banner-refresh"]');
    refreshButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('new');
    expect(feedHarness.applyRefreshMarkers).toHaveBeenCalledWith([]);

    app.stop();
  });

  it('keeps the current feed and applies markers when refresh is clicked from a non-New tab', async () => {
    // Refresh should preserve active feed context and only update marker state.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const router = createFakeRouter(feedRoute);
    const mount = vi.fn();
    const unmount = vi.fn();
    const destroy = vi.fn();
    let updatesSubscriber = null;

    const createLiveBannerViewImpl = vi.fn(({ onRefresh }) => {
      const element = document.createElement('section');
      const refreshButton = document.createElement('button');
      refreshButton.type = 'button';
      refreshButton.dataset.testid = 'live-banner-refresh';
      refreshButton.addEventListener('click', () => {
        onRefresh();
      });
      element.append(refreshButton);

      return {
        element,
        render: vi.fn(),
        destroy,
      };
    });

    const createLiveBannerControllerImpl = vi.fn(({ onRefresh }) => ({
      mount,
      unmount,
      refresh: vi.fn(async () => onRefresh()),
      clear: vi.fn(async () => ({ ok: true, data: undefined })),
    }));

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: {
        poll: vi.fn(async () => ({ ok: true, data: {} })),
        subscribe: vi.fn((callback) => {
          updatesSubscriber = callback;
          return () => {
            updatesSubscriber = null;
          };
        }),
      },
      createLiveBannerViewImpl,
      createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    updatesSubscriber?.({
      newIds: [901, 902],
      currentIds: [901, 902],
      previousIds: [900],
      polledAtMs: 10_000,
      isFirstPoll: false,
    });

    const refreshButton = featureContainer.querySelector('[data-testid="live-banner-refresh"]');
    refreshButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenNthCalledWith(1, 'top');
    expect(feedHarness.applyRefreshMarkers).toHaveBeenCalledWith([901, 902]);

    app.stop();
  });

  it('clears pending updates from banner and removes active refresh markers when clear is clicked', async () => {
    // Clear should dismiss popup state and clear marker badges without switching tabs.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const router = createFakeRouter(feedRoute);
    let updatesSubscriber = null;

    const createLiveBannerViewImpl = vi.fn(({ onClear }) => {
      const element = document.createElement('section');
      const clearButton = document.createElement('button');
      clearButton.type = 'button';
      clearButton.dataset.testid = 'live-banner-clear';
      clearButton.addEventListener('click', () => {
        onClear();
      });
      element.append(clearButton);

      return {
        element,
        render: vi.fn(),
        destroy: vi.fn(),
      };
    });

    const clear = vi.fn(async () => ({ ok: true, data: undefined }));
    const createLiveBannerControllerImpl = vi.fn(({ onClear }) => ({
      mount: vi.fn(),
      unmount: vi.fn(),
      refresh: vi.fn(async () => ({ ok: true, data: undefined })),
      clear: vi.fn(async () => {
        await onClear();
        return clear();
      }),
    }));

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: {
        poll: vi.fn(async () => ({ ok: true, data: {} })),
        subscribe: vi.fn((callback) => {
          updatesSubscriber = callback;
          return () => {
            updatesSubscriber = null;
          };
        }),
      },
      createLiveBannerViewImpl,
      createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    updatesSubscriber?.({
      newIds: [701],
      currentIds: [701],
      previousIds: [],
      polledAtMs: 12_000,
      isFirstPoll: false,
    });

    const clearButton = featureContainer.querySelector('[data-testid="live-banner-clear"]');
    clearButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    await flushMicrotasks();

    expect(clear).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenCalledTimes(1);
    expect(feedHarness.clearRefreshMarkers).toHaveBeenCalledTimes(1);
    expect(feedHarness.applyRefreshMarkers).not.toHaveBeenCalled();

    app.stop();
  });

  it('resets back to top feed when routing from another feed to /feed/top', async () => {
    // Returning to top must actively re-select top so controller state cannot stay pinned to another tab.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('new');

    router.emit(feedRoute);
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenCalledTimes(2);
    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('top');

    app.stop();
  });

  it('returns from detail to the previously active feed route on back', async () => {
    // Back from detail should preserve route context (for example /feed/new) instead of forcing top.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });

    let backHandler = null;
    const createPostDetailViewImpl = vi.fn(({ onBack }) => {
      backHandler = onBack;
      const element = document.createElement('section');
      element.dataset.testid = 'fake-detail-view';

      return {
        element,
        render: vi.fn(),
        destroy: vi.fn(),
      };
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: {
        load: vi.fn(async () => ({ ok: false, error: 'Not found', reason: 'not-found' })),
      },
      createPostDetailViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('new');

    router.emit(itemRoute);
    await flushMicrotasks();

    expect(typeof backHandler).toBe('function');

    backHandler?.();

    expect(router.navigate).toHaveBeenCalledWith('/feed/new');

    app.stop();
  });

  it('returns from poll detail to poll feed on back', async () => {
    // Poll detail back should respect the originating poll feed route.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter({
      name: 'feed',
      path: '/feed/poll',
      params: { feedType: 'poll' },
    });

    let backHandler = null;
    const createPostDetailViewImpl = vi.fn(({ onBack }) => {
      backHandler = onBack;
      const element = document.createElement('section');
      element.dataset.testid = 'fake-detail-view';

      return {
        element,
        render: vi.fn(),
        destroy: vi.fn(),
      };
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: {
        load: vi.fn(async () => ({ ok: false, error: 'Not found', reason: 'not-found' })),
      },
      createPostDetailViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.selectTab).toHaveBeenLastCalledWith('poll');

    router.emit(itemRoute);
    await flushMicrotasks();

    backHandler?.();

    expect(router.navigate).toHaveBeenCalledWith('/feed/poll');

    app.stop();
  });

  it('restores prior feed scroll position when returning from detail to the same feed route', async () => {
    // Returning to the same feed route should reuse the preserved feed DOM and restore the previous viewport offset.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });
    const scrollToSpy = vi.fn();

    Object.defineProperty(window, 'scrollTo', {
      value: scrollToSpy,
      configurable: true,
      writable: true,
    });

    Object.defineProperty(window, 'scrollY', {
      value: 640,
      configurable: true,
      writable: true,
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: {
        load: vi.fn(async () => ({ ok: false, error: 'Not found', reason: 'not-found' })),
      },
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(feedHarness.createFeedViewImpl).toHaveBeenCalledTimes(1);

    router.emit(itemRoute);
    await flushMicrotasks();

    router.emit({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });
    await flushMicrotasks();

    expect(feedHarness.createFeedViewImpl).toHaveBeenCalledTimes(1);
    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(scrollToSpy).toHaveBeenCalledWith({ top: 640, left: 0, behavior: 'auto' });

    app.stop();
  });

  it('switches from detail to feed by unmounting the detail view and mounting the feed shell', async () => {
    // Routing from item to feed should clean up the detail view before the banner and feed are mounted.
    const featureContainer = createTestContainer();
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter(itemRoute);
    const load = vi.fn(async () => ({
      ok: true,
      data: {
        item: { id: 42, type: 'job' },
        comments: [],
        viewModel: {
          id: 42,
          type: 'job',
          title: 'Route-wired job',
          author: 'team',
          time: null,
          score: null,
          url: null,
          hasUrl: false,
          text: '',
          hasText: false,
        },
      },
    }));
    const destroy = vi.fn();
    const createPostDetailViewImpl = vi.fn(() => {
      const element = document.createElement('section');
      element.dataset.testid = 'fake-detail-view';
      return {
        element,
        render: vi.fn(),
        destroy,
      };
    });

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: { load },
      createPostDetailViewImpl,
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    router.emit(feedRoute);
    await flushMicrotasks();

    expect(destroy).toHaveBeenCalledTimes(1);
    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.mount).toHaveBeenCalledTimes(1);
    expect(featureContainer.querySelector('[data-testid="app-feed-route"]')).not.toBeNull();
    expect(featureContainer.querySelector('[data-testid="feed-view"]')).not.toBeNull();
    expect(featureContainer.querySelector('[data-testid="live-banner"]')).not.toBeNull();

    app.stop();

    expect(feedHarness.unmount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.unmount).toHaveBeenCalledTimes(1);
    expect(liveBannerHarness.destroy).toHaveBeenCalledTimes(1);
  });

  it('lets sticky header tabs route from detail and apply the selected feed tab', async () => {
    // Header tabs must remain interactive on detail routes and carry the chosen feed type into feed mount.
    document.body.innerHTML = `
      <div id="app">
        <header>
          <button type="button" data-feed-type="top" data-testid="feed-tab-top">Top</button>
          <button type="button" data-feed-type="new" data-testid="feed-tab-new">New</button>
          <button type="button" data-feed-type="job" data-testid="feed-tab-jobs">Jobs</button>
          <button type="button" data-feed-type="poll" data-testid="feed-tab-polls">Polls</button>
          <button type="button" data-feed-type="ask" data-testid="feed-tab-ask">Ask HN</button>
          <button type="button" data-feed-type="show" data-testid="feed-tab-show">Show HN</button>
        </header>
        <div class="app-feature"></div>
      </div>
    `;

    const featureContainer = document.querySelector('.app-feature');
    const feedHarness = createFeedHarness();
    const liveBannerHarness = createLiveBannerHarness();
    const router = createFakeRouter(itemRoute);
    const load = vi.fn(async () => ({
      ok: true,
      data: {
        item: { id: 42, type: 'story' },
        comments: [],
        viewModel: {
          id: 42,
          type: 'story',
          title: 'Route-wired story',
          author: 'alice',
          time: null,
          score: 5,
          url: null,
          hasUrl: false,
          text: '',
          hasText: false,
        },
      },
    }));

    const app = createApp({
      featureContainer,
      router,
      feedController: feedHarness.feedController,
      createFeedViewImpl: feedHarness.createFeedViewImpl,
      postDetailController: { load },
      pollUpdates: liveBannerHarness.pollUpdates,
      createLiveBannerViewImpl: liveBannerHarness.createLiveBannerViewImpl,
      createLiveBannerControllerImpl: liveBannerHarness.createLiveBannerControllerImpl,
    });

    app.start();
    await flushMicrotasks();

    const newTab = document.querySelector('[data-testid="feed-tab-new"]');
    newTab?.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(router.navigate).toHaveBeenCalledWith('/feed/new');

    router.emit({
      name: 'feed',
      path: '/feed/new',
      params: { feedType: 'new' },
    });
    await flushMicrotasks();

    expect(feedHarness.mount).toHaveBeenCalledTimes(1);
    expect(feedHarness.selectTab).toHaveBeenCalledWith('new');

    app.stop();
  });
});
