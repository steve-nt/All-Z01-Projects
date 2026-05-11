/*
 * Purpose: Render a recursive Hacker News comments tree as a standalone feature view.
 * Public API: createCommentsTreeView() -> { element, render(comments), hide(), destroy() } plus COMMENTS_TREE_TEST_IDS.
 * Constraints: This module owns comments DOM creation, sanitizes only API HTML fields with DOMPurify, and consumes the upstream TreeNode order and parent contract as-is.
 */

import DOMPurify from 'dompurify';

import './comments.css';
import { clearElement, createElement } from '../../shared/dom-helpers.js';
import { formatRelativeTime } from '../../shared/time-format.js';

/** @typedef {import('../../core/entities/item.js').HnItem} HnItem */

/**
 * @typedef {Object} CommentsTreeNode
 * @property {HnItem} item
 * @property {CommentsTreeNode[]} comments
 */

// Stable test IDs give TC-4 and later audit coverage a single public DOM contract to target.
export const COMMENTS_TREE_TEST_IDS = Object.freeze({
  root: 'comments-tree',
  heading: 'comments-tree-heading',
  list: 'comments-tree-list',
  comment: 'comment-item',
  author: 'comments-tree-author',
  parent: 'comment-parent',
  time: 'comment-time',
  body: 'comments-tree-body',
});

// Visibility is centralized so render and hide paths stay consistent.
const setVisibility = (element, isVisible) => {
  element.hidden = !isVisible;
  return element;
};

// Invalid or missing author names collapse to a stable fallback instead of rendering empty metadata.
const toDisplayAuthor = (value) =>
  typeof value === 'string' && value.trim().length > 0 ? value : 'Unknown author';

// Missing timestamps fall back cleanly so the view never throws while formatting metadata.
const toDisplayTime = (value) =>
  Number.isInteger(value) ? formatRelativeTime(value) : 'Unknown time';

// Integer guards keep data attributes and CSS depth variables clean and intentional.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Numeric fallback keeps local sorting stable when upstream data is partially missing.
const toSortableTime = (value) => (Number.isInteger(value) ? value : 0);

// Numeric fallback keeps same-time comments deterministic so the list never flickers between renders.
const toSortableId = (value) => (isPositiveInteger(value) ? value : 0);

// Newest-first ordering is enforced locally as defense-in-depth even though core currently provides sorted trees.
const sortNodesNewestFirst = (nodes) =>
  nodes.toSorted((leftNode, rightNode) => {
    const timeDelta = toSortableTime(rightNode?.item?.time) - toSortableTime(leftNode?.item?.time);

    if (timeDelta !== 0) {
      return timeDelta;
    }

    return toSortableId(rightNode?.item?.id) - toSortableId(leftNode?.item?.id);
  });

// Comment bodies only render sanitized HTML when the upstream text field actually contains content.
const hasRenderableCommentText = (value) => typeof value === 'string' && value.trim().length > 0;

// Sanitized API HTML is appended directly and never mutated afterward.
const renderSanitizedBody = (element, value) => {
  // Each render starts from an empty body so prior content never leaks across states.
  clearElement(element);

  // Empty text fields stay hidden instead of producing empty comment cards.
  if (!hasRenderableCommentText(value)) {
    return false;
  }

  // DOMPurify returns a safe fragment that can be appended without using innerHTML.
  element.append(DOMPurify.sanitize(value, { RETURN_DOM_FRAGMENT: true }));

  // Child count is used as the final visibility signal after sanitization strips unsafe markup.
  return element.childNodes.length > 0;
};

/**
 * @param {CommentsTreeNode} node
 * @param {number} depth
 * @returns {HTMLElement}
 */
