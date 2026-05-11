/*
 * Purpose: Normalize the getItem use-case output into a small story/job/poll detail payload for the view layer.
 * Public API: createPostDetailController({ getItem }) -> { load(id) }.
 * Constraints: This module stays pure, performs no DOM work, and leaves HTML sanitization to the view.
 */

/** @typedef {import('../../core/entities/item.js').HnItem} HnItem */
import { createPollController } from '../polls/poll.controller.js';

/**
 * @typedef {'invalid-item' | 'not-found' | 'unsupported-item' | 'load-error' | 'misconfigured-controller'} PostDetailLoadFailureReason
 */

/**
 * @typedef {{ ok: true, data: PostDetailLoadOutput } | { ok: false, error: string, reason: PostDetailLoadFailureReason }} PostDetailLoadResult
 */

/**
 * @typedef {Object} PostDetailTreeNode
 * @property {HnItem} item
 * @property {PostDetailTreeNode[]} comments
 */

/**
 * @typedef {'story' | 'job' | 'poll'} SupportedPostDetailType
 */

/**
 * @typedef {Object} StoryOrJobPostDetailViewModel
 * @property {number} id
 * @property {'story' | 'job'} type
 * @property {string} title
 * @property {string} author
 * @property {number | null} time
 * @property {number | null} score
 * @property {string | null} url
 * @property {boolean} hasUrl
 * @property {string} text
 * @property {boolean} hasText
 */

/**
 * @typedef {Object} PollPostDetailOptionViewModel
 * @property {number} id
 * @property {number} position
 * @property {string} text
 * @property {boolean} hasText
 * @property {number} score
 * @property {number} voteRatio
 * @property {number} barWidthPercent
 */

/**
 * @typedef {Object} PollPostDetailViewModel
 * @property {number} id
 * @property {'poll'} type
 * @property {string} title
 * @property {string} author
 * @property {number | null} time
 * @property {string} text
 * @property {boolean} hasText
 * @property {number} totalVotes
 * @property {number} optionCount
 * @property {number} leadingScore
 * @property {boolean} hasOptions
 * @property {PollPostDetailOptionViewModel[]} options
 */

/**
 * @typedef {StoryOrJobPostDetailViewModel | PollPostDetailViewModel} PostDetailViewModel
 */

/**
 * @typedef {Object} PostDetailLoadOutput
 * @property {HnItem} item
 * @property {PostDetailTreeNode[]} comments
 * @property {PostDetailViewModel} viewModel
 */

/**
 * @typedef {Object} PostDetailControllerDependencies
 * @property {(id: number) => Promise<{ ok: true, data: { item: HnItem, comments: PostDetailTreeNode[] } } | { ok: false, error: string }>} getItem
 */

// Local Result helpers keep the controller contract aligned with the rest of the project without throwing.
const ok = (data) => ({ ok: true, data });
const err = (error, reason = 'load-error') => ({ ok: false, error, reason });

// Positive integer validation fails early so the use-case contract is consumed deliberately.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Only item types with a dedicated post-detail rendering path are accepted here.
const isSupportedPostDetailType = (value) =>
  value === 'story' || value === 'job' || value === 'poll';

// The controller recognizes Track A's stable not-found message so the UI can render a dedicated empty state.
const isNotFoundError = (value) => /^Item \d+ was not found\.$/.test(value);

// Nullable string normalization keeps the view free from repetitive undefined and empty-string checks.
const toOptionalString = (value) =>
  typeof value === 'string' && value.trim().length > 0 ? value : null;

// Fallback labels keep incomplete API payloads renderable without inventing new runtime branches in the view.
const toDisplayString = (value, fallback) =>
  typeof value === 'string' && value.trim().length > 0 ? value : fallback;

/**
 * @param {HnItem & { type: SupportedPostDetailType }} item
 * @param {string} fallbackTitle
 * @returns {StoryOrJobPostDetailViewModel}
 */
const toBaseViewModel = (item, fallbackTitle) => {
  const url = toOptionalString(item.url);
  const text = typeof item.text === 'string' ? item.text : '';

  return {
    id: item.id,
    type: item.type,
    title: toDisplayString(item.title, fallbackTitle),
    author: toDisplayString(item.by, 'Unknown author'),
    time: Number.isInteger(item.time) ? item.time : null,
    score: Number.isInteger(item.score) ? item.score : null,
    url,
    hasUrl: url !== null,
    text,
    hasText: text.trim().length > 0,
  };
};

/**
 * @param {HnItem & { type: 'story' }} item
 * @returns {StoryOrJobPostDetailViewModel}
 */
const toStoryViewModel = (item) => toBaseViewModel(item, 'Untitled story');

/**
 * @param {HnItem & { type: 'job' }} item
 * @returns {StoryOrJobPostDetailViewModel}
 */
const toJobViewModel = (item) => toBaseViewModel(item, 'Untitled job');

/**
 * @param {HnItem & { type: 'story' | 'job' }} item
 * @returns {StoryOrJobPostDetailViewModel}
 */
const toViewModel = (item) => {
  if (item.type === 'story') {
    return toStoryViewModel(item);
  }

  return toJobViewModel(item);
};

/**
 * @param {PostDetailControllerDependencies=} dependencies
 * @returns {{ load(id: number): Promise<PostDetailLoadResult> }}
 */
export const createPostDetailController = ({ getItem } = {}) => {
  // Dependency validation keeps the controller honest about its reliance on Track A's use-case.
  if (typeof getItem !== 'function') {
    return {
      load: async () =>
        err(
          'Post detail controller requires a getItem use-case function.',
          'misconfigured-controller',
        ),
    };
  }

  // Poll normalization stays delegated to the poll feature so ratio math and option shaping live in one place.
  const pollController = createPollController();

  return {
    // Loading remains the only public action for now so the view can stay dumb and data-driven.
    async load(id) {
      if (!isPositiveInteger(id)) {
        return err('Post detail controller requires a positive integer item ID.', 'invalid-item');
      }

      const getItemResult = await getItem(id);

      if (!getItemResult.ok) {
        return err(
          getItemResult.error,
          isNotFoundError(getItemResult.error) ? 'not-found' : 'load-error',
        );
      }

      const { item, comments } = getItemResult.data;

      if (!isSupportedPostDetailType(item.type)) {
        return err(`Unsupported item type "${item.type}" for post detail.`, 'unsupported-item');
      }

      if (item.type === 'poll') {
        const pollResult = pollController.load(item);

        if (!pollResult.ok) {
          return err(pollResult.error, 'load-error');
        }

        return ok({
          item,
          comments,
          viewModel: {
            ...pollResult.data.viewModel,
            type: 'poll',
          },
        });
      }

      return ok({
        item,
        comments,
        viewModel: toViewModel(item),
      });
    },
  };
};
