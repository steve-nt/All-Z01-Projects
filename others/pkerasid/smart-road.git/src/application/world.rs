#![allow(
    clippy::cast_possible_truncation,
    clippy::cast_precision_loss,
    clippy::cast_sign_loss
)]

use std::collections::{HashMap, HashSet};

use crate::application::commands::Command;
use crate::application::spawn::SpawnSystem;
use crate::config::{
    INTER_BOTTOM, INTER_LEFT, INTER_RIGHT, INTER_TOP, LANE_W, SAFETY_BUFFER, VEHICLE_LENGTH,
};
use crate::domain::intersection::{EntryDir, RouteId, TileId, TurnDir};
use crate::domain::lane::{Direction, Route};
use crate::domain::path::Vec2;
use crate::domain::reservation::ReservationManager;
use crate::domain::safety::adjusted_follower_speed;
use crate::domain::stats::Statistics;
use crate::domain::vehicle::{SpeedTier, Vehicle, VehicleState};

/// Top-level simulation state.
pub struct World {
    /// Total elapsed simulation time in seconds.
    pub sim_time: f32,
    /// Active vehicles currently simulated and rendered.
    vehicles: Vec<Vehicle>,
    spawner: SpawnSystem,
    reservation_manager: ReservationManager,
    stats: Statistics,
    detection_started_at: HashMap<u64, f32>,
    counted_close_calls: HashSet<(u64, u64)>,
    /// Sim-time of the last reservation grant per approach direction.
    /// Prevents multiple vehicles from the same direction flooding the
    /// intersection simultaneously — each grant enforces a minimum gap.
    last_entry_time: [f32; 4],
}

impl World {
    #[must_use]
    pub fn new() -> Self {
        Self {
            sim_time: 0.0,
            vehicles: Vec::new(),
            spawner: SpawnSystem::new(),
            reservation_manager: ReservationManager::new(),
            stats: Statistics::default(),
            detection_started_at: HashMap::new(),
            counted_close_calls: HashSet::new(),
            last_entry_time: [-RESERVATION_ENTRY_INTERVAL; 4],
        }
    }

    /// Apply an application command produced by infrastructure input adapters.
    pub fn apply_command(&mut self, command: Command) {
        match command {
            Command::SpawnVehicle(request) => self.spawner.submit(&mut self.vehicles, request),
        }
    }

    /// Advance the simulation by one fixed step.
    pub fn tick(&mut self, dt: f32) {
        let tick_start = self.sim_time;
        self.sim_time += dt;

        // Phase 5: lane-following safety (sets target_speed before tick).
        apply_lane_safety(&mut self.vehicles);

        // Phase 6: intersection reservation control (may override target_speed to 0).
        apply_reservation_control(
            &mut self.vehicles,
            &mut self.reservation_manager,
            &mut self.detection_started_at,
            &mut self.last_entry_time,
            tick_start,
            dt,
        );

        // Advance every vehicle with the resolved target_speed.
        for vehicle in &mut self.vehicles {
            vehicle.tick(dt);
        }

        backfill_detection_times(&self.vehicles, &mut self.detection_started_at, tick_start);

        // Phase 6: post-tick — update entered_intersection flag and release
        // reservations for vehicles that have fully exited the conflict zone.
        release_exited_reservations(&mut self.vehicles, &mut self.reservation_manager);

        update_statistics(
            &self.vehicles,
            &mut self.stats,
            &mut self.detection_started_at,
            &mut self.counted_close_calls,
            self.sim_time,
        );

        self.vehicles.retain(|v| !v.is_done());
        self.spawner.tick(dt, &mut self.vehicles);
    }

    #[must_use]
    pub fn vehicles(&self) -> &[Vehicle] {
        &self.vehicles
    }

    #[must_use]
    pub fn queue_len(&self, direction: Direction) -> usize {
        self.spawner.queue_len(direction)
    }

    #[must_use]
    pub fn statistics(&self) -> &Statistics {
        &self.stats
    }
}

// ── Phase 5: lane safety ──────────────────────────────────────────────────────

