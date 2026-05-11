



# Smart Road

Rust + SDL2 smart-intersection simulation for autonomous vehicles.

The project is built as a small simulation engine with a rendering layer on top:

- `domain/` contains the traffic rules, path geometry, safety helpers, and statistics
- `application/` contains the world tick, spawn system, and reservation flow
- `infrastructure/` contains SDL input/app setup and random vehicle generation
- `presentation/` draws the road, cars, and statistics screen

## What The Code Does

Each vehicle has:

- a fixed origin: `North`, `South`, `East`, `West`
- a fixed route: `Right`, `Straight`, `Left`
- a path built once at spawn time
- a scalar progress value `s` that moves along that path

The simulation does not move vehicles by hand-written `x/y` rules for every turn.
Instead, each tick does:

1. choose a target speed
2. advance `progress.s` by `speed * dt`
3. sample the path to get position and heading

## Vehicle Flow

The current implementation works in three map parts.

### 1. Approach Lane

Vehicle behavior before the intersection:

- vehicles spawn at the lane entry
- they start at `Low` speed
- same-lane safety is active here
- the first car may wait before the intersection if reservation is denied

Same-lane following is based on path distance, not raw screen distance.

For two cars in the same lane:

```text
gap = leader_progress - follower_progress - vehicle_length
```

That `gap` is the empty space between the rear of the front car and the front of
the car behind it.

Current spacing values:

- `VEHICLE_LENGTH = 36 px`
- `SAFETY_BUFFER = 30 px`
- target empty gap in queues = `30 px`

### 2. Intersection

Vehicle behavior while entering and crossing:

- vehicles request permission before entering the conflict zone
- right turns do not use the inner conflict map
- straight and left turns use the reservation manager
- cars inside the intersection run at `High` speed

The intersection is represented as a 6x6 road grid with an inner 4x4 conflict zone.

- right turns use outer corner tiles
- straight and left routes pass through the inner 4x4
- conflicts are resolved by route/tile overlap, not by traffic lights

### 3. Leaving

Vehicle behavior after clearing the conflict zone:

- reservation is released once the rear clears the inner conflict box
- leaving vehicles run at `Medium` speed
- lane-following safety is no longer applied here

## Speed Model

The simulation now uses phase-based speeds instead of per-car cruise speeds.

Current speed values:

- `Low = 50 px/s`
- `Medium = 100 px/s`
- `High = 150 px/s`

Current rules:

- `Approaching` -> `Low`
- `WaitingReservation` -> `Low`
- `Entering` -> `High`
- `Inside` -> `High`
- `Leaving` -> `Medium`

There is one extra detail at the stop line:

- if a vehicle already has a reservation and its front reaches the intersection
  edge, it switches to `High` before the movement step into the intersection

## Spawning

Spawning is controlled by keyboard input and a queue.

Controls:

- `Arrow Up` -> spawn from `South`
- `Arrow Down` -> spawn from `North`
- `Arrow Left` -> spawn from `East`
- `Arrow Right` -> spawn from `West`
- `R` -> toggle random spawning
- `Esc` -> show statistics screen

Spawn safety rules:

- each direction has a cooldown
- blocked requests are queued
- a new car is spawned only if the lane entry is clear

Spawn clearance uses center-to-center distance:

```text
min_spawn_center_distance = vehicle_length + safety_buffer = 66 px
```

This is separate from the queue gap used by same-lane following.

## Reservation System

The project uses a reservation-based intersection manager.

For straight and left turns:

1. a car reaches the detection zone
2. it requests a reservation
3. if there is no active route conflict, it is granted
4. otherwise it keeps approaching slowly and is capped at the stop line

Right-turn vehicles are outside the inner conflict map, so they skip the normal
route-conflict reservation batch.

Reservation requests are also ordered fairly:

- earlier arrivals are served first
- long waits get a small aging bonus

## Safety Rules

There are three different safety ideas in the code.

### Same-Lane Safety

Used only before a vehicle has committed through the intersection:

- active in `Approaching`, `WaitingReservation`, and `Entering`
- not active in `Inside` or `Leaving`

This keeps queued cars from rear-ending each other while still allowing the
intersection phase to be controlled mainly by reservations.

### Spawn Safety

Used only at the lane entry:

- prevents cars from being created on top of each other
- uses center-to-center distance from the spawn point

### Intersection Conflict Safety

Used only for crossing routes:

- based on route conflicts in the inner 4x4 tile map
- not based on Euclidean distance

## Statistics

The simulation tracks:

- vehicles completed
- max speed
- min speed
- max time to pass
- min time to pass
- close calls

Time-to-pass is measured from first detection by the smart-intersection logic
until the vehicle is removed from the scene.

Close calls are currently counted when two vehicles from different routes occupy
the same inner conflict tile at the same time.

## Tests

The project has active unit and integration tests covering:

- path geometry and progress
- lane safety
- spawn queueing
- reservation conflicts
- statistics updates

The current codebase is expected to pass:

```bash
cargo test
cargo fmt --check
cargo clippy --all-targets --all-features -- -D warnings
```

## Summary

In simple words:

- a car spawns at the start of a lane
- it follows a fixed route
- it approaches slowly
- it keeps distance from the car ahead
- it asks permission before entering the intersection
- if denied, it waits at the stop line
- if allowed, it enters at high speed
- once it clears the conflict zone, it leaves at medium speed

That is the current behavior implemented in the code.
