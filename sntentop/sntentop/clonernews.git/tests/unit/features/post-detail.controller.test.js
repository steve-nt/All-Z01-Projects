// Post-detail controller tests pin the Track C normalization contract before route wiring lands.
// Public API under test: createPostDetailController and its load(id) behavior for story/job/poll detail flows.
// Constraints: tests stay pure, mock getItem directly, and verify Result contracts without touching the DOM.
import { describe, expect, it, vi } from 'vitest';

import { createPostDetailController } from '../../../src/features/post-detail/post-detail.controller.js';

// Result helpers mirror the production contract so test fixtures stay concise and explicit.
const ok = (data) => ({ ok: true, data });
const err = (error) => ({ ok: false, error });

describe('post-detail controller', () => {
  it('returns a misconfigured-controller failure when getItem is missing', async () => {
    // Missing dependencies must fail cleanly so view wiring errors do not turn into thrown exceptions.
    const controller = createPostDetailController();

    const result = await controller.load(1);

    expect(result).toEqual({
      ok: false,
      error: 'Post detail controller requires a getItem use-case function.',
      reason: 'misconfigured-controller',
    });
  });

  it('rejects invalid item ids before calling getItem', async () => {
    // Early validation keeps impossible route params from leaking into the shared use-case contract.
    const getItem = vi.fn(async () => ok({ item: { id: 1, type: 'story' }, comments: [] }));
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(0);

    expect(result).toEqual({
      ok: false,
      error: 'Post detail controller requires a positive integer item ID.',
      reason: 'invalid-item',
    });
    expect(getItem).not.toHaveBeenCalled();
  });

  it('normalizes story items into a story view model', async () => {
    // Story normalization should preserve optional fields and expose flags the view can branch on safely.
    const getItem = vi.fn(async () =>
      ok({
        item: {
          id: 11,
          type: 'story',
          title: 'Launch post',
          by: 'alice',
          time: 1_234,
          score: 42,
          url: 'https://example.com/story',
          text: '<p>Story text</p>',
        },
        comments: [],
      }),
    );
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(11);

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          item: expect.objectContaining({ id: 11, type: 'story' }),
          comments: [],
          viewModel: {
            id: 11,
            type: 'story',
            title: 'Launch post',
            author: 'alice',
            time: 1_234,
            score: 42,
            url: 'https://example.com/story',
            hasUrl: true,
            text: '<p>Story text</p>',
            hasText: true,
          },
        }),
      }),
    );
    expect(getItem).toHaveBeenCalledWith(11);
  });

  it('normalizes job items with missing fields without breaking the view contract', async () => {
    // Job posts often omit urls and metadata, so the controller must supply safe fallbacks and flags.
    const getItem = vi.fn(async () =>
      ok({
        item: {
          id: 22,
          type: 'job',
          text: '<p>Now hiring</p>',
          url: '',
        },
        comments: [],
      }),
    );
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(22);

    expect(result).toEqual(
      expect.objectContaining({
        ok: true,
        data: expect.objectContaining({
          item: expect.objectContaining({ id: 22, type: 'job' }),
          comments: [],
          viewModel: {
            id: 22,
            type: 'job',
            title: 'Untitled job',
            author: 'Unknown author',
            time: null,
            score: null,
            url: null,
            hasUrl: false,
            text: '<p>Now hiring</p>',
            hasText: true,
          },
        }),
      }),
    );
  });

  it('maps missing items to a not-found failure reason', async () => {
    // The view-state mapper depends on this reason to render an empty state instead of a generic error.
    const getItem = vi.fn(async () => err('Item 404 was not found.'));
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(404);

    expect(result).toEqual({
      ok: false,
      error: 'Item 404 was not found.',
      reason: 'not-found',
    });
  });

  it('normalizes poll items into a poll view model for the poll subview', async () => {
    // Poll normalization should preserve resolved option order and precomputed bar-chart data.
    const getItem = vi.fn(async () =>
      ok({
        item: {
          id: 33,
          type: 'poll',
          title: 'Best runtime?',
          by: 'dang',
          time: 3_300,
          text: '<p>Vote now</p>',
          parts: [
            {
              id: 331,
              type: 'pollopt',
              text: 'Node.js',
              score: 10,
            },
            {
              id: 332,
              type: 'pollopt',
              text: 'Deno',
              score: 8,
            },
          ],
        },
        comments: [],
      }),
    );
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(33);

    expect(result.ok).toBe(true);

    if (!result.ok) {
      throw new Error('Expected poll items to normalize successfully.');
    }

    expect(result.data).toEqual(
      expect.objectContaining({
        item: expect.objectContaining({ id: 33, type: 'poll' }),
        comments: [],
        viewModel: expect.objectContaining({
          id: 33,
          type: 'poll',
          title: 'Best runtime?',
          author: 'dang',
          time: 3_300,
          text: '<p>Vote now</p>',
          hasText: true,
          totalVotes: 18,
          optionCount: 2,
          leadingScore: 10,
          hasOptions: true,
        }),
      }),
    );
    expect(result.data.viewModel.options).toHaveLength(2);
    expect(result.data.viewModel.options[0]).toEqual(
      expect.objectContaining({
        id: 331,
        position: 1,
        text: 'Node.js',
        score: 10,
      }),
    );
    expect(result.data.viewModel.options[0].voteRatio).toBeCloseTo(10 / 18);
    expect(result.data.viewModel.options[0].barWidthPercent).toBeCloseTo(55.5555555556);
    expect(result.data.viewModel.options[1]).toEqual(
      expect.objectContaining({
        id: 332,
        position: 2,
        text: 'Deno',
        score: 8,
      }),
    );
    expect(result.data.viewModel.options[1].voteRatio).toBeCloseTo(8 / 18);
    expect(result.data.viewModel.options[1].barWidthPercent).toBeCloseTo(44.4444444444);
  });

  it('rejects unsupported post types that still fall outside the detail workflow', async () => {
    // Comment nodes still belong to the comments feature rather than the post-detail success path.
    const getItem = vi.fn(async () =>
      ok({
        item: {
          id: 33,
          type: 'comment',
        },
        comments: [],
      }),
    );
    const controller = createPostDetailController({ getItem });

    const result = await controller.load(33);

    expect(result).toEqual({
      ok: false,
      error: 'Unsupported item type "comment" for post detail.',
      reason: 'unsupported-item',
    });
  });
});