/// Enforce lane-following safety constraints across all active vehicles.
///
/// Lane-following applies only before the vehicle has committed through the
/// intersection.  Once a vehicle is inside the intersection it is governed by
/// the reservation system; once it is leaving it has already cleared the
/// conflict path.
fn apply_lane_safety(vehicles: &mut [Vehicle]) {
    // Step 1 — restore every vehicle to its phase speed. Safety and reservation
    // rules below can only reduce this value.
    for v in vehicles.iter_mut() {
        v.target_speed = phase_speed(v);
    }

    for origin in Direction::ALL {
        for route in Route::ALL {
            let mut lane: Vec<(usize, f32)> = vehicles
                .iter()
                .enumerate()
                .filter(|(_, v)| {
                    v.origin == origin
                        && v.route == route
                        && !v.is_done()
                        && needs_lane_following_safety(v)
                })
                .map(|(i, v)| (i, v.progress.s))
                .collect();

            if lane.len() < 2 {
                continue;
            }

            // Highest s = leader.
            lane.sort_unstable_by(|a, b| {
                b.1.partial_cmp(&a.1).unwrap_or(std::cmp::Ordering::Equal)
            });

            for pair in lane.windows(2) {
                let (leader_idx, leader_s) = pair[0];
                let (follower_idx, follower_s) = pair[1];

                let gap = leader_s - follower_s - VEHICLE_LENGTH;
                let leader_speed = vehicles[leader_idx].current_speed;
                let natural = phase_speed(&vehicles[follower_idx]);

                vehicles[follower_idx].target_speed = adjusted_follower_speed(
                    gap,
                    natural,
                    leader_speed,
                    VEHICLE_LENGTH,
                    SAFETY_BUFFER,
                    SpeedTier::Low.units_per_second(),
                    SpeedTier::Medium.units_per_second(),
                );
            }

            // Leader is never constrained by lane safety.
            let (leader_idx, _) = lane[0];
            vehicles[leader_idx].target_speed = phase_speed(&vehicles[leader_idx]);
        }
    }
}

fn needs_lane_following_safety(vehicle: &Vehicle) -> bool {
    matches!(
        vehicle.state,
        VehicleState::Approaching | VehicleState::WaitingReservation | VehicleState::Entering
    )
}

fn phase_speed(vehicle: &Vehicle) -> f32 {
    match vehicle.state {
        VehicleState::Leaving => SpeedTier::Medium.units_per_second(),
        VehicleState::Entering | VehicleState::Inside => SpeedTier::High.units_per_second(),
        VehicleState::Approaching | VehicleState::WaitingReservation | VehicleState::Done => {
            if vehicle.reservation_granted && front_has_reached_intersection(vehicle) {
                SpeedTier::High.units_per_second()
            } else {
                SpeedTier::Low.units_per_second()
            }
        }
    }
}

fn front_has_reached_intersection(vehicle: &Vehicle) -> bool {
    let pos = vehicle.sample_pose().position;
    let half_len = VEHICLE_LENGTH / 2.0;

    match vehicle.origin {
        Direction::North => pos.y + half_len >= INTER_TOP as f32,
        Direction::South => pos.y - half_len <= INTER_BOTTOM as f32,
        Direction::East => pos.x - half_len <= INTER_RIGHT as f32,
        Direction::West => pos.x + half_len >= INTER_LEFT as f32,
    }
}

// ── Phase 6: reservation control ─────────────────────────────────────────────

/// Detection zone depth (px) in front of each intersection edge.
/// Vehicles entering this zone request a reservation; if denied they slow down.
/// 320 px gives ~4.1 s at Low tier before reaching the intersection —
/// enough buffer for even the slowest left-turn crossing (~3.6 s) to clear.
const DETECTION_MARGIN: f32 = LANE_W as f32 * 8.0; // 320 px = 8 lane widths

/// Minimum seconds between consecutive reservation grants from the same approach
/// direction.  Prevents a group of bunched-up waiting vehicles from all rushing
/// through the intersection simultaneously when the conflicting car clears.
const RESERVATION_ENTRY_INTERVAL: f32 = 1.0;

// ── Phase 8c: fairness when granting reservations ─────────────────────────────
//
// Requests are ordered by a synthetic priority so the raw `Vec` order does not
// decide who wins when several vehicles compete the same tick.

