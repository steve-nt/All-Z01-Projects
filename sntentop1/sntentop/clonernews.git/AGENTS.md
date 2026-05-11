# Agent Instructions: Modern JavaScript (2026) Guidelines

## Core Directives
- **Action-Oriented**: Focus purely on writing clean, idiomatic code and performing rigorous robust static analysis.
-  **Bug Fixing**: When I report a bug, dont try to fix it immediately. Instead, try to understand the root cause of the bug. Then write a test that reproduces the bug. Then have subagents try to fix the bug and prove it with passing test. 
- **Explanations**: Always have single-line comments for every non-trivial line of code. Comments should explain the "why" and "how", not the "what". Also in the begining of each file, have a comment block that explains the file's purpose, its public API (if applicable), and any important implementation notes or constraints.

## Modern JS & Idioms (ES2026 Standard)
- **Immutability by Default**: 
  - ALWAYS use immutable array methods (`toSorted()`, `toReversed()`, `toSpliced()`, `with()`) instead of mutating counterparts (`sort()`, `reverse()`, `splice()`).
  - Rely on destructuring and the spread operator (`...`) instead of direct object/array mutation.
- **Modern APIs & Primitives**:
  - **Temporal API**: Exhaustively use `Temporal` for accurate date/time tracking. The legacy `Date` object is strictly forbidden.
  - **Sets**: Utilize native ES2025 Set operations (`union()`, `intersection()`, `difference()`, `isSubsetOf()`, etc.).
  - **Grouping**: Leverage `Object.groupBy()` and `Map.groupBy()` for grouping data iteratively.
- **Promises & Error Flow**: 
  - Use `async/await` exclusively with top-level await where applicable. 
  - Rely on `Promise.withResolvers()` and `Promise.try()` for unified sync/async wrapping.
  - **Explicit Result Pattern**: For business logic rejections, prefer returning a Result object (e.g., `{ ok: true, data } | { ok: false, error }`) instead of throwing exceptions. Reserve `throw` for truly exceptional, unrecoverable system failures.
- **Resource Management**: Exhaustively use the **`using`** and **`await using`** declarations for deterministic cleanup of resources (file handles, event listeners, streams). Avoid manual `try...finally` cleanup.
- **Variables & Functions**:
  - `const` is the default. Only use `let` when variable reassignment is essential.
  - Prefer arrow functions for standard logic and to preserve lexical `this`.

## Code Organization & Architecture (2026)
- **Screaming Architecture**: The directory structure MUST communicate the application's domain and purpose, not the technical framework used.
  - Prefer a **Feature-First** structure (e.g., `src/features/billing/`) over a Type-First structure (e.g., `src/components/`).
  - **Colocation**: Keep logic, UI components, styles, types, and tests together within the feature folder they serve.
- **clonernews Project Mapping**:
  - `src/core/` — Pure domain: entities, interfaces (JSDoc interfaces), use-cases. **Zero DOM, zero `fetch`, zero side effects.**
  - `src/infra/` — Infrastructure adapters: `hn-api-adapter.js`, `cache-adapter.js`, `throttle.js`. I/O lives here.
  - `src/features/{feed,post-detail,live-banner,polls}/` — Feature-first UI slices with colocated views, controllers, styles, and tests.
  - `src/shared/` — Cross-feature utilities: router, DOM helpers, signals, time formatting.
- **Layered Decoupling (Clean Architecture)**:
  - **Core/Domain**: This layer MUST be pure JavaScript, housing business rules and entities. It has zero knowledge of the database, UI, or external APIs.
  - **Application/Use Cases**: Orchestrates flow. It depends on abstractions (interfaces) and not concrete implementations.
  - **Infrastructure/Adapters**: Implements the technical details (e.g., `fetch` calls, `Temporal` formatting for UI, storage). Use Dependency Inversion to provide these to the application layer.
- **Encapsulated Modules**: 
  - Use `package.json` `exports` to define strict public boundaries. Prevent unauthorized deep-imports into a feature's internal modules.
  - Maintain a strict hierarchy: Feature A should never reach into Feature B's `internal/` folder.
- **Functional Core, Imperative Shell**: Centralize business logic in pure, deterministic functions; push side effects (I/O) to the outermost edges of the architecture.
- **Fine-Grained Reactivity (Signals)**: For state management, prefer **Signals** (standardized observer primitives) over heavy global state containers. Ensure UI updates are targeted and granular.

## Code Quality & Static Analysis (2026)
- **State-of-the-Art Tooling**: Use **Biome** as the standard for ultra-fast unified linting and formatting (or ESLint v10 Flat Config). Actively perform static analysis to detect code smells.
- **Type-Aware JS**: Provide rich JSDocs for type-aware linting in vanilla JavaScript for maximum safety.
- **Code Health**: Prevent "AI slop", over-engineering, and boilerplate code. Write tightly scoped functions and observe complexity limits.
- When a ticket is complete, use [.github/pull_request_template.md](.github/pull_request_template.md) for the PR description and fill every required section before requesting review.

