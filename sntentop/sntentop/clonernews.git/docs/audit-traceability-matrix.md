# TA-6 Audit Traceability Matrix

This matrix maps every audit question occurrence in `docs/audit.md` to implementing ticket(s) and a planned/target e2e test file.

## Matrix

| Audit ID | Audit Section | Audit Prompt Context | Audit Question (verbatim from `docs/audit.md`) | Implementing Ticket(s) | Planned/Target Test File |
|---|---|---|---|---|---|
| AUDIT-F-01 | Functional | Try to open a story post | Does this post open without any errors? | TC-3 | tests/e2e/audit/stories.spec.js |
| AUDIT-F-02 | Functional | Try to open a job post | Does this post open without any errors? | TC-3 | tests/e2e/audit/jobs.spec.js |
| AUDIT-F-03 | Functional | Try to open a poll post | Does this post open without any errors? | TC-6 | tests/e2e/audit/polls.spec.js |
| AUDIT-F-04 | Functional | Try to load more posts | Did the posts load without error and without spamming the user? | TB-3 | tests/e2e/audit/load-more.spec.js |
| AUDIT-F-05 | Functional | Try to open a post with comments | Are the comments being displayed in the correct order (from newest to oldest)? | TC-4 | tests/e2e/audit/comments.spec.js |
| AUDIT-G-01 | General | General checklist | Does the UI have at least stories, jobs and polls? | TB-2, TC-3, TC-6 | tests/e2e/audit/stories.spec.js; tests/e2e/audit/polls.spec.js |
| AUDIT-G-02 | General | General checklist | Are the posts displayed in the correct order (from newest to oldest)? | TB-2 | tests/e2e/audit/stories.spec.js |
| AUDIT-G-03 | General | General checklist | Does each comment present the right parent post? | TC-4 | tests/e2e/audit/comments.spec.js |
| AUDIT-G-04 | General | General checklist | Does the UI notify the user when a certain post is updated? | TC-5 | tests/e2e/audit/live-data.spec.js |
| AUDIT-G-05 | General | General checklist | Is the project using throttling to regulate the number of requests (every 5 seconds)? | TA-4, TC-5 | tests/e2e/audit/live-data.spec.js |
| AUDIT-B-01 | Bonus | Bonus checklist | +Does the UI have more types of posts than stories, jobs and polls? | TB-2 | tests/e2e/audit/stories.spec.js |
| AUDIT-B-02 | Bonus | Bonus checklist | +Have sub-comments (nested comments) been implemented? | TC-4 | tests/e2e/audit/comments.spec.js |

## Verification

- ID coverage is contiguous and complete for all required ranges (Functional 5, General 5, Bonus 2).
- Every audit question occurrence from `docs/audit.md` is represented exactly once in this matrix.
- No duplicate IDs are present in this file.
