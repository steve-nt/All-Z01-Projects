# Smart Road Plan

## Goal

Build a Rust + SDL2 simulation of a four-way smart intersection for autonomous vehicles.
The system must:

- spawn vehicles from keyboard input
- move them through a cross intersection without collisions
- manage at least 3 velocities
- enforce a strictly positive safety distance
- animate turning vehicles
- collect and display statistics on exit

## Core Approach

Use a **reservation-based intersection manager**.

Instead of traffic lights, each vehicle asks the intersection manager for permission to enter the conflict zone.
The manager decides:

- when the vehicle may enter
- what target speed it should maintain
- whether it must slow down before the intersection

This separates the problem into two layers:

1. **Lane following and safety distance**
   Vehicles follow their fixed lane and keep safe spacing from the vehicle ahead.
2. **Intersection conflict control**
   Vehicles only enter the intersection when their path reservation does not conflict with another active reservation.

This is simpler and safer than trying to make all cars negotiate continuously.

## Architecture Style

Use **Clean Architecture** so the core simulation rules do not depend on SDL2, input devices, or rendering.

Dependency direction must always point inward:

- presentation depends on application
- application depends on domain
- infrastructure depends on application and domain
- domain depends on nothing external

This matters here because:

- the traffic algorithm should be testable without a window
- SDL2 should be replaceable without rewriting the simulation
- statistics, spawning, and reservation logic should run in unit and integration tests

## Clean Architecture Layers

### Domain

Pure business rules and core models:

- vehicle entity
- lane and route definitions
- path progress model
- safety distance rules
- reservation/conflict rules
- statistics domain objects

Rules:

- no SDL2 types
- no file I/O
- no random generation directly inside entities
- deterministic where possible

### Application

Use cases that orchestrate domain behavior:

- advance simulation tick
- spawn vehicle request
- toggle random spawning
- request intersection reservation
- remove completed vehicle
- finalize statistics on exit

This layer coordinates domain objects and exposes clear APIs to the outer layers.

### Infrastructure

Technical adapters and concrete implementations:

- SDL2 window creation
- texture loading
- keyboard event adapter
- random vehicle generator
- clock/timestep adapter
- asset loading

Infrastructure should translate external data into application commands.

### Presentation

User-facing rendering and screens:

- road rendering
- vehicle sprite rendering
- stats window rendering
- HUD or overlays if added later

Presentation reads application state and renders it, but it should not contain simulation rules.

## Main Systems

### 1. App / Game Loop

Responsibility:

- initialize SDL2
- load textures/assets
- run fixed-timestep simulation updates
- process keyboard input
- render the world
- show statistics window on exit

Recommended loop:

- poll events
- accumulate delta time
- update simulation in fixed steps, e.g. `1/60`
- render once per frame

### 2. World State

Create a `World` struct containing:

- all active vehicles
- spawn queues per incoming direction
- intersection manager
- statistics collector
- random generator state
- simulation clock

### 3. Vehicle Model

Each vehicle should contain:

- `id`
- `origin` (`North`, `South`, `East`, `West`)
- `route` (`Right`, `Straight`, `Left`)
- `state` (`Approaching`, `WaitingReservation`, `Entering`, `Inside`, `Leaving`, `Done`)
- `position`
- `heading`
- `distance_travelled`
- `current_speed`
- `target_speed`
- `spawn_time`
- `detection_time`
- `exit_time`
- `reservation_id` optional

Physics values:

- `speed_low`
- `speed_medium`
- `speed_high`
- optional acceleration/deceleration later as bonus

Start with constant-speed tiers and simple speed transitions.

### 4. Lane and Path Definitions

Define 12 inbound lanes:

- 4 directions × 3 routes

Each lane has exactly one route:

- right
- straight
- left

Represent each route as a geometric path through the world:

- straight segments for approach/exit
- arc segment or sampled curve for turning

Recommended representation:

- a path as sampled waypoints or a parametric spline
- each vehicle moves by advancing a scalar `s` along its path

This avoids hardcoding special movement logic for each frame.

### 5. Intersection Manager

This is the core smart strategy.

#### Conflict Model

For each route, define:

- approach zone
- conflict zone inside the intersection
- exit zone

Precompute which route pairs conflict.

Example:

- two right turns from compatible sides may be non-conflicting
- straight vs left from crossing directions usually conflicts
- opposite straights may or may not conflict depending on path geometry

Represent this with a conflict matrix:

- `conflicts[route_a][route_b] -> bool`

#### Reservation Logic

When a vehicle reaches a detection point before the intersection:

1. it requests entry
2. the manager checks active reservations
3. if no conflict exists for the intended entry interval, it grants reservation
4. otherwise the car slows or stops before the stop line

Simple first version:

- only allow one conflicting vehicle group at a time
- allow multiple vehicles simultaneously only if their routes do not conflict

Reservation data:

- `vehicle_id`
- `route_key`
- `entry_time`
- `expected_exit_time`

