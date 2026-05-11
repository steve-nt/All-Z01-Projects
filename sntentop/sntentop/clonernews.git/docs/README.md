# clonernews Documentation Hub

Start here when you need the project docs, the contribution flow, or the PR requirements.

## Source of Truth

If docs conflict, use this order:
1. [AGENTS.md](../AGENTS.md)
2. [docs/requirements.md](requirements.md)
3. [docs/audit.md](audit.md)
4. [docs/implementation-plan.md](implementation-plan.md)
5. [docs/tickets.md](tickets.md)
6. [docs/agentic-workflow-guide.md](agentic-workflow-guide.md)
7. [README.md](../README.md)
8. [docs/hackernewsAPI.md](hackernewsAPI.md)

## What to Read

- Quick start: [README.md](../README.md) -> [AGENTS.md](../AGENTS.md) -> [docs/requirements.md](requirements.md) -> [docs/tickets.md](tickets.md)
- Full onboarding: add [docs/implementation-plan.md](implementation-plan.md), [docs/audit.md](audit.md), [docs/agentic-workflow-guide.md](agentic-workflow-guide.md), and [docs/hackernewsAPI.md](hackernewsAPI.md)
- Track A: [docs/implementation-plan.md](implementation-plan.md), [docs/audit.md](audit.md), [docs/hackernewsAPI.md](hackernewsAPI.md)
- Track B: [docs/implementation-plan.md](implementation-plan.md), [docs/requirements.md](requirements.md), [docs/audit.md](audit.md)
- Track C: [docs/implementation-plan.md](implementation-plan.md), [docs/hackernewsAPI.md](hackernewsAPI.md), [docs/audit.md](audit.md)

## Onboarding

1. Run setup checks: `npm ci`, `npm run check`, `npm test`, `npm run dev`.
2. Keep the boundaries in mind: `src/core/` has no DOM/fetch, `src/infra/` has no DOM, and API HTML must be sanitized with DOMPurify.
3. Claim a ticket in [docs/tickets.md](tickets.md) by changing `[ ]` to `[-]`.
4. Confirm the ticket requirements in [docs/implementation-plan.md](implementation-plan.md) and the matching audit items in [docs/audit.md](audit.md).

## Ticket to PR

1. Verify dependencies in [docs/implementation-plan.md](implementation-plan.md).
2. Create a focused branch like `feature/<ticket-id>-<name>` or `fix/<ticket-id>-<name>`.
3. Implement the smallest viable change and add or update tests in the same branch.
4. Run the relevant checks: `npm run check`, `npm test`, `npm run build`, `npm run preview`.
5. Fill [.github/pull_request_template.md](../.github/pull_request_template.md), list affected audit IDs, and request review.
6. After merge, update [docs/tickets.md](tickets.md) from `[-]` to `[x]`.

## Review Checklist

- [ ] Scope matches the ticket.
- [ ] Forbidden APIs are not introduced (`innerHTML`, `eval`, `var`, CommonJS, legacy Date).
- [ ] Layer boundaries are intact.
- [ ] Fetches use AbortController with timeout.
- [ ] Live-data polling stays at 5 seconds or more.
- [ ] API HTML is sanitized with DOMPurify.
- [ ] Tests pass and audit items are covered.
- [ ] PR template is complete.

## Templates

### Ticket kickoff

```md
## Ticket Kickoff
- Ticket ID:
- Title:
- Owner:
- Branch:
- Scope:
- Dependencies:
- Audit IDs:
```

### Implementation plan

```md
## Implementation Plan
1. Context
2. Files to change
3. Tests
4. Risks
```

### Bug fix

```md
## Bug Fix Plan
1. Reproduce
2. Add failing test
3. Fix
4. Verify
```

### PR body

Use this with [.github/pull_request_template.md](../.github/pull_request_template.md):

```md
## Scope
- Ticket:
- Files:

## Verification
- Checks run:
- Tests run:
- Audit IDs:

## Risk
- Main risk:
- Mitigation:
```

## Keep It Current

- Update [docs/tickets.md](tickets.md) when ticket status changes.
- Update [docs/implementation-plan.md](implementation-plan.md) when ownership or dependencies change.
- Update [docs/audit.md](audit.md) only when acceptance criteria change.
- Keep this file as the short navigation layer, not the source of detailed policy.
