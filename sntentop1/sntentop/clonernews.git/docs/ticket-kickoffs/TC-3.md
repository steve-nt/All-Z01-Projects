# Ticket Kickoff
- Ticket ID: `TC-3`
- Title: `Post Detail View — Story & Job`
- Owner: `Dev 3 / Track C`
- Branch: `TBD`
- Scope:
  - Build the base post-detail feature in `src/features/post-detail/`.
  - Support `story` and `job` items only.
  - Render title, URL, author, time, score, and sanitized `text` when present.
  - Add stable back navigation and descriptive `data-testid` attributes.
- Out of scope:
  - `TC-4` comments tree rendering.
  - `TC-5` live-data banner behavior.
  - `TC-6` poll-specific rendering.
  - Changes to Track A use-case behavior beyond consuming the existing `getItem(id)` contract.
- Dependencies:
  - `T0-3` shared utilities: complete.
  - `TB-1` design system: complete.
- Audit IDs:
  - `AUDIT-F-01`
  - `AUDIT-F-02`

## Integration Boundary
- The route contract for `#/item/:id` already exists in `src/shared/router.js`, so `TC-3` should consume that contract rather than redefine routing.
- The current app shell in `src/main.js` is still scaffold-only, so detail-page mounting will require a small integration change there unless Track B lands broader shell wiring first.
- Any cross-track edits should stay minimal and explicit:
  - Prefer leaving `src/shared/router.js` unchanged.
  - Limit Track B-owned touch points to composition wiring in `src/main.js` only if necessary.
  - Keep feature behavior, DOM rendering, and styling inside `src/features/post-detail/`.
- The feature should consume the existing Track A `getItem(id)` use-case result and must not move fetching logic into the view layer.

## Step 4 — `getItem(id)` Contract for TC-3
- Source:
  - `src/core/use-cases/get-item.js`
  - `tests/unit/core/get-item.test.js`
- Public contract:
  - Input: positive integer item ID.
  - Output: `Promise<{ ok: true, data: { item, comments } } | { ok: false, error: string }>`
- Verified behavior relevant to `TC-3`:
  - Success returns `data.item` as the root HN item and `data.comments` as an array.
  - Story items return `type: 'story'` and may include a sorted nested comments tree.
  - Job items return `type: 'job'` and an empty `comments` array when no discussion exists.
  - Invalid IDs and missing items return failed `Result` objects instead of throwing.
  - Comment nodes are already sorted newest-to-oldest per depth and carry `parent` references, but `TC-3` should not render them yet.
  - Poll `parts` are resolved for poll items, but that is future `TC-6` work and out of scope for the base detail page.
- Root item fields are optional except for `id` and `type`, so `TC-3` must tolerate missing `by`, `time`, `text`, `url`, `score`, and `title`.

## Step 5 — Proposed Controller API
- File: `src/features/post-detail/post-detail.controller.js`
- Public API:
  - `createPostDetailController({ getItem })`
  - Returns: `{ load(id) }`
- `load(id)` responsibilities:
  - Call Track A's `getItem(id)` use-case.
  - Reject unsupported item types for this ticket (`poll`, `pollopt`, `comment`).
  - Normalize the result into a small view model for the renderer.
  - Keep raw `text` unchanged so sanitization still happens in the view layer with DOMPurify.
- `load(id)` return shape:
  - Success: `{ ok: true, data: { item, comments, viewModel } }`
  - Failure: `{ ok: false, error, reason }`
- Failure reasons:
  - `invalid-item`
  - `not-found`
  - `unsupported-item`
  - `load-error`
  - `misconfigured-controller`
- `viewModel` fields:
  - `id`
  - `type`
  - `title`
  - `author`
  - `time`
  - `score`
  - `url`
  - `hasUrl`
  - `text`
  - `hasText`

## Step 6 — Render States
- File: `src/features/post-detail/post-detail.render-state.js`
- Explicit states:
  - `loading`
  - `success`
  - `error`
  - `not-found`
- Mapping rules:
  - Before `load(id)` resolves: `loading`
  - Successful controller result: `success`
  - `invalid-item`, `not-found`, and `unsupported-item` controller failures: `not-found`
  - `load-error` and `misconfigured-controller` controller failures: `error`
