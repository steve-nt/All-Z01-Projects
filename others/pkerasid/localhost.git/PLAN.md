# Localhost — HTTP/1.1 Server in Rust

A single-process, single-thread, non-blocking HTTP/1.1 server using `epoll` (via `libc`).
Architecture: **Clean Architecture** (Domain / Application / Infrastructure / Interface)
Principles: **SOLID**, **DRY**, small composable modules, no panics on the hot path.
Scope: required features only — no bonus items.

---

## Architectural Overview

```
crates/
└── localhost/
    ├── src/
    │   ├── main.rs                     # composition root
    │   ├── domain/                     # pure types, no I/O
    │   │   ├── http/                   # Request, Response, Method, Status, Headers, Version
    │   │   ├── config/                 # ServerConfig, RouteConfig, CgiConfig (validated)
    │   │   ├── session/                # Session, Cookie, SessionId
    │   │   └── error.rs                # DomainError
    │   ├── application/                # use cases, orchestration (pure where possible)
    │   │   ├── router.rs               # match host+port+path -> Route
    │   │   ├── request_pipeline.rs     # parse -> validate -> dispatch -> respond
    │   │   ├── handlers/               # static_file, upload, delete, redirect, listing, cgi
    │   │   ├── error_pages.rs          # status -> body resolution
    │   │   └── ports.rs                # traits: FileSystem, Clock, ProcessRunner, SessionStore
    │   ├── infrastructure/             # syscalls, OS, real I/O
    │   │   ├── net/                    # listener, connection, non-blocking sockets
    │   │   ├── reactor/                # epoll wrapper, event loop, timers
    │   │   ├── fs/                     # std::fs adapter implementing FileSystem
    │   │   ├── cgi/                    # fork/exec, pipes, env (PATH_INFO, etc.)
    │   │   └── session_store/          # in-memory store implementing SessionStore
    │   ├── interface/                  # parsers/serializers (boundary)
    │   │   ├── http_parser.rs          # bytes -> Request (chunked + unchunked)
    │   │   ├── http_writer.rs          # Response -> bytes
    │   │   └── config_parser.rs        # TOML/JSON -> ServerConfig
    │   └── lib.rs
    ├── tests/                          # integration tests
    └── config/
        └── default.toml                # sample configs
```

**SOLID mapping**
- **S**: each module owns one concern (parser ≠ router ≠ reactor).
- **O**: handlers are trait objects behind `Handler`; new behaviors plug in without edits.
- **L**: `FileSystem`, `SessionStore`, `Clock`, `ProcessRunner` traits with swappable impls (real / fake for tests).
- **I**: small, role-based traits (no fat interfaces).
- **D**: `application/` depends on traits in `application/ports.rs`; `infrastructure/` implements them.

**DRY**: header serialization, error-page rendering, MIME guessing, and path resolution each live in exactly one place.

---

## Phase 0 — Bootstrap & Tooling
- [x] `cargo init` workspace; pin Rust toolchain (`rust-toolchain.toml`).
- [x] Add deps: `libc` only (runtime). Dev: `tempfile`, `assert_matches`, `serde`/`serde_derive` for config (allowed — not a server crate).
- [x] Add `clippy` + `rustfmt` configs; CI-friendly `cargo fmt --check && cargo clippy -- -D warnings`.
- [x] Add `Makefile`: `build`, `run`, `test`, `siege`, `lint`.
- [x] Set up `tracing`-free lightweight logger (or stderr macro) — no async runtimes.
- [x] Skeleton crate compiles with empty `main`.

## Phase 1 — Domain Model (pure, no I/O)
- [x] `Method` enum (GET/POST/DELETE + Other).
- [x] `Status` (with canonical reason phrases for 200, 201, 204, 301, 302, 400, 403, 404, 405, 408, 411, 413, 500, 501, 505).
- [x] `Headers` map (case-insensitive keys, preserves insertion order, single source of truth for header ops).
- [x] `Request` / `Response` value types (immutable builders).
- [x] `Url` / `Path` value object (normalizes `..`, percent-decoding, prevents traversal).
- [x] `ServerConfig`, `HostConfig`, `RouteConfig`, `CgiConfig` with **validated constructors** (invalid configs cannot exist).
- [x] `Session`, `Cookie`, `SessionId` types.
- [x] `DomainError` (no `unwrap` outside tests).
- [x] Unit tests for every domain type (round-trip parse/serialize where applicable).