Expected exit time can be estimated from:

- remaining path length inside intersection
- assigned speed

This is enough for a solid first implementation.

### 6. Safety Distance System

For each lane, vehicles must maintain a positive gap from the vehicle ahead.

Implement:

- sort or track vehicles by progress along the same lane
- compare follower-to-leader distance
- if gap < safety distance, reduce target speed
- if necessary, stop the follower before overlap

Safe initial rule:

- `safe_distance = vehicle_length + buffer`

Use one constant buffer first.

### 7. Vehicle Spawning

Controls:

- `Up`: spawn from south to north side
- `Down`: spawn from north to south side
- `Left`: spawn from east to west side
- `Right`: spawn from west to east side
- `R`: toggle continuous random spawning

Important requirement:

- key spam must not spawn overlapping vehicles

Solution:

- keep a per-direction spawn cooldown
- before spawning, verify the lane start zone is clear
- if blocked, enqueue the spawn request instead of forcing creation

When generating a vehicle:

- choose a random route valid for that direction
- choose an initial speed tier
- assign the correct path

### 8. Animation and Rendering

Use SDL2 textures for:

- road background
- lane markings
- vehicles

Rendering plan:

- draw static map first
- draw vehicles after updating transforms

Vehicle animation:

- rotate sprite based on current heading tangent to the path
- for turns, heading must change continuously along the curve

This satisfies the requirement that turning vehicles visually rotate while moving.

Simple rendering strategy:

- top-down 2D world
- logical world coordinates mapped directly to screen space
- fixed window size, e.g. `1280x720`

### 9. Statistics

Track globally:

- total vehicles completed
- max completed count observed
- max speed reached
- min speed reached
- max time to pass
- min time to pass
- close calls

Definitions:

- time to pass starts when detected by the smart intersection manager
- time ends when vehicle leaves the intersection and is removed

Close call rule:

- if two vehicles inside or near the conflict zone pass within less than safety distance, increment close calls
- collisions should never happen; close calls detect near-failures

On `Esc`:

- stop simulation
- open a statistics window
- render the final metrics clearly

## Project Structure

```text
DOCS/
  instructions.md          ← project subject
  audit_questions.md       ← audit checklist
src/
  main.rs                  ✅ bootstraps App::build + App::run
  lib.rs                   ✅ exposes all modules
  config.rs                ✅ window/road/timing/vehicle constants
  domain/
    mod.rs                 ✅
    vehicle.rs             ✅ Vehicle entity, speed tiers, state transitions
    lane.rs                ✅ Direction, Route
    path.rs                ✅ 12 strict L-shaped lane paths, col_x/row_y helpers
    safety.rs              ✅ min_safe_gap, is_too_close, adjusted_follower_speed
    intersection.rs        ✅ conflict map (inner 4×4 tile grid)
    reservation.rs         ✅ Phase 6 ReservationManager
    stats.rs               ✅ Statistics accumulator + final report model
  application/
    mod.rs                 ✅
    world.rs               ✅ World tick + lane safety + reservation control + statistics
    commands.rs            ✅ Phase 4
    spawn.rs               ✅ Phase 4
  infrastructure/
    mod.rs                 ✅
    sdl/
      mod.rs               ✅
      app.rs               ✅ SDL2 init, fixed-timestep loop, Esc to statistics screen
      input.rs             ✅ Phase 4
      texture_store.rs     ✅ generated vehicle texture
    random/
      mod.rs               ✅ Phase 4
      vehicle_generator.rs ✅ Phase 4
  presentation/
    mod.rs                 ✅
    scene.rs               ✅ top-level render() + render_statistics() dispatcher
    map_view.rs            ✅ static road layout (grass, asphalt, lane marks)
    vehicle_view.rs        ✅ rotated vehicle sprite rendering
    stats_view.rs          ✅ Phase 7 statistics screen
tests/
  reservation_flow.rs      ✅ Phase 6 — all tests active
  lane_safety.rs           ✅ Phase 5 — all tests active
  spawn_flow.rs            ✅ Phase 4
  stats_flow.rs            ✅ Phase 7
```

## Dependency Rules

- `domain` must not import from `application`, `presentation`, or `infrastructure`
- `application` may import `domain` and `shared`
- `presentation` may import `application` and `shared`
- `infrastructure` may import `application`, `domain`, and `shared`

If SDL2 types start appearing in domain structs, the architecture is already drifting in the wrong direction.

## Data Types

Useful enums and structs:

```rust
enum Direction {
    North,
    South,
    East,
    West,
}

enum Route {
    Right,
    Straight,
    Left,
}

enum SpeedTier {
    Low,
    Medium,
    High,
}

enum VehicleState {
    Approaching,
    WaitingReservation,
    Entering,
    Inside,
    Leaving,
    Done,
}
```

## Implementation Phases

### Phase 1: Test Setup

