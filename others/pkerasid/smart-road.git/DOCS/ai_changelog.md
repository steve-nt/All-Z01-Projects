# AI Changelog

Records every instruction given to the AI agent (Claude Code), the decisions made, and the
reasoning behind them. Teammates can read this to understand *why* the code exists, not just
what it does.

---

## Session 1 — Phase 1: Test Setup & Domain Foundation

**Date:** April 2026
**AI Agent:** Claude Code (claude-sonnet-4-6)

### Instructions given

- Define how logic modules will be tested outside SDL2.
- Set up unit test structure for pure simulation code.
- Set up integration test structure for multi-vehicle and intersection scenarios.
- Add first tests for path progress, safety distance, and reservation conflict checks.
- Add first integration tests for spawn flow, lane safety, and reservation conflicts.
- Decide which behaviors stay manual visual checks only.
- Build the conflict map for the intersection domain.

### Architectural decisions made

| Decision | Reason |
|---|---|
| Use Clean Architecture (domain / application / infrastructure / presentation) | Keeps simulation rules testable without SDL2; SDL2 can be replaced without rewriting the algorithm |
| Domain has zero external dependencies | Conflict map, path progress, and safety distance must be runnable in pure unit tests with no window |
| Integration tests use `#[ignore]` with phase labels for unimplemented systems | Prevents false green tests while clearly marking what is planned and when it will be activated |
| `build_conflict_map()` is a pure function returning a `HashMap` | Called once at startup; deterministic; fully testable without any game state |
| `PathProgress` advances by `speed * dt` and clamps at `total_length` | Simple scalar model avoids hardcoding per-route movement logic |
| `min_safe_gap = vehicle_length + buffer` | Guarantees a strictly positive safety distance as required by the project spec |

### Constraints given to the agent

- No SDL2 types anywhere in `domain/` or `application/`.
- Prefer `Result` and `?` over `unwrap()` in library code.
- Every `unsafe` block must have a `// SAFETY:` comment.
- Rust naming conventions: `snake_case`, `PascalCase`, `SCREAMING_SNAKE_CASE`.
- Edition 2024, rust-version 1.85, clippy pedantic enabled.

### Files created

- `src/domain/lane.rs` — `Direction`, `Route`
- `src/domain/vehicle.rs` — `VehicleState`, `SpeedTier`
- `src/domain/path.rs` — `PathProgress`
- `src/domain/safety.rs` — `min_safe_gap`, `is_too_close`
- `src/domain/intersection.rs` — `build_conflict_map()`, 6×6 tile grid conflict model
- `src/domain/stats.rs` — `Statistics` struct
- `tests/reservation_flow.rs` — conflict map tests + ignored Phase 6 tests
- `tests/lane_safety.rs` — safety distance tests + ignored Phase 5 tests
- `tests/spawn_flow.rs` — ignored Phase 4 tests

### Quality gate

15 unit + integration tests pass. `cargo test` clean.

---

## Session 2 — Phase 2: SDL2 Window + Static Road Rendering

**Date:** April 2026
**AI Agent:** Claude Code (claude-sonnet-4-6)

### Instructions given

- Add `sdl2` and `rand` as dependencies.
- Create a fixed-timestep game loop (60 fps, vsync) inside the infrastructure layer.
- Open a 1280×720 window through infrastructure and presentation only — domain and
  application must remain SDL2-free.
- Draw a static road layout: grass background, asphalt roads, lane markings, stop lines.
- Verify all Phase 1 tests still pass after adding SDL2.

### Architectural decisions made

| Decision | Reason |
|---|---|
| Fixed timestep with `accumulator += delta` and spiral-of-death guard (`delta.min(0.25)`) | Simulation advances at a deterministic rate regardless of render frame timing |
| `present_vsync()` on the canvas | Locks render to display refresh rate; no manual sleep needed |
| `App::build()` returns `Result<App, String>` | SDL2 init can fail; errors surface cleanly to `main` without `unwrap` |
| Road width = 240 px (3 lanes × 40 px × 2 directions) | Fits the 6×6 tile conflict model; leaves enough approach length at 1280×720 |
| All geometry constants in `src/config.rs` | Single source of truth; road bounds will be reused by path definitions (Phase 3), vehicle rendering (Phase 3), and the reservation manager (Phase 6) |
| Stop lines drawn only on inbound lanes | Matches the real-world convention and the ASCII diagram in `instructions.md` |
| Presentation layer imports SDL2 directly | Rendering is inherently coupled to the graphics API; the boundary that matters is domain/application being SDL2-free |

