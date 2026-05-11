/*
 * Purpose: Render a Hacker News poll as a standalone child view that can be mounted inside post-detail.
 * Public API: createPollView() -> { element, render(viewModel), hide(), destroy() } plus POLL_TEST_IDS.
 * Constraints: This module owns poll DOM creation, sanitizes only API HTML fields with DOMPurify, and uses safe DOM sinks for all plain-text surfaces.
 */

import DOMPurify from 'dompurify';

import './poll.css';
import { clearElement, createElement, setText } from '../../shared/dom-helpers.js';
import { formatRelativeTime } from '../../shared/time-format.js';

// Stable test IDs give the poll feature one public DOM contract for unit and e2e targeting.
export const POLL_TEST_IDS = Object.freeze({
  root: 'poll-view',
  title: 'poll-title',
  metadata: 'poll-metadata',
  textSection: 'poll-text-section',
  textBody: 'poll-text-body',
  totalVotes: 'poll-total-votes',
  optionsList: 'poll-options-list',
  option: 'poll-option',
  optionLabel: 'poll-option-label',
  optionVotes: 'poll-option-votes',
  optionBar: 'poll-option-bar',
});

// Visibility toggling stays centralized so render and hide behavior share one rule.
const setVisibility = (element, isVisible) => {
  element.hidden = !isVisible;
  return element;
};

// Poll metadata falls back cleanly so partial API payloads still render without blank labels.
const toDisplayString = (value, fallback) =>
  typeof value === 'string' && value.trim().length > 0 ? value : fallback;

// Time formatting stays null-safe so invalid timestamps do not leak into the shared formatter.
const toDisplayTime = (value) =>
  Number.isInteger(value) ? formatRelativeTime(value) : 'Unknown time';

// Vote totals stay numeric and non-negative so ratio bars never render with misleading labels.
const toVoteCount = (value) => (Number.isInteger(value) && value >= 0 ? value : 0);

// Bar widths are clamped defensively so malformed input can never escape the visual bounds.
const toBarWidthPercent = (value) => {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return 0;
  }

  return Math.max(0, Math.min(100, value));
};

// Singular and plural vote labels are normalized in one place so all rows stay consistent.
const formatVoteLabel = (value) => {
  const voteCount = toVoteCount(value);
  const voteWord = voteCount === 1 ? 'vote' : 'votes';

  return `${voteCount} ${voteWord}`;
};

// Metadata stays compact so the poll card matches the rest of the detail-page information density.
const formatMetadata = (viewModel) =>
  `${toDisplayString(viewModel?.author, 'Unknown author')} • ${toDisplayTime(viewModel?.time)}`;

// Title entities are decoded off-DOM so user-facing text stays readable without using innerHTML.
const decodeHtmlEntities = (value) => {
  const parser = new DOMParser();
  const parsedDocument = parser.parseFromString(String(value ?? ''), 'text/html');

  return parsedDocument.body.textContent ?? '';
};

// Only API rich-text fields are sanitized, and the sanitized fragment is appended directly afterward.
const renderSanitizedHtml = (element, value) => {
  clearElement(element);

  if (typeof value !== 'string' || value.trim().length === 0) {
    return false;
  }

  element.append(DOMPurify.sanitize(value, { RETURN_DOM_FRAGMENT: true }));

  return element.childNodes.length > 0;
};

/**
 * @param {{ id?: number, position?: number, text?: string, hasText?: boolean, score?: number, barWidthPercent?: number }} optionViewModel
 * @returns {HTMLElement}
 */
const createOptionNode = (optionViewModel) => {
  // The controller owns ordering, so each row only reflects the supplied position and values as-is.
  const optionId =
    Number.isInteger(optionViewModel?.id) && optionViewModel.id > 0
      ? optionViewModel.id
      : undefined;
  // Sequential positions remain visible even when an option's text is empty after sanitization.
  const position =
    Number.isInteger(optionViewModel?.position) && optionViewModel.position > 0
      ? optionViewModel.position
      : 1;
  // Safe bar widths prevent malformed ratios from exceeding the visual track.
  const barWidthPercent = toBarWidthPercent(optionViewModel?.barWidthPercent);
  // Vote labels reuse one formatter so totals stay human-readable and deterministic.
  const voteLabel = formatVoteLabel(optionViewModel?.score);

  // Each option row acts as the stable root for one choice and its bar graph.
  const option = createElement('li', {
    className: 'poll__option',
    attributes: {
      'data-testid': POLL_TEST_IDS.option,
      'data-option-id': optionId,
      'data-option-position': position,
    },
  });

  // The header groups the label and vote total so the row reads predictably before the bar track.
  const optionHeader = createElement('div', {
    className: 'poll__option-header',
  });

  // Option labels are the second rich-text surface in this feature, so they are sanitized before insertion.
  const optionLabel = createElement('div', {
    className: 'poll__option-label',
    attributes: {
      'data-testid': POLL_TEST_IDS.optionLabel,
    },
  });

  // Vote counts remain plain text because they are numeric controller output, not API HTML.
  const optionVotes = createElement('span', {
    className: 'poll__option-votes',
    text: voteLabel,
    attributes: {
      'data-testid': POLL_TEST_IDS.optionVotes,
    },
  });

  // The bar track provides a consistent visual rail even when the fill width is zero.
  const optionTrack = createElement('div', {
    className: 'poll__option-track',
  });

  // The fill width is driven entirely by controller-provided ratio math so the view stays dumb.
  const optionBar = createElement('div', {
    className: 'poll__option-bar',
    attributes: {
      'data-testid': POLL_TEST_IDS.optionBar,
    },
  });

  // Sanitized option content determines whether the label shows rich text or a stable fallback label.
  const hasVisibleLabel =
    optionViewModel?.hasText === true && renderSanitizedHtml(optionLabel, optionViewModel?.text);

  // Empty or fully stripped option labels fall back to a deterministic sequence name for accessibility.
  if (!hasVisibleLabel) {
    setText(optionLabel, `Option ${position}`);
  }

  // Inline width assignment keeps the bar graph contract explicit and easy for tests to inspect.
  optionBar.style.width = `${barWidthPercent}%`;

  // Static assembly keeps the rendered row order identical to the provided option order.
  optionHeader.append(optionLabel, optionVotes);
  optionTrack.append(optionBar);
  option.append(optionHeader, optionTrack);

  return option;
};