/// Seconds of artificial earliness granted per real second spent waiting in the
/// detection zone (capped by [`RESERVATION_AGING_CAP_SECS`]).  Reduces starvation
/// when one approach would otherwise monopolise the conflict batch.
const RESERVATION_AGING_PER_WAIT_SEC: f32 = 0.12;
/// Maximum total aging bonus (seconds subtracted from the recorded detection time).
const RESERVATION_AGING_CAP_SECS: f32 = 2.5;

/// Lower values are served first.  `detection_time` is sim time when the vehicle
/// first entered the zone; `now` is the start of the current tick.
#[must_use]
fn reservation_request_priority(detection_time: f32, now: f32) -> f32 {
    let wait = (now - detection_time).max(0.0);
    let aging = (RESERVATION_AGING_PER_WAIT_SEC * wait).min(RESERVATION_AGING_CAP_SECS);
    detection_time - aging
}

/// For each vehicle at the detection line that does not yet have a reservation,
/// attempt to grant one.  If denied, cap `target_speed` so the vehicle can roll
/// up to the stop line but cannot enter the intersection.  Right-turn vehicles
/// are granted immediately because their corner tiles are outside the inner
/// conflict zone.
///
/// Phase 8c: non-right requests are sorted by [`reservation_request_priority`]
/// (earlier detection wins; long waits gain extra priority) before calling the
/// manager, so ordering matches FIFO-by-detection rather than spawn order.
fn apply_reservation_control(
    vehicles: &mut [Vehicle],
    manager: &mut ReservationManager,
    detection_started_at: &mut HashMap<u64, f32>,
    last_entry_time: &mut [f32; 4],
    tick_start: f32,
    dt: f32,
) {
    // Stamp first-seen detection times for anyone currently in the zone.
    for vehicle in vehicles.iter() {
        if vehicle.is_done() {
            continue;
        }
        let pos = vehicle.sample_pose().position;
        if is_in_detection_zone(pos, vehicle.origin) {
            detection_started_at.entry(vehicle.id).or_insert(tick_start);
        }
    }

    // Right-turn routes bypass the conflict map entirely.
    for vehicle in vehicles.iter_mut() {
        if vehicle.is_done() || vehicle.reservation_granted {
            continue;
        }
        let pos = vehicle.sample_pose().position;
        if !is_in_detection_zone(pos, vehicle.origin) {
            continue;
        }
        if vehicle.route == Route::Right {
            vehicle.reservation_granted = true;
        }
    }

    let mut batch: Vec<(usize, RouteId, f32, u64)> = vehicles
        .iter()
        .enumerate()
        .filter_map(|(idx, vehicle)| {
            if vehicle.is_done() || vehicle.reservation_granted || vehicle.route == Route::Right {
                return None;
            }
            let pos = vehicle.sample_pose().position;
            if !is_in_detection_zone(pos, vehicle.origin) {
                return None;
            }
            let route_id = to_route_id(vehicle.origin, vehicle.route);
            let t = *detection_started_at.get(&vehicle.id).unwrap_or(&tick_start);
            let prio = reservation_request_priority(t, tick_start);
            Some((idx, route_id, prio, vehicle.id))
        })
        .collect();

    batch.sort_unstable_by(|a, b| a.2.total_cmp(&b.2).then_with(|| a.3.cmp(&b.3)));

    for (idx, route_id, _, _) in batch {
        let vehicle = &mut vehicles[idx];
        let dir_idx = vehicle.origin.index();
        let entry_gap_ok = tick_start - last_entry_time[dir_idx] >= RESERVATION_ENTRY_INTERVAL;
        if entry_gap_ok && manager.request(vehicle.id, route_id) {
            vehicle.reservation_granted = true;
            last_entry_time[dir_idx] = tick_start;
        } else {
            let stop_cap = stop_line_speed_cap(vehicle.sample_pose().position, vehicle.origin, dt);
            vehicle.target_speed = vehicle.target_speed.min(stop_cap);
            vehicle.state = VehicleState::WaitingReservation;
        }
    }
}

