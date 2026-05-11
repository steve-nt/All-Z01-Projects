// The entity tests verify that the pure contract layer exports the expected types and field names.
// Public API under test: HN_ITEM_TYPES, HN_ITEM_FIELD_NAMES, and isHnItemType from the item entity module.
// Constraints: assertions stay pure and deterministic without DOM, fetch, or time-dependent behavior.
import { describe, expect, it } from 'vitest';

import {
  HN_ITEM_FIELD_NAMES,
  HN_ITEM_TYPES,
  isHnItemType,
} from '../../../src/core/entities/item.js';

describe('hn item entity', () => {
  it('exposes the expected item types', () => {
    // These checks ensure the downstream validators keep the same set of supported item kinds.
    expect(HN_ITEM_TYPES).toEqual(['story', 'job', 'poll', 'pollopt', 'comment']);
    expect(isHnItemType('story')).toBe(true);
    expect(isHnItemType('unknown')).toBe(false);
  });

  it('lists the expected HN field names', () => {
    // The field list is a contract boundary, so tests protect it from accidental drift.
    expect(HN_ITEM_FIELD_NAMES).toContain('title');
    expect(HN_ITEM_FIELD_NAMES).toContain('parent');
  });
});