- [x] define how logic modules will be tested outside SDL2
- [x] set up unit test structure for pure simulation code
- [x] set up integration test structure for multi-vehicle and intersection scenarios
- [x] add first tests for path progress, safety distance, and reservation checks
- [x] add first integration tests for spawn flow, lane safety, and reservation conflicts
- [x] decide which behaviors stay manual visual checks only

Deliverable:

- basic unit and integration test foundation exists before core implementation grows

### Phase 2: Clean Architecture Skeleton ✅

- [x] create Rust project
- [x] create `domain`, `application`, `infrastructure`, and `presentation` modules
- [x] define dependency rules between layers
- [x] add `sdl2` dependency
- [x] create fixed timestep loop
- [x] open a window through infrastructure and presentation only
- [x] draw static road layout through presentation only

Deliverable:

- empty intersection scene running at stable FPS with clean layer boundaries ✅

### Phase 3: Vehicle Paths ✅

- [x] define the 12 lane paths
- [x] keep path and motion rules inside `domain`
- [x] spawn one vehicle manually
- [x] move it along its route
- [x] rotate sprite to match heading in `presentation`

Deliverable:

- one vehicle can drive right, straight, or left correctly ✅

### Phase 4: Input and Spawning ✅

- [x] implement arrow key input as application commands
- [x] implement random spawning with `R`
- [x] add spawn cooldown and queueing
- [x] keep random generation in infrastructure, not in domain

Deliverable:

- multiple cars can enter the map without overlapping on spawn ✅

### Phase 5: Lane Safety ✅

- [x] detect leader/follower per lane
- [x] apply safe distance rule in domain services
- [x] enforce speed reduction or stop behavior through application tick flow

Deliverable:

- no rear-end collisions on the same lane ✅

### Phase 6: Smart Intersection Manager ✅

- [x] add detection point before intersection
- [x] implement reservation requests in application use cases
- [x] build conflict matrix in domain
- [x] only allow non-conflicting entries
- [x] cover reservation behavior with integration tests

Deliverable:

- vehicles cross the intersection without collisions ✅

### Phase 7: Statistics ✅

- [x] track speed/time metrics in domain/application
- [x] count completed vehicles
- [x] detect close calls
- [x] add exit statistics window in presentation

Deliverable:

- full stats shown after `Esc` ✅

### Phase 8: Polish ✅

- [x] improve vehicle sprites and turning visuals
- [x] tune speeds and safety distance
- [x] refine reservation fairness to reduce congestion

Deliverable:

- smoother and more believable simulation ✅

## Fairness and Congestion Strategy

To avoid starvation:

- maintain waiting queues per approach
- if one direction dominates, raise priority for long-waiting vehicles
- use FIFO among vehicles with conflicting claims

Simple first rule:

- earliest detection time wins among conflicting requests

Then improve if needed with:

- aging bonus for vehicles waiting too long

## Recommended Math

Keep the first version simple:

- use `Vec2` positions as `f32`
- use path progress `s` in world units
- compute `speed = distance / time`
- update by `s += current_speed * dt`

Turning:

- either use cubic Bezier curves
- or pre-sampled path points with interpolation

Pre-sampled paths are easier to debug.

## Testing Strategy

### Manual tests

- spam each arrow key
- run random mode for several minutes
- verify no visual overlap at spawn
- verify no collisions in the intersection
- verify turning animation rotates correctly
- verify stats values update plausibly

### Deterministic tests

Extract pure logic from SDL2 and test:

- conflict matrix correctness
- reservation granting
- safe distance checks
- travel time calculation
- close call detection

This is important because SDL2 rendering should not hold all logic.

## Crates

Start with:

- `sdl2`
- `rand`

Optional later:

- `anyhow` for app setup errors
- `thiserror` for domain errors

Keep dependencies small.

## Risks

### Risk 1: Overcomplicated physics

Avoid full real-world physics first.
Start with:

- constant speed tiers
- instant speed switching

Then add acceleration only if there is time.

### Risk 2: Hardcoded turning behavior

Avoid per-route special-case movement in update code.
Use reusable path definitions.

### Risk 3: Coupling logic to rendering

Keep simulation data independent from SDL2 types wherever possible.

## Minimum Viable Version

The minimum version that satisfies the project well is:

- SDL2 intersection map
- arrow-key and random vehicle spawning
- 3 speed levels
- safe following distance
- reservation-based no-collision intersection manager
- turning animation with heading rotation
- statistics window on exit

## Nice Bonuses

- acceleration and braking curves
- custom-made car sprites
- heatmap or lane occupancy overlay
- per-route throughput stats
- average wait time
- fairness score by direction

## Final Recommendation

Build this in layers and do not start with optimization.

The right order is:

1. fixed timestep + map
2. reusable route paths
3. spawn/input
4. lane safety
5. reservation manager
6. statistics
7. visual polish

If the architecture stays clean, the difficult part is not SDL2. The difficult part is the intersection policy, so keep that logic isolated and testable.
