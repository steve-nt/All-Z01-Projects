// Render-state tests pin the small UI-state vocabulary that the post-detail view consumes.
// Public API under test: post-detail render-state helpers and toPostDetailRenderState(loadResult).
// Constraints: tests stay pure and verify mapping rules without constructing DOM nodes.
import { describe, expect, it } from 'vitest';

import {
  createErrorPostDetailState,
  createLoadingPostDetailState,
  createNotFoundPostDetailState,
  createSuccessPostDetailState,
  POST_DETAIL_RENDER_STATES,
  toPostDetailRenderState,
} from '../../../src/features/post-detail/post-detail.render-state.js';

describe('post-detail render state', () => {
  it('creates the loading helper state', () => {
    // Loading is constructed independently because it exists before any controller result is available.
    expect(createLoadingPostDetailState()).toEqual({
      status: POST_DETAIL_RENDER_STATES.loading,
    });
  });

  it('maps successful load results to a success render state', () => {
    // Success should pass the normalized controller payload straight through to the renderer.
    const data = { viewModel: { id: 1, type: 'story' } };

    expect(toPostDetailRenderState({ ok: true, data })).toEqual(createSuccessPostDetailState(data));
  });

  it('maps not-found style controller failures to the not-found render state', () => {
    // Invalid ids, missing items, and still-unsupported types all intentionally share one empty-state UI branch.
    const invalidItemState = toPostDetailRenderState({
      ok: false,
      error: 'Bad route param',
      reason: 'invalid-item',
    });
    const notFoundState = toPostDetailRenderState({
      ok: false,
      error: 'Item 99 was not found.',
      reason: 'not-found',
    });
    const unsupportedItemState = toPostDetailRenderState({
      ok: false,
      error: 'Unsupported item type "comment" for post detail.',
      reason: 'unsupported-item',
    });

    expect(invalidItemState).toEqual(createNotFoundPostDetailState('Bad route param'));
    expect(notFoundState).toEqual(createNotFoundPostDetailState('Item 99 was not found.'));
    expect(unsupportedItemState).toEqual(
      createNotFoundPostDetailState('Unsupported item type "comment" for post detail.'),
    );
  });

  it('maps general controller failures to the error render state', () => {
    // Operational failures should remain distinct so the UI can present a generic failure message.
    expect(
      toPostDetailRenderState({
        ok: false,
        error: 'Network timeout',
        reason: 'load-error',
      }),
    ).toEqual(createErrorPostDetailState('Network timeout'));
  });
});
