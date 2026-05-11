// Core item contracts stay pure so feature and infra layers can validate Hacker News payloads consistently.
// The exported type metadata is shared across Track A and later feature layers without importing UI code.
// Keeping the contract centralized reduces drift between adapters, use cases, and tests.
/**
 * @typedef {'story' | 'job' | 'poll' | 'pollopt' | 'comment'} HnItemType
 */

/**
 * @typedef {Object} HnItem
 * @property {number} id
 * @property {HnItemType} type
 * @property {string=} by
 * @property {number=} time
 * @property {string=} text
 * @property {number[]=} kids
 * @property {string=} url
 * @property {number=} score
 * @property {string=} title
 * @property {number[]=} parts
 * @property {number=} descendants
 * @property {number=} parent
 */

/**
 * @typedef {ReadonlyArray<'id' | 'type' | 'by' | 'time' | 'text' | 'kids' | 'url' | 'score' | 'title' | 'parts' | 'descendants' | 'parent'>} HnItemFieldNames
 */

// Centralized item-type constants keep adapter and use-case validation rules synchronized.
export const HN_ITEM_TYPES = Object.freeze(['story', 'job', 'poll', 'pollopt', 'comment']);

// Field-name constants provide a stable contract snapshot for tests and downstream tooling.
export const HN_ITEM_FIELD_NAMES = Object.freeze([
  'id',
  'type',
  'by',
  'time',
  'text',
  'kids',
  'url',
  'score',
  'title',
  'parts',
  'descendants',
  'parent',
]);

// Type checks are centralized to prevent duplicate string-literal checks across layers.
export const isHnItemType = (value) => HN_ITEM_TYPES.includes(value);
