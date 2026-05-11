/*
 * Purpose: Define the explicit UI render states for the post-detail feature before the DOM view is implemented.
 * Public API: POST_DETAIL_RENDER_STATES, state factory helpers, and toPostDetailRenderState(loadResult).
 * Constraints: This module stays pure, does no DOM work, and only maps controller output into UI-friendly state shapes.
 */

/**
 * @typedef {'loading' | 'success' | 'error' | 'not-found'} PostDetailRenderStatus
 */

/**
 * @typedef {Object} PostDetailLoadingState
 * @property {'loading'} status
 */

/**
 * @typedef {Object} PostDetailSuccessState
 * @property {'success'} status
 * @property {object} data
 */

/**
 * @typedef {Object} PostDetailErrorState
 * @property {'error'} status
 * @property {string} message
 */

/**
 * @typedef {Object} PostDetailNotFoundState
 * @property {'not-found'} status
 * @property {string} message
 */

/**
 * @typedef {PostDetailLoadingState | PostDetailSuccessState | PostDetailErrorState | PostDetailNotFoundState} PostDetailRenderState
 */

// Shared status constants keep the view and controller-state mapping aligned on one vocabulary.
export const POST_DETAIL_RENDER_STATES = Object.freeze({
  loading: 'loading',
  success: 'success',
  error: 'error',
  notFound: 'not-found',
});

// Loading state stays intentionally minimal because the view only needs to know that work is in progress.
export const createLoadingPostDetailState = () => ({
  status: POST_DETAIL_RENDER_STATES.loading,
});

// Success state carries the normalized controller payload straight through to the future renderer.
export const createSuccessPostDetailState = (data) => ({
  status: POST_DETAIL_RENDER_STATES.success,
  data,
});

// Error normalization prevents blank UI states when upstream failures omit or mangle their messages.
export const createErrorPostDetailState = (message) => ({
  status: POST_DETAIL_RENDER_STATES.error,
  message:
    typeof message === 'string' && message.trim().length > 0 ? message : 'Failed to load post.',
});

// Missing and invalid detail routes share one empty-state bucket so the view can stay small and explicit.
export const createNotFoundPostDetailState = (message = 'Post not found.') => ({
  status: POST_DETAIL_RENDER_STATES.notFound,
  message,
});

// Invalid and unsupported item cases map to the empty state instead of the generic error state by design.
const isNotFoundRenderReason = (reason) =>
  reason === 'not-found' || reason === 'invalid-item' || reason === 'unsupported-item';

/**
 * @param {{ ok: true, data: object } | { ok: false, error: string, reason?: string }} loadResult
 * @returns {PostDetailRenderState}
 */
export const toPostDetailRenderState = (loadResult) => {
  if (loadResult?.ok) {
    return createSuccessPostDetailState(loadResult.data);
  }

  if (isNotFoundRenderReason(loadResult?.reason)) {
    return createNotFoundPostDetailState(loadResult?.error);
  }

  return createErrorPostDetailState(loadResult?.error);
};