## clonernews: HTML Sanitization Policy
- The HN API returns HTML in `text` fields (story text, comment text, job descriptions, poll text).
- When rendering API-sourced HTML, **sanitize with DOMPurify** before DOM insertion. Never modify the sanitized output afterward.
- For plain-text fields (`by`, `url`, `title`), use `textContent` exclusively.
- The `title` field may contain HTML entities — decode safely, never via `innerHTML`.

## clonernews: API Client Discipline
- All HN API live-data polling **must** use a minimum **5-second throttle** interval.
- All item fetches **must** check the in-memory cache before hitting the network.
- Batch-fetch concurrency **must** be bounded (max 6 concurrent requests).
- All fetch calls **must** use `AbortController` with a timeout (8s default).
- Return `Result` objects (`{ ok, data }` or `{ ok, error }`) — never throw on API failures.

## clonernews: DOM Rendering Contract
- DOM manipulation is **only permitted** inside `src/features/*/` view files and `src/shared/dom-helpers.js`.
- `src/core/` and `src/infra/` **must never** reference `document`, `window`, or any DOM API.
- All DOM construction **must** use `createElement` / `textContent` / `appendChild` (safe sinks).
- Use `DocumentFragment` for batch insertions. Use event delegation over per-element listeners.

## JavaScript Security Guidelines (2026)
- **XSS Prevention by Construction**:
  - Prefer safe DOM sinks (`textContent`, `createTextNode`, `setAttribute` with hardcoded safe attribute names).
  - Treat `innerHTML`, `insertAdjacentHTML`, `outerHTML`, `document.write`, and string-to-code APIs as high-risk injection sinks.
  - If HTML rendering is required, sanitize with a maintained sanitizer (for example, DOMPurify) and never mutate sanitized markup afterward.
  - Never put untrusted values into event-handler attributes, inline scripts/styles, or `javascript:` URLs.
- **Token & Session Handling**:
  - Prefer server-managed sessions in `HttpOnly` + `Secure` + `SameSite` cookies for authentication state.
  - If bearer tokens must be used client-side, keep TTL short, rotate regularly, scope audience/claims minimally, and validate `iss`, `aud`, `exp`, and algorithm expectations server-side.
  - Do not place long-lived secrets in `localStorage`; if temporary browser storage is unavoidable, prefer shorter-lived strategies and strict CSP hardening.
  - Renew session identifiers after authentication or privilege changes and enforce server-side invalidation on logout/expiration.
- **CSP and Trusted Types Enforcement**:
  - Use a strict CSP (`script-src` nonce/hash-based, avoid `unsafe-inline` and `unsafe-eval`, set `object-src 'none'`, `base-uri 'none'`).
  - Use `frame-ancestors` to prevent untrusted embedding/clickjacking where appropriate.
  - Enforce Trusted Types with `require-trusted-types-for 'script'` and explicit `trusted-types` policy allowlists.
  - Roll out policy changes with report-only mode first, then enforce.
- **Dependency and Supply-Chain Controls**:
  - Pin and enforce lockfiles (`npm ci` in CI).
  - Continuously audit dependencies and transitive dependencies; block known exploited/high-severity issues before release.
  - Minimize install-time script risk (prefer `--ignore-scripts` or allowlisted lifecycle scripts where feasible).
  - Use provenance/SBOM/signing workflows for release artifacts and prefer trusted publishers with short-lived credentials.
  - Enforce 2FA and least privilege for package publishing and CI tokens.
- **Secure Network and Fetch Handling**:
  - Enforce HTTPS-only endpoints; never send credentials or tokens over plaintext transport.
  - Use explicit `fetch` options (`method`, `mode`, `credentials`, `cache`, `redirect`) and fail closed on unexpected response status/content-type.
  - Bound network operations with timeouts/cancellation (`AbortController`) and implement retry with backoff only for idempotent operations.
  - Validate response shape before use; reject ambiguous or mixed-content responses.
- **Safe DOM and Browser API Usage**:
  - Build UI trees with `createElement`/`createElementNS` and explicit property assignment.
  - Avoid dynamic code execution (`eval`, `new Function`, string arguments to `setTimeout`/`setInterval`).
  - Keep third-party script inclusion minimal and integrity-protected where applicable.
- **Logging and Privacy Hygiene**:
  - Never log credentials, tokens, session identifiers, auth headers, or PII.
  - Redact sensitive fields before logging and prefer structured logs with minimal, purpose-bound fields.
  - Use correlation IDs instead of user secrets for traceability.
  - Treat client telemetry as potentially sensitive; collect the least data necessary and retain for the shortest practical duration.

## Legacy Anti-Patterns (NEVER USE)
- **`var` is completely forbidden**.
- **`Date` constructor is forbidden**: Use `Temporal` exclusively. No `new Date()`, `Date.now()`, or `Date.parse()`.
- **CommonJS is obsolete**: Exclusively use ES Modules (`import`/`export`). No `require`.
- **Frameworks are forbidden**: No React, Vue, Angular, Svelte, Phaser, jQuery, or similar. This is a vanilla JavaScript project.
- **Legacy objects & properties**: Avoid the `arguments` object (use rest parameter `...args`), avoid `_private` naming conventions (use `#private` class fields), and do not use `XMLHttpRequest` (use `fetch`).
- **Forbidden Patterns**: NEVER modify built-in prototypes. Avoid "clever" one-liners that sacrifice readability for terseness.