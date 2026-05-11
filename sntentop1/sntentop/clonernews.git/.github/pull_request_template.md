# PR Gate Checklist

## Required checks

- [ ] I read AGENTS.md and the agentic workflow guide.
- [ ] I ran npm run ci:quality locally.
- [ ] I ran npm run ci:policy locally.
- [ ] I verified my branch commits reference at least one ticket ID from docs/tickets.md.
- [ ] I confirmed changed files stay within the declared ticket track ownership scope.
- [ ] I ran the applicable local checks for this change.
- [ ] I listed the audit IDs affected by this change.
- [ ] I checked security sinks and trust boundaries.
- [ ] I checked architecture boundaries.
- [ ] I checked dependency and lockfile impact.
- [ ] I requested human review.

## Layer boundary confirmation

- [ ] `src/core/` has no DOM or fetch references.
- [ ] `src/infra/` has no DOM references.
- [ ] All HN API HTML content is sanitized via DOMPurify before insertion.
- [ ] All fetch calls use AbortController with timeout.
- [ ] Live-data polling respects 5-second throttle minimum.

## What changed
- 

## Why
- 

## Tests
- [ ] npm run check
- [ ] npm run test
- [ ] npm run test:e2e
- [ ] npm run build
- [ ] npm run ci:quality
- [ ] npm run ci:policy

## Audit questions affected
- 

## Security notes
- 

## Architecture / dependency notes
- 

## Risks
- 
