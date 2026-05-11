// @vitest-environment jsdom

// Post-detail view tests verify the DOM skeleton, state rendering, and sanitizer-backed text insertion.
// Public API under test: createPostDetailView, POST_DETAIL_TEST_IDS, and render-state integration at the DOM level.
// Constraints: tests stay within jsdom, use stable data-testid hooks, and avoid coupling to future route wiring.
import { describe, expect, it, vi } from 'vitest';

import { COMMENTS_TREE_TEST_IDS } from '../../../src/features/comments/comments-tree.view.js';
import { POLL_TEST_IDS } from '../../../src/features/polls/poll.view.js';
import {
  createErrorPostDetailState,
  createLoadingPostDetailState,
  createNotFoundPostDetailState,
  createSuccessPostDetailState,
} from '../../../src/features/post-detail/post-detail.render-state.js';
import {
  createPostDetailView,
  POST_DETAIL_TEST_IDS,
} from '../../../src/features/post-detail/post-detail.view.js';

// Data-testid lookup keeps assertions aligned with the feature's public DOM contract.
const getByTestId = (root, testId) => root.querySelector(`[data-testid="${testId}"]`);

// Repeated selectors are centralized so nested comment assertions stay concise and explicit.
const getAllByTestId = (root, testId) => [...root.querySelectorAll(`[data-testid="${testId}"]`)];

// Story-comment fixtures stay small but still exercise nested rendering and hidden parent markers.
const createNestedComments = () => [
  {
    item: {
      id: 501,
      type: 'comment',
      by: 'root-commenter',
      time: 501,
      text: '<p>Root comment</p>',
      parent: 5,
    },
    comments: [
      {
        item: {
          id: 502,
          type: 'comment',
          by: 'reply-commenter',
          time: 502,
          text: '<p>Reply comment</p>',
          parent: 501,
        },
        comments: [],
      },
    ],
  },
];

