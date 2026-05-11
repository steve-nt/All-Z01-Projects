/*
 * Purpose: Wire the routed shell together with feed, detail, and live-banner feature controllers.
 * Public API: createApp({ ...dependencies }) -> { start(), stop() }.
 * Constraints: route composition stays in the shell, DOM creation stays in shared helpers and feature views, and cleanup must be deterministic.
 */

import './app.css';

import { createGetItemUseCase } from './core/use-cases/get-item.js';
import { createListItemsUseCase } from './core/use-cases/list-items.js';
import { createPollUpdatesUseCase } from './core/use-cases/poll-updates.js';
import { createFeedController } from './features/feed/feed.controller.js';
import { createFeedView } from './features/feed/feed.view.js';
import { createLiveBannerController } from './features/live-banner/live-banner.controller.js';
import { createLiveBannerView } from './features/live-banner/live-banner.view.js';
import { createPostDetailController } from './features/post-detail/post-detail.controller.js';
import {
  createLoadingPostDetailState,
  toPostDetailRenderState,
} from './features/post-detail/post-detail.render-state.js';
import { createPostDetailView } from './features/post-detail/post-detail.view.js';
import { createHnApiAdapter } from './infra/hn-api-adapter.js';
import {
  appendChildren,
  createElement,
  getElementById,
  isHtmlElement,
  replaceElementChildren,
  setDocumentTitle,
} from './shared/dom-helpers.js';
import { createHashRouter } from './shared/router.js';
import { Temporal } from '@js-temporal/polyfill';

// The app title resets on route changes so both feed and detail screens share one stable baseline.
const APP_TITLE = 'clonernews';
const FEED_TAB_SELECTOR = '[data-feed-type]';
const FEED_TYPES = new Set(['top', 'new', 'job', 'poll', 'ask', 'show']);
const TOP_FEED_ROUTE_PATH = '/feed/top';
const toFeedRoutePath = (feedType) => `/feed/${feedType}`;
const toRouteFeedType = (route) => {
  const routeFeedType = route?.params?.feedType;
  return typeof routeFeedType === 'string' && FEED_TYPES.has(routeFeedType) ? routeFeedType : 'top';
};

// The feed controller keeps Track B isolated while still letting the shell supply concrete contracts.
const createDefaultFeedController = () =>
  createFeedController({
    listItems: createListItemsUseCase({
      api: createHnApiAdapter(),
    }),
  });

// The post-detail use case stays behind a tiny adapter factory so the shell never owns transport details.
const createDefaultPostDetailController = () =>
  createPostDetailController({
    getItem: createGetItemUseCase({
      api: createHnApiAdapter(),
    }),
  });

// Live polling stays behind a core use case so the banner only consumes a Result-driven contract.
const createDefaultPollUpdates = () =>
  createPollUpdatesUseCase({
    api: createHnApiAdapter(),
  });

const ok = (data) => ({ ok: true, data });
const toNowMs = () => Temporal.Now.instant().epochMilliseconds;

/**
 * @param {{
 *   featureContainer?: HTMLElement | null,
 *   router?: ReturnType<typeof createHashRouter>,
 *   feedController?: ReturnType<typeof createDefaultFeedController>,
 *   createFeedViewImpl?: typeof createFeedView,
 *   postDetailController?: { load(id: number): Promise<{ ok: true, data: object } | { ok: false, error: string, reason?: string }> },
 *   pollUpdates?: { poll(): Promise<{ ok: true, data: object } | { ok: false, error: string }>, subscribe(callback: (payload: object) => void): () => void },
 *   createPostDetailViewImpl?: typeof createPostDetailView,
 *   createLiveBannerViewImpl?: typeof createLiveBannerView,
 *   createLiveBannerControllerImpl?: typeof createLiveBannerController,
 * }=} dependencies
 * @returns {{ start(): void, stop(): void }}
 */