/// After each tick, update `entered_intersection` and release reservations for
/// vehicles that have passed fully through the conflict box.
fn release_exited_reservations(vehicles: &mut [Vehicle], manager: &mut ReservationManager) {
    for vehicle in vehicles.iter_mut() {
        if vehicle.route == Route::Right {
            if vehicle.is_done() {
                vehicle.reservation_granted = false;
            }
            continue;
        }

        let pos = vehicle.sample_pose().position;
        let cleared_conflict_zone = has_cleared_conflict_zone(pos, vehicle.origin, vehicle.route);

        // Mark the first time the vehicle is known to have crossed the box.
        if vehicle.state == VehicleState::Inside || cleared_conflict_zone {
            vehicle.entered_intersection = true;
        }

        if vehicle.reservation_granted && cleared_conflict_zone {
            manager.release(vehicle.id, to_route_id(vehicle.origin, vehicle.route));
            vehicle.reservation_granted = false;
        }
    }
}

fn backfill_detection_times(
    vehicles: &[Vehicle],
    detection_started_at: &mut HashMap<u64, f32>,
    tick_start: f32,
) {
    for vehicle in vehicles {
        if vehicle.state != VehicleState::Approaching {
            detection_started_at.entry(vehicle.id).or_insert(tick_start);
        }
    }
}

fn update_statistics(
    vehicles: &[Vehicle],
    stats: &mut Statistics,
    detection_started_at: &mut HashMap<u64, f32>,
    counted_close_calls: &mut HashSet<(u64, u64)>,
    sim_time: f32,
) {
    for vehicle in vehicles {
        stats.record_speed_sample(vehicle.current_speed);

        if vehicle.is_done() {
            let detected_at = detection_started_at.remove(&vehicle.id).unwrap_or(sim_time);
            stats.record_completion(sim_time - detected_at);
        }
    }

    // Close call: two vehicles from different routes occupying the same inner
    // conflict tile simultaneously.  Tile-based detection is more meaningful
    // than Euclidean distance — it directly checks whether two paths literally
    // share the same 40×40 px section of the conflict zone.
    for (idx, vehicle) in vehicles.iter().enumerate() {
        if vehicle.is_done() {
            continue;
        }
        let Some(tile_a) = vehicle_current_tile(vehicle.sample_pose().position) else {
            continue;
        };

        for other in &vehicles[idx + 1..] {
            if other.is_done() {
                continue;
            }
            if vehicle.origin == other.origin && vehicle.route == other.route {
                continue;
            }
            let Some(tile_b) = vehicle_current_tile(other.sample_pose().position) else {
                continue;
            };
            if tile_a != tile_b {
                continue;
            }

            let pair = ordered_pair(vehicle.id, other.id);
            if counted_close_calls.insert(pair) {
                stats.record_close_call();
            }
        }
    }
}

// ── Geometry helpers ──────────────────────────────────────────────────────────

/// Map a pixel position to the inner 4×4 conflict-tile that contains it, or
/// `None` if the position is outside the inner zone.
///
/// Inner tile `(c, r)` corresponds to outer lane column `c+1` / row `r+1`:
/// ```text
///   x ∈ [INTER_LEFT + (c+1)*LANE_W,  INTER_LEFT + (c+2)*LANE_W)
///   y ∈ [INTER_TOP  + (r+1)*LANE_W,  INTER_TOP  + (r+2)*LANE_W)
/// ```
/// Each tile is 40 × 40 px (one lane width).  Right-turn vehicles use the
/// outer corner tiles and never land here, so they are excluded naturally.
fn vehicle_current_tile(pos: Vec2) -> Option<TileId> {
    let x0 = INTER_LEFT as f32 + LANE_W as f32; // 560
    let y0 = INTER_TOP as f32 + LANE_W as f32; // 280
    let x1 = INTER_LEFT as f32 + 5.0 * LANE_W as f32; // 720
    let y1 = INTER_TOP as f32 + 5.0 * LANE_W as f32; // 440

    if pos.x < x0 || pos.x >= x1 || pos.y < y0 || pos.y >= y1 {
        return None;
    }

    let col = ((pos.x - x0) / LANE_W as f32) as u8;
    let row = ((pos.y - y0) / LANE_W as f32) as u8;
    Some((col.min(3), row.min(3)))
}

