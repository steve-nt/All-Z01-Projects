// Hash routing keeps the app compatible with static hosting and avoids server-side rewrite requirements.
const ROUTE_EVENT = 'routechange';
const FEED_TYPES = new Set(['top', 'new', 'job', 'poll', 'ask', 'show']);
const TOP_FEED_PATH = '/feed/top';

// Normalize loose hash inputs into a consistent path form for the rest of the router.
const normalizeHash = (value) => {
  if (!value || value === '#' || value === '#/') {
    return TOP_FEED_PATH;
  }

  return value.startsWith('#') ? value.slice(1) || '/' : value;
};

// Route parsing stays centralized so the app shell can mount views from a single route contract.
const parseRoute = (hash) => {
  const path = normalizeHash(hash).replace(/^\/?/, '/');

  if (path === '/' || path === '') {
    return {
      name: 'feed',
      path: TOP_FEED_PATH,
      params: { feedType: 'top' },
    };
  }

  const feedMatch = path.match(/^\/feed\/(top|new|job|poll|ask|show)$/);

  if (feedMatch && FEED_TYPES.has(feedMatch[1])) {
    return {
      name: 'feed',
      path: `/feed/${feedMatch[1]}`,
      params: {
        feedType: feedMatch[1],
      },
    };
  }

  const itemMatch = path.match(/^\/item\/(\d+)$/);

  if (itemMatch) {
    return {
      name: 'item',
      path: `/item/${itemMatch[1]}`,
      params: {
        id: Number(itemMatch[1]),
      },
    };
  }

  return {
    name: 'not-found',
    path,
    params: {},
  };
};

// The router exposes a tiny event-driven API so the shell can react without a framework.
export const createHashRouter = ({ target = window } = {}) => {
  let lastEmittedHash = null;

  // Emitting only on actual changes keeps view teardown/mount cycles predictable.
  const emitRoute = (hash = target.location.hash, force = false) => {
    const normalizedHash = normalizeHash(hash);

    if (!force && normalizedHash === lastEmittedHash) {
      return parseRoute(normalizedHash);
    }

    lastEmittedHash = normalizedHash;
    const route = parseRoute(normalizedHash);

    target.dispatchEvent(
      new CustomEvent(ROUTE_EVENT, {
        detail: route,
      }),
    );

    return route;
  };

  // The native hashchange event is enough because the app only needs client-side fragment routing.
  const handleHashChange = () => {
    emitRoute(target.location.hash);
  };

  return {
    start() {
      if (!target.location.hash || target.location.hash === '#' || target.location.hash === '#/') {
        target.location.hash = `#${TOP_FEED_PATH}`;
      }

      target.addEventListener('hashchange', handleHashChange);
      return emitRoute(target.location.hash, true);
    },
    stop() {
      target.removeEventListener('hashchange', handleHashChange);
    },
    navigate(path) {
      const normalizedPath = path === '/' ? TOP_FEED_PATH : path;
      const nextHash = normalizedPath.startsWith('#')
        ? normalizedPath
        : `#${normalizedPath.startsWith('/') ? normalizedPath : `/${normalizedPath}`}`;
      target.location.hash = nextHash;
      return emitRoute(nextHash, true);
    },
    getCurrentRoute() {
      return parseRoute(target.location.hash);
    },
    subscribe(handler) {
      target.addEventListener(ROUTE_EVENT, handler);
      return () => target.removeEventListener(ROUTE_EVENT, handler);
    },
    parseRoute,
  };
};
