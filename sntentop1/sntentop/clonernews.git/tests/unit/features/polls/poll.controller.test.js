/*
 * Purpose: Verify poll-controller normalization for resolved poll items before the poll view is wired in.
 * Public API: Exercises createPollController().load(item) only.
 * Constraints: Tests stay pure, assert Result contracts explicitly, and cover ratio math plus invalid-parts guards.
 */

import { describe, expect, it } from 'vitest';

import { createPollController } from '../../../../src/features/polls/poll.controller.js';

describe('poll controller', () => {
  it('normalizes a loaded poll item into ordered option bar-chart data', () => {
    // A realistic poll fixture proves that totals, ratios, and fallback flags are derived in one pass.
    const pollItem = {
      id: 90,
      type: 'poll',
      title: 'Best JavaScript runtime?',
      by: 'pg',
      time: 1_234,
      text: 'Choose one',
      parts: [
        {
          id: 901,
          type: 'pollopt',
          text: 'Node.js',
          score: 12,
        },
        {
          id: 902,
          type: 'pollopt',
          text: '',
          score: 3,
        },
        {
          id: 903,
          type: 'pollopt',
          text: 'Deno',
          score: 5,
        },
      ],
    };
    const controller = createPollController();

    const result = controller.load(pollItem);

    expect(result.ok).toBe(true);

    if (!result.ok) {
      throw new Error('Expected a successful poll normalization result.');
    }

    // Order must remain stable because the view renders options sequentially, not score-sorted.
    expect(result.data.item).toBe(pollItem);
    expect(result.data.viewModel).toEqual(
      expect.objectContaining({
        id: 90,
        title: 'Best JavaScript runtime?',
        author: 'pg',
        time: 1_234,
        text: 'Choose one',
        hasText: true,
        totalVotes: 20,
        optionCount: 3,
        leadingScore: 12,
        hasOptions: true,
      }),
    );
    expect(result.data.viewModel.options.map((option) => option.id)).toEqual([901, 902, 903]);
    expect(result.data.viewModel.options.map((option) => option.position)).toEqual([1, 2, 3]);
    expect(result.data.viewModel.options[0]).toEqual(
      expect.objectContaining({
        id: 901,
        position: 1,
        text: 'Node.js',
        hasText: true,
        score: 12,
        barWidthPercent: 60,
      }),
    );
    expect(result.data.viewModel.options[0].voteRatio).toBeCloseTo(0.6);
    expect(result.data.viewModel.options[1]).toEqual(
      expect.objectContaining({
        id: 902,
        position: 2,
        text: '',
        hasText: false,
        score: 3,
        barWidthPercent: 15,
      }),
    );
    expect(result.data.viewModel.options[1].voteRatio).toBeCloseTo(0.15);
    expect(result.data.viewModel.options[2]).toEqual(
      expect.objectContaining({
        id: 903,
        position: 3,
        text: 'Deno',
        hasText: true,
        score: 5,
        barWidthPercent: 25,
      }),
    );
    expect(result.data.viewModel.options[2].voteRatio).toBeCloseTo(0.25);
  });

  it('returns zero-width bars when every poll option has zero normalized votes', () => {
    // This protects the ratio math from division-by-zero and malformed negative scores.
    const pollItem = {
      id: 91,
      type: 'poll',
      parts: [
        {
          id: 911,
          type: 'pollopt',
          text: 'Option A',
          score: undefined,
        },
        {
          id: 912,
          type: 'pollopt',
          text: 'Option B',
          score: -4,
        },
      ],
    };
    const controller = createPollController();

    const result = controller.load(pollItem);

    expect(result.ok).toBe(true);

    if (!result.ok) {
      throw new Error('Expected a successful zero-vote poll normalization result.');
    }

    // Both the aggregate and each option should collapse cleanly to zero so the view can render safely.
    expect(result.data.viewModel).toEqual(
      expect.objectContaining({
        title: 'Untitled poll',
        author: 'Unknown author',
        totalVotes: 0,
        leadingScore: 0,
        hasOptions: true,
      }),
    );
    expect(result.data.viewModel.options).toEqual([
      expect.objectContaining({
        id: 911,
        score: 0,
        voteRatio: 0,
        barWidthPercent: 0,
      }),
      expect.objectContaining({
        id: 912,
        score: 0,
        voteRatio: 0,
        barWidthPercent: 0,
      }),
    ]);
  });

  it('rejects non-poll items explicitly', () => {
    // A dedicated failure reason keeps the later post-detail handoff branch easy to diagnose.
    const controller = createPollController();

    const result = controller.load({
      id: 92,
      type: 'story',
      title: 'Not a poll',
    });

    expect(result).toEqual({
      ok: false,
      error: 'Unsupported item type "story" for poll view.',
      reason: 'unsupported-item',
    });
  });

  it('rejects poll items whose parts are not resolved pollopt objects', () => {
    // Numeric IDs here would mean the upstream use-case contract was not honored for poll rendering.
    const controller = createPollController();

    const result = controller.load({
      id: 93,
      type: 'poll',
      parts: [931, 932],
    });

    expect(result).toEqual({
      ok: false,
      error: 'Poll item parts must contain resolved pollopt items only.',
      reason: 'invalid-parts',
    });
  });

  it('rejects malformed root items before any poll-specific normalization runs', () => {
    // Early validation keeps impossible route payloads from leaking into ratio and text handling.
    const controller = createPollController();

    const result = controller.load({
      id: 0,
      type: 'poll',
      parts: [],
    });

    expect(result).toEqual({
      ok: false,
      error: 'Poll controller requires a loaded item with a positive integer ID.',
      reason: 'invalid-item',
    });
  });
});