/// True when `pos` is in the approach detection zone for the given entry direction
/// (outside the intersection box, within `DETECTION_MARGIN` of the edge).
fn is_in_detection_zone(pos: Vec2, origin: Direction) -> bool {
    match origin {
        Direction::North => {
            pos.y >= INTER_TOP as f32 - DETECTION_MARGIN && pos.y < INTER_TOP as f32
        }
        Direction::South => {
            pos.y > INTER_BOTTOM as f32 && pos.y <= INTER_BOTTOM as f32 + DETECTION_MARGIN
        }
        Direction::East => {
            pos.x > INTER_RIGHT as f32 && pos.x <= INTER_RIGHT as f32 + DETECTION_MARGIN
        }
        Direction::West => {
            pos.x >= INTER_LEFT as f32 - DETECTION_MARGIN && pos.x < INTER_LEFT as f32
        }
    }
}

/// Highest speed that lets a denied vehicle move toward, but not beyond, the
/// approach stop line during this tick.  Position is the vehicle center, so the
/// center stops half a vehicle length before the intersection edge.
fn stop_line_speed_cap(pos: Vec2, origin: Direction, dt: f32) -> f32 {
    if dt <= f32::EPSILON {
        return 0.0;
    }

    let remaining = match origin {
        Direction::North => INTER_TOP as f32 - VEHICLE_LENGTH / 2.0 - pos.y,
        Direction::South => pos.y - (INTER_BOTTOM as f32 + VEHICLE_LENGTH / 2.0),
        Direction::East => pos.x - (INTER_RIGHT as f32 + VEHICLE_LENGTH / 2.0),
        Direction::West => INTER_LEFT as f32 - VEHICLE_LENGTH / 2.0 - pos.x,
    };

    (remaining.max(0.0)) / dt
}

/// Returns true once the vehicle's rear has crossed the inner conflict box and
/// is now outside it on the route's exit side.
fn has_cleared_conflict_zone(pos: Vec2, origin: Direction, route: Route) -> bool {
    let inner_left = INTER_LEFT as f32 + LANE_W as f32;
    let inner_top = INTER_TOP as f32 + LANE_W as f32;
    let inner_right = INTER_LEFT as f32 + 5.0 * LANE_W as f32;
    let inner_bottom = INTER_TOP as f32 + 5.0 * LANE_W as f32;
    let rear_clearance = VEHICLE_LENGTH / 2.0;

    match (origin, route) {
        (Direction::North, Route::Straight) | (Direction::East, Route::Left) => {
            pos.y > inner_bottom + rear_clearance
        }
        (Direction::South, Route::Straight) | (Direction::West, Route::Left) => {
            pos.y < inner_top - rear_clearance
        }
        (Direction::East, Route::Straight) | (Direction::South, Route::Left) => {
            pos.x < inner_left - rear_clearance
        }
        (Direction::West, Route::Straight) | (Direction::North, Route::Left) => {
            pos.x > inner_right + rear_clearance
        }
        (_, Route::Right) => false,
    }
}

fn ordered_pair(a: u64, b: u64) -> (u64, u64) {
    if a < b { (a, b) } else { (b, a) }
}

/// Map a vehicle's (origin, route) to the conflict-map `RouteId` used by
/// `ReservationManager`.  Returns the canonical `(EntryDir, TurnDir)` pair.
///
/// Right turns have no `RouteId` — callers must guard against `Route::Right`
/// before calling this function.
fn to_route_id(origin: Direction, route: Route) -> RouteId {
    let entry = match origin {
        Direction::North => EntryDir::North,
        Direction::South => EntryDir::South,
        Direction::East => EntryDir::East,
        Direction::West => EntryDir::West,
    };
    let turn = match route {
        Route::Straight => TurnDir::Straight,
        Route::Left => TurnDir::Left,
        Route::Right => unreachable!("right turns do not use the conflict map"),
    };
    (entry, turn)
}

// ── Default ───────────────────────────────────────────────────────────────────

impl Default for World {
    fn default() -> Self {
        Self::new()
    }
}