- State helpers:
  - `createLoadingPostDetailState()`
  - `createSuccessPostDetailState(data)`
  - `createErrorPostDetailState(message)`
  - `createNotFoundPostDetailState(message)`
  - `toPostDetailRenderState(loadResult)`

## Step 7 — Static View Skeleton
- File: `src/features/post-detail/post-detail.view.js`
- Public API:
  - `createPostDetailView({ onBack })`
  - Returns: `{ element, render(state), destroy() }`
- Always-present skeleton nodes:
  - back button
  - title
  - status message
  - metadata line
  - URL section wrapper
  - text section wrapper
- Stable `data-testid` hooks:
  - `post-detail-view`
  - `post-detail-back`
  - `post-detail-title`
  - `post-detail-status`
  - `post-detail-metadata`
  - `post-detail-url-section`
  - `post-detail-url-link`
  - `post-detail-text-section`
  - `post-detail-text-body`
- Render behavior:
  - `loading` fills the title and status placeholders and hides optional URL/text sections.
  - `success` fills title, metadata, optional URL, and optional text in the existing skeleton.
  - `not-found` reuses the skeleton with an empty-state message.
  - `error` reuses the skeleton with a generic failure message.

## Step 8 — Safe Text Rendering Rules
- Plain fields stay on safe text sinks:
  - title uses entity decoding plus `textContent`
  - metadata uses text-only formatting
  - URL label uses `textContent`
- API HTML handling:
  - only `viewModel.text` is sanitized with `DOMPurify`
  - sanitized output is appended directly as a `DocumentFragment`
  - the view does not mutate sanitized nodes after insertion
- URL safety:
  - outbound `href` values are restricted to `http` and `https`
  - invalid or unsafe URLs keep the URL section hidden

## Step 9 — Explicit Story vs Job Handling
- Controller normalization:
  - story items normalize through a dedicated `toStoryViewModel(...)` branch
  - job items normalize through a dedicated `toJobViewModel(...)` branch
  - both branches still tolerate missing `url`, missing `text`, and missing metadata fields
- View rendering:
  - success rendering now branches explicitly on `viewModel.type`
  - story rendering keeps the optional source link and optional story text path
  - job rendering keeps the optional link path but expects the body text to be the primary content when present
- Resilience rules:
  - no URL keeps the URL section hidden
  - no text keeps the text section hidden
  - unsupported types fall back to the error state instead of breaking the screen

## Step 10 — Back Navigation
- Current implementation:
  - the Back button uses a custom `onBack` callback when one is provided by the app shell
  - otherwise it falls back to `#/`, which is the router's stable feed route
- Why this is the smallest compatible version:
  - Track B has not landed feed-state restoration yet, so the detail view cannot reliably restore the exact active feed tab, page, or scroll position on its own
  - returning to `#/` still guarantees a stable path back to the feed surface without depending on unfinished shell state
- Future dependency:
  - restoring the exact previously active feed view depends on Track B app-shell/feed state wiring

## Step 11 — Detail Page Styling
- File: `src/features/post-detail/post-detail.css`
- Scope:
  - styling stays local to the post-detail feature and is imported from `post-detail.view.js`
  - layout is optimized for long-form reading rather than feed-card density
- Covered styling concerns:
  - readable page width and vertical rhythm
  - editorial title and byline presentation
  - safe, clearly interactive external link styling
  - padded rich-text body container for story/job text
  - responsive spacing and typography for tablet and mobile widths
- Visual direction:
  - follows the existing warm-accent, editorial shell instead of introducing a conflicting feature theme

## Step 12 — Route Wiring
- Files:
  - `src/main.js`
  - `tests/unit/main.test.js`
- Thin composition now lives in `main.js`:
  - build one shared router
  - build one post-detail controller from the HN adapter and `getItem` use-case
  - on `item` routes, mount the post-detail view and render `loading` immediately
  - when the controller resolves, map the result into a render state and repaint the mounted view
- Current non-item behavior:
  - `feed` and unknown routes fall back to the existing boot/feed placeholder content
  - this keeps route wiring minimal until Track B replaces the placeholder with the real feed shell
