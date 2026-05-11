/*
 * Purpose: Render the live-data notification banner as a standalone presentation-only feature view.
 * Public API: createLiveBannerView({ onRefresh? }) -> { element, render(state), destroy() } plus LIVE_BANNER_TEST_IDS.
 * Constraints: This module owns only DOM creation and event wiring, accepts simple render state, and contains no polling or app-shell logic.
 */

import './live-banner.css';
import { createElement, setText } from '../../shared/dom-helpers.js';

// Stable test IDs define the public DOM contract for unit tests and later integration coverage.
export const LIVE_BANNER_TEST_IDS = Object.freeze({
  root: 'live-banner',
  message: 'live-banner-message',
  refreshButton: 'live-banner-refresh',
  clearButton: 'live-banner-clear',
});

// Positive integer counts are the only values that should produce a visible update banner.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Hidden-state rendering is centralized so initial, invalid, and reset paths stay identical.
const applyHiddenState = ({ element, message, refreshButton, clearButton }) => {
  // The feature boundary stays mounted but hidden so callers can render eagerly without stale announcements.
  element.hidden = true;
  // The live-region text is cleared so prior update counts are not re-announced accidentally.
  setText(message, '');
  // The refresh action hides alongside the banner because there is nothing actionable while hidden.
  refreshButton.hidden = true;
  // The clear action hides alongside refresh to avoid showing controls on an empty banner.
  clearButton.hidden = true;
};

// Visible-state rendering keeps count and elapsed copy formatting centralized for consistency.
const applyVisibleState = (
  { element, message, refreshButton, clearButton },
  {
    updateCount,
    elapsedLabel,
    messageText,
    showRefresh = true,
    showClear = true,
  },
) => {
  // The banner becomes visible only after the message has a valid actionable count.
  element.hidden = false;
  const normalizedElapsedLabel =
    typeof elapsedLabel === 'string' && elapsedLabel.length > 0 ? elapsedLabel : '0 seconds';

  // Explicit messages support post-refresh summaries while preserving count-based copy as the default.
  if (typeof messageText === 'string' && messageText.length > 0) {
    setText(message, messageText);
  } else {
    setText(
      message,
      `Newest information: ${updateCount} update${updateCount === 1 ? '' : 's'} in the last ${normalizedElapsedLabel}`,
    );
  }

  // Action visibility is caller-driven so summary states can keep only clear without exposing refresh.
  refreshButton.hidden = showRefresh !== true;
  clearButton.hidden = showClear !== true;
};

/**
 * @param {{ onRefresh?: () => void, onClear?: () => void }} [dependencies]
 * @returns {{ element: HTMLElement, render(state?: { isVisible?: boolean, updateCount?: number, elapsedLabel?: string, messageText?: string, showRefresh?: boolean, showClear?: boolean } | undefined): void, destroy(): void }}
 */
export const createLiveBannerView = ({ onRefresh = () => {}, onClear = () => {} } = {}) => {
  // The root section marks the feature boundary and acts as the live region announced to assistive tech.
  const element = createElement('section', {
    className: 'live-banner',
    attributes: {
      'data-testid': LIVE_BANNER_TEST_IDS.root,
      role: 'status',
      'aria-live': 'polite',
      'aria-atomic': 'true',
    },
  });

  // Message text stays in its own node so render() can update copy without touching button state.
  const message = createElement('p', {
    className: 'live-banner__message',
    attributes: {
      'data-testid': LIVE_BANNER_TEST_IDS.message,
    },
  });

  // The explicit button keeps refresh behavior keyboard-accessible and decoupled from the surrounding region.
  const refreshButton = createElement('button', {
    className: 'live-banner__refresh',
    text: 'Refresh',
    attributes: {
      type: 'button',
      'data-testid': LIVE_BANNER_TEST_IDS.refreshButton,
    },
  });

  // Clear action mirrors refresh styling so both actions read as a related control pair.
  const clearButton = createElement('button', {
    className: 'live-banner__refresh live-banner__refresh--clear',
    text: 'Clear',
    attributes: {
      type: 'button',
      'data-testid': LIVE_BANNER_TEST_IDS.clearButton,
    },
  });

  // Concurrent action clicks are serialized so async shell callbacks cannot race banner state transitions.
  let isActionPending = false;

  // Both controls share pending-disable behavior so refresh and clear always stay in sync.
  const setActionsDisabled = (isDisabled) => {
    refreshButton.disabled = isDisabled;
    clearButton.disabled = isDisabled;
  };

  // Action execution is centralized so refresh/clear callbacks share the same reentrancy and cleanup rules.
  const runAction = async (action) => {
    if (isActionPending) {
      return;
    }

    isActionPending = true;
    setActionsDisabled(true);

    try {
      await action();
    } finally {
      setActionsDisabled(false);
      isActionPending = false;
    }
  };

  // Click handling is injected so the view remains reusable regardless of how refresh is orchestrated later.
  const handleRefreshClick = () => {
    void runAction(onRefresh);
  };

  // Clear is injected so dismissal behavior stays in the controller/app shell layer.
  const handleClearClick = () => {
    void runAction(onClear);
  };

  // Static assembly happens once so later renders only flip visibility and message text.
  const actions = createElement('div', {
    className: 'live-banner__actions',
    children: [clearButton, refreshButton],
  });

  element.append(message, actions);
  // The button listener is attached once and removed during destroy() to keep the view leak-free.
  refreshButton.addEventListener('click', handleRefreshClick);
  clearButton.addEventListener('click', handleClearClick);
  // The banner starts hidden because mounting it should not announce updates before polling has run.
  applyHiddenState({ element, message, refreshButton, clearButton });

  return {
    element,
    // Rendering accepts simple state so polling and diff semantics stay outside the view boundary.
    render(state) {
      // Invalid or hidden state collapses to the shared hidden rendering contract.
      const hasExplicitMessage =
        typeof state?.messageText === 'string' && state.messageText.length > 0;

      if (state?.isVisible !== true || (!hasExplicitMessage && !isPositiveInteger(state?.updateCount))) {
        applyHiddenState({ element, message, refreshButton, clearButton });
        return;
      }

      // Valid visible state renders the formatted count and reveals the refresh action.
      applyVisibleState({ element, message, refreshButton, clearButton }, state);
    },
    // Cleanup removes the button listener so repeated mounts do not accumulate duplicate callbacks.
    destroy() {
      refreshButton.removeEventListener('click', handleRefreshClick);
      clearButton.removeEventListener('click', handleClearClick);
    },
  };
};
