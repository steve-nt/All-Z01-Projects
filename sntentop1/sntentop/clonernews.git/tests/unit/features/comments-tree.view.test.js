// @vitest-environment jsdom

// Comments-tree view tests verify recursive rendering, safe rich-text handling, and visibility semantics.
// Public API under test: createCommentsTreeView and COMMENTS_TREE_TEST_IDS at the standalone feature level.
// Constraints: tests stay within jsdom, mock relative-time formatting for determinism, and avoid post-detail integration concerns.
import { describe, expect, it, vi } from 'vitest';

// Relative-time labels are mocked so assertions stay deterministic regardless of wall-clock time.
vi.mock('../../../src/shared/time-format.js', () => ({
  formatRelativeTime: (value) => `relative-${value}`,
}));

import {
  COMMENTS_TREE_TEST_IDS,
  createCommentsTreeView,
} from '../../../src/features/comments/comments-tree.view.js';

// Data-testid lookup includes the root node itself so feature-section assertions can target the mounted boundary.
const getByTestId = (root, testId) =>
  root?.matches?.(`[data-testid="${testId}"]`)
    ? root
    : root?.querySelector(`[data-testid="${testId}"]`);

// Repeated selectors are centralized so list and metadata assertions stay compact.
const getAllByTestId = (root, testId) => [...root.querySelectorAll(`[data-testid="${testId}"]`)];

// Comment IDs are surfaced as stable data attributes so tests can assert tree structure directly.
const getCommentById = (root, commentId) => root.querySelector(`[data-comment-id="${commentId}"]`);

// Direct-child ordering matters for TC-4, so this helper ignores deeper descendants.
const getDirectChildCommentIds = (listElement) =>
  Array.from(listElement?.children ?? []).map((child) => child.getAttribute('data-comment-id'));

// Scoped lookup keeps repeated metadata assertions tied to one rendered comment node.
const getScopedByTestId = (root, testId) => root?.querySelector(`[data-testid="${testId}"]`);

// A reusable multi-level fixture keeps tree-order and nesting assertions easy to read.
const createMultiLevelComments = () => [
  {
    item: {
      id: 302,
      type: 'comment',
      by: 'ava',
      time: 302,
      text: '<p>Newest top-level comment</p>',
      parent: 99,
    },
    comments: [
      {
        item: {
          id: 3022,
          type: 'comment',
          by: 'cora',
          time: 3022,
          text: '<p>Newest nested reply</p>',
          parent: 302,
        },
        comments: [
          {
            item: {
              id: 30221,
              type: 'comment',
              by: 'drew',
              time: 30221,
              text: '<p>Deep reply</p>',
              parent: 3022,
            },
            comments: [],
          },
        ],
      },
      {
        item: {
          id: 3021,
          type: 'comment',
          by: 'ben',
          time: 3021,
          text: '<p>Older nested reply</p>',
          parent: 302,
        },
        comments: [],
      },
    ],
  },
  {
    item: {
      id: 301,
      type: 'comment',
      by: 'eli',
      time: 301,
      text: '<p>Older top-level comment</p>',
      parent: 99,
    },
    comments: [],
  },
];