### Constraints given to the agent

- SDL2 must not appear in `domain/` or `application/` — verified with `grep -r "sdl2" src/domain/ src/application/`.
- `#[allow(clippy::cast_possible_wrap)]` scoped to `map_view.rs` only — all geometry values are well within i32 range.
- Rust toolchain and `libsdl2-dev` installed during this session (system did not have them).

### Files created

- `src/config.rs` — `WINDOW_W/H`, `LANE_W`, `ROAD_W`, `INTER_LEFT/TOP/RIGHT/BOTTOM`, `FIXED_DT`
- `src/application/world.rs` — `World` stub with `sim_time` and `tick(dt)`
- `src/infrastructure/sdl/app.rs` — `App::build()`, `App::run()` with fixed-timestep loop
- `src/infrastructure/sdl/input.rs` — stub for Phase 4
- `src/presentation/map_view.rs` — full static road renderer
- `src/presentation/scene.rs` — `render()` dispatcher

### Files modified

- `src/application/mod.rs` — added `pub mod world`
- `src/infrastructure/mod.rs` — added `pub mod sdl`
- `src/presentation/mod.rs` — added `pub mod map_view`, `pub mod scene`
- `src/lib.rs` — added `pub mod config`
- `src/main.rs` — `App::build().and_then(App::run)`, exits cleanly on error
- `Cargo.toml` — added `sdl2 = "0.37"`, `rand = "0.8"`
- `plan.md` — Phase 2 ticked off; project structure section updated to reflect real state
- `README.md` — progress line updated

### Quality gate

15/15 tests pass. `cargo build` clean. Window opens, road renders, Esc closes cleanly.

---

## Session 3 — Phase 3: Vehicle Paths

**Date:** April 2026
**AI Agent:** Codex (gpt-5.2)

### Instructions given

- Check all project markdown files and confirm the current build phase.
- Start Phase 3 after verifying that Phases 1 and 2 were already complete.
- Keep the implementation aligned with the existing agent-driven plan and changelog.

### Architectural decisions made

| Decision | Reason |
|---|---|
| Represent each lane route as a sampled polyline in `src/domain/path.rs` | Keeps movement deterministic and testable while supporting continuous heading changes for turns |
| Use quadratic sampled turns for left/right paths | Simpler than splines, but smooth enough for visible sprite rotation |
| Keep a scalar `PathProgress` and sample pose from `RoutePath` | Preserves the Phase 1 motion model while upgrading it with geometry and heading |
| Spawn one manual `North + Left` vehicle in `World::new()` | Satisfies the Phase 3 requirement without leaking Phase 4 input/spawn concerns into the implementation |
| Generate the vehicle texture procedurally in SDL infrastructure | Avoids introducing external asset files before the asset pipeline exists, while still using a rotatable sprite |
| Scope geometry-cast lint allowances to the new path/vehicle modules only | Maintains `clippy::pedantic` cleanliness without polluting unrelated modules |

### Constraints given to the agent

- Phase structure had to stay consistent with `README.md`, `plan.md`, and the existing changelog.
- Domain/application remained SDL2-free.
- Turning animation had to rotate continuously along the path heading.
- Use non-destructive edits only; the worktree started clean.

### Files created

- `src/infrastructure/sdl/texture_store.rs` — procedural vehicle texture for rotated sprite rendering
- `src/presentation/vehicle_view.rs` — vehicle drawing with `copy_ex` rotation

### Files modified

- `src/domain/path.rs` — `Vec2`, `PathSample`, `RoutePath`, sampled geometry for all 12 lanes, new path tests
- `src/domain/vehicle.rs` — `Vehicle` entity, speed tiers, movement tick, pose sampling
- `src/application/world.rs` — one manual vehicle plus tick integration
- `src/presentation/scene.rs` — render map + vehicles
- `src/infrastructure/sdl/app.rs` — propagate render errors
- `src/infrastructure/sdl/mod.rs` — export `texture_store`
- `src/presentation/mod.rs` — export `vehicle_view`
- `README.md` — Phase 3 marked complete, implementation summary added
- `DOCS/plan.md` — Phase 3 checklist marked complete, structure updated
- `src/domain/intersection.rs`, `src/domain/safety.rs` — small clippy cleanup while validating the milestone

### Quality gate

- `cargo test` clean
- `cargo build` clean
- `cargo clippy --all-targets --all-features` clean

---

## Session 4 — Phase 4: Input and Spawning


**Date:** April 2026
**AI Agent:** Codex (gpt-5.4)

### Instructions given

