---
title: "PR Review — {Ticket ID}: {Short Title}"
---

# {Ticket ID} Review — {Area/Module}

**Date**: {YYYY-MM-DD}  
**PR**: {link or branch name}  
**Files reviewed**: {e.g. `src/...`, `tests/...`}  
**Ticket coverage**: {e.g. RQ5 (event handling)}

---

## Ticket Definition

> {Paste the ticket’s definition / acceptance criteria here.}

---

## Implementation Review

### `{primary file path}`

| Requirement / Behavior | Status | Notes |
|---|---:|---|
| {AC #1} | ⬜ | {notes} |
| {AC #2} | ⬜ | {notes} |
| {AC #3} | ⬜ | {notes} |

**Verdict**: {✅ correct / ⚠️ mostly correct / ❌ incorrect}. {1-sentence rationale.}

---

## Test Coverage Review

### `{test file path(s)}`

| Scenario | Covered? | Notes |
|---|---:|---|
| {happy path} | ⬜ |  |
| {edge case} | ⬜ |  |
| {error / no-op behavior} | ⬜ |  |
| {idempotency / dedupe semantics} | ⬜ |  |
| {multi-subscriber / multi-event behavior} | ⬜ |  |

**Notes**:
- {Call out whether tests validate behavior vs “cleanup-only” usage.}

---

## Issues Found

### 1. {Issue title}

- **Severity**: {Blocker / Strongly recommend / Nice-to-have}
- **What/Where**: `{file}` {line/range if helpful}
- **Why it matters**: {impact}
- **Suggested fix**: {concrete change}

### 2. {Issue title}

- **Severity**: {Blocker / Strongly recommend / Nice-to-have}
- **What/Where**: `{file}` {line/range if helpful}
- **Why it matters**: {impact}
- **Suggested fix**: {concrete change}

---

## Overall Verdict

| Category | Result |
|---|---:|
| Ticket spec satisfied | {✅ / ⚠️ / ❌} |
| Lint | {✅ / ⚠️ / ❌} |
| Tests pass | {✅ / ⚠️ / ❌} |
| Test coverage completeness | {✅ / ⚠️ / ❌} |
| Follow-ups required | {✅ / ⚠️ / ❌} |

**Decision**: {✅ Approve / ⚠️ Approve with follow-ups / ❌ Request changes}

**Follow-ups**:
1. {follow-up #1}
2. {follow-up #2}

