/*
 * Purpose: Build the post-detail DOM skeleton and render explicit detail-page states without owning data fetching.
 * Public API: createPostDetailView({ onBack }) -> { element, render(state), destroy() }.
 * Constraints: This file owns DOM creation for the feature, exposes stable data-testid hooks, and keeps rendering simple until richer styling and sanitization land.
 */

import DOMPurify from 'dompurify';

import './post-detail.css';
import { clearElement, createElement, setText } from '../../shared/dom-helpers.js';
import { formatRelativeTime } from '../../shared/time-format.js';
import { createCommentsTreeView } from '../comments/comments-tree.view.js';
import { createPollView, POLL_TEST_IDS } from '../polls/poll.view.js';
import { POST_DETAIL_RENDER_STATES } from './post-detail.render-state.js';

/**
 * @typedef {Object} PostDetailViewDependencies
 * @property {(() => void)=} onBack
 */

// The feed route is the smallest stable back target until Track B adds richer feed-state restoration.
const FEED_ROUTE_HASH = '#/feed/top';

// Stable test IDs give future e2e coverage one authoritative place to target this feature.
export const POST_DETAIL_TEST_IDS = Object.freeze({
  root: 'post-detail-view',
  backButton: 'post-detail-back',
  scrollTopButton: 'post-detail-scroll-top-button',
  title: 'post-detail-title',
  status: 'post-detail-status',
  loadingIndicator: 'post-detail-loading-indicator',
  metadata: 'post-detail-metadata',
  urlSection: 'post-detail-url-section',
  urlLabel: 'post-detail-url-label',
  urlLink: 'post-detail-url-link',
  textSection: 'post-detail-text-section',
  textBody: 'post-detail-text-body',
});

// Visibility is centralized so the render paths stay declarative and avoid repeated hidden toggles.
const setVisibility = (element, isVisible) => {
  element.hidden = !isVisible;
  return element;
};

// Metadata pieces are filtered compactly so missing optional fields do not leave awkward separators.
const compactValues = (values) => values.filter((value) => value !== null);

// Null-safe score formatting keeps the metadata string readable when some API fields are missing.
const formatScore = (score) => (Number.isInteger(score) ? `${score} points` : null);

// Null-safe relative-time formatting keeps this helper pure from view branching details.
const formatTime = (time) => (Number.isInteger(time) ? formatRelativeTime(time) : null);

// Metadata stays as one line for the skeleton so the upcoming CSS step can style it predictably.
const formatMetadata = (viewModel) =>
  compactValues([viewModel.author, formatTime(viewModel.time), formatScore(viewModel.score)]).join(
    ' • ',
  );

// Title entities are decoded off-DOM so user-facing text reads naturally without using innerHTML.
const decodeHtmlEntities = (value) => {
  const parser = new DOMParser();
  const parsedDocument = parser.parseFromString(value, 'text/html');
  return parsedDocument.body.textContent ?? '';
};

// Untrusted link targets are restricted to http/https so the view never writes javascript: URLs into href.
const toSafeHref = (value) => {
  try {
    const parsedUrl = new URL(value);

    if (parsedUrl.protocol !== 'http:' && parsedUrl.protocol !== 'https:') {
      return null;
    }

    return parsedUrl.toString();
  } catch {
    return null;
  }
};

// Only the API text field is inserted as sanitized HTML, and the sanitized fragment is appended directly.
const renderSanitizedText = (element, value) => {
  clearElement(element);
  element.append(DOMPurify.sanitize(value, { RETURN_DOM_FRAGMENT: true }));
  return element;
};

/**
 * @param {PostDetailViewDependencies=} dependencies
 * @returns {{ element: HTMLElement, render(state: { status: string, data?: object, message?: string }): void, destroy(): void }}
 */