- Implement Phase 4.
- Keep the implementation aligned with the existing plan and architecture.

### Architectural decisions made

| Decision | Reason |
|---|---|
| Add `application::commands::Command` and `application::spawn::SpawnRequest` | Gives infrastructure a narrow application-facing API instead of mutating world state directly |
| Move spawn queueing and cooldown logic into `SpawnSystem` | Keeps `World` orchestration simple and makes spawn behavior unit-testable |
| Keep route/speed randomness in `infrastructure::random::VehicleGenerator` | Preserves the rule that randomness belongs to infrastructure, not domain |
| Treat `R` as a continuous random-mode toggle in the SDL app loop | Matches the project instructions for game-loop-driven spawning |
| Remove the Phase 3 hardcoded startup vehicle | Phase 4 requires keyboard/random spawning to be the entry path for new vehicles |
| Release queued vehicles only after movement/removal in `World::tick()` | Ensures a blocked lane can clear naturally before the next queued spawn is retried |

### Constraints given to the agent

- Domain and application had to remain SDL-free.
- Arrow keys had to map to the project’s travel directions (`Up` = south-to-north, etc.).
- Spawn spam could not create overlapping vehicles.
- Existing ignored Phase 5 and Phase 6 tests had to remain untouched.

### Files created

- `src/application/commands.rs` — application command enum for spawn requests
- `src/application/spawn.rs` — spawn queues, per-direction cooldowns, lane-entry checks
- `src/infrastructure/random/mod.rs` — random module export
- `src/infrastructure/random/vehicle_generator.rs` — randomized route/speed generation

### Files modified

- `src/application/mod.rs` — export Phase 4 modules
- `src/application/world.rs` — empty startup world, command handling, queue-driven spawning
- `src/domain/lane.rs` — enum helpers for indexed direction/route iteration
- `src/domain/vehicle.rs` — `SpeedTier::ALL` for generator selection
- `src/infrastructure/mod.rs` — export random infrastructure
- `src/infrastructure/sdl/app.rs` — input command handling and `R` random-spawn loop
- `src/infrastructure/sdl/input.rs` — SDL keyboard-to-action adapter
- `tests/spawn_flow.rs` — real Phase 4 integration coverage
- `README.md`, `DOCS/plan.md` — Phase 4 marked complete

### Quality gate

- `cargo test` clean

---

## Session 5 — Phase 5: Lane Safety

**Date:** April 2026
**AI Agent:** Claude Code (claude-sonnet-4-6)

### Instructions given

- Implement Phase 5: detect leader/follower per lane, apply safe distance rule in domain services, and enforce speed reduction or stop behavior through the application tick flow.
- Read all relevant project files before writing any code.
- Implement the three previously ignored integration tests in `tests/lane_safety.rs`.

### Architectural decisions made

| Decision | Reason |
|---|---|
| Promote `VEHICLE_LENGTH` and `SAFETY_BUFFER` to `config.rs` | Both constants were duplicated (private in `spawn.rs`); moving them to `config.rs` makes them the single source of truth for spawn checks and lane safety alike |
| Add `adjusted_follower_speed()` to `domain/safety.rs` | Lane speed adjustment is a pure business rule (no SDL2, no world state); belongs in the domain alongside the existing `min_safe_gap` and `is_too_close` functions |
| `apply_lane_safety` runs before `vehicle.tick()` each step | Vehicles must see the corrected `target_speed` before they advance; running it after would delay enforcement by one tick and could allow momentary overlap |
| Step 1 restores all vehicles to their natural tier speed before applying constraints | Unblocks followers whose leader has since moved away without requiring each vehicle to track whether it was previously constrained |
| Group vehicles by `(origin, route)` and sort by `progress.s` descending | Vehicles on the same lane share an identical path; progress scalar comparison is exact and avoids Euclidean distance approximation across curved segments |
| Cascading stop: follower speed uses `leader.current_speed` (previous tick) | If a middle vehicle is already stopped, the vehicle behind it inherits speed 0 — no special multi-vehicle case needed |
| Leader's `target_speed` is explicitly restored after the follower loop | Prevents the leader from being accidentally zeroed; leaves a clean hook for Phase 6 to apply intersection-waiting constraints on top |

### Files modified

- `src/config.rs` — added `VEHICLE_LENGTH` and `SAFETY_BUFFER` as public constants
- `src/domain/safety.rs` — added `adjusted_follower_speed()` domain service
- `src/application/spawn.rs` — replaced local constants with `config::VEHICLE_LENGTH` / `config::SAFETY_BUFFER`
- `src/application/world.rs` — added `apply_lane_safety()` and wired it into `World::tick()` before vehicle updates
- `tests/lane_safety.rs` — implemented all three previously ignored Phase 5 tests

