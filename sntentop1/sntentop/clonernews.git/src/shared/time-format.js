// Relative time formatting stays here so feature views can reuse Temporal-based labels consistently.
import { Temporal } from '@js-temporal/polyfill';

// The browser formatter handles pluralization and localized phrasing for the final label.
const relativeTimeFormatter = new Intl.RelativeTimeFormat('en', {
  numeric: 'auto',
});

// Units are ordered from largest to smallest so the formatter chooses the most human-friendly bucket first.
const timeUnits = [
  { unit: 'year', seconds: 60 * 60 * 24 * 365 },
  { unit: 'month', seconds: 60 * 60 * 24 * 30 },
  { unit: 'day', seconds: 60 * 60 * 24 },
  { unit: 'hour', seconds: 60 * 60 },
  { unit: 'minute', seconds: 60 },
  { unit: 'second', seconds: 1 },
];

// Normalize supported inputs into Temporal.Instant so the public formatter can stay flexible.
const toInstant = (value) => {
  if (value instanceof Temporal.Instant) {
    return value;
  }

  if (typeof value === 'number') {
    return Temporal.Instant.fromEpochMilliseconds(value * 1000);
  }

  if (typeof value === 'string') {
    return Temporal.Instant.from(value);
  }

  throw new TypeError('Expected an epoch second number, ISO string, or Temporal.Instant.');
};

/**
 * Format an HN timestamp as relative time without relying on legacy Date APIs.
 *
 * The optional `now` parameter keeps tests deterministic and avoids hidden wall-clock dependence.
 *
 * @param {number | string | Temporal.Instant} value
 * @param {Temporal.Instant=} now
 * @returns {string}
 */
export const formatRelativeTime = (value, now = Temporal.Now.instant()) => {
  const instant = toInstant(value);
  // Convert to seconds early so the bucket loop can stay simple and integer-based.
  const differenceInSeconds = Math.trunc(
    (now.epochMilliseconds - instant.epochMilliseconds) / 1000,
  );

  if (differenceInSeconds === 0) {
    return 'now';
  }

  const absoluteDifference = Math.abs(differenceInSeconds);

  // Walk the buckets from large to small so the first matching unit reads naturally.
  for (const { unit, seconds } of timeUnits) {
    if (absoluteDifference >= seconds || unit === 'second') {
      const amount = -Math.trunc(differenceInSeconds / seconds);

      return relativeTimeFormatter.format(amount, unit);
    }
  }

  return 'now';
};