// ── Tests ─────────────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;
    use crate::application::commands::Command;
    use crate::application::spawn::SpawnRequest;
    use crate::config::FIXED_DT;
    use crate::domain::path::RoutePath;
    use crate::domain::vehicle::{SpeedTier, Vehicle, VehicleState};

    #[test]
    fn world_starts_empty_until_input_arrives() {
        let world = World::new();
        assert!(world.vehicles().is_empty());
    }

    #[test]
    fn tick_advances_time() {
        let mut world = World::new();
        world.tick(1.0);
        assert!(world.sim_time > 0.0);
    }

    fn spawn_request(origin: Direction, route: Route) -> SpawnRequest {
        SpawnRequest {
            origin,
            route,
            speed_tier: SpeedTier::Medium,
        }
    }

    #[test]
    fn reservation_is_released_once_vehicle_is_on_the_exit_side() {
        let mut world = World::new();
        world.apply_command(Command::SpawnVehicle(spawn_request(
            Direction::North,
            Route::Straight,
        )));
        world.apply_command(Command::SpawnVehicle(spawn_request(
            Direction::East,
            Route::Straight,
        )));

        // North starts inside the detection zone and receives the first reservation.
        // East starts just outside its detection zone and should remain blocked until
        // North's rear clears the inner conflict zone.
        world.vehicles[0].progress.s = 210.0;
        world.vehicles[1].progress.s = 200.0;

        for _ in 0..600 {
            world.tick(FIXED_DT);
            if world.reservation_manager.is_empty() {
                break;
            }
        }

        assert!(
            world.reservation_manager.is_empty(),
            "northbound reservation should release after the vehicle clears the conflict box"
        );

        for _ in 0..600 {
            world.tick(FIXED_DT);
            if world.vehicles[1].reservation_granted {
                break;
            }
        }

        assert!(
            world.vehicles[1].reservation_granted,
            "eastbound vehicle should acquire the reservation on the next tick after release"
        );
    }

    #[test]
    fn completed_vehicle_updates_phase_seven_statistics() {
        let mut world = World::new();
        world.apply_command(Command::SpawnVehicle(spawn_request(
            Direction::North,
            Route::Straight,
        )));

        while !world.vehicles().is_empty() {
            world.tick(FIXED_DT);
        }

        let stats = world.statistics();
        assert_eq!(stats.vehicles_completed, 1);
        assert!((stats.max_speed - SpeedTier::High.units_per_second()).abs() < f32::EPSILON);
        assert!((stats.min_speed - SpeedTier::Low.units_per_second()).abs() < f32::EPSILON);

        let path = RoutePath::for_lane(Direction::North, Route::Straight);
        let speed = SpeedTier::Low.units_per_second();
        // Time-to-pass is from intersection detection until removal — always less
        // than traversing the entire route at cruise speed; bound avoids magic
        // numbers when `SpeedTier` or path geometry is retuned.
        let full_path_secs = path.total_length / speed;
        assert!(
            stats.max_time_to_pass > 0.5
                && stats.max_time_to_pass < full_path_secs + FIXED_DT * 2.0,
            "unexpected pass time: {} (full path at tier ≈ {:.3}s)",
            stats.max_time_to_pass,
            full_path_secs
        );
        assert!((stats.max_time_to_pass - stats.min_time_to_pass).abs() < f32::EPSILON);
    }

    #[test]
    fn phase_speed_rule_uses_low_high_then_medium() {
        let mut vehicle = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Low);

        vehicle.state = VehicleState::Approaching;
        vehicle.reservation_granted = false;
        assert!((phase_speed(&vehicle) - SpeedTier::Low.units_per_second()).abs() < f32::EPSILON);

        vehicle.reservation_granted = true;
        assert!((phase_speed(&vehicle) - SpeedTier::Low.units_per_second()).abs() < f32::EPSILON);

        vehicle.progress.s = 262.0; // front reaches the north stop line: y + len/2 == INTER_TOP.
        assert!((phase_speed(&vehicle) - SpeedTier::High.units_per_second()).abs() < f32::EPSILON);

        vehicle.state = VehicleState::Entering;
        assert!((phase_speed(&vehicle) - SpeedTier::High.units_per_second()).abs() < f32::EPSILON);

        vehicle.state = VehicleState::Inside;
        assert!((phase_speed(&vehicle) - SpeedTier::High.units_per_second()).abs() < f32::EPSILON);

        vehicle.state = VehicleState::Leaving;
        assert!(
            (phase_speed(&vehicle) - SpeedTier::Medium.units_per_second()).abs() < f32::EPSILON
        );
    }

    #[test]
    fn lane_safety_only_applies_before_intersection_commitment() {
        let mut inside_leader =
            Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        inside_leader.state = VehicleState::Inside;
        inside_leader.progress.s = 340.0;

        let mut inside_follower =
            Vehicle::new(2, Direction::North, Route::Straight, SpeedTier::Medium);
        inside_follower.state = VehicleState::Inside;
        inside_follower.progress.s = 330.0;

        let mut leaving_leader =
            Vehicle::new(3, Direction::South, Route::Straight, SpeedTier::Medium);
        leaving_leader.state = VehicleState::Leaving;
        leaving_leader.progress.s = 540.0;

        let mut leaving_follower =
            Vehicle::new(4, Direction::South, Route::Straight, SpeedTier::Medium);
        leaving_follower.state = VehicleState::Leaving;
        leaving_follower.progress.s = 530.0;

        let mut vehicles = vec![
            inside_leader,
            inside_follower,
            leaving_leader,
            leaving_follower,
        ];

        apply_lane_safety(&mut vehicles);

        assert!(
            (vehicles[0].target_speed - SpeedTier::High.units_per_second()).abs() < f32::EPSILON
        );
        assert!(
            (vehicles[1].target_speed - SpeedTier::High.units_per_second()).abs() < f32::EPSILON
        );
        assert!(
            (vehicles[2].target_speed - SpeedTier::Medium.units_per_second()).abs() < f32::EPSILON
        );
        assert!(
            (vehicles[3].target_speed - SpeedTier::Medium.units_per_second()).abs() < f32::EPSILON
        );
    }

    #[test]
    fn close_call_is_counted_once_per_vehicle_pair() {
        // NS-Straight travels at x=580.  Inner tile (0,0) spans
        // x ∈ [560,600) × y ∈ [280,320).
        // North at s=325 → y = -40+325 = 285  → tile (0,0).
        // East  at s=740 → x = 1320-740 = 580, y=300 → tile (0,0).
        // Both in the same inner tile → close call.
        let mut world = World::new();

        let mut north = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        north.progress.s = 325.0;
        north.state = VehicleState::Inside;

        let mut east = Vehicle::new(2, Direction::East, Route::Straight, SpeedTier::Medium);
        east.progress.s = 740.0;
        east.state = VehicleState::Inside;

        world.vehicles.push(north);
        world.vehicles.push(east);

        world.tick(0.0);
        world.tick(0.0);

        assert_eq!(world.statistics().close_calls, 1);
    }

    #[test]
    fn denied_conflicting_vehicle_waits_without_advancing() {
        let mut world = World::new();

        let mut north = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        north.progress.s = 300.0;
        north.reservation_granted = true;
        assert!(
            world
                .reservation_manager
                .request(north.id, to_route_id(north.origin, north.route))
        );

        let mut east = Vehicle::new(2, Direction::East, Route::Straight, SpeedTier::Medium);
        east.progress.s = 542.0; // x = 778, center is at the East stop line.
        let east_start = east.progress.s;

        world.vehicles.push(north);
        world.vehicles.push(east);

        world.tick(1.0);

        let east = world
            .vehicles()
            .iter()
            .find(|vehicle| vehicle.id == 2)
            .expect("east vehicle should still be active");

        assert_eq!(east.state, VehicleState::WaitingReservation);
        assert!(
            (east.progress.s - east_start).abs() < f32::EPSILON,
            "denied vehicle must not move toward a conflicting reserved path"
        );
    }

    #[test]
    fn reservation_priority_prefers_earlier_detection() {
        let now = 100.0;
        assert!(
            super::reservation_request_priority(5.0, now)
                < super::reservation_request_priority(8.0, now)
        );
    }

    #[test]
    fn reservation_priority_aging_boosts_long_wait_over_recent_arrival() {
        let now = 20.0;
        let long_wait = super::reservation_request_priority(0.0, now);
        let recent = super::reservation_request_priority(17.0, now);
        assert!(
            long_wait < recent,
            "long wait should sort before a vehicle that only recently reached the zone"
        );
    }
}