export const createApp = ({
  featureContainer = getElementById('app')?.querySelector('.app-feature'),
  router = createHashRouter(),
  feedController = createDefaultFeedController(),
  createFeedViewImpl = createFeedView,
  postDetailController = createDefaultPostDetailController(),
  createPostDetailViewImpl = createPostDetailView,
  pollUpdates = createDefaultPollUpdates(),
  createLiveBannerViewImpl = createLiveBannerView,
  createLiveBannerControllerImpl = createLiveBannerController,
} = {}) => {
  // Missing feature containers are tolerated so tests and partial shells can import the module safely.
  if (!isHtmlElement(featureContainer)) {
    return {
      start() {},
      stop() {},
    };
  }

  let currentView = null;
  let unsubscribeRoute = null;
  let cleanupHeaderTabs = () => {};
  let activeRouteToken = 0;
  let pendingFeedType = null;
  let lastFeedRoutePath = TOP_FEED_ROUTE_PATH;
  let preservedFeedRoute = null;
  let pendingLiveUpdateItemIds = new Set();
  let pendingLiveUpdateSinceMs = null;
  let refreshMarkersByFeedType = new Map();
  let unsubscribeGlobalLiveUpdateTracking = null;

  // Live-update tracking is app-global so feed route remounts preserve banner counters and pending ids.
  const trackPendingLiveUpdates = (payload) => {
    if (payload?.isFirstPoll === true) {
      return;
    }

    const newIds = Array.isArray(payload?.newIds) ? payload.newIds : [];
    const hadNoPendingUpdates = pendingLiveUpdateItemIds.size === 0;

    for (const id of newIds) {
      if (Number.isInteger(id) && id > 0) {
        pendingLiveUpdateItemIds.add(id);
      }
    }

    if (hadNoPendingUpdates && pendingLiveUpdateItemIds.size > 0) {
      pendingLiveUpdateSinceMs = Number.isInteger(payload?.polledAtMs) ? payload.polledAtMs : toNowMs();
      // Marker persistence expires once a new refresh popup cycle starts.
      refreshMarkersByFeedType = new Map();
    }
  };

  // Route teardown stays centralized so feed, detail, and banner cleanup always happen together.
  const clearCurrentRoute = ({ preserveFeed = false } = {}) => {
    if (preserveFeed && currentView?.kind === 'feed' && currentView.element) {
      preservedFeedRoute = {
        path: lastFeedRoutePath,
        view: currentView,
        element: currentView.element,
        scrollY: typeof window !== 'undefined' ? window.scrollY ?? 0 : 0,
      };

      currentView = null;
      replaceElementChildren(featureContainer);
      return;
    }

    currentView?.destroy?.();
    currentView?.unmount?.();
    currentView = null;
    replaceElementChildren(featureContainer);
  };

  // Sticky-header feed tabs must work from any route, so feed selection is coordinated at shell level.
  const bindHeaderFeedTabs = () => {
    const tabButtons = [...document.querySelectorAll(FEED_TAB_SELECTOR)].filter((button) => {
      const feedType = button.getAttribute('data-feed-type');
      return feedType !== null && FEED_TYPES.has(feedType);
    });

    const listenerCleanups = tabButtons.map((button) => {
      const handleClick = () => {
        const feedType = button.getAttribute('data-feed-type');

        if (feedType === null || !FEED_TYPES.has(feedType)) {
          return;
        }

        pendingFeedType = feedType;
        router.navigate(toFeedRoutePath(feedType));
      };

      button.addEventListener('click', handleClick);

      return () => {
        button.removeEventListener('click', handleClick);
      };
    });

    cleanupHeaderTabs = () => {
      for (const cleanup of listenerCleanups) {
        cleanup();
      }
    };
  };

  // Feed rendering composes the banner above the feed content without forcing either feature to know about the other.
  const renderFeedRoute = (route) => {
    activeRouteToken += 1;
    setDocumentTitle(APP_TITLE);

    const routeFeedType = toRouteFeedType(route);
    const nextFeedType = pendingFeedType ?? routeFeedType;
    pendingFeedType = null;
    const targetFeedRoutePath = toFeedRoutePath(nextFeedType);
    // Preserve the active feed path so detail back-navigation can restore where the user came from.
    lastFeedRoutePath = targetFeedRoutePath;

    if (preservedFeedRoute && preservedFeedRoute.path === targetFeedRoutePath) {
      clearCurrentRoute();
      replaceElementChildren(featureContainer, preservedFeedRoute.element);
      currentView = preservedFeedRoute.view;

      const { scrollY } = preservedFeedRoute;
      preservedFeedRoute = null;

      if (typeof window !== 'undefined' && typeof window.scrollTo === 'function') {
        queueMicrotask(() => {
          try {
            window.scrollTo({ top: scrollY, left: 0, behavior: 'auto' });
          } catch {
            window.scrollTo(0, scrollY);
          }
        });
      }

      return;
    }

    if (preservedFeedRoute) {
      preservedFeedRoute.view.destroy?.();
      preservedFeedRoute = null;
    }

    clearCurrentRoute();

    let activeFeedType = nextFeedType ?? 'top';

    const feedRouteElement = createElement('section', {
      className: 'app-route app-route--feed',
      attributes: {
        'data-testid': 'app-feed-route',
      },
    });

    const bannerContainer = createElement('div', {
      className: 'app-route__banner',
      attributes: {
        'data-testid': 'app-feed-banner',
      },
    });

    const feedContainer = createElement('div', {
      className: 'app-route__feed',
      attributes: {
        'data-testid': 'app-feed-content',
      },
    });

    appendChildren(feedRouteElement, [bannerContainer, feedContainer]);
    replaceElementChildren(featureContainer, feedRouteElement);

    const feedView = createFeedViewImpl({
      container: feedContainer,
      controller: feedController,
    });

    const selectFeedTab = async (feedType, options) => {
      activeFeedType = feedType;

      if (options === undefined) {
        return feedView.selectTab?.(feedType);
      }

      return feedView.selectTab?.(feedType, options);
    };

    let bannerController = null;

    // The banner button simply reuses the controller refresh contract so the view stays decoupled from shell logic.
    const bannerView = createLiveBannerViewImpl({
      onRefresh: async () => {
        if (typeof bannerController?.refresh !== 'function') {
          return ok(undefined);
        }

        return bannerController.refresh();
      },
      onClear: async () => {
        if (typeof bannerController?.clear !== 'function') {
          return ok(undefined);
        }

        return bannerController.clear();
      },
    });

    bannerController = createLiveBannerControllerImpl({
      pollUpdates,
      view: bannerView,
      initialPendingUpdateIds: [...pendingLiveUpdateItemIds],
      initialPendingSinceMs: pendingLiveUpdateSinceMs,
      onRefresh: async () => {
        const normalizedRefreshedIds = [
          ...new Set(
            [...pendingLiveUpdateItemIds].filter((itemId) => Number.isInteger(itemId) && itemId > 0),
          ),
        ];

        // A new refresh cycle replaces prior badge state so symbols only reflect the latest acknowledged updates.
        refreshMarkersByFeedType = new Map();

        if (normalizedRefreshedIds.length > 0) {
          refreshMarkersByFeedType.set(activeFeedType, {
            updatedItemIds: normalizedRefreshedIds,
          });
        }

        if (typeof feedView.applyRefreshMarkers === 'function') {
          await feedView.applyRefreshMarkers(normalizedRefreshedIds);
        }

        pendingLiveUpdateItemIds = new Set();
        pendingLiveUpdateSinceMs = null;

        return undefined;
      },
      onClear: async () => {
        pendingLiveUpdateItemIds = new Set();
        pendingLiveUpdateSinceMs = null;
        refreshMarkersByFeedType = new Map();

        if (typeof feedView.clearRefreshMarkers === 'function') {
          feedView.clearRefreshMarkers();
        }

        return undefined;
      },
    });

    currentView = {
      kind: 'feed',
      element: feedRouteElement,
      destroy() {
        bannerController?.unmount?.();
        bannerView.destroy?.();
        feedView.unmount?.();
      },
    };

    bannerContainer.append(bannerView.element);
    bannerController.mount();
    feedView.mount();

    if (nextFeedType !== null) {
      const persistedRefreshMarkers = refreshMarkersByFeedType.get(nextFeedType);

      void (async () => {
        await selectFeedTab(nextFeedType);

        if (
          persistedRefreshMarkers &&
          Array.isArray(persistedRefreshMarkers.updatedItemIds) &&
          typeof feedView.applyRefreshMarkers === 'function'
        ) {
          await feedView.applyRefreshMarkers(persistedRefreshMarkers.updatedItemIds, {
            isMarkerReplay: true,
          });
        }
      })();
    }
  };

  const renderItemRoute = async (itemId) => {
    const routeToken = activeRouteToken + 1;
    activeRouteToken = routeToken;
    clearCurrentRoute({ preserveFeed: true });

    const view = createPostDetailViewImpl({
      onBack: () => {
        router.navigate(lastFeedRoutePath);
      },
    });

    currentView = view;
    replaceElementChildren(featureContainer, view.element);
    view.render(createLoadingPostDetailState());

    const loadResult = await postDetailController.load(itemId);

    // Stale async results are ignored so rapid route changes do not repaint an outdated detail screen.
    if (routeToken !== activeRouteToken || currentView !== view) {
      return;
    }

    const renderState = toPostDetailRenderState(loadResult);
    view.render(renderState);
    setDocumentTitle(
      renderState.status === 'success'
        ? `${renderState.data.viewModel.title} · ${APP_TITLE}`
        : APP_TITLE,
    );
  };

  // Route composition stays thin by mapping the shared router contract directly to feed or detail rendering.
  const handleRoute = (route) => {
    if (route?.name === 'item' && Number.isInteger(route.params?.id)) {
      void renderItemRoute(route.params.id);
      return;
    }

    renderFeedRoute(route);
  };

  return {
    // Startup renders the initial route once and then subscribes to later routechange events.
    start() {
      setDocumentTitle(APP_TITLE);
      bindHeaderFeedTabs();

      if (typeof pollUpdates?.subscribe === 'function') {
        unsubscribeGlobalLiveUpdateTracking = pollUpdates.subscribe(trackPendingLiveUpdates);
      }

      const initialRoute = router.start();
      unsubscribeRoute = router.subscribe((event) => {
        handleRoute(event.detail);
      });
      handleRoute(initialRoute);
    },
    // Stop removes route listeners and active feature views so repeated test runs clean up safely.
    stop() {
      activeRouteToken += 1;
      cleanupHeaderTabs();
      cleanupHeaderTabs = () => {};
      unsubscribeRoute?.();
      unsubscribeRoute = null;
      unsubscribeGlobalLiveUpdateTracking?.();
      unsubscribeGlobalLiveUpdateTracking = null;
      router.stop();
      clearCurrentRoute();

      if (preservedFeedRoute) {
        preservedFeedRoute.view.destroy?.();
        preservedFeedRoute = null;
      }
    },
  };
};

// The browser entrypoint auto-starts, while tests import createApp directly and control startup themselves.
if (import.meta.env.MODE !== 'test') {
  createApp().start();
}
