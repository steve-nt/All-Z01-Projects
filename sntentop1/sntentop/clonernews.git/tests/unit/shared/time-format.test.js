// Time formatting tests stay deterministic by pinning both the sample time and the reference instant.
// Public API under test: formatRelativeTime from the shared time-format module.
// Constraints: tests use Temporal instants only and avoid locale-unstable assumptions.
import { Temporal } from '@js-temporal/polyfill';
import { describe, expect, it } from 'vitest';

import { formatRelativeTime } from '../../../src/shared/time-format.js';

describe('time format', () => {
  it('formats a recent past timestamp', () => {
    // Past-tense formatting is the common case for HN timestamps, so it gets the first assertion.
    const now = Temporal.Instant.from('2026-03-28T12:00:00Z');
    const tenMinutesAgo = Temporal.Instant.from('2026-03-28T11:50:00Z');

    expect(formatRelativeTime(tenMinutesAgo, now)).toBe('10 minutes ago');
  });

  it('formats a future timestamp', () => {
    // Future dates verify that the formatter is symmetric instead of only handling elapsed time.
    const now = Temporal.Instant.from('2026-03-28T12:00:00Z');
    const inTwoHours = Temporal.Instant.from('2026-03-28T14:00:00Z');

    expect(formatRelativeTime(inTwoHours, now)).toBe('in 2 hours');
  });
});
