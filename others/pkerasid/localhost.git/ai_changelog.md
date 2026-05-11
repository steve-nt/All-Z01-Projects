# AI Changelog

This file is required by the Engineering SOP (memo from C. Avramoudis, 2026-04-02,
"AI-Assisted Engineering and Spec-Driven Development"). Every AI-assisted change
must be reproducible from this log: the **spec section** that drove the change,
the **prompt/instruction** the agent received, and a short **rationale** for the
resulting diff. Reviewers audit this file alongside `PLAN.md` and the diff.

Conventions:

- One entry per phase or sub-task, newest at the bottom.
- Spec ref points into `PLAN.md` (phase + checkbox text).
- Prompt is a faithful summary of what the human asked the agent to do.
- Outcome lists the files touched and the gate that was met before sign-off.

## Quality gate (every phase must pass before the next starts)

- `cargo build` clean
- `cargo clippy --all-targets --all-features -- -D warnings` clean
- `cargo fmt --all -- --check` clean
- `cargo test` — all green
- The phase's PLAN.md checkboxes flipped to `[x]`

---

## 2026-05-06 — PLAN.md (project bootstrap spec)

- **Spec ref:** N/A (this entry creates the spec).
- **Prompt:** "Read the localhost README + audit README from 01-edu/public; produce a phased plan with checkboxes. Stack: Rust. Constraints: clean architecture, SOLID, DRY."
- **Rationale:** Per SOP Guideline 1, the plan is the source of truth. PLAN.md was authored before any code, capturing layered architecture (domain / application / infrastructure / interface), SOLID mapping, and 13 phases (0–12) with definition-of-done.
- **Outcome:** `PLAN.md` created. No code yet.

## 2026-05-07 — Phase 0: Bootstrap & tooling

- **Spec ref:** PLAN.md "Phase 0 — Bootstrap & Tooling" (6 checkboxes).
- **Prompt:** "Do phase 0."
- **Rationale:** Establish the quality gates SOP Guideline 3 mandates *before* any feature code: pinned toolchain, lint table denying `unwrap_used`/`expect_used`/`panic`/`dbg_macro`, fmt config, Makefile with `lint`/`siege`/`test` targets, lightweight stderr logger, library + binary skeleton.
- **Outcome:** `Cargo.toml`, `rust-toolchain.toml`, `rustfmt.toml`, `clippy.toml`, `Makefile`, `.gitignore`, `src/lib.rs`, `src/logger.rs`, `src/main.rs`. Gate green: build, clippy `-D warnings`, fmt, 1/1 test.

## 2026-05-07 — Phase 1: Domain model

- **Spec ref:** PLAN.md "Phase 1 — Domain Model (pure, no I/O)" (9 checkboxes).
- **Prompt:** "Do phase 1."
- **Rationale:** Inner layer of clean architecture. Validated constructors so invalid `Method`, `Status`, `HeaderName/Value`, `RequestPath`, `RouteConfig`, `HostConfig`, `ServerConfig`, `Cookie`, `SessionId` cannot exist at runtime — the rest of the system trusts them. Audit-relevant invariants (path traversal blocked, CRLF/NUL rejected in headers, duplicate listener detection, redirect must be 3xx) live here.
- **Outcome:** `src/domain/{error.rs, http/*, config/*, session/*}` (15 files). 64 unit tests. Gate green.

## 2026-05-07 — Phase 2: Configuration parser

- **Spec ref:** PLAN.md "Phase 2 — Configuration" (5 checkboxes).
- **Prompt:** "Go phase 2."
- **Rationale:** Boundary translator. TOML deserializes into thin `Raw*` shapes, then folds through the validated domain constructors — no validation logic duplicated. `ConfigError` unions IO/TOML/Domain/BadShape/BadAddress/BadStatusCode with `From` impls so `?` carries domain errors transparently. Five sample configs cover audit scenarios (single-port, multi-port, multi-host on same IP:port, CGI-enabled, upload+redirect).
- **Outcome:** `src/interface/{mod.rs, config_parser.rs}`, `config/{default,multi-port,multi-host,cgi,upload}.toml`, `src/main.rs` wired to argv config path. 17 new tests (one positive + one negative per validation rule). Gate green: clippy `-D warnings` clean, fmt clean, **81/81 tests**, both `default.toml` and `multi-host.toml` load and report listener resolution.

