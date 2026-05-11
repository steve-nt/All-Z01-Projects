# Project Structure Overview

This document mirrors the layout described in `README.md` but focuses on how modules depend on each other inside the runtime. Treat it as the map for onboarding new developers or planning architectural changes.

## Entry Points & Static Assets

- `index.html` defines the layered DOM (background, entities, HUD, overlays) and loads `main.js` as an ES module.
- `main.js` is the only script tag entry point; it bootstraps every subsystem, creates the shared state container, and exposes `window.startNewRun` so the story overlay can trigger gameplay.
- `style/` holds all CSS and menu imagery:
  - `style/styles.css` covers the aquarium scene, HUD, and gameplay layers.
  - `style/menu.css` styles start/pause/story overlays plus celebratory FX.
  - `style/images/` stores background JPEGs/PNGs referenced by CSS.

## Core Layer (`src/core`)

- `constants.js` centralizes configuration like durations, max lives, and world bounds. Any module needing shared tuning imports from here.
- `settings.js` exports `createSettingsHandlers` so audio toggles persist via `localStorage` without leaking storage calls across the codebase.
- `gameLoop.js` builds the `createGameLoop` factory used by `main.js` to update physics, spawning, and HUD every frame.
- `lifecycle.js` exposes `createLifecycleSystem`, orchestrating pause/resume/restart and main menu transitions.
- `input.js` provides `createInputSystem`, handling mouse/keyboard events plus the dedicated crosshair RAF loop.

## Entity Layer (`src/entities`)

- `entities.js` supplies low-level DOM helpers (`positionElement`, `removeEntity`, `clearEntities`).
- `lifeFish.js` implements `LifeFishManager`, which tracks spawn timers and behaviors for the life-granting fish.
- `turtle.js` implements `TurtleManager`, controlling the hazard spawn cycle and reactions when hit or escaped.

## Gameplay Layer (`src/gameplay`)

- `spawning.js` contains `spawnBubble` for ambient visuals and `spawnParrot` for main fish (and indirectly life fish) creation.
- `shooting.js` exposes `createShootingSystem`, responsible for hit detection, scoring, combo logic, and miss penalties.

## UI Layer (`src/ui`)

- `highScores.js` encapsulates localStorage access for leaderboard data plus DOM update helpers.
- `hud.js` exports `createHudSystem`, syncing HUD and pause menu elements with live state.
- `menuEffects.js` exports `createMenuEffects`, handling animated menu bubbles, pause bubbles, and celebration visuals.
- `feedback.js` exports `createFeedbackSystem`, managing score popups, combo banners, center notifications, and Game Over.

## Story Layer (`src/story`)

- `story.js` defines `StoryManager`, which runs the optional story overlay, handles navigation events, and triggers `window.startNewRun()` when the player begins a session.

## Dependency Flow

1. `main.js` imports constants/settings first, then instantiates the game loop, lifecycle, input, entities, gameplay, and UI factories.
2. Entity managers (`LifeFishManager`, `TurtleManager`) and DOM helpers feed into gameplay factories (`spawnParrot`, `createShootingSystem`).
3. Gameplay modules mutate shared state; UI factories receive getter/setter closures so they stay stateless and focus on rendering.
4. Lifecycle and story modules call back into the shared state/pause helpers exposed from `main.js`, keeping overlays decoupled from gameplay logic.

## Reuse & Extension Notes

- Constants and settings modules isolate tuning/persistence, making them safe extension points for new mechanics or backend syncing.
- Entity helpers ensure all DOM manipulation flows through one utility, which simplifies future refactors (e.g., swapping to Canvas/WebGL).
- UI factories are reusable: the same HUD update logic powers both the in-game HUD and pause menu by sharing getter/setter bindings.

Consult this document alongside the README when reorganizing files or introducing new subsystems to ensure imports remain acyclic and responsibilities stay narrow.