export const createPostDetailView = ({ onBack } = {}) => {
  // The root section defines the feature boundary and gives tests one stable mount target.
  const element = createElement('section', {
    className: 'post-detail',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.root,
      'aria-labelledby': POST_DETAIL_TEST_IDS.title,
    },
  });

  // Custom back behavior can restore richer feed state later, while the fallback always returns to the feed route.
  const navigateBack = () => {
    if (typeof onBack === 'function') {
      onBack();
      return;
    }

    window.location.hash = FEED_ROUTE_HASH;
  };

  // Scroll-to-top behavior mirrors feed ergonomics so long detail pages remain easy to navigate.
  const scrollViewportToTop = () => {
    if (typeof window === 'undefined') {
      return;
    }

    const isJsDomEnvironment =
      typeof window.navigator !== 'undefined' &&
      typeof window.navigator.userAgent === 'string' &&
      window.navigator.userAgent.toLowerCase().includes('jsdom');

    if (!isJsDomEnvironment && typeof window.scrollTo === 'function') {
      try {
        window.scrollTo({ top: 0, left: 0, behavior: 'auto' });
      } catch {
        // Browsers that reject object-form scroll options still receive the fallback below.
      }
    }

    if (document.documentElement) {
      document.documentElement.scrollTop = 0;
    }

    if (document.body) {
      document.body.scrollTop = 0;
    }
  };

  // A single delegated click handler keeps interaction cleanup simple and aligns with the repo's DOM rules.
  const handleRootClick = (event) => {
    const clickTarget = event.target;

    if (!(clickTarget instanceof Element)) {
      return;
    }

    if (clickTarget.closest(`[data-testid="${POST_DETAIL_TEST_IDS.scrollTopButton}"]`)) {
      scrollViewportToTop();
      return;
    }

    // The back action is exposed by dependency injection so routing stays outside the view.
    if (clickTarget.closest(`[data-testid="${POST_DETAIL_TEST_IDS.backButton}"]`)) {
      navigateBack();
    }
  };

  // The view header groups navigation and primary identity content into one predictable block.
  const header = createElement('header', {
    className: 'post-detail__header',
  });

  // The back button exists in every state so users always have a consistent escape button.
  const backButton = createElement('button', {
    className: 'post-detail__back-button',
    text: 'Back',
    attributes: {
      type: 'button',
      'data-testid': POST_DETAIL_TEST_IDS.backButton,
    },
  });

  // The title node is stable across states so loading, errors, and success all write into one place.
  const title = createElement('h1', {
    className: 'post-detail__title',
    text: 'Loading post…',
    attributes: {
      id: POST_DETAIL_TEST_IDS.title,
      'data-testid': POST_DETAIL_TEST_IDS.title,
    },
  });

  // Status text handles loading, not-found, and generic error messages without replacing the overall skeleton.
  const status = createElement('p', {
    className: 'post-detail__status',
    text: 'Fetching post details…',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.status,
    },
  });

  // Loading dots add clear activity feedback while the detail request is in-flight.
  const loadingIndicator = createElement('div', {
    className: 'post-detail__loading-indicator',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.loadingIndicator,
      'aria-hidden': 'true',
    },
    children: [
      createElement('span', {
        className: 'post-detail__loading-dot post-detail__loading-dot--first',
      }),
      createElement('span', {
        className: 'post-detail__loading-dot post-detail__loading-dot--second',
      }),
      createElement('span', {
        className: 'post-detail__loading-dot post-detail__loading-dot--third',
      }),
    ],
  });

  // The status row keeps loading/done indicators aligned with status copy on one horizontal line.
  const statusRow = createElement('div', {
    className: 'post-detail__status-row',
  });

  // Metadata stays mounted in all states so the final CSS can reserve one consistent layout slot.
  const metadata = createElement('p', {
    className: 'post-detail__metadata',
    text: '',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.metadata,
    },
  });

  // The URL wrapper is optional because many job posts have no outbound link.
  const urlSection = createElement('section', {
    className: 'post-detail__url',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.urlSection,
    },
  });

  // The URL label clarifies that the anchor points to the source article outside the app.
  const urlLabel = createElement('span', {
    className: 'post-detail__url-label',
    text: 'Original Link:',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.urlLabel,
    },
  });

  // The outbound link uses explicit safe attributes because the target URL comes from API data.
  const urlLink = createElement('a', {
    className: 'post-detail__url-link',
    text: '',
    attributes: {
      href: '#',
      rel: 'noreferrer noopener',
      target: '_blank',
      'data-testid': POST_DETAIL_TEST_IDS.urlLink,
    },
  });

  // The text wrapper is always part of the DOM tree even when hidden so tests can target it stably.
  const textSection = createElement('section', {
    className: 'post-detail__text',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.textSection,
    },
  });

  // The body container receives sanitized HTML because Hacker News sends markup in the text field.
  const textBody = createElement('div', {
    className: 'post-detail__text-body',
    text: '',
    attributes: {
      'data-testid': POST_DETAIL_TEST_IDS.textBody,
    },
  });

  // The comments subview is created once so state transitions can reuse one mounted feature boundary.
  const commentsView = createCommentsTreeView();
  // The poll subview is created once so poll route changes can reuse one mounted feature boundary.
  const pollView = createPollView();

  // Floating action keeps long detail pages consistent with feed-page top navigation affordances.
  const scrollTopButton = createElement('button', {
    className: 'post-detail__scroll-top-button',
    text: 'Top',
    attributes: {
      type: 'button',
      'data-testid': POST_DETAIL_TEST_IDS.scrollTopButton,
      'aria-label': 'Back to top',
      title: 'Back to top',
    },
  });

  // Static assembly happens once so later renders only update text and visibility instead of rebuilding nodes.
  urlSection.append(urlLabel, urlLink);
  textSection.append(textBody);
  statusRow.append(status, loadingIndicator);
  header.append(backButton, title, statusRow, metadata);
  element.append(header, urlSection, textSection, pollView.element, commentsView.element, scrollTopButton);
  element.addEventListener('click', handleRootClick);

  // The initial skeleton hides optional content until a success state provides real values.
  setVisibility(urlSection, false);
  setVisibility(textSection, false);

  // Small reset logic keeps each render branch focused on the fields that actually differ.
  const resetContent = () => {
    element.setAttribute('aria-labelledby', POST_DETAIL_TEST_IDS.title);
    element.removeAttribute('data-post-type');
    element.classList.remove('post-detail--loading');
    loadingIndicator.classList.remove('post-detail__loading-indicator--done');
    setText(metadata, '');
    setText(urlLink, '');
    urlLink.setAttribute('href', '#');
    clearElement(textBody);
    commentsView.hide();
    pollView.hide();
    setVisibility(title, true);
    setVisibility(metadata, true);
    setVisibility(statusRow, true);
    setVisibility(status, true);
    setVisibility(loadingIndicator, false);
    setVisibility(urlSection, false);
    setVisibility(textSection, false);
  };

  // Loading keeps the skeleton mounted while making the in-flight state obvious to users and tests.
  const renderLoading = () => {
    resetContent();
    element.classList.add('post-detail--loading');
    setText(title, 'Loading post…');
    setText(status, 'Fetching post details…');
    loadingIndicator.classList.remove('post-detail__loading-indicator--done');
    setVisibility(status, true);
    setVisibility(loadingIndicator, true);
  };

  // Not-found shares the skeleton but uses empty-state language instead of implying a network failure.
  const renderNotFound = (message) => {
    resetContent();
    setText(title, 'Post unavailable');
    setText(status, message ?? 'Post not found.');
    setVisibility(status, true);
  };

  // Generic errors remain distinct from not-found so future retries can be presented clearly.
  const renderError = (message) => {
    resetContent();
    setText(title, 'Unable to load post');
    setText(status, message ?? 'Failed to load post.');
    setVisibility(status, true);
  };

  // Shared URL rendering keeps the story and job branches explicit without duplicating sink logic.
  const renderOptionalUrl = (viewModel) => {
    const safeHref = viewModel.hasUrl && viewModel.url !== null ? toSafeHref(viewModel.url) : null;

    if (safeHref === null) {
      return;
    }

    setText(urlLink, viewModel.url);
    urlLink.setAttribute('href', safeHref);
    setVisibility(urlSection, true);
  };

  // Shared text rendering keeps DOMPurify usage centralized while still allowing type-specific branches.
  const renderOptionalText = (viewModel) => {
    if (!viewModel.hasText) {
      return;
    }

    renderSanitizedText(textBody, viewModel.text);
    setVisibility(textSection, textBody.childNodes.length > 0);
  };

  // Story posts may show an outbound source link and may also include supplemental story text.
  const renderStoryContent = (viewModel, comments) => {
    renderOptionalUrl(viewModel);
    renderOptionalText(viewModel);
    commentsView.render(comments);
  };

  // Job posts commonly have no link and instead lean on a larger text body, but both fields remain optional.
  const renderJobContent = (viewModel) => {
    renderOptionalUrl(viewModel);
    renderOptionalText(viewModel);
  };

  // Poll posts hand off to the dedicated poll feature so vote bars and option text stay encapsulated there.
  const renderPollContent = (viewModel) => {
    element.setAttribute('aria-labelledby', POLL_TEST_IDS.title);
    setVisibility(title, false);
    setVisibility(metadata, false);
    pollView.render(viewModel);
  };

  // Success fills the existing skeleton in place so the final UI can stay predictable for tests and CSS.
  const renderSuccess = (data) => {
    const viewModel = data?.viewModel;

    if (!viewModel) {
      renderError('Post details were incomplete.');
      return;
    }

    resetContent();
    element.dataset.postType = viewModel.type;
    setText(title, decodeHtmlEntities(viewModel.title));
    setText(status, 'Done loading post details.');
    setVisibility(status, true);
    loadingIndicator.classList.add('post-detail__loading-indicator--done');
    setVisibility(loadingIndicator, true);
    setText(metadata, formatMetadata(viewModel));

    if (viewModel.type === 'story') {
      renderStoryContent(viewModel, data?.comments);
      return;
    }

    if (viewModel.type === 'job') {
      renderJobContent(viewModel);
      return;
    }

    if (viewModel.type === 'poll') {
      renderPollContent(viewModel);
      return;
    }

    renderError('Unsupported post detail view model type.');
  };

  return {
    element,
    // The public render API accepts the explicit render-state contract defined in the previous step.
    render(state) {
      switch (state?.status) {
        case POST_DETAIL_RENDER_STATES.loading:
          renderLoading();
          return;
        case POST_DETAIL_RENDER_STATES.success:
          renderSuccess(state.data);
          return;
        case POST_DETAIL_RENDER_STATES.notFound:
          renderNotFound(state.message);
          return;
        default:
          renderError(state?.message);
      }
    },
    // Cleanup removes the delegated listener so future unmounting does not leak handlers.
    destroy() {
      element.removeEventListener('click', handleRootClick);
      commentsView.destroy();
      pollView.destroy();
    },
  };
};
