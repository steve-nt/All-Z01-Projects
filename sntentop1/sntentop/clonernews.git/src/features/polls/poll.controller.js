/*
 * Purpose: Normalize a loaded poll item into poll-option bar-chart data for the poll view.
 * Public API: createPollController() -> { load(item) }.
 * Constraints: This module stays pure, expects a getItem-resolved poll payload, preserves option order, and leaves HTML sanitization to the view.
 */

/** @typedef {import('../../core/entities/item.js').HnItem} HnItem */

/**
 * @typedef {Omit<HnItem, 'type'> & { type: 'pollopt' }} PollOptionItem
 */

/**
 * @typedef {Omit<HnItem, 'type' | 'parts'> & { type: 'poll', parts?: PollOptionItem[] }} PollItem
 */

/**
 * @typedef {'invalid-item' | 'unsupported-item' | 'invalid-parts'} PollLoadFailureReason
 */

/**
 * @typedef {{ ok: true, data: PollLoadOutput } | { ok: false, error: string, reason: PollLoadFailureReason }} PollLoadResult
 */

/**
 * @typedef {Object} PollOptionViewModel
 * @property {number} id
 * @property {number} position
 * @property {string} text
 * @property {boolean} hasText
 * @property {number} score
 * @property {number} voteRatio
 * @property {number} barWidthPercent
 */

/**
 * @typedef {Object} PollViewModel
 * @property {number} id
 * @property {string} title
 * @property {string} author
 * @property {number | null} time
 * @property {string} text
 * @property {boolean} hasText
 * @property {number} totalVotes
 * @property {number} optionCount
 * @property {number} leadingScore
 * @property {boolean} hasOptions
 * @property {PollOptionViewModel[]} options
 */

/**
 * @typedef {Object} PollLoadOutput
 * @property {PollItem} item
 * @property {PollViewModel} viewModel
 */

// Small Result constructors keep the controller aligned with the rest of the codebase without throwing.
const ok = (data) => ({ ok: true, data });
const err = (error, reason) => ({ ok: false, error, reason });

// Poll IDs must stay valid so later routing and DOM test hooks can safely rely on them.
const isPositiveInteger = (value) => Number.isInteger(value) && value > 0;

// Object guards keep malformed payloads from leaking into downstream normalization code.
const isRecord = (value) => typeof value === 'object' && value !== null;

// Poll rendering only applies to root poll items, so other item types are rejected explicitly.
const isPollItem = (value) =>
  isRecord(value) && value.type === 'poll' && isPositiveInteger(value.id);

// Resolved poll options must arrive as concrete pollopt objects rather than raw numeric IDs.
const isResolvedPollOption = (value) =>
  isRecord(value) && value.type === 'pollopt' && isPositiveInteger(value.id);

// Human-readable fallbacks keep incomplete API payloads renderable without extra view branching.
const toDisplayString = (value, fallback) =>
  typeof value === 'string' && value.trim().length > 0 ? value : fallback;

// Poll and option text stays raw here because the feature view owns DOMPurify and DOM insertion.
const toRawText = (value) => (typeof value === 'string' ? value : '');

// Missing vote tallies should collapse to zero so ratio math remains deterministic and exception-free.
const toVoteScore = (value) => (Number.isInteger(value) && value >= 0 ? value : 0);

// Zero-total polls should render empty bars instead of producing NaN or Infinity ratios.
const toVoteRatio = (score, totalVotes) => (totalVotes > 0 ? score / totalVotes : 0);

/**
 * @param {PollOptionItem[]} parts
 * @returns {PollOptionViewModel[]}
 */
const toOptionViewModels = (parts) => {
  // Raw score normalization happens first so totals and individual options use the same canonical values.
  const normalizedOptions = parts.map((part, index) => {
    const score = toVoteScore(part.score);
    const text = toRawText(part.text);

    return {
      id: part.id,
      position: index + 1,
      text,
      hasText: text.trim().length > 0,
      score,
      voteRatio: 0,
      barWidthPercent: 0,
    };
  });

  // Total votes are derived from normalized scores so bad API values cannot skew bar math.
  const totalVotes = normalizedOptions.reduce(
    (runningTotal, option) => runningTotal + option.score,
    0,
  );

  // Final option objects carry precomputed ratios so the view can stay focused on rendering only.
  return normalizedOptions.map((option) => {
    const voteRatio = toVoteRatio(option.score, totalVotes);

    return {
      ...option,
      voteRatio,
      barWidthPercent: voteRatio * 100,
    };
  });
};

/**
 * @param {PollItem} item
 * @returns {PollViewModel}
 */
const toViewModel = (item) => {
  // Missing parts are treated as an empty option list so edge-case polls still render safely.
  const parts = Array.isArray(item.parts) ? item.parts : [];
  const options = toOptionViewModels(parts);
  // The aggregate tally is recomputed from normalized options so totals always match rendered bars.
  const totalVotes = options.reduce((runningTotal, option) => runningTotal + option.score, 0);
  // The leading score is useful for future styling and testing without forcing the view to rescan options.
  const leadingScore = options.reduce(
    (currentLeadingScore, option) => Math.max(currentLeadingScore, option.score),
    0,
  );
  const text = toRawText(item.text);

  return {
    id: item.id,
    title: toDisplayString(item.title, 'Untitled poll'),
    author: toDisplayString(item.by, 'Unknown author'),
    time: Number.isInteger(item.time) ? item.time : null,
    text,
    hasText: text.trim().length > 0,
    totalVotes,
    optionCount: options.length,
    leadingScore,
    hasOptions: options.length > 0,
    options,
  };
};

/**
 * @returns {{ load(item: HnItem): PollLoadResult }}
 */
export const createPollController = () => ({
  // The controller accepts a fully loaded item so poll normalization stays decoupled from data fetching.
  load(item) {
    if (!isRecord(item) || !isPositiveInteger(item.id)) {
      return err(
        'Poll controller requires a loaded item with a positive integer ID.',
        'invalid-item',
      );
    }

    if (!isPollItem(item)) {
      return err(`Unsupported item type "${item.type}" for poll view.`, 'unsupported-item');
    }

    if (item.parts !== undefined && !Array.isArray(item.parts)) {
      return err('Poll item parts must be an array of resolved poll options.', 'invalid-parts');
    }

    if (Array.isArray(item.parts) && !item.parts.every(isResolvedPollOption)) {
      return err('Poll item parts must contain resolved pollopt items only.', 'invalid-parts');
    }

    return ok({
      item,
      viewModel: toViewModel(item),
    });
  },
});

export default createPollController;