## 2026-05-07 — Spec amendment: drop bonus scope

- **Spec ref:** PLAN.md "Phase 11 — Bonus" (removed) and Phase 7 CGI checkbox.
- **Prompt:** "Remove phase 11 we don't want the bonus."
- **Rationale:** SOP Guideline 1 — fix the spec, not downstream code. The team is descoping bonus items (extra CGI runtimes, secondary-language server, basic auth, range requests). Phase 7's "Bonus: Perl, PHP, C++ binary" suffix was also removed since it implied scope we are no longer pursuing. Old "Phase 12 — Documentation & Submission" renumbered to Phase 11. Plan header now states explicitly: "required features only — no bonus items."
- **Outcome:** `PLAN.md` updated. No code changed. Gate not re-run (no source changes).

## 2026-05-07 — Spec amendment: split Phase 3 into 3a / 3b

- **Spec ref:** PLAN.md "Phase 3 — Reactor".
- **Prompt:** "Do phase 3."
- **Rationale:** SOP Guideline 3 — original Phase 3 bundled epoll wrapper + listener + connection state machine + event loop + timeouts (~1k+ LOC). That violates "small, logical chunks" and would be unreviewable. Split: **3a** = primitives (port traits, epoll wrapper, sockets, mock); **3b** = orchestration (listener, connection state machine, event loop, timer wheel). Spec now also notes that `Connection`'s parser body lands in Phase 4 — the state machine here will be a skeleton with raw byte I/O. Linux-only epoll module is cfg-gated; `MockReactor` makes upper layers testable on any platform.
- **Outcome:** PLAN.md restructured. No code yet.

## 2026-05-07 — Phase 3a: reactor primitives

- **Spec ref:** PLAN.md "Phase 3 — Reactor / 3a — Reactor primitives" (6 checkboxes).
- **Prompt:** "Do phase 3."  → executed Phase 3a only per the spec amendment above.
- **Rationale:** Establishes the dependency-inversion seam for I/O. `Reactor` and `Clock` traits live in `application/ports/`, so handlers (Phase 5+) and the event loop (Phase 3b) depend only on the trait. `EpollReactor` is the audit-targeted impl (Linux, edge-triggered, `EPOLLRDHUP`); `MockReactor` lets tests assert single-`poll`-per-iteration without OS state. `bind_listener` configures `O_NONBLOCK` + `SO_REUSEADDR`; `make_nonblocking` mirrors `accept4`'s flag inheritance for `std`'s `accept`. `setsockopt` size cast went through `try_from` to avoid `cast_possible_truncation`. The 4-bool `Event` struct carries a one-line `#[allow(struct_excessive_bools)]` with rationale: each bool maps 1:1 to an epoll bit.
- **Outcome:** New files: `src/application/{mod.rs, ports/{mod.rs, clock.rs, reactor.rs}}`, `src/infrastructure/{mod.rs, clock.rs, net/{mod.rs, socket.rs}, reactor/{mod.rs, epoll.rs, mock.rs}}`. `src/lib.rs` exposes `application` + `infrastructure`. **94/94 tests** (13 new in 3a). Gate green: `cargo build`, `clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`. Linux-gated tests in `epoll.rs` will run on the audit host.

## 2026-05-07 — Phase 3b: reactor orchestration

