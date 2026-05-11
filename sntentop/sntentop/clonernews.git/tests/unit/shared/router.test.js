// @vitest-environment jsdom

// Router tests need a browser-like target because hash navigation and events are DOM-driven.
// Public API under test: createHashRouter parseRoute, subscribe, navigate, start, and stop behavior.
// Constraints: route assertions stay hash-based and avoid coupling to feature-view rendering details.
import { describe, expect, it } from 'vitest';

import { createHashRouter } from '../../../src/shared/router.js';

describe('router', () => {
  it('parses the root and item routes', () => {
    // Parsing is the core contract because the shell depends on route shape, not raw hashes.
    const router = createHashRouter();

    expect(router.parseRoute('#/')).toEqual({
      name: 'feed',
      path: '/feed/top',
      params: { feedType: 'top' },
    });

    expect(router.parseRoute('#/item/42')).toEqual({
      name: 'item',
      path: '/item/42',
      params: { id: 42 },
    });

    expect(router.parseRoute('#/feed/poll')).toEqual({
      name: 'feed',
      path: '/feed/poll',
      params: { feedType: 'poll' },
    });
  });

  it('emits routechange when navigating', () => {
    // The event stream verifies that mounted features can react without polling.
    const router = createHashRouter();
    const routes = [];

    const stopListening = router.subscribe((event) => {
      routes.push(event.detail);
    });

    router.start();
    router.navigate('/item/9');
    router.navigate('/');
    stopListening();
    router.stop();

    expect(routes.at(-1)).toEqual({
      name: 'feed',
      path: '/feed/top',
      params: { feedType: 'top' },
    });
    expect(routes.some((route) => route.name === 'item' && route.params.id === 9)).toBe(true);
  });
});