/**
 * @returns {{ element: HTMLElement, render(viewModel: object | undefined): void, hide(): void, destroy(): void }}
 */
export const createPollView = () => {
  // The root section defines the feature boundary and stays hidden until a poll is rendered.
  const element = createElement('section', {
    className: 'poll',
    attributes: {
      'data-testid': POLL_TEST_IDS.root,
      'aria-labelledby': POLL_TEST_IDS.title,
    },
  });

  // The header groups the poll identity and metadata into one stable card-like block.
  const header = createElement('header', {
    className: 'poll__header',
  });

  // The title remains a plain-text sink because poll titles are text fields that may contain entities only.
  const title = createElement('h2', {
    className: 'poll__title',
    text: '',
    attributes: {
      id: POLL_TEST_IDS.title,
      'data-testid': POLL_TEST_IDS.title,
    },
  });

  // Metadata exposes the author and relative time in one predictable surface for tests and readers.
  const metadata = createElement('p', {
    className: 'poll__metadata',
    text: '',
    attributes: {
      'data-testid': POLL_TEST_IDS.metadata,
    },
  });

  // Total votes are called out separately so users can read the aggregate without scanning the full list.
  const totalVotes = createElement('p', {
    className: 'poll__total-votes',
    text: '',
    attributes: {
      'data-testid': POLL_TEST_IDS.totalVotes,
    },
  });

  // The poll text section holds the sanitized rich-text question or description when one exists.
  const textSection = createElement('section', {
    className: 'poll__text',
    attributes: {
      'data-testid': POLL_TEST_IDS.textSection,
    },
  });

  // The text body is the feature's primary rich-text sink and is always cleared before rerendering.
  const textBody = createElement('div', {
    className: 'poll__text-body',
    attributes: {
      'data-testid': POLL_TEST_IDS.textBody,
    },
  });

  // The ordered list keeps option order identical to the controller-provided sequence.
  const optionsList = createElement('ol', {
    className: 'poll__options',
    attributes: {
      'data-testid': POLL_TEST_IDS.optionsList,
    },
  });

  // A root-level delegated handler preserves the same cleanup pattern used by other child views.
  const handleRootClick = (event) => {
    if (!(event.target instanceof Element)) {
      return;
    }
  };

  // Static assembly happens once so render only needs to update content and visibility.
  textSection.append(textBody);
  header.append(title, metadata, totalVotes);
  element.append(header, textSection, optionsList);
  element.addEventListener('click', handleRootClick);

  // The view starts hidden because post-detail can mount it before the route resolves to a poll.
  setVisibility(element, false);
  setVisibility(textSection, false);

  // Reset logic clears stale rich text and option rows before any new render path runs.
  const reset = () => {
    setText(title, '');
    setText(metadata, '');
    setText(totalVotes, '');
    clearElement(textBody);
    clearElement(optionsList);
    setVisibility(textSection, false);
  };

  return {
    element,
    // Rendering consumes a controller-normalized poll view model and reflects it directly into the DOM.
    render(viewModel) {
      reset();

      // Missing view models hide the entire feature so callers can fail closed without extra guards.
      if (!viewModel || typeof viewModel !== 'object') {
        setVisibility(element, false);
        return;
      }

      // Titles are decoded through a text parser so entities display naturally while remaining safe.
      setText(title, decodeHtmlEntities(viewModel.title));
      // Metadata and totals stay on plain-text sinks because they are derived values, not API HTML.
      setText(metadata, formatMetadata(viewModel));
      setText(totalVotes, formatVoteLabel(viewModel.totalVotes));

      // Poll text only becomes visible when sanitized content survives purification.
      const hasVisibleText =
        viewModel.hasText === true && renderSanitizedHtml(textBody, viewModel.text);

      setVisibility(textSection, hasVisibleText);

      // Options are appended through one fragment so large polls do not trigger repeated reflow.
      const optionRows = Array.isArray(viewModel.options) ? viewModel.options : [];
      const listFragment = document.createDocumentFragment();

      for (const optionViewModel of optionRows) {
        listFragment.append(createOptionNode(optionViewModel));
      }

      // One append keeps the list update compact and preserves the supplied controller order.
      optionsList.append(listFragment);
      // The feature becomes visible only after the DOM has been fully populated.
      setVisibility(element, true);
    },
    // Hiding clears the subview so stale poll data never survives route or type changes.
    hide() {
      reset();
      setVisibility(element, false);
    },
    // Cleanup removes the delegated listener so unmounting stays leak-free.
    destroy() {
      element.removeEventListener('click', handleRootClick);
    },
  };
};

export default createPollView;