## Phase 2 — Configuration
- [x] Define config schema (TOML preferred for readability).
- [x] `interface/config_parser.rs`: file -> raw -> validated `ServerConfig`.
- [x] Support: multiple `[[server]]`, `listen` (host:port, multi), `server_name`, `client_max_body_size`, `error_pages`, `[[server.route]]` with `path`, `methods`, `root`, `index`, `autoindex`, `redirect`, `upload_dir`, `cgi { extension, interpreter, working_dir }`.
- [x] Validation rules: duplicate listener detection across servers (reject), at least one route, sane size limits.
- [x] Sample configs in `config/`: minimal, multi-port, multi-host (same IP:port), CGI-enabled, upload-enabled.
- [x] Tests: every validation rule has a positive + negative case.

## Phase 3 — Reactor (epoll, single-threaded)
Split into 3a (primitives) and 3b (orchestration) per SOP Guideline 3 (incremental gates).
Linux-only epoll module is `cfg`-gated so the crate still builds on non-Linux dev hosts;
a `MockReactor` makes the application logic testable on any platform.

### 3a — Reactor primitives
- [x] `Reactor` port trait + `Token`, `Interest`, `Event` types in `application/ports/reactor.rs`.
- [x] `Clock` port trait in `application/ports/clock.rs`; `SystemClock` impl in infrastructure.
- [x] Non-blocking socket helpers (`O_NONBLOCK`, `SO_REUSEADDR`).
- [x] `infrastructure/reactor/epoll.rs`: safe wrapper around `epoll_create1`/`epoll_ctl`/`epoll_wait` (cfg = linux).
- [x] `MockReactor` for cross-platform unit tests of upper layers.
- [x] Tests for each.

### 3b — Reactor orchestration
- [x] `Listener` registers `READ`; on accept, spawns `Connection` registered `READ`.
- [x] `Connection` state-machine skeleton: `ReadingHeaders -> ReadingBody -> Dispatching -> WritingResponse -> (KeepAlive | Close)` (parser body filled in Phase 4).
- [x] **Single `epoll_wait` per loop iteration** drives both reads and writes (audit requirement).
- [x] Per-connection timeout (configurable, e.g. 30s); idle connections evicted.
- [x] Backpressure: bounded read/write buffers; partial writes re-arm `EPOLLOUT`.
- [x] On any socket error, connection is removed cleanly (no leaks).

## Phase 4 — HTTP Parser & Writer
- [x] Streaming request parser: request-line, headers, body.
- [x] Body modes: `Content-Length`, `Transfer-Encoding: chunked`, none.
- [x] Strict limits: max header size, max line length, max body (per route `client_max_body_size` -> 413).
- [x] Malformed input -> 400 (never panic).
- [x] Response writer: status line, headers (Date, Server, Content-Length/Transfer-Encoding, Connection), body.
- [x] Connection: keep-alive support for HTTP/1.1, `Connection: close` honored.
- [x] Tests: golden-file tests for parser; fuzz-ish corpus of malformed requests.

## Phase 5 — Routing & Request Pipeline
- [x] `Router`: (host, port, path) -> `RouteConfig` with longest-prefix match; `server_name` fallback to default server on the listener.
- [x] `RequestPipeline`: parse -> validate (method allowed? size? path safe?) -> resolve handler -> produce `Response`.
- [x] Method-not-allowed produces 405 with `Allow` header.
- [x] All errors funnel through `error_pages.rs` (custom path if configured, otherwise built-in template).
- [x] Tests: matrix of (host, path, method) -> expected route/status.

## Phase 6 — Static Handlers
- [x] `StaticFileHandler`: resolves under `root`, blocks traversal, serves with correct `Content-Type` (MIME table), `Last-Modified`, `Content-Length`.
- [x] `DirectoryHandler`: serves `index` if present; else `autoindex` HTML listing if enabled; else 403.
- [x] `RedirectHandler`: 301/302 from config.
- [x] `DeleteHandler`: deletes file under route root if method allowed; 204 on success.
- [x] `UploadHandler` (POST): writes body to `upload_dir`; supports multipart and raw; respects body limit.
- [x] Tests with a fake `FileSystem` impl.

