# Ticket Tracker — clonernews

> **Status Key:**
> `[ ]` To Do
> `[-]` In Progress / Review
> `[x]` Done

This file strictly tracks the execution of the implementation plan and handles the DRI (Directly Responsible Individual) bindings.

## Phase 0 — Shared Kick-off
- [x] **T0-1:** Project Scaffold (Dev 1)
- [x] **T0-2:** Core Contracts & Entities (Dev 1)
- [x] **T0-3:** Shared Utilities (Dev 2)

## Track A — Core, Use Cases, All Tests & Delivery
*Owner: Dev 1*
- [x] **TA-1:** HN API Adapter (validated and integrated with `src/infra/cache-adapter.js`) | Blocked by: T0-2
- [x] **TA-2:** List Items Use Case | Blocked by: TA-1, TC-1 (cache), TC-2 (throttle)
- [x] **TA-3:** Get Item Use Case (comment tree; implemented and validated; integrated with TA-1 and TC-1) | Blocked by: TA-1, TC-1 (cache)
- [x] **TA-4:** Poll Updates Use Case (implemented and validated) | Blocked by: TA-1, TC-2 (throttle)
- [x] **TA-5:** Integration Smoke Tests (implemented and verified with merged dependencies) | Blocked by: TA-2, TA-3, TA-4; TC-1 (cache) and TC-2 (throttle) must be merged first
- [x] **TA-6:** Audit Traceability Matrix | Blocked by: None (can start immediately)
- [x] **TA-7:** E2E Audit Test Suite — Functional & General | Blocked by: TB-2, TB-3, TC-1, TC-2, TC-3, TC-4 (features must be merged first)
- [x] **TA-8:** E2E Audit Test Suite — Bonus | Blocked by: TB-2, TC-3 (comments tree)
- [ ] **TA-9:** Performance & Accessibility QA | Blocked by: TA-7
- [ ] **TA-10:** GitHub Pages Deployment | Blocked by: T0-1

## Track B — Shared Utilities, Design System, Feed & App Shell
*Owner: Dev 2*
- [x] **TB-1:** Design System (CSS Tokens & Global Styles) | Blocked by: T0-1
- [x] **TB-2:** Feed View — Tab Navigation & Story List | Blocked by: T0-3, TB-1
- [x] **TB-3:** Infinite Scroll / Load More | Blocked by: TB-2
- [x] **TB-4:** App Shell, Routing & Responsive Layout | Blocked by: TB-1, TB-2

## Track C — Infra Utilities & Feature Views
*Owner: Dev 3*
- [x] **TC-1:** Cache Adapter | Blocked by: T0-2
- [x] **TC-2:** Throttle Utility | Blocked by: None
- [x] **TC-3:** Post Detail View — Story & Job | Blocked by: T0-3, TB-1
- [x] **TC-4:** Comments Tree (Nested/Recursive) | Blocked by: T0-3, TC-3
- [x] **TC-5:** Live-Data Notification Banner | Blocked by: T0-3, TC-1, TC-2
- [x] **TC-6:** Poll View (Pollopts) | Blocked by: T0-3, TC-3

---

**Notes:**
- Before starting a ticket, change `[ ]` to `[-]` and optionally note the working branch next to the ticket name.
- Upon merging to main, change `[-]` to `[x]`.
- Phase 0 scaffold has been bootstrapped in the workspace, including the starter README and shared utility tests.
- TA-1 now uses the concrete cache adapter from Track C, with cache-hit behavior covered in infra unit tests.
