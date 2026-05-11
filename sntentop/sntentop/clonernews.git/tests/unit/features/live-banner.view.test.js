// @vitest-environment jsdom

// Live-banner view tests define the standalone DOM contract before app-shell wiring exists.
// Public API under test: createLiveBannerView and LIVE_BANNER_TEST_IDS from the feature view module.
// Constraints: tests stay within jsdom, use stable data-testid hooks, and avoid controller or polling integration.
import { describe, expect, it, vi } from 'vitest';

import {
  createLiveBannerView,
  LIVE_BANNER_TEST_IDS,
} from '../../../src/features/live-banner/live-banner.view.js';

// Data-testid lookup includes the root node itself so wrapper-level assertions stay direct.
const getByTestId = (root, testId) =>
  root?.matches?.(`[data-testid="${testId}"]`)
    ? root
    : root?.querySelector(`[data-testid="${testId}"]`);

describe('live-banner view', () => {
  it('renders an initially hidden polite status region', () => {
    // The banner starts hidden so the shell can mount it eagerly without announcing stale updates.
    const view = createLiveBannerView();

    const root = getByTestId(view.element, LIVE_BANNER_TEST_IDS.root);
    const message = getByTestId(view.element, LIVE_BANNER_TEST_IDS.message);
    const refreshButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.refreshButton);
    const clearButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.clearButton);

    expect(root?.hidden).toBe(true);
    expect(root?.getAttribute('role')).toBe('status');
    expect(root?.getAttribute('aria-live')).toBe('polite');
    expect(message?.textContent).toBe('');
    expect(refreshButton?.hidden).toBe(true);
    expect(clearButton?.hidden).toBe(true);

    view.destroy();
  });

  it('renders the update count and exposes the refresh action when visible', () => {
    // Visible state must communicate the count clearly and expose the explicit refresh affordance.
    const view = createLiveBannerView();

    view.render({
      isVisible: true,
      updateCount: 3,
      elapsedLabel: '12 seconds',
    });

    const root = getByTestId(view.element, LIVE_BANNER_TEST_IDS.root);
    const message = getByTestId(view.element, LIVE_BANNER_TEST_IDS.message);
    const refreshButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.refreshButton);
    const clearButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.clearButton);

    expect(root?.hidden).toBe(false);
    expect(message?.textContent).toContain('Newest information: 3 updates in the last 12 seconds');
    expect(refreshButton?.hidden).toBe(false);
    expect(refreshButton?.textContent).toBe('Refresh');
    expect(clearButton?.hidden).toBe(false);
    expect(clearButton?.textContent).toBe('Clear');

    view.destroy();
  });

  it('renders an explicit summary message with a clear-only action state', () => {
    // Post-refresh summaries should keep banner styling while hiding refresh until the user clears the popup.
    const view = createLiveBannerView();

    view.render({
      isVisible: true,
      messageText: 'Moved 1 refreshed post to top (0 new, 1 updated).',
      showRefresh: false,
      showClear: true,
    });

    const message = getByTestId(view.element, LIVE_BANNER_TEST_IDS.message);
    const refreshButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.refreshButton);
    const clearButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.clearButton);

    expect(message?.textContent).toBe('Moved 1 refreshed post to top (0 new, 1 updated).');
    expect(refreshButton?.hidden).toBe(true);
    expect(clearButton?.hidden).toBe(false);

    view.destroy();
  });

  it('calls onRefresh when the refresh control is clicked', () => {
    // Refresh remains an injected callback so the view stays decoupled from feed and app-shell logic.
    const onRefresh = vi.fn();
    const view = createLiveBannerView({ onRefresh });

    view.render({
      isVisible: true,
      updateCount: 2,
    });

    const refreshButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.refreshButton);

    refreshButton?.dispatchEvent(
      new MouseEvent('click', {
        bubbles: true,
      }),
    );

    expect(onRefresh).toHaveBeenCalledTimes(1);

    view.destroy();
  });

  it('calls onClear when the clear control is clicked', () => {
    // Clear remains independent from refresh so users can dismiss the popup without reloading feed content.
    const onClear = vi.fn();
    const view = createLiveBannerView({ onClear });

    view.render({
      isVisible: true,
      updateCount: 2,
    });

    const clearButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.clearButton);

    clearButton?.dispatchEvent(
      new MouseEvent('click', {
        bubbles: true,
      }),
    );

    expect(onClear).toHaveBeenCalledTimes(1);

    view.destroy();
  });

  it('rehides and clears stale banner content when rerendered as hidden', () => {
    // Hidden rerenders must clear prior text so old update counts are never re-announced accidentally.
    const view = createLiveBannerView();

    view.render({
      isVisible: true,
      updateCount: 4,
    });
    view.render({
      isVisible: false,
      updateCount: 0,
    });

    const root = getByTestId(view.element, LIVE_BANNER_TEST_IDS.root);
    const message = getByTestId(view.element, LIVE_BANNER_TEST_IDS.message);
    const refreshButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.refreshButton);
    const clearButton = getByTestId(view.element, LIVE_BANNER_TEST_IDS.clearButton);

    expect(root?.hidden).toBe(true);
    expect(message?.textContent).toBe('');
    expect(refreshButton?.hidden).toBe(true);
    expect(clearButton?.hidden).toBe(true);

    view.destroy();
  });
});
