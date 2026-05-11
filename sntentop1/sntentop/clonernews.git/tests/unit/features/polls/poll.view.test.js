// @vitest-environment jsdom

// Poll-view tests define the standalone DOM contract before the feature is wired into post-detail.
// Public API under test: createPollView, POLL_TEST_IDS, and the child-view render/hide lifecycle.
// Constraints: tests stay within jsdom, mock relative-time formatting for determinism, and assert sanitized rich-text behavior without depending on routing.
import { describe, expect, it, vi } from 'vitest';

// Relative-time labels are mocked so header assertions stay stable across all runs.
vi.mock('../../../../src/shared/time-format.js', () => ({
  formatRelativeTime: (value) => `relative-${value}`,
}));

import { createPollView, POLL_TEST_IDS } from '../../../../src/features/polls/poll.view.js';

// Data-testid lookup includes the root node itself so the feature boundary can be asserted directly.
const getByTestId = (root, testId) =>
  root?.matches?.(`[data-testid="${testId}"]`)
    ? root
    : root?.querySelector(`[data-testid="${testId}"]`);

// Repeated selectors stay centralized so option-list assertions read as feature behavior, not selector plumbing.
const getAllByTestId = (root, testId) => [...root.querySelectorAll(`[data-testid="${testId}"]`)];

// Scoped lookups keep per-option assertions tied to one rendered row.
const getScopedByTestId = (root, testId) => root?.querySelector(`[data-testid="${testId}"]`);

// A realistic fixture keeps the acceptance contract explicit for header text, vote counts, and bar widths.
const createPollViewModel = () => ({
  id: 70,
  title: 'Best JavaScript runtime?',
  author: 'dang',
  time: 70,
  text: '<p>Choose a <strong>winner</strong>.</p><script>boom()</script>',
  hasText: true,
  totalVotes: 20,
  optionCount: 3,
  leadingScore: 12,
  hasOptions: true,
  options: [
    {
      id: 701,
      position: 1,
      text: '<span>Node.js</span><img src="x" onerror="alert(1)">',
      hasText: true,
      score: 12,
      voteRatio: 0.6,
      barWidthPercent: 60,
    },
    {
      id: 702,
      position: 2,
      text: '',
      hasText: false,
      score: 3,
      voteRatio: 0.15,
      barWidthPercent: 15,
    },
    {
      id: 703,
      position: 3,
      text: '<em>Deno</em>',
      hasText: true,
      score: 5,
      voteRatio: 0.25,
      barWidthPercent: 25,
    },
  ],
});

describe('poll view', () => {
  it('starts hidden so post-detail can mount the subview without showing stale poll content', () => {
    // Hidden-by-default behavior matches the existing child-view pattern used for other post-detail subviews.
    const view = createPollView();

    expect(getByTestId(view.element, POLL_TEST_IDS.root)?.hidden).toBe(true);
    expect(getAllByTestId(view.element, POLL_TEST_IDS.option)).toHaveLength(0);

    view.destroy();
  });

  it('renders the poll header, ordered options, vote counts, and ratio-based bars', () => {
    // The main success case defines the DOM contract the upcoming implementation must satisfy end to end.
    const view = createPollView();
    const viewModel = createPollViewModel();

    view.render(viewModel);

    const root = getByTestId(view.element, POLL_TEST_IDS.root);
    const title = getByTestId(view.element, POLL_TEST_IDS.title);
    const metadata = getByTestId(view.element, POLL_TEST_IDS.metadata);
    const textBody = getByTestId(view.element, POLL_TEST_IDS.textBody);
    const totalVotes = getByTestId(view.element, POLL_TEST_IDS.totalVotes);
    const optionRows = getAllByTestId(view.element, POLL_TEST_IDS.option);
    const optionBars = getAllByTestId(view.element, POLL_TEST_IDS.optionBar);

    expect(root?.hidden).toBe(false);
    expect(title?.textContent).toBe('Best JavaScript runtime?');
    expect(metadata?.textContent).toContain('dang');
    expect(metadata?.textContent).toContain('relative-70');
    expect(totalVotes?.textContent).toBe('20 votes');
    expect(textBody?.querySelector('strong')?.textContent).toBe('winner');
    expect(textBody?.querySelector('script')).toBeNull();
    expect(optionRows).toHaveLength(3);
    expect(optionBars).toHaveLength(3);

    // Option rows must preserve controller order so the view does not silently reshuffle poll choices.
    const firstOptionLabel = getScopedByTestId(optionRows[0], POLL_TEST_IDS.optionLabel);
    const secondOptionLabel = getScopedByTestId(optionRows[1], POLL_TEST_IDS.optionLabel);
    const thirdOptionLabel = getScopedByTestId(optionRows[2], POLL_TEST_IDS.optionLabel);
    const firstOptionVotes = getScopedByTestId(optionRows[0], POLL_TEST_IDS.optionVotes);
    const secondOptionVotes = getScopedByTestId(optionRows[1], POLL_TEST_IDS.optionVotes);
    const thirdOptionVotes = getScopedByTestId(optionRows[2], POLL_TEST_IDS.optionVotes);

    expect(firstOptionLabel?.textContent).toContain('Node.js');
    expect(secondOptionLabel?.textContent).toBe('Option 2');
    expect(thirdOptionLabel?.textContent).toContain('Deno');
    expect(firstOptionVotes?.textContent).toBe('12 votes');
    expect(secondOptionVotes?.textContent).toBe('3 votes');
    expect(thirdOptionVotes?.textContent).toBe('5 votes');
    expect(optionBars[0]?.style.width).toBe('60%');
    expect(optionBars[1]?.style.width).toBe('15%');
    expect(optionBars[2]?.style.width).toBe('25%');

    view.destroy();
  });

  it('sanitizes both poll text and option text while preserving safe markup', () => {
    // Both the root poll text and each option label come from API text fields, so both need purification.
    const view = createPollView();

    view.render(createPollViewModel());

    const textBody = getByTestId(view.element, POLL_TEST_IDS.textBody);
    const firstOption = getAllByTestId(view.element, POLL_TEST_IDS.option)[0];
    const firstOptionLabel = getScopedByTestId(firstOption, POLL_TEST_IDS.optionLabel);
    const sanitizedImage = firstOptionLabel?.querySelector('img');

    expect(textBody?.querySelector('strong')?.textContent).toBe('winner');
    expect(textBody?.querySelector('script')).toBeNull();
    expect(firstOptionLabel?.querySelector('span')?.textContent).toBe('Node.js');
    expect(firstOptionLabel?.querySelector('script')).toBeNull();
    expect(sanitizedImage?.getAttribute('onerror')).toBeNull();

    view.destroy();
  });

  it('hide clears rendered rows and hides the feature boundary again', () => {
    // The reset path matters because post-detail will reuse one mounted poll subview across route changes.
    const view = createPollView();

    view.render(createPollViewModel());

    expect(getByTestId(view.element, POLL_TEST_IDS.root)?.hidden).toBe(false);
    expect(getAllByTestId(view.element, POLL_TEST_IDS.option)).toHaveLength(3);

    view.hide();

    expect(getByTestId(view.element, POLL_TEST_IDS.root)?.hidden).toBe(true);
    expect(getAllByTestId(view.element, POLL_TEST_IDS.option)).toHaveLength(0);

    view.destroy();
  });
});