- **Spec ref:** PLAN.md "Phase 3 — Reactor / 3b — Reactor orchestration" (6 checkboxes).
- **Prompt:** "go 3b".
- **Rationale:** Closes the audit invariant that "the server uses only one select (or equivalent) to read the client requests and write answers": `EventLoop::tick` calls `Reactor::poll` exactly once per iteration, then dispatches both reads and writes from the events that one call returned (`MockReactor::poll_calls` test asserts the count). Listeners register `READABLE` only; on accept a new `Connection` is registered `READABLE`, then re-armed to `WRITABLE` when a response is queued. Skeleton response is a stub that recognises `\r\n\r\n` — Phase 4 replaces it with the parser/router/handler pipeline; the four state-machine variants (`ReadingHeaders`/`ReadingBody`/`Dispatching`/`WritingResponse`) are encoded now so Phase 4 only fills bodies. Backpressure: `READ_BUF_CAP` caps inbound buffering, `on_writable` returns `Rearm(WRITABLE)` on a partial write so the unwritten tail is re-driven by the next epoll event. Idle eviction runs at the end of every tick using the `Clock` port (`FixedClock` in tests, `SystemClock` in prod). Socket errors / `EPOLLERR` / `EPOLLHUP` funnel through `close_connection`, which deregisters then drops the stream — no fd leaks. `main.rs` is `cfg(target_os = "linux")`-gated so non-Linux dev hosts still build the library + run unit tests via `MockReactor` while the binary cleanly errors out.
- **Outcome:** New files: `src/application/{connection.rs, event_loop.rs}`. Edits: `src/application/mod.rs`, `src/main.rs` (event-loop wiring). **108/108 tests** (14 new in 3b: 7 connection state-machine, 7 event-loop). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.

## 2026-05-07 — Phase 4: HTTP parser & writer