describe('comments-tree view', () => {
  it('renders a recursive tree and preserves the provided order at every depth', () => {
    // Recursive rendering should preserve the upstream newest-first order rather than rebuilding the tree.
    const view = createCommentsTreeView();
    const comments = createMultiLevelComments();

    view.render(comments);

    // Successful renders should unhide the feature section and expose the heading and root list.
    const root = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);
    const heading = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.heading);
    const list = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.list);

    expect(root?.hidden).toBe(false);
    expect(heading?.textContent).toBe('Comments');
    expect(getDirectChildCommentIds(list)).toEqual(['302', '301']);

    // Nested lists should keep the provided order at each level and still render deeper descendants.
    const topLevelComment = getCommentById(view.element, 302);
    const nestedList = topLevelComment?.querySelector('.comments-tree__children');

    expect(getDirectChildCommentIds(nestedList)).toEqual(['3022', '3021']);
    expect(getCommentById(view.element, 30221)).not.toBeNull();

    // Relative-time labels should stay visible for both top-level and nested comments.
    const topLevelTime = getScopedByTestId(topLevelComment, COMMENTS_TREE_TEST_IDS.time);
    const nestedComment = getCommentById(view.element, 3022);
    const nestedTime = getScopedByTestId(nestedComment, COMMENTS_TREE_TEST_IDS.time);

    expect(topLevelTime?.textContent).toBe('relative-302');
    expect(nestedTime?.textContent).toBe('relative-3022');

    view.destroy();
  });

  it('sanitizes comment HTML while preserving safe rich-text content', () => {
    // The comment body should preserve safe markup but strip executable or dangerous attributes.
    const view = createCommentsTreeView();

    view.render([
      {
        item: {
          id: 410,
          type: 'comment',
          by: 'safe-html',
          time: 410,
          text: '<p>Hello <strong>world</strong></p><img src="x" onerror="alert(1)"><script>boom()</script>',
          parent: 99,
        },
        comments: [],
      },
    ]);

    // Sanitized bodies should stay visible when safe content remains after purification.
    const commentNode = getCommentById(view.element, 410);
    const body = getScopedByTestId(commentNode, COMMENTS_TREE_TEST_IDS.body);
    const sanitizedImage = body?.querySelector('img');

    expect(body?.hidden).toBe(false);
    expect(body?.querySelector('strong')?.textContent).toBe('world');
    expect(body?.querySelector('script')).toBeNull();
    expect(sanitizedImage?.getAttribute('onerror')).toBeNull();

    view.destroy();
  });

  it('hides comment bodies when sanitization strips all unsafe markup', () => {
    // Fully unsafe comment HTML should not leave an empty visible body container behind.
    const view = createCommentsTreeView();

    view.render([
      {
        item: {
          id: 420,
          type: 'comment',
          by: 'unsafe-only',
          time: 420,
          text: '<script>boom()</script>',
          parent: 99,
        },
        comments: [],
      },
    ]);

    // The comment still renders for metadata purposes, but the stripped body stays hidden.
    const commentNode = getCommentById(view.element, 420);
    const body = getScopedByTestId(commentNode, COMMENTS_TREE_TEST_IDS.body);

    expect(commentNode).not.toBeNull();
    expect(body?.hidden).toBe(true);
    expect(body?.childNodes.length).toBe(0);

    view.destroy();
  });

  it('hides and clears the section when render receives an empty comments array', () => {
    // Re-rendering with no comments should remove stale nodes and hide the feature boundary entirely.
    const view = createCommentsTreeView();

    view.render(createMultiLevelComments());

    // A populated render establishes the precondition that the feature is visible and contains nodes.
    const root = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.root);
    const list = getByTestId(view.element, COMMENTS_TREE_TEST_IDS.list);

    expect(root?.hidden).toBe(false);
    expect(list?.childElementCount).toBe(2);

    // Empty comment input should clear prior content and hide the section for the caller.
    view.render([]);

    expect(root?.hidden).toBe(true);
    expect(list?.childElementCount).toBe(0);

    view.destroy();
  });

  it('adds stable depth and parent markers that match the rendered hierarchy', () => {
    // Depth and hidden parent data attributes remain available as structural hooks for tests and audits.
    const view = createCommentsTreeView();

    view.render(createMultiLevelComments());

    // Top-level comments start at depth zero and point back to the root post as their parent.
    const topLevelComment = getCommentById(view.element, 302);
    const siblingTopLevelComment = getCommentById(view.element, 301);

    expect(topLevelComment?.getAttribute('data-depth')).toBe('0');
    expect(topLevelComment?.style.getPropertyValue('--comment-depth')).toBe('0');
    expect(topLevelComment?.getAttribute('data-parent-id')).toBe('99');
    expect(siblingTopLevelComment?.getAttribute('data-depth')).toBe('0');

    // Nested replies should increment depth and expose their direct comment parent in hidden data only.
    const nestedComment = getCommentById(view.element, 3022);
    const deepComment = getCommentById(view.element, 30221);

    expect(nestedComment?.getAttribute('data-depth')).toBe('1');
    expect(nestedComment?.style.getPropertyValue('--comment-depth')).toBe('1');
    expect(nestedComment?.getAttribute('data-parent-id')).toBe('302');
    expect(deepComment?.getAttribute('data-depth')).toBe('2');
    expect(deepComment?.style.getPropertyValue('--comment-depth')).toBe('2');
    expect(deepComment?.getAttribute('data-parent-id')).toBe('3022');

    // Public test-id hooks should exist for each rendered comment metadata surface.
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.comment)).toHaveLength(5);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.author)).toHaveLength(5);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.time)).toHaveLength(5);
    expect(getAllByTestId(view.element, COMMENTS_TREE_TEST_IDS.body)).toHaveLength(5);

    view.destroy();
  });
});
