# Fish Tank Hunt
Fast-paced, pointer-driven aquarium target game built with plain JavaScript modules and DOM layers.

## 🎮 Overview
- Aim the crosshair (mouse or arrow keys) and click/space to tag fish before they exit the tank.
- Chain catches to build combos, boost the score multiplier, and chase the local high-score table.
- Catch the rare red life fish for extra lives and dodge the turtle hazards that break combos and drain points.
- Levels advance on a fixed timer; spawn rates, fish variants, and difficulty scale with each stage.

## 🧩 System Highlights
### Game Loop & Lifecycle (`src/core/gameLoop.js`, `src/core/lifecycle.js`)
- `createGameLoop` owns `update → render → HUD` and keeps frame work under `requestAnimationFrame`. It advances timers, transitions levels, enforces spawn caps, and prunes out-of-bounds entities before dispatching DOM positioning.
- `createLifecycleSystem` centralizes start/pause/restart flows. It resets shared state, clears DOM entities, syncs overlays, and exposes helpers like `startNewRun()` (also used by the story overlay).

### Entities & Gameplay (`src/entities/*`, `src/gameplay/*`)
- `LifeFishManager` and `TurtleManager` coordinate rare entity scheduling, motion curves, and callbacks when caught or missed.
- `spawnParrot` (main fish spawner) and `spawnBubble` (ambient layer) encapsulate DOM creation and variant selection; `entities.js` provides `positionElement`, `removeEntity`, and mass cleanup for render safety.
- `createShootingSystem` performs input hit tests, applies combo/score math, animates “caught” fish, and routes special cases (life fish heals, turtle penalties) back into shared state.

### Input & Controls (`src/core/input.js`)
- Mouse movement instantly drives the crosshair via `translate3d` for GPU compositing; `crosshairLoop` is a dedicated `requestAnimationFrame` that handles keyboard navigation independent of the main loop for minimal latency.
- Shooting is rate-limited (`SHOT_RATE_LIMIT_MS`) and normalized to game-space coordinates so future camera changes stay localized.

### HUD, Feedback, and Meta UI (`src/ui/*`, `src/story/story.js`, `index.html`)
- `createHudSystem` mirrors runtime stats in both the floating HUD and pause menu, including heart indicators and low-time styling.
- `createFeedbackSystem` handles combo banners, score popups, center notifications, and the Game Over overlay. High scores persist through `localStorage` (`src/ui/highScores.js`).
- `createMenuEffects` animates start/pause menus, celebratory bubble bursts, and shared overlays. The optional story sequence in `src/story/story.js` injects narrative beats before gameplay starts.

## 🧱 Architecture Deep Dive
- **Runtime state graph:** `main.js` owns canonical state (time, score, lives, combo, RNG seed, entity array). Smaller modules receive getter/setter lambdas instead of global imports, keeping side effects centralized and making it easy to replace the state container with Redux, Zustand, or a custom proxy later.
- **Render pipeline:** The game loop’s `render()` calls `positionElement` for every alive entity, ensuring layout work stays on `translate3d` transforms. HUD updates run after physics each frame, while keyboard-driven crosshair updates execute on a parallel RAF so input latency stays predictable even if simulation work spikes.
- **Entity lifecycle:** Every spawned DOM node is wrapped in an entity record (id, x/y, velocities, flags). `removeEntity` plus `clearEntities` guarantee DOM cleanup when levels reset or players pause, and managers such as `LifeFishManager`/`TurtleManager` hold higher-level spawn semantics so gameplay tuning lives in one place.
- **UI orchestration:** Overlay elements (start, pause, game over, story) are toggled through the lifecycle system. Menu FX rely on layered DOM nodes (`plants`, `bubbles`, `entities`, HUD) so designer tweaks do not require code changes, and the story manager can call `window.startNewRun()` without knowing anything about game internals.
- **Data surfaces:** High scores live in `localStorage`, audio settings are persisted via `settings.js`, and `responsibilities.md` plus the `docs/` folder track product decisions. These files effectively serve as lightweight “backend” artifacts the browser runtime reads/writes without external services.

## 📁 Folder Structure
```text
.
├── index.html            # Static shell: layered tank DOM, HUD, overlays
├── main.js               # Entry point wiring every system together
├── style                 # All CSS plus theme-specific images
│   ├── styles.css        # Core gameplay layout and HUD styling
│   ├── menu.css          # Start/pause/story overlay presentation
│   └── images            # Backgrounds, coral layers, decorative sprites for CSS
├── src
│   ├── core              # Game loop, lifecycle, input, settings, constants
│   ├── entities          # DOM helpers plus LifeFish & Turtle managers
│   ├── gameplay          # Spawning logic and shooting/combat rules
│   ├── ui                # HUD, feedback, high-score persistence, menu FX
│   └── story             # Narrative overlay controller
├── docs                  # Specs, feature logs, UX notes shared by the team
├── images                # Sprite sheets referenced directly by JS (fish, turtle, etc.)
└── LICENSE               # MIT terms applied to the entire repository
```

