# Agentic Workflow Guide for clonernews

This document is the working guide for a 3-developer team using coding agents in this repository. It turns the repo constraints into an operating model for planning, coding, review, testing, security, and release readiness.

If this guide conflicts with [AGENTS.md](../AGENTS.md), [docs/requirements.md](requirements.md), or [docs/audit.md](audit.md), those files win.

## 1. Operating Principles

1. One human owns one slice of work. The agent drafts code, but a human is accountable for the result.
2. Keep every task small enough to review in one pass. Prefer one feature, one use-case, or one bug fix per branch.
3. Treat the HN API as untrusted. Validate every response shape. Sanitize every HTML field before DOM insertion.
4. Keep domain logic pure and DOM side effects isolated. `src/core/` and `src/infra/` must never touch the DOM.
5. Treat agent output as untrusted until reviewed and tested.
6. Optimize for mergeability. Small PRs and clear ownership reduce conflicts more than parallelism does.

## 2. Team Model for 3 Developers

> **Canonical track ownership is defined in [`docs/implementation-plan.md` §3](implementation-plan.md#3-workflow-tracks-3-developers)**. The tracks are: **Track A** (Core, Tests & Delivery), **Track B** (Feed & App Shell), **Track C** (Feature Views). When this guide conflicts with that document on task ownership, the implementation plan wins.

This guide describes the *process layer* on top of those tracks:

- **Track A owner** (Dev 1): `src/core/`, `src/infra/hn-api-adapter.js`, `tests/` (all unit, integration, e2e), `.gitea/workflows/` (canonical CI) and `.github/workflows/` (parity/policy), project scaffold, CI, and static deployment workflow ownership. **Zero UI work.**
- **Track B owner** (Dev 2): `src/shared/`, `src/features/feed/`, `src/styles/`, `index.html`, `public/`, `src/app.css`, `src/main.js`, plus `tests/unit/features/feed/` for feed-only assertions and `data-testid` contract updates (Track A review required). Design system, feed list, infinite scroll, app shell wiring and responsive layout.
- **Track C owner** (Dev 3): `src/infra/cache-adapter.js`, `src/infra/throttle.js`, `src/features/post-detail/`, `src/features/comments/`, `src/features/live-banner/`, `src/features/polls/`. All four feature views plus infra utilities — self-contained, colocated, no interdependency between them.

That split is not rigid, but each task must have a single DRI. If a task crosses ownership boundaries, write down the boundary before the agent starts.

Recommended rule for all three devs:

- Do not let two people or agents edit the same subsystem at the same time unless the work is intentionally paired.
- Keep branches short-lived.
- Rebase or sync early and often.
- Keep `docs/tickets.md` updated with status, owner, and review state.

## 3. How to Use Agents Well

Use an agent for bounded work, not open-ended exploration.

A good task brief includes:

- Objective: what must change.
- Scope: exact files or subsystems allowed.
- Out of scope: what the agent must not touch.
- Constraints: API discipline, DOM safety, sanitization rules, and style rules.
- Acceptance: tests and manual checks that define done.
- Stop condition: the smallest proof that the task is complete.

Good examples:

- Implement the cache adapter with TTL expiry and add unit tests.
- Build the feed view with tab navigation for stories/jobs and skeleton loading.
- Add DOMPurify sanitization to post-detail comment rendering and test it.

Bad examples:

- Make the UI better.
- Fix performance.
- Improve architecture everywhere.

If the task is risky, require the agent to work in a draft PR and stop at the first verified pass of the relevant tests.

## 4. Workflow for Each Task

1. Define the slice.
2. Assign one human owner.
3. Give the agent a bounded prompt.
4. Have the agent implement the smallest viable change.
5. Add or update tests in the same branch.
6. Run the relevant checks.
7. Review the diff as a human.
8. Merge only after the PR gate passes.

For bug fixes, follow the repo bug-fix workflow:

1. Reproduce the issue.
2. Add a failing test.
3. Implement the minimal fix.
4. Prove the fix passes.
5. Check nearby systems for regressions.

## 5. Branch and PR Rules

Every branch should represent one logical change.

- Keep the branch focused on one feature, bug, or refactor.
- Avoid mixing cleanup and behavior changes unless they are inseparable.
- Delete the branch after merge.
- Do not use a feature branch as a shared workspace.
- Rebase or sync before the PR gets large.

PRs should be easy to scan in under 15 minutes.

A good PR description answers:

- What changed?
- Why was it needed?
- How was it tested?
- What is risky or still unknown?

## 6. Pre-PR Gate

A PR is not ready until the following are true.

### Required checks

- Full quality gate passes locally (`npm run ci:quality`).
- PR policy gate passes locally (`npm run ci:policy`).
- Branch name or commit messages include at least one ticket ID from `docs/tickets.md` (for example, `TA-7`).
- Audit-related e2e coverage exists for any affected audit question.
- The diff does not introduce forbidden APIs or unsafe DOM patterns.
- The change respects layer boundaries (`src/core/` has no DOM, `src/infra/` has no DOM).
- Documentation is updated if behavior, constraints, or testing expectations changed.

### Required evidence for API-critical changes

Attach a short note with:

- Endpoints tested.
- Throttle/cache behavior observed.
- Error handling verified (timeout, 404, malformed response).
- Memory or allocation notes if relevant.

## 7. Audit Queries to Check Before PR

Use `docs/audit.md` as the acceptance checklist. For any change that touches the feed, post-detail, comments, polls, live-data, or pagination, ask these questions before opening the PR:

### Functional

- Does a story post open without errors?
- Does a job post open without errors?
- Does a poll post open without errors?
- Do posts load more without error and without spamming the user?
- Are comments displayed in correct order (newest to oldest)?

### General

- Does the UI have at least stories, jobs, and polls?
- Are posts displayed in the correct order (newest to oldest)?
- Does each comment present the right parent post?
- Does the UI notify the user when a certain post is updated?
- Is the project using throttling to regulate requests (every 5 seconds)?

### Bonus

- Does the UI have more types of posts than stories, jobs, and polls?
- Have sub-comments (nested comments) been implemented?

If a PR touches one of these areas, the author should state which audit IDs changed and how they were verified. If a change affects several audit items, list them explicitly in the PR description.

## 8. Security Rules

Security is part of code review, not a separate afterthought.

### DOM and browser safety

- Prefer `textContent`, `createElement`, `appendChild`, and explicit attribute APIs.
- Avoid `innerHTML`, `outerHTML`, `insertAdjacentHTML`, `document.write`, and string-based event handlers.
- Never route untrusted data into `eval`, `Function`, or string timers.
- Sanitize all HN API HTML content with DOMPurify before DOM insertion.
- Keep domain logic (`src/core/`) and infrastructure (`src/infra/`) free of all DOM references.

### Input and validation

- Validate HN API response shapes before use.
- Reject malformed item data early (missing `id`, unexpected `type`).
- Guard against circular references in comment trees.
- Enforce max depth on recursive comment fetching.

### Dependency and supply-chain hygiene

- Keep lockfiles committed and current.
- Review new dependencies before adding them.
- Minimize lifecycle scripts and other installation-time surprises.
- Do not commit secrets, tokens, or private keys.
- Do not print sensitive data in logs, tests, or debug output.

## 9. Code Review Checklist

Reviewers should check the following before approving:

- The change is small enough to understand quickly.
- The implementation matches the stated requirement.
- The code keeps layer boundaries intact (`src/core/` = no DOM, `src/infra/` = no DOM).
- The code does not add unsafe DOM access.
- The code does not add unnecessary allocations in hot paths.
- The code has tests that fail before the fix and pass after it.
- The PR description explains the impact and the verification.
- The change does not break audit coverage.
- HN API HTML content is sanitized before insertion.
- All fetch calls use `AbortController` and validate response shape.

## 10. Suggested Review Questions

Ask these questions on every non-trivial PR:

- What is the smallest behavior change this PR makes?
- What could regress if this lands?
- How is the behavior verified automatically?
- Which audit questions does this affect?
- Does the PR add any new trust boundary?
- Does the PR respect the 5-second throttle for live data?
- Is there any safer API or simpler approach?
- Could this be split into two smaller PRs?

## 11. Team Cadence

For a 3-dev team, this cadence works well:

- Morning: claim or confirm one task each.
- During work: keep short status updates on blockers and handoffs.
- Before PR: run the local gate and attach evidence.
- During review: review one PR at a time per developer whenever possible.
- After merge: clean up the branch and mark the ticket `[x]` in `docs/tickets.md`.

If a task stalls, stop adding scope. Either finish the slice or split it.

## 12. Quick PR Template

Use this structure in PR descriptions.

Important: this section is only a compact writing aid. To pass the repository policy gate, use the full checklist and layer confirmations from `.github/pull_request_template.md`.

```md
## What changed
- 

## Why
- 

## Tests
- 

## Audit questions affected
- 

## Security notes
- 

## Architecture / dependency notes
- 

## Risks
- 
```

## 13. Practical Standard

If you only remember one rule, use this one:

> Every agent task must be small, owned, testable, and reviewable before it becomes a PR.

That single rule keeps the team fast without turning the repository into a pile of unreviewable agent output.

## 14. Automated Enforcement

Use repository automation to block unsafe or incomplete PRs before merge.

### What the gate enforces

- `npm run ci:quality` passes (`check`, `ci:guards`, unit tests, e2e tests, build).
- PR template sections/checklists are validated by `scripts/validate-pr-template.mjs` through `npm run ci:policy`.
- Architecture boundaries are validated for `src/core/` and `src/infra/` by `scripts/enforce-architecture-boundaries.mjs`.
- Legacy guardrails are validated by `scripts/ci-guards.mjs` (`var`, legacy `Date` API, `require`).
- In this repository, canonical automation runs in Gitea CI at `.gitea/workflows/ci.yml`.

### What the gate should block next (hardening backlog)

- `docs/audit.md` changes without matching traceability and test updates.
- Source changes that introduce `innerHTML`, `outerHTML`, `insertAdjacentHTML`, `document.write`, `eval`, string timers, `var`, or CommonJS imports.
- `src/core/` or `src/infra/` changes that reference `document`, `window`, `querySelector`, or other DOM APIs.
- Dependency edits that skip the lockfile.
- PRs that are still missing the required human review approval.
- Framework imports (React, Vue, Angular, Svelte, Phaser, jQuery, etc.).

### How to use it

- Keep the PR template filled out.
- Ensure branch commits include the ticket ID and validate locally with `npm run ci:policy` before opening the PR.
- Treat a green gate as the minimum for review readiness, not the finish line.
- If the gate fails, fix the root cause in the same branch before asking for another review.

## 15. Commit and PR Submission Process

Use this sequence for every branch before requesting review:

1. Commit your code changes.
2. Ensure your branch name or at least one commit message includes a valid ticket ID from `docs/tickets.md`.
3. Run `npm run ci:quality` and wait for all checks to pass.
4. Run `npm run ci:policy` and ensure it passes.
5. Fill the PR description using `.github/pull_request_template.md` and request human review.

If any step fails, fix the issue in the same branch and rerun both commands before opening or updating the PR.

## 16. Layer Boundary Quick Reference

| Layer | Directory | Allowed | Forbidden |
|---|---|---|---|
| Core/Domain | `src/core/` | Pure JS, JSDoc, Result objects | `document`, `window`, `fetch`, DOM APIs |
| Infrastructure | `src/infra/` | `fetch`, `AbortController`, `Map`, `Temporal` | `document`, `window`, DOM APIs |
| Features | `src/features/*/` | DOM (via dom-helpers), `DOMPurify`, CSS | Direct `innerHTML`, cross-feature imports |
| Shared | `src/shared/` | DOM helpers, signals, router, formatters | Business logic, API calls |
| Tests | `tests/` | Playwright, Vitest, assertions | Production code |

This table is the definitive boundary reference. If any code violates these boundaries, the CI gate should catch it. If it does not, fix the gate.