### Quality gate

- 40 tests pass (`cargo test` clean), 0 ignored among Phase 5 tests
- `cargo build` clean, no warnings

---

## Session 6 — Phase 6: Smart Intersection Manager + Path Geometry Overhaul

**Date:** April 2026
**AI Agent:** Claude Code (claude-sonnet-4-6)

### Instructions given

- Fix lane paths: cars from East and West were moving incorrectly; turning points should be one tile only, all other movement straight.
- Right turns must use the outer corner tile on the vehicle's **entry side** (first tile entering), not the far corner (last tile before exit).
- Implement Phase 6: detection point before intersection, reservation requests, conflict-based entry control, integration tests.

### Part A — Path geometry overhaul (`path.rs`)

The old quadratic Bezier turn system was replaced entirely with strict 2- or 3-point polylines.

| Route type | Waypoints | Notes |
|---|---|---|
| Straight | `[entry_offscreen, exit_offscreen]` | 2 points, constant heading |
| Left turn | `[entry_offscreen, turn_tile_center, exit_offscreen]` | 3 points; turn tile matches inner 4×4 conflict map exactly |
| Right turn | `[entry_offscreen, entry_corner_tile, exit_offscreen]` | 3 points; corner tile is at the **entry edge** of the intersection |

#### Bug fixed: East-Right and West-Right corner tiles

The original paths crossed the entire intersection before turning:

| Route | Old (wrong) corner | Fixed corner |
|---|---|---|
| East-Right (going west) | `(540, 260)` — last column | `(740, 260)` — first column from East |
| West-Right (going east) | `(740, 460)` — last column | `(540, 460)` — first column from West |

Each right turn now uses the corner tile on its own entry side:
- North-Right → top-left `(540, 260)`
- East-Right → top-right `(740, 260)`
- South-Right → bottom-right `(740, 460)`
- West-Right → bottom-left `(540, 460)`

#### Left turn waypoints verified against the conflict map

Inner tile `(c, r)` maps to pixel `(col_x(c+1), row_y(r+1))`:

| Route | Inner tile | Pixel waypoint |
|---|---|---|
| North-Left (NE) | (1,2) | (620, 380) |
| South-Left (SW) | (2,1) | (660, 340) |
| East-Left (ES) | (1,1) | (620, 340) |
| West-Left (WN) | (2,2) | (660, 380) |

### Part B — Phase 6 reservation manager

#### Architectural decisions made

| Decision | Reason |
|---|---|
| `ReservationManager` in `domain/reservation.rs` | Pure business logic (no SDL2, no world state); holds the static conflict map and a live `HashMap<RouteId, Vec<u64>>` of active holders |
| Multiple vehicles allowed on the same route simultaneously | Same-route vehicles never conflict in the tile grid; lane safety (Phase 5) handles their spacing |
| Right turns skip reservation entirely | Their corner tiles are outside the inner 4×4 conflict zone; they receive `reservation_granted = true` on first detection |
| `apply_reservation_control` runs before `vehicle.tick()` | Sets `target_speed = 0` for denied vehicles so the stop takes effect in the same step |
| `release_exited_reservations` runs after `vehicle.tick()` | Reads post-tick geometric state; releases reservation once the vehicle has been `Inside` and is no longer inside the box |
| `entered_intersection: bool` field on `Vehicle` | Needed to distinguish "waiting at detection line" (Inside not yet seen) from "exiting after crossing" (Inside seen, now outside) |
| Detection zone = 2 lane widths (80 px) before the intersection edge | Wide enough to stop a vehicle at medium speed within one fixed step; narrow enough not to trigger for vehicles on crossing roads |
| `VehicleState::WaitingReservation` set by application layer, not vehicle.tick() | Keeps domain vehicle logic geometric; reservation policy lives in the application |

### Files created

- `src/domain/reservation.rs` — `ReservationManager` with `request` / `release` API and 8 unit tests

### Files modified

- `src/domain/path.rs` — complete rewrite of `for_lane()`; removed Bezier system; added `col_x` / `row_y` helpers; updated and expanded path tests
- `src/domain/mod.rs` — added `pub mod reservation`
- `src/domain/vehicle.rs` — added `reservation_granted: bool` and `entered_intersection: bool` fields
- `src/application/world.rs` — added `ReservationManager` to `World`; added `apply_reservation_control` and `release_exited_reservations`; updated tick order; added `to_route_id` and `is_in_detection_zone` helpers
- `tests/reservation_flow.rs` — implemented all four previously ignored Phase 6 tests