const createCommentNode = (node, depth, rootPostId) => {
  // Tree traversal accepts only the existing Track A node shape and falls back safely when partial data appears.
  const commentItem = node?.item ?? null;
  // Child comments are sorted defensively so each depth level always renders newest-first.
  const childNodes = sortNodesNewestFirst(Array.isArray(node?.comments) ? node.comments : []);
  // Positive depth values power both test assertions and indentation styling.
  const normalizedDepth = Number.isInteger(depth) && depth >= 0 ? depth : 0;
  // Stable integer IDs are surfaced for audit assertions without creating a new data contract.
  const commentId = isPositiveInteger(commentItem?.id) ? commentItem.id : undefined;
  // Parent IDs are surfaced exactly as supplied by the get-item use case.
  const parentId = isPositiveInteger(commentItem?.parent) ? commentItem.parent : undefined;
  const parentPostId =
    isPositiveInteger(rootPostId) || normalizedDepth > 0
      ? rootPostId
      : isPositiveInteger(parentId)
        ? parentId
        : undefined;

  // Each list item is the authoritative root for one rendered comment node.
  const commentNode = createElement('li', {
    className: 'comments-tree__item',
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.comment,
      'data-comment-id': commentId,
      'data-parent-id': parentId,
      'data-parent-post-id': parentPostId,
      'data-parent': parentPostId,
      'data-depth': normalizedDepth,
      'data-time': Number.isInteger(commentItem?.time) ? commentItem.time : undefined,
    },
  });

  // A depth variable lets CSS indent nested nodes without coupling tests to layout math.
  commentNode.style.setProperty('--comment-depth', String(normalizedDepth));

  // The article groups metadata and body content into a single semantic comment card.
  const commentArticle = createElement('article', {
    className: 'comments-tree__comment',
  });

  // Metadata is grouped in a header so author and time stay visually tied.
  const commentHeader = createElement('header', {
    className: 'comments-tree__header',
  });

  // Author labels always use text sinks because they are plain-text API fields.
  const author = createElement('span', {
    className: 'comments-tree__author',
    text: toDisplayAuthor(commentItem?.by),
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.author,
    },
  });

  // Parent references stay visible so audit checks can confirm each comment points to the expected parent item.
  const parent = createElement('span', {
    className: 'comments-tree__parent',
    text: parentPostId !== undefined ? `Parent post #${parentPostId}` : 'Reply target unavailable',
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.parent,
      'data-parent-post-id': parentPostId,
      'data-parent-id': parentId,
      'data-parent': parentPostId,
    },
  });

  // Relative-time labels stay consistent with the shared Temporal formatter used elsewhere in the app.
  const time = createElement('span', {
    className: 'comments-tree__time',
    text: toDisplayTime(commentItem?.time),
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.time,
      'data-time': Number.isInteger(commentItem?.time) ? commentItem.time : undefined,
    },
  });

  // Comment bodies isolate the only sanitized rich-text surface in this feature.
  const body = createElement('div', {
    className: 'comments-tree__body',
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.body,
    },
  });

  // Sanitized-body visibility depends on whether safe content survived purification.
  const hasVisibleBody = renderSanitizedBody(body, commentItem?.text);

  // Empty or fully stripped bodies are hidden so metadata stays compact and readable.
  setVisibility(body, hasVisibleBody);

  // Nested comments stay in an ordered list so the upstream newest-first order remains visible.
  const childList = createElement('ol', {
    className: 'comments-tree__children',
  });

  // Child levels are built in a fragment to keep recursive DOM insertion efficient.
  const childListFragment = document.createDocumentFragment();

  // Recursive rendering consumes the existing tree shape directly with one extra depth level.
  for (const childNode of childNodes) {
    childListFragment.append(createCommentNode(childNode, normalizedDepth + 1, parentPostId));
  }

  // The fragment is appended in one operation to avoid repeated reflow while recursing.
  childList.append(childListFragment);

  // Empty child lists are hidden so leaf comments do not reserve unnecessary space.
  setVisibility(childList, childNodes.length > 0);

  // Static assembly keeps the rendered order obvious and avoids later mutation branches.
  commentHeader.append(author, parent, time);
  commentArticle.append(commentHeader, body, childList);
  commentNode.append(commentArticle);

  return commentNode;
};

/**
 * @returns {{ element: HTMLElement, render(comments: CommentsTreeNode[] | undefined): void, hide(): void, destroy(): void }}
 */
export const createCommentsTreeView = () => {
  // The root section defines the feature boundary and stays hidden until comments exist.
  const element = createElement('section', {
    className: 'comments-tree',
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.root,
      'aria-labelledby': COMMENTS_TREE_TEST_IDS.heading,
    },
  });

  // The heading makes the section self-describing when it is mounted into post-detail later.
  const heading = createElement('h2', {
    className: 'comments-tree__heading',
    text: 'Comments',
    attributes: {
      id: COMMENTS_TREE_TEST_IDS.heading,
      'data-testid': COMMENTS_TREE_TEST_IDS.heading,
    },
  });

  // The ordered list preserves the provided newest-first comment order at the top level.
  const list = createElement('ol', {
    className: 'comments-tree__list',
    attributes: {
      'data-testid': COMMENTS_TREE_TEST_IDS.list,
    },
  });

  // A root-level delegated click handler keeps future interaction ownership on the section boundary.
  const handleRootClick = (event) => {
    // Non-element targets are ignored so future delegated logic has one defensive entry point.
    if (!(event.target instanceof Element)) {
      return;
    }
  };

  // Static assembly happens once so later renders only swap list content and visibility.
  element.append(heading, list);
  element.addEventListener('click', handleRootClick);

  // The comments feature starts hidden because post-detail may render states without discussions.
  setVisibility(element, false);

  return {
    element,
    // Rendering accepts the upstream tree directly and never mutates its structure.
    render(comments) {
      // Every render starts from an empty list so previous trees cannot linger across route updates.
      clearElement(list);

      // Empty inputs keep the entire feature hidden until there is something meaningful to show.
      if (!Array.isArray(comments) || comments.length === 0) {
        setVisibility(element, false);
        return;
      }

      // Top-level comments are sorted defensively so the section always displays newest threads first.
      const sortedComments = sortNodesNewestFirst(comments);

      // Top-level comments are appended through a fragment so large trees render efficiently.
      const listFragment = document.createDocumentFragment();

      // Each top-level node starts at depth zero so nested indentation remains predictable.
      for (const commentNode of sortedComments) {
        listFragment.append(createCommentNode(commentNode, 0, undefined));
      }

      // One append keeps the final DOM update compact and deterministic.
      list.append(listFragment);
      // Visibility flips only after the list is fully populated to avoid empty flashes.
      setVisibility(element, true);
    },
    // Hiding clears the current tree so stale comments never survive non-success states.
    hide() {
      clearElement(list);
      setVisibility(element, false);
    },
    // Cleanup removes the delegated root listener so future unmounting stays leak-free.
    destroy() {
      element.removeEventListener('click', handleRootClick);
    },
  };
};