- **Spec ref:** PLAN.md "Phase 4 — HTTP Parser & Writer" (7 checkboxes).
- **Prompt:** "do phase 4".
- **Rationale:** `interface/http_parser.rs` is the request-side boundary translator: a streaming state machine (`RequestLine -> Headers -> Body -> Done`) fed from the connection's `read_buf` slice; bytes consumed are reported so the caller drains exactly what the parser absorbed. Body modes cover the audit's three required framings — none, `Content-Length`, and `Transfer-Encoding: chunked` (with extension and trailer tolerance, both discarded per RFC 7230 §4.1.2). Limits are encoded as `ParserLimits` (`max_request_line` -> 414, `max_headers_size` -> 400, `max_body_size` -> 413); `POST` without framing -> 411; obs-fold rejected; non-`chunked` Transfer-Encoding -> 501; unknown HTTP version -> 505; everything else malformed -> 400. The parser never panics — every fallible path returns `ParseError`. `interface/http_writer.rs` is the response-side translator: it auto-supplies `Date` (RFC 7231 IMF-fixdate computed in-house via Howard Hinnant's `civil_from_days`, no `chrono` dep), `Server: localhost/<VERSION>`, `Content-Length` (suppressed when caller sets `Transfer-Encoding`), and `Connection: keep-alive|close` based on the server's decision (caller's `Connection` header is overridden — server has the final say). `application/connection.rs` was rewired: read path feeds the parser; on `Complete` we produce a stub 200 (Phase 5 plugs in the router); on `ParseError` we produce a status-mapped error response with `Connection: close`. Keep-alive decision follows RFC 7230: HTTP/1.1 default is keep-alive unless `Connection: close`; HTTP/1.0 default is close unless `Connection: keep-alive`. On full flush with keep-alive we reset the parser and re-arm `READABLE` (the `read_buf` is preserved so a pipelined request would already be in hand). `event_loop.rs` test had to send `Connection: close` because the new keep-alive behavior keeps connections around past the response — that's the intended HTTP/1.1 default.
- **Outcome:** New files: `src/interface/{http_parser.rs, http_writer.rs}`. Edits: `src/interface/mod.rs`, `src/application/connection.rs` (full rewrite of state machine to drive the parser/writer), one event-loop test updated for keep-alive default. **135/135 tests** (27 new in Phase 4: 19 parser, 7 writer, 1 reading-body state). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.

## 2026-05-07 — Phase 5: Routing & Request Pipeline

- **Spec ref:** PLAN.md "Phase 5 — Routing & Request Pipeline" (5 checkboxes).
- **Prompt:** "We have this project with my teammates for zone01 athens. We have made the plan and the changelog. I want you to read the code and help me do the phases 5 & 6."
- **Rationale:** Phase 5 wires the parsed `Request` to a `Response` through a three-step pipeline: host resolution → route selection → dispatch. `application/router.rs` implements `resolve_host` (named `server_name` match first, unnamed-default fallback, then first host) and `longest_prefix_match` with `path_has_prefix` to prevent `/api` spuriously matching `/apiv2`. `application/request_pipeline.rs` introduces `PipelineContext { config: Rc<ServerConfig>, fs: Rc<dyn FileSystem> }` — `Rc` (not `Arc`) because the server is intentionally single-threaded. `PipelineContext` required a manual `Debug` impl because `Rc<dyn FileSystem>` is not `Debug`. `application/error_pages.rs` tries the host's configured custom error-page path via `FileSystem`, falls back to a built-in HTML template; the let-chain idiom (`if let Some(h) = … && let Some(p) = …`) keeps it lint-clean. `connection.rs` gained `local_addr: SocketAddr` + `pipeline: Option<Rc<PipelineContext>>` and a `with_pipeline` constructor; backward-compat `Connection::new` keeps all 135 pre-existing tests passing in stub mode. `event_loop.rs` gained `listener_addrs: HashMap<Token, SocketAddr>`, records `listener.local_addr()` on bind, and uses `Connection::with_pipeline` when both are present.
- **Outcome:** New files: `src/application/{router.rs, request_pipeline.rs, error_pages.rs, ports/filesystem.rs}`. Edits: `src/application/{mod.rs, connection.rs, event_loop.rs, ports/mod.rs}`, `src/infrastructure/mod.rs`. **152/152 tests** (17 new: 8 router, 7 pipeline, 2 error-pages). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.

## 2026-05-07 — Phase 6: Static Handlers

- **Spec ref:** PLAN.md "Phase 6 — Static Handlers" (6 checkboxes).
- **Prompt:** Same session as Phase 5 — phases 5 and 6 were implemented together.
- **Rationale:** All handlers receive a `&dyn FileSystem` and are tested exclusively against `FakeFileSystem`, satisfying the "tests with a fake `FileSystem` impl" checkbox. `format_http_date` was extracted from `interface/http_writer.rs` to `domain/http/date.rs` as a pure `u64 -> String` function; this avoids an application → interface dependency violation while letting both `http_writer` and `static_file` reuse it. `static_file::handle` stats the resolved path (constructed by `resolve_path` which strips the route prefix then joins with `root`): if it is a file, it streams bytes with `Content-Type` from a 24-extension MIME table, `Last-Modified`, and `Content-Length`; if it is a directory, it tries the configured index file, then generates an HTML autoindex listing (HTML-escaped via `push_str` + `writeln!` to satisfy `format_push_string` and `write_with_newline` lint rules), then 403. `redirect::handle` sets `Location`. `delete::handle` stats first (404 if missing, 403 if directory), then removes. `upload::handle` detects `multipart/form-data` from the `Content-Type` header, parses the first part's `Content-Disposition` for a filename, and writes the part body; raw (non-multipart) mode uses the last URL path segment. `sanitize_filename` strips path separators and NUL bytes. `main.rs` wires `OsFileSystem` (new `infrastructure/fs/mod.rs`) into a `PipelineContext` and passes it to `EventLoop::with_pipeline`.
- **Outcome:** New files: `src/domain/http/date.rs`, `src/application/handlers/{mod.rs, static_file.rs, delete.rs, upload.rs, redirect.rs}`, `src/infrastructure/fs/mod.rs`. Edits: `src/domain/http/mod.rs`, `src/interface/http_writer.rs`, `src/application/mod.rs`, `src/main.rs`, `src/infrastructure/reactor/epoll.rs` (targeted `#[allow]` for pre-existing FFI casts). **170/170 tests** (18 new: 5 static-file, 4 delete, 5 upload, 2 pipeline, 2 redirect). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.

## 2026-05-07 — Bootstrap www/ test content

- **Spec ref:** PLAN.md "Phase 10 — Integration & Audit Scenarios" (browser smoke test checkbox).
- **Prompt:** Server returned nothing at localhost:8080 — investigated and created minimal www/ content to verify the server responds end-to-end.
- **Rationale:** `config/default.toml` points `root = "www"` and `error_pages` at `www/errors/`. Neither directory existed, so every request returned an empty error response. Created the minimum needed to confirm the server is wired correctly: `www/index.html` (smoke test), `www/errors/404.html`, `www/errors/500.html`. This is **not** the final audit demo content — a teammate must expand `www/` before Phase 10 to cover every audited feature (autoindex, upload, redirect, CGI, sessions). See the new checkbox added to Phase 10 in PLAN.md.
- **Outcome:** New files: `www/index.html`, `www/errors/404.html`, `www/errors/500.html`. No code changes. Server confirmed serving HTML at `http://127.0.0.1:8080/`.

## 2026-05-07 — Phase 7 (partial): CGI handler, process runner, and tests

- **Spec ref:** PLAN.md "Phase 7 — CGI" (10 checkboxes).
- **Prompt:** "Read changelog and plan.md; do phase 7; commit one by one; update md for the next step."
- **Rationale:** Added a dedicated `ProcessRunner` port so application logic depends on an interface rather than process APIs. Implemented `OsProcessRunner` in infrastructure and a new `application::handlers::cgi` path that builds CGI env vars (`REQUEST_METHOD`, `QUERY_STRING`, `SCRIPT_FILENAME`, `PATH_INFO`, `SERVER_*`, `HTTP_*`), sets CWD to script dir, streams request body to stdin, enforces timeout with child kill/reap, parses CGI `Status`/headers/body, and maps non-zero exit to 502. `RequestPipeline` now dispatches to CGI when route has `cgi` config and path extension matches. Added unit tests with fake runner plus a real Python script execution test through the OS runner.
- **Outcome:** New files: `src/application/ports/process_runner.rs`, `src/application/handlers/cgi.rs`, `src/infrastructure/cgi/mod.rs`. Edits: `src/application/request_pipeline.rs`, `src/main.rs`, `src/infrastructure/mod.rs`, plus supporting lint-cleanups and status 504 constant in `src/domain/http/status.rs`. Gate green: `cargo fmt --all`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test` (**175/175**).
- **Next step:** Finish remaining Phase 7 audit parity by wiring CGI stdin/stdout pipes into the epoll loop and replacing the current process spawn path with explicit libc-level `fork`/`execve` integration.

## 2026-05-07 — Phase 7 completion: libc fork/execve + epoll CGI pipes

- **Spec ref:** PLAN.md "Phase 7 — CGI" (remaining unchecked items: `fork/execve` and epoll-registered pipes).
- **Prompt:** "Yes, I want Phase 7 fully closed."
- **Rationale:** Replaced the temporary `std::process::Command` runner with a Linux syscall-based implementation in `infrastructure/cgi/mod.rs`: `pipe2(O_NONBLOCK|O_CLOEXEC)` for stdin/stdout/stderr, `fork`, child-side `dup2`/`chdir`/`execve`, parent-side epoll loop over CGI pipes (`EPOLLIN`/`EPOLLOUT`/`EPOLLHUP`) to stream request body in and collect stdout/stderr out without threads. Child lifecycle now uses `waitpid(..., WNOHANG)` on each loop tick and timeout enforcement via `SIGKILL` + reap, matching the race-safe requirement.
- **Outcome:** `src/infrastructure/cgi/mod.rs` rewritten around libc + epoll. Existing CGI unit/integration tests remain green, and full quality gate is green: `cargo fmt --all`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test` (**175/175**).
- **Next step:** Start **Phase 8 — Sessions & Cookies** (`SessionStore`, cookie parsing/serialization, pipeline session attachment, demo counter route).

## 2026-05-07 — Phase 8: Sessions & Cookies

- **Spec ref:** PLAN.md "Phase 8 — Sessions & Cookies" (4 checkboxes).
- **Prompt:** "I made a git pull because a teammate pushed commits. Start implementing Phase 8, one thing at a time; commit then push, then move on."
- **Rationale:** Domain types (`Session`, `SessionId`, `Cookie`, `parse_cookie_header`) were already in place from Phase 1. Phase 8 added the remaining layers. **`SessionStore` port trait** (`application/ports/session_store.rs`) — object-safe interface with `get`/`put`/`remove`/`evict_expired`; `get` returns a clone so callers can mutate and `put` back without holding a borrow on the store. A `FakeSessionStore` ships alongside for unit tests. **In-memory store** (`infrastructure/session_store/memory.rs`) — `HashMap<String, Session>` with TTL eviction via `Session::is_expired`. **Session ID generator** (`infrastructure/session_store/id_gen.rs`) — reads 32 bytes from `/dev/urandom`, hex-encodes to a 64-char string valid as a `SessionId`. **`RouteKind::SessionCounter`** added to the domain route config and config parser (`session_counter = true` in TOML); a dedicated `session_counter` handler builds an HTML visit-count page given a `u64` visits argument. **Session middleware** wired into `request_pipeline::handle`: before dispatch it reads the `SID` cookie, loads an existing session or creates a new one (with `evict_expired` called opportunistically on each request), mutates the session in the `SessionCounter` arm, saves it back via `store_rc.borrow_mut().put(session)`, and injects a `Set-Cookie: SID=…; Path=/; HttpOnly; SameSite=Lax` header only when the session is new. `PipelineContext` gained `session_store: Option<Rc<RefCell<dyn SessionStore>>>` — optional so all prior tests work unchanged with `session_store: None`. `main.rs` constructs a `MemorySessionStore` and passes it in. `Response::add_header` was added to the domain type to let the middleware append a header without rebuilding the response. `config/default.toml` gained a `/session` demo route.
- **Outcome:** New files: `src/application/ports/session_store.rs`, `src/application/handlers/session_counter.rs`, `src/infrastructure/session_store/{mod.rs, memory.rs, id_gen.rs}`. Edits: `src/application/ports/mod.rs`, `src/application/handlers/mod.rs`, `src/application/request_pipeline.rs`, `src/domain/config/route.rs`, `src/domain/http/response.rs`, `src/interface/config_parser.rs`, `src/infrastructure/mod.rs`, `src/main.rs`, `config/default.toml`, `PLAN.md`. **192/192 tests** (12 new: 5 port fake, 6 infra store+id-gen, 2 pipeline session counter, 2 config parser). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.

## 2026-05-07 — Phase 9: Error Pages & Resilience

- **Spec ref:** PLAN.md "Phase 9 — Error Pages & Resilience" (5 checkboxes).
- **Prompt:** "Start implementing Phase 9."
- **Rationale:** Three items needed real work; two were already satisfied. **Already done:** per-server error-page overrides with fallback (implemented in Phase 5 via `error_pages::error_response`); hot-path `unwrap`/`expect`/`panic` audit found all such calls are inside `#[cfg(test)]` blocks or are infallible `unwrap_or` variants. **Improved error page templates:** `default_body` rewritten with a proper HTML5 skeleton (`<meta charset>`, inline CSS, per-status description text for the 11 most common error codes). Explicit tests added for all 7 audit-required codes (400/403/404/405/408/413/500) plus an HTML5-structure check — 8 new tests total. **408 on idle timeout:** Previously `evict_idle` silently dropped idle connections. Now `Connection::start_timeout(now_ms)` is called first: for connections in `ReadingHeaders` or `ReadingBody` state it queues a `408 Request Timeout` response, sets `keep_alive = false`, advances `last_activity_ms` to prevent immediate re-eviction, and returns `Rearm(WRITABLE)`; the event loop reregisters the fd for WRITABLE so the response is flushed on the next reactor event, after which `on_writable` returns `Close`. For any other state (connection was already responding) it returns `Close` directly. The existing `idle_connections_are_evicted` test was updated to drive the extra WRITABLE tick. **Bounded buffers:** `READ_BUF_CAP = 16 KiB` (enforced in `on_readable`) and one-response-at-a-time `write_buf` already bounded memory per connection; confirmed no unbounded growth path on the hot path.
- **Outcome:** Edits: `src/application/error_pages.rs`, `src/application/connection.rs`, `src/application/event_loop.rs`, `PLAN.md`. **202/202 tests** (10 new: 8 error-page coverage, 2 connection timeout). Gate green: `cargo build`, `cargo clippy --all-targets -- -D warnings`, `cargo fmt --check`, `cargo test`.