### Quality gate

- 55 tests pass (`cargo test` clean), 0 ignored
- `cargo build` clean, no warnings

---

## Session 9 — Speed & Collision Bugfix Pass

**Date:** April 2026
**AI Agent:** Claude Code (claude-sonnet-4-6)

### Instructions given

- Vehicles were stopping completely (reaching speed 0); they must always move at one of the three fixed speed tiers — never zero.
- After removing stops, vehicles were entering the intersection without a reservation and colliding; fix collision avoidance.
- Right-turn vehicles were not always travelling at full speed; they have no conflict path and the spawn system guarantees clearance — they must always go at High speed.
- Vehicles were being assigned random speed tiers at spawn; all vehicles should spawn at High speed and the collision avoidance logic is the only thing that should change speed.

### Root causes identified

| Bug | Root cause |
|---|---|
| Vehicles stopping | `adjusted_follower_speed` returned `0.0` when gap < safe gap; reservation denial set `target_speed = 0.0` |
| Intersection collisions after removing stops | `DETECTION_MARGIN` was only 80 px — at Low speed (78 px/s) a denied vehicle crossed the zone in ~1 s, far less than a left-turn crossing (~3.6 s) |
| Right-turn speed reduction | `apply_lane_safety` iterated over `Route::ALL` including `Route::Right`; a follower right-turner behind a slower one got its speed reduced |
| Right-turn same-lane collisions | Random speed tier assignment meant a High-tier right-turner could spawn behind a Low-tier one; with lane safety disabled for right turns, catch-up was inevitable |
| Random spawn speeds | `VehicleGenerator::spawn_request_for_origin` called `choose_random` on `SpeedTier::ALL`, assigning Low/Medium/High arbitrarily |

### Architectural decisions made

| Decision | Reason |
|---|---|
| Minimum speed is always `SpeedTier::Low` (78 px/s), never zero | Project specification: vehicles must always move at one of three fixed speeds |
| `adjusted_follower_speed` returns `leader_speed.min(natural_speed)` for all close-gap cases | Matching leader speed prevents closing the gap further while ensuring the vehicle never stops; the two previously separate close-gap branches collapse to one |
| Reservation denial sets `target_speed = target_speed.min(SpeedTier::Low)` | Denied vehicles creep instead of stopping; `min` preserves any lane-safety constraint already applied |
| `DETECTION_MARGIN` increased from 80 px (2 lane widths) to 320 px (8 lane widths) | At Low speed, 320 px takes ~4.1 s — enough buffer for worst-case left-turn crossing at Low speed (~3.6 s) to clear before the denied vehicle reaches the conflict zone |
| `Route::Right` skipped entirely in `apply_lane_safety` loop | Right turns have no conflict path; the spawn system guarantees entry clearance; lane safety constraints are irrelevant and caused unnecessary slowdowns |
| Right-turn vehicles forced to `SpeedTier::High` in speed-reset step | All right-turners travel at the same speed, eliminating catch-up between them and the need for any lane safety; assigned tier becomes irrelevant |
| All vehicles spawn at `SpeedTier::High` | Speed changes are the responsibility of the collision avoidance systems; the spawn assignment should not pre-decide a vehicle's cruise speed |

### Files modified

- `src/domain/safety.rs` — `adjusted_follower_speed`: removed `0.0` branch; gap < 2× safe gap now always returns `leader_speed.min(natural_speed)`
- `src/application/world.rs` — reservation denial uses `SpeedTier::Low` instead of `0.0`; `DETECTION_MARGIN` 80 px → 320 px; `Route::Right` skipped in lane safety loop; right-turn vehicles forced to `SpeedTier::High` in speed-reset step; `SpeedTier` added to imports
- `src/infrastructure/random/vehicle_generator.rs` — `speed_tier` hardcoded to `SpeedTier::High` instead of random
- `tests/lane_safety.rs` — `follower_slows_when_too_close_to_leader` assertion updated: checks `current_speed <= leader.current_speed && current_speed > 0.0` instead of `current_speed == 0.0`
- `src/application/world.rs` (test) — `reservation_is_released_once_vehicle_is_on_the_exit_side`: East vehicle repositioned from `s=520` (40 px from intersection, designed for stop-and-wait) to `s=200` (outside new 320 px detection zone; drifts in naturally during the 2.5 s tick)

### Quality gate

- 66 tests pass (`cargo test` clean), 0 ignored
- `cargo build` clean, no warnings