describe('post-detail view', () => {
  it('renders the loading skeleton with optional sections hidden', () => {
    // Loading should mount the base structure immediately without exposing stale optional content.
    const view = createPostDetailView();

    view.render(createLoadingPostDetailState());

    const loadingIndicator = getByTestId(view.element, POST_DETAIL_TEST_IDS.loadingIndicator);

    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.title)?.textContent).toBe(
      'Loading post…',
    );
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.status)?.textContent).toBe(
      'Fetching post details…',
    );
    expect(view.element.classList.contains('post-detail--loading')).toBe(true);
    expect(loadingIndicator?.hidden).toBe(false);
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection)?.hidden).toBe(true);
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.textSection)?.hidden).toBe(true);

    view.destroy();
  });

  it('calls onBack when the back button is clicked', () => {
    // The back callback is the view's only interaction responsibility before route wiring is added.
    const onBack = vi.fn();
    const view = createPostDetailView({ onBack });
    const backButton = getByTestId(view.element, POST_DETAIL_TEST_IDS.backButton);

    backButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(onBack).toHaveBeenCalledTimes(1);

    view.destroy();
  });

  it('falls back to the feed route when no custom back callback is provided', () => {
    // Returning to #/feed/top keeps top-feed routing explicit and aligned with the tab route contract.
    window.location.hash = '#/item/42';
    const view = createPostDetailView();
    const backButton = getByTestId(view.element, POST_DETAIL_TEST_IDS.backButton);

    backButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(window.location.hash).toBe('#/feed/top');

    view.destroy();
  });

  it('renders a floating top button that resets scroll positions', () => {
    // Detail views should expose the same top-navigation affordance as feed pages.
    const view = createPostDetailView();
    const topButton = getByTestId(view.element, POST_DETAIL_TEST_IDS.scrollTopButton);

    document.documentElement.scrollTop = 520;
    document.body.scrollTop = 520;

    topButton?.dispatchEvent(new MouseEvent('click', { bubbles: true }));

    expect(topButton).not.toBeNull();
    expect(document.documentElement.scrollTop).toBe(0);
    expect(document.body.scrollTop).toBe(0);

    view.destroy();
  });

  it('renders a story success state with a safe link and sanitized HTML text', () => {
    // Story rendering should keep plain fields on text sinks while sanitizing only the API text field.
    const view = createPostDetailView();
    const state = createSuccessPostDetailState({
      comments: createNestedComments(),
      viewModel: {
        id: 5,
        type: 'story',
        title: 'Fish &amp; Chips',
        author: 'alice',
        time: null,
        score: 7,
        url: 'https://example.com/story',
        hasUrl: true,
        text: '<p>Hello <strong>world</strong></p><img src="x" onerror="alert(1)"><script>boom()</script>',
        hasText: true,
      },
    });

    view.render(state);

    const title = getByTestId(view.element, POST_DETAIL_TEST_IDS.title);
    const metadata = getByTestId(view.element, POST_DETAIL_TEST_IDS.metadata);
    const urlSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection);
    const urlLink = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlLink);
    const status = getByTestId(view.element, POST_DETAIL_TEST_IDS.status);
    const loadingIndicator = getByTestId(view.element, POST_DETAIL_TEST_IDS.loadingIndicator);
    const textSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.textSection);
    const textBody = getByTestId(view.element, POST_DETAIL_TEST_IDS.textBody);
    const sanitizedImage = textBody?.querySelector('img');
    const commentsRoot = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);
    const renderedComments = getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment);

    expect(title?.textContent).toBe('Fish & Chips');
    expect(metadata?.textContent).toBe('alice • 7 points');
    expect(view.element.dataset.postType).toBe('story');
    expect(urlSection?.hidden).toBe(false);
    expect(urlLink?.textContent).toBe('https://example.com/story');
    expect(urlLink?.getAttribute('href')).toBe('https://example.com/story');
    expect(status?.textContent).toBe('Done loading post details.');
    expect(loadingIndicator?.classList.contains('post-detail__loading-indicator--done')).toBe(true);
    expect(textSection?.hidden).toBe(false);
    expect(textBody?.querySelector('strong')?.textContent).toBe('world');
    expect(textBody?.querySelector('script')).toBeNull();
    expect(sanitizedImage?.getAttribute('onerror')).toBeNull();
    expect(commentsRoot?.hidden).toBe(false);
    expect(renderedComments).toHaveLength(2);
    expect(renderedComments[0]?.getAttribute('data-comment-id')).toBe('501');
    expect(renderedComments[1]?.getAttribute('data-comment-id')).toBe('502');
    expect(renderedComments[1]?.getAttribute('data-parent-id')).toBe('501');
    expect(renderedComments[0]?.getAttribute('data-depth')).toBe('0');
    expect(renderedComments[1]?.getAttribute('data-depth')).toBe('1');

    view.destroy();
  });

  it('keeps comments hidden for job success states', () => {
    // Job posts should not leave a comments tree visible even if a previous story render populated one.
    const view = createPostDetailView();
    const storyState = createSuccessPostDetailState({
      comments: createNestedComments(),
      viewModel: {
        id: 5,
        type: 'story',
        title: 'Story before job',
        author: 'alice',
        time: null,
        score: 7,
        url: null,
        hasUrl: false,
        text: '<p>Story body</p>',
        hasText: true,
      },
    });
    const jobState = createSuccessPostDetailState({
      viewModel: {
        id: 6,
        type: 'job',
        title: 'Remote role',
        author: 'team',
        time: null,
        score: null,
        url: null,
        hasUrl: false,
        text: '<p>Now hiring</p>',
        hasText: true,
      },
    });

    view.render(storyState);
    view.render(jobState);

    const metadata = getByTestId(view.element, POST_DETAIL_TEST_IDS.metadata);
    const urlSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection);
    const textSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.textSection);
    const textBody = getByTestId(view.element, POST_DETAIL_TEST_IDS.textBody);
    const commentsRoot = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);

    expect(view.element.dataset.postType).toBe('job');
    expect(metadata?.textContent).toBe('team');
    expect(urlSection?.hidden).toBe(true);
    expect(textSection?.hidden).toBe(false);
    expect(textBody?.textContent).toContain('Now hiring');
    expect(commentsRoot?.hidden).toBe(true);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(0);

    view.destroy();
  });

  it('renders poll success states through the poll subview and hides story-only surfaces', () => {
    // Poll posts should hand off to the poll feature instead of reusing the story/job content blocks.
    const view = createPostDetailView();
    const state = createSuccessPostDetailState({
      comments: [],
      viewModel: {
        id: 70,
        type: 'poll',
        title: 'Best JavaScript runtime?',
        author: 'dang',
        time: null,
        text: '<p>Choose a <strong>winner</strong>.</p><script>boom()</script>',
        hasText: true,
        totalVotes: 20,
        optionCount: 3,
        leadingScore: 12,
        hasOptions: true,
        options: [
          {
            id: 701,
            position: 1,
            text: '<span>Node.js</span>',
            hasText: true,
            score: 12,
            voteRatio: 0.6,
            barWidthPercent: 60,
          },
          {
            id: 702,
            position: 2,
            text: '',
            hasText: false,
            score: 3,
            voteRatio: 0.15,
            barWidthPercent: 15,
          },
          {
            id: 703,
            position: 3,
            text: '<em>Deno</em>',
            hasText: true,
            score: 5,
            voteRatio: 0.25,
            barWidthPercent: 25,
          },
        ],
      },
    });

    view.render(state);

    const title = getByTestId(view.element, POST_DETAIL_TEST_IDS.title);
    const metadata = getByTestId(view.element, POST_DETAIL_TEST_IDS.metadata);
    const urlSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection);
    const textSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.textSection);
    const commentsRoot = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);
    const pollRoot = getByTestId(view.element, POLL_TEST_IDS.root);
    const pollTitle = getByTestId(view.element, POLL_TEST_IDS.title);
    const pollTotalVotes = getByTestId(view.element, POLL_TEST_IDS.totalVotes);
    const pollBars = getAllByTestId(view.element, POLL_TEST_IDS.optionBar);

    expect(view.element.dataset.postType).toBe('poll');
    expect(title?.hidden).toBe(true);
    expect(metadata?.hidden).toBe(true);
    expect(urlSection?.hidden).toBe(true);
    expect(textSection?.hidden).toBe(true);
    expect(commentsRoot?.hidden).toBe(true);
    expect(pollRoot?.hidden).toBe(false);
    expect(pollTitle?.textContent).toBe('Best JavaScript runtime?');
    expect(pollTotalVotes?.textContent).toBe('20 votes');
    expect(pollBars).toHaveLength(3);
    expect(pollBars[0]?.style.width).toBe('60%');
    expect(pollBars[1]?.style.width).toBe('15%');
    expect(pollBars[2]?.style.width).toBe('25%');

    view.destroy();
  });

  it('clears previously rendered poll content when a later story success state renders', () => {
    // Poll subview content should not linger once the detail workflow switches back to a story.
    const view = createPostDetailView();
    const pollState = createSuccessPostDetailState({
      comments: [],
      viewModel: {
        id: 70,
        type: 'poll',
        title: 'Best JavaScript runtime?',
        author: 'dang',
        time: null,
        text: '<p>Choose one</p>',
        hasText: true,
        totalVotes: 20,
        optionCount: 2,
        leadingScore: 12,
        hasOptions: true,
        options: [
          {
            id: 701,
            position: 1,
            text: 'Node.js',
            hasText: true,
            score: 12,
            voteRatio: 0.6,
            barWidthPercent: 60,
          },
          {
            id: 702,
            position: 2,
            text: 'Deno',
            hasText: true,
            score: 8,
            voteRatio: 0.4,
            barWidthPercent: 40,
          },
        ],
      },
    });
    const storyState = createSuccessPostDetailState({
      comments: createNestedComments(),
      viewModel: {
        id: 71,
        type: 'story',
        title: 'Story after poll',
        author: 'alice',
        time: null,
        score: 9,
        url: null,
        hasUrl: false,
        text: '<p>Story body</p>',
        hasText: true,
      },
    });

    view.render(pollState);
    view.render(storyState);

    const title = getByTestId(view.element, POST_DETAIL_TEST_IDS.title);
    const metadata = getByTestId(view.element, POST_DETAIL_TEST_IDS.metadata);
    const commentsRoot = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);
    const pollRoot = getByTestId(view.element, POLL_TEST_IDS.root);

    expect(view.element.dataset.postType).toBe('story');
    expect(title?.hidden).toBe(false);
    expect(title?.textContent).toBe('Story after poll');
    expect(metadata?.hidden).toBe(false);
    expect(metadata?.textContent).toBe('alice • 9 points');
    expect(commentsRoot?.hidden).toBe(false);
    expect(pollRoot?.hidden).toBe(true);

    view.destroy();
  });

  it('keeps the url section hidden when the provided url is unsafe', () => {
    // Unsafe protocols must never reach href even when the controller marks the field as present.
    const view = createPostDetailView();
    const state = createSuccessPostDetailState({
      viewModel: {
        id: 7,
        type: 'story',
        title: 'Unsafe url story',
        author: 'alice',
        time: null,
        score: null,
        url: 'javascript:alert(1)',
        hasUrl: true,
        text: '',
        hasText: false,
      },
    });

    view.render(state);

    const urlSection = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection);
    const urlLink = getByTestId(view.element, POST_DETAIL_TEST_IDS.urlLink);

    expect(urlSection?.hidden).toBe(true);
    expect(urlLink?.getAttribute('href')).toBe('#');

    view.destroy();
  });

  it('clears previously rendered comments for loading, not-found, and error states', () => {
    // Non-success states should never leave stale comments visible after a prior story render.
    const view = createPostDetailView();
    const successState = createSuccessPostDetailState({
      comments: createNestedComments(),
      viewModel: {
        id: 5,
        type: 'story',
        title: 'Story with comments',
        author: 'alice',
        time: null,
        score: 7,
        url: null,
        hasUrl: false,
        text: '<p>Story body</p>',
        hasText: true,
      },
    });
    const commentsRoot = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);

    view.render(successState);

    expect(commentsRoot?.hidden).toBe(false);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(2);

    view.render(createLoadingPostDetailState());

    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.title)?.textContent).toBe(
      'Loading post…',
    );
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.status)?.textContent).toBe(
      'Fetching post details…',
    );
    expect(commentsRoot?.hidden).toBe(true);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(0);

    view.render(successState);
    view.render(createNotFoundPostDetailState('Missing item'));

    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.title)?.textContent).toBe(
      'Post unavailable',
    );
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.status)?.textContent).toBe(
      'Missing item',
    );
    expect(commentsRoot?.hidden).toBe(true);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(0);

    view.render(successState);
    view.render(createErrorPostDetailState('Network timeout'));

    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.title)?.textContent).toBe(
      'Unable to load post',
    );
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.status)?.textContent).toBe(
      'Network timeout',
    );
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.urlSection)?.hidden).toBe(true);
    expect(getByTestId(view.element, POST_DETAIL_TEST_IDS.textSection)?.hidden).toBe(true);
    expect(commentsRoot?.hidden).toBe(true);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(0);

    view.destroy();
  });
});
