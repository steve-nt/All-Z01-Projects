---
title: "feat({Track}): {Ticket Name}"
---

# PR Implementation Report: {Ticket ID}
<!-- Save as: docs/pr-messages/{TicketID}-{ShortDescription}-pr.md -->

## Summary
{Brief summary of the change and how it satisfies the ticket and linked RQ/AQ IDs.}

## Key Changes
- **{Module or file}**: {What changed}
- **{Module or file}**: {What changed}

## Technical Decisions
- **{Decision}**: {Why this approach (protocol handling, concurrency model, parsing strategy, etc.).}

## Verification Results

### Automated checks
- [ ] Unit or integration tests for this ticket pass (`go test`, `cargo test`, or the project’s documented test command).
- [ ] Build succeeds for the chosen language and entrypoint used as `./wget` (or documented equivalent).

### Manual audit (against `docs/audit.md`)
- [ ] {AQ ID}: {Pass / Fail / N/A — short note}
- [ ] {AQ ID}: {Pass / Fail / N/A — short note}

## Artifacts
- **Commands run**: {Exact commands and key output snippets if useful}
- **Sample run**: {Optional: example URL or flag combination exercised}

---

## Next Steps
{Follow-up ticket, known limitations, or deferred scope.}
