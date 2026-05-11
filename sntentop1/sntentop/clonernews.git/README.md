# 🌊 clonernews

> A polished, vanilla JavaScript Hacker News client with a feature-first architecture, strong UI foundations, and a strict safety posture.

## ✨ What this is

`clonernews` is a static Hacker News app built with Vite and plain JavaScript. The codebase is structured around clean boundaries so the domain stays pure, the infrastructure stays focused on I/O, and the UI stays accessible and easy to extend.

## 🎨 Design direction

The visual direction is deliberate: a dark, editorial-style canvas, warm accents, expressive typography, and lightweight motion. The goal is to feel sharp and readable instead of generic or overfitted to a framework.

## 🧱 Architecture

- `src/core/` for pure entities and use cases.
- `src/infra/` for adapters and runtime I/O.
- `src/shared/` for reusable UI primitives.
- `src/features/` for feature-local views and controllers.

## 📚 Documentation

- Start here for contribution flow, reading order, and templates: [docs/README.md](docs/README.md)
- Team operating model and PR discipline: [docs/agentic-workflow-guide.md](docs/agentic-workflow-guide.md)
- Work breakdown, tracks, and dependencies: [docs/implementation-plan.md](docs/implementation-plan.md)
- Live ticket status and ownership: [docs/tickets.md](docs/tickets.md)

## 🛠️ Scripts

- `npm run dev` starts the Vite dev server.
- `npm run build` creates the production bundle.
- `npm run preview` serves the production build locally.
- `npm run check` runs Biome checks.
- `npm test` runs the Vitest suite.
- `npm run ci:quality` runs the full quality gate (`check`, AGENTS guards, unit tests, e2e, build).
- `npm run ci:policy` validates policy gates (PR metadata in CI, or branch ticket metadata locally) and architecture boundaries.
- `npm run ci:all` runs both quality and policy gates.

## 🚧 Phase 0 status

The shared kick-off, initial infra layer, and core styling foundation are in place:

- Project config is in place.
- Core contract files exist.
- Shared utilities are implemented and tested.
- HN API adapter, cache adapter, throttle utility, and core use-cases are implemented and test-backed.
- The starter UI, typography system, and design tokens are wired into the bootstrap shell.

## 🔎 Notes

- HN API HTML fields will be sanitized before insertion in later feature work.
- Hash routing is used so the app can run cleanly on GitHub Pages.
- The implementation plan and ticket tracker remain the source of truth for remaining work.

## 📌 Next steps

1. Build feed/detail/comments/poll feature views and the app shell wiring.
2. Finish the integration smoke tests and audit traceability matrix.
3. Expand audit coverage and finalize the deployment workflow.