## Phase 7 — CGI
- [x] `CgiHandler`: `fork` + `execve` via `libc`, two pipes (stdin = body, stdout = response).
- [x] Env: `REQUEST_METHOD`, `CONTENT_LENGTH`, `CONTENT_TYPE`, `QUERY_STRING`, `SCRIPT_FILENAME`, `SCRIPT_NAME`, `PATH_INFO` (full path), `SERVER_PROTOCOL`, `SERVER_NAME`, `SERVER_PORT`, `HTTP_*` headers.
- [x] **Working dir set to script dir** for correct relative paths.
- [x] Pipes registered on epoll (still single-threaded).
- [x] Body streamed in; EOF on stdin signals end-of-request.
- [x] Parse CGI response headers (`Status:`, `Content-Type:`) -> HTTP response.
- [x] Per-CGI timeout; child reaped without `SIGCHLD` races (use `waitpid` with `WNOHANG`).
- [x] Support both chunked + unchunked input bodies.
- [x] At least one interpreter wired (Python).
- [x] Tests: fake `ProcessRunner` for unit tests + real script integration test.

## Phase 8 — Sessions & Cookies
- [x] `SessionStore` trait; in-memory impl (HashMap with TTL eviction on access).
- [x] `Cookie` parser/serializer (Set-Cookie incl. `HttpOnly`, `Path`, `Max-Age`).
- [x] `SessionMiddleware`-style step in pipeline: read `Cookie: SID=`, attach session, set on response if new.
- [x] Demo route exercising session counter for the audit (`/session` in default.toml).

## Phase 9 — Error Pages & Resilience
- [x] Default templated pages for 400/403/404/405/408/413/500.
- [x] Per-server overrides from config; missing override file falls back to default (never 500-loops).
- [x] All `Result` paths reviewed: no `unwrap`/`expect` on hot path; panics replaced by 500.
- [x] Global timeout: idle connections receive a 408 response then are closed.
- [x] Memory: read buffer capped at 16 KiB; write buffer holds one response at a time (no unbounded growth).

## Phase 10 — Integration & Audit Scenarios

> **Note for teammates:** A minimal `www/` folder was created (index.html + errors/404.html + errors/500.html)
> just to verify the server starts and responds. Before the audit, this folder needs to be expanded into
> a proper demo tree that exercises every feature below. See the ai_changelog.md entry dated 2026-05-07
> "Bootstrap www/ test content" for context.

- [ ] Build out `www/` demo tree: static files, a subdirectory for autoindex, an upload target, a redirect route, a CGI script, and a session demo page — enough to walk through every audit question in the browser.
- [ ] Curl scripts under `tests/scripts/` for: GET, POST upload, DELETE, malformed, oversized body (413), method not allowed, redirect, autoindex, custom error page, chunked upload, CGI chunked, CGI unchunked.
- [ ] Multi-port config + server_name resolution test (same IP:port, two `server_name`s).
- [ ] Duplicate-port config rejected with clear error.
- [ ] Browser smoke test checklist (DevTools): headers, status codes, redirects, cookies, autoindex, upload.
- [ ] `make siege` -> `siege -b 127.0.0.1:8080` for ~1 min: **>= 99.5% availability**, no leaks (`top` flat).
- [ ] No hanging connections after siege exits.

## Phase 11 — Documentation & Submission
- [ ] `README.md`: build, run, config schema, examples, audit walkthrough.
- [ ] `ARCHITECTURE.md`: layer diagram + dependency rule (inward only).
- [ ] Annotated answers to audit questions (single epoll, why, error handling).
- [ ] Final pass: `cargo fmt`, `cargo clippy -D warnings`, all tests green, siege ≥ 99.5%.

---

## Definition of Done (per phase)
A phase is complete only when: code compiles, `clippy -D warnings` is clean, unit tests cover happy + error paths, and the relevant audit item from `audit/README.md` can be demonstrated end-to-end.