Key module responsibilities:
- `src/core/constants.js` – Shared tuning knobs (duration, spawn cadence, level caps).
- `src/core/settings.js` – Sound/volume persistence helpers.
- `src/core/lifecycle.js` – Start/pause/resume orchestration and resets.
- `src/gameplay/spawning.js` – Fish/bubble creation and world placement.
- `src/gameplay/shooting.js` – Shot handling, combos, penalties, and capture animations.
- `src/ui/hud.js` – Scoreboard rendering for HUD and pause stats.
- `src/ui/feedback.js` – Combo messaging, popups, and Game Over flow.

## 🚀 Setup & Run
The project is 100% static assets plus ES modules, so any HTTP server works. Use Python’s built-in server while developing:

```bash
python3 -m http.server 8080
```

Then open `http://localhost:8080/` in a modern browser. The start overlay appears first; press **Start Game** (or progress through the story overlay) to allow `main.js` to bootstrap the loop. No build step or bundler is required—just keep the server running so module imports resolve correctly.

## ⚡ Performance Notes
- **Layered compositing:** Separate DOM layers for plants, bubbles, entities, and HUD keep repaint regions tight. Entities are positioned via `translate3d(...)` to stay on the GPU compositor.
- **Dual RAF loops:** The main loop handles simulation/render while the crosshair/input loop runs independently, preventing long update steps from introducing cursor lag.
- **Adaptive spawning:** Spawn timers and max fish counts scale with level, preventing DOM overload while still raising difficulty.
- **FPS instrumentation:** Instantaneous FPS is sampled each frame (`hud-fps`) so feature work can watch for regressions. Rate limiting on shots and capped bubble counts protect against DOM spam.

## 🛠️ Extending the Codebase
For deeper changes, start by tracing control flow in `main.js`: it is the central handshake that wires modules together, so augmenting the state map or providing new UI hooks there keeps the rest of the code decoupled.

## 🗄️ Data & Backend Considerations
- **Local-first persistence:** High scores and audio preferences are written to `localStorage`. The code isolates persistence inside `highScores.js` and `settings.js`, so migrating to IndexedDB or a remote API requires stubbing those modules rather than rewriting gameplay.
- **Backend integration points:** To mirror scores or telemetry on a server, replace `addHighScore`/`updateHighScoreListUI` with fetch calls. Lifecycle hooks (`showGameOver`, `restartGame`, `startNewRun`) provide clean points to emit analytics events or WebSocket messages.
- **Session state serialization:** Because all mutable state is exposed via getters/setters in `main.js`, exporting a snapshot (for saving replays or supporting backend-driven resumes) is a matter of cloning that state bag and the `entities` array—no hidden globals exist.
- **Future services:** A Node/Express (or Deno) backend could receive POSTed score packets, perform validation, and broadcast leaderboards. Nothing in the frontend assumes network availability, so backend rollout can happen incrementally without breaking offline play.

## 🔧 Development Workflow
- **Tooling:** No build chain is required. Running `python3 -m http.server 8080` (or `npm serve`) is sufficient, which keeps debugging straightforward through browser devtools. Source maps are irrelevant because every module is already human-readable.
- **Docs & planning:** `docs/` record design discussions, art direction, and work agreements. Treat them as a living product spec; contributors should update these files when introducing mechanics or visual language changes.
- **Testing & profiling:** Use the HUD FPS readout plus Chrome DevTools’ performance tab to evaluate new effects. The dual-loop architecture makes it easy to isolate regressions: if crosshair lag appears, inspect `input.js`; if entity stutter occurs, profile `gameLoop.update`.
- **Story-driven flows:** The `story` module demonstrates how to build modal interactions that gate game start. Future tutorials, quests, or backend-driven missions can reuse that pattern.

---

## 🤝 Credits

**Team Fish Tank Hunt**
Built with care by a small crew of creators and tinkerers.

* **Andy** — Core gameplay systems, high-score flow, and FX polish. 

    [Let's chat on Discord](https://discordapp.com/users/780150798927134740) · [Connect on LinkedIn](https://www.linkedin.com/in/andriana-stas-419437329/) 🐚

* **Georgia** — Menu & HUD styling, celebratory effects, and UX polish. 

    [Say hey on Discord](https://discordapp.com/users/1277216244910522371) · [Find me on LinkedIn](https://www.linkedin.com/in/georgia-marouli/) 🌊



* **Sofia** — Entity design, difficulty tuning, and celebratory interactions.

    [Catch me on Discord](https://discordapp.com/users/1276592724979613697) · [Let's connect on LinkedIn](https://www.linkedin.com/in/sofia-busho-626433201/) 🐙



* **Xaroula** — Input responsiveness, accessibility tweaks, and menu interactions

    [Drop by on Discord!](https://discordapp.com/users/1242540766879023160) · [Network on LinkedIn](https://www.linkedin.com/in/theocharoula-tarara-650017200/) 🐠


> *Fish Tank Hunt* is maintained collectively by the team above.
> See [`LICENSE`](./LICENSE) for usage rights.

---

## 🌐 Community

Have improvement ideas, want to study a subsystem, or plan to integrate a backend leaderboard? Reach out through the channels above—we’re happy to collaborate.