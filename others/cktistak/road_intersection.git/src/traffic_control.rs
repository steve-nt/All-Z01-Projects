//! Traffic phases (20 s) and per-approach stop: each vehicle only uses its own signal.
//! When a lane is packed to capacity, that axis gets green early. If both axes are packed, we compare demand fairly (vertical total is doubled vs horizontal so one full N/S arm is not starved by two full E/W arms); exact ties alternate every `BOTH_JAM_ALTERNATE_SECONDS`.

use crate::render::{TrafficLightState, CX, CY, ROAD_HALF};
use crate::vehicle::{
    Direction, Vehicle, SAFETY_GAP, SPAWN_OFFSET, VEHICLE_HEIGHT, VEHICLE_WIDTH,
};

/// Duration of each phase (one axis green) in seconds.
pub const LIGHT_PHASE_SECONDS: f64 = 10.0;

/// When both axes are at capacity, alternate green this fast (seconds) if scores tie.
const BOTH_JAM_ALTERNATE_SECONDS: f64 = 5.0;

/// Usable queue segment along one lane: spawn ↔ stop line (matches `Vehicle::spawn_pos` vs intersection edge).
const LANE_QUEUE_LENGTH: f32 = SPAWN_OFFSET - ROAD_HALF - VEHICLE_HEIGHT;

const STOP_DIST: f32 = 30.0;

/// Max vehicles that fit in `lane_queue_length` with bumper-to-bumper spacing `extent + SAFETY_GAP` (no spare slot).
fn max_vehicles_packed(lane_queue_length: f32, extent: f32) -> usize {
    let pitch = extent + SAFETY_GAP;
    if lane_queue_length < extent {
        return 0;
    }
    // n cars need n*extent + (n-1)*SAFETY_GAP <= lane_queue_length  ⇒  n <= (L + SAFETY_GAP) / pitch
    (((lane_queue_length + SAFETY_GAP) / pitch).floor() as usize).max(1)
}

fn lane_capacity(dir: Direction) -> usize {
    let extent = match dir {
        Direction::North | Direction::South => VEHICLE_HEIGHT,
        Direction::East | Direction::West => VEHICLE_WIDTH,
    };
    max_vehicles_packed(LANE_QUEUE_LENGTH, extent)
}

/// State from simulation time: e.g. 0–20 s N/S green, 20–40 s E/W green, repeating.
pub fn current_traffic_light_state(now_secs: f64) -> TrafficLightState {
    let phase = (now_secs / LIGHT_PHASE_SECONDS).floor() as i64;
    let half = phase.rem_euclid(2);
    if half == 0 {
        TrafficLightState {
            north: true,
            south: true,
            east: false,
            west: false,
        }
    } else {
        TrafficLightState {
            north: false,
            south: false,
            east: true,
            west: true,
        }
    }
}

/// True while the vehicle has not fully left the intersection on its exit side (counts toward queue pressure).
fn still_on_approach_or_intersection(v: &Vehicle, dir: Direction) -> bool {
    match dir {
        Direction::North => v.y + VEHICLE_HEIGHT > CY - ROAD_HALF,
        Direction::South => v.y < CY + ROAD_HALF,
        Direction::East => v.x < CX + ROAD_HALF + VEHICLE_WIDTH,
        Direction::West => v.x + VEHICLE_WIDTH > CX - ROAD_HALF,
    }
}

fn count_queued_on_direction(vehicles: &[Vehicle], dir: Direction) -> usize {
    vehicles
        .iter()
        .filter(|v| v.direction == dir && still_on_approach_or_intersection(v, dir))
        .count()
}

fn is_lane_physically_full(vehicles: &[Vehicle], dir: Direction) -> bool {
    let cap = lane_capacity(dir);
    cap > 0 && count_queued_on_direction(vehicles, dir) >= cap
}

/// Timer-based state, overridden when at least one lane hits capacity; if both axes are full, the busier axis wins.
pub fn traffic_light_state(now_secs: f64, vehicles: &[Vehicle]) -> TrafficLightState {
    let scheduled = current_traffic_light_state(now_secs);
    let ns_jam = is_lane_physically_full(vehicles, Direction::North)
        || is_lane_physically_full(vehicles, Direction::South);
    let ew_jam = is_lane_physically_full(vehicles, Direction::East)
        || is_lane_physically_full(vehicles, Direction::West);

    let n = count_queued_on_direction(vehicles, Direction::North);
    let s = count_queued_on_direction(vehicles, Direction::South);
    let e = count_queued_on_direction(vehicles, Direction::East);
    let w = count_queued_on_direction(vehicles, Direction::West);
    let ns_sum = n + s;
    let ew_sum = e + w;
    // One vertical queue vs two horizontal arms: without weighting, E+W sum almost always beats N or S alone.
    let ns_demand = ns_sum.saturating_mul(2);
    let ew_demand = ew_sum;

    let ns_green = TrafficLightState {
        north: true,
        south: true,
        east: false,
        west: false,
    };
    let ew_green = TrafficLightState {
        north: false,
        south: false,
        east: true,
        west: true,
    };

    if ns_jam && !ew_jam {
        ns_green
    } else if ew_jam && !ns_jam {
        ew_green
    } else if ns_jam && ew_jam {
        if ns_demand > ew_demand {
            ns_green
        } else if ew_demand > ns_demand {
            ew_green
        } else {
            let half = (now_secs / BOTH_JAM_ALTERNATE_SECONDS).floor() as i64 % 2;
            if half == 0 {
                ns_green
            } else {
                ew_green
            }
        }
    } else {
        scheduled
    }
}

/// Green only for this vehicle's approach; other `TrafficLightState` fields are ignored.
pub fn is_green_for_vehicle(v: &Vehicle, state: TrafficLightState) -> bool {
    match v.direction {
        Direction::North => state.north,
        Direction::South => state.south,
        Direction::East => state.east,
        Direction::West => state.west,
    }
}

/// True when the vehicle is in the approach band *before* entering the intersection box.
/// Uses the **leading bumper** so we do not hold with the car already inside the white square.
/// Intersection interior: x ∈ [CX ± ROAD_HALF], y ∈ [CY ± ROAD_HALF].
fn is_at_stop_zone(v: &Vehicle) -> bool {
    match v.direction {
        // North: front is top edge; approach from south (large y).
        Direction::North => {
            let front_y = v.y;
            front_y >= CY + ROAD_HALF && front_y <= CY + ROAD_HALF + STOP_DIST
        }
        // South: front is bottom; approach from north (small y).
        Direction::South => {
            let front_y = v.y + VEHICLE_HEIGHT;
            front_y >= CY - ROAD_HALF - STOP_DIST && front_y <= CY - ROAD_HALF
        }
        // East: front is right edge; approach from west (small x).
        Direction::East => {
            let front_x = v.x + VEHICLE_WIDTH;
            front_x >= CX - ROAD_HALF - STOP_DIST && front_x <= CX - ROAD_HALF
        }
        // West: front is left edge; approach from east (large x).
        Direction::West => {
            let front_x = v.x;
            front_x >= CX + ROAD_HALF && front_x <= CX + ROAD_HALF + STOP_DIST
        }
    }
}

/// Move unless we are in the entry stop band with a red light for this approach.
pub fn may_move_through_light(v: &Vehicle, state: TrafficLightState) -> bool {
    !is_at_stop_zone(v) || is_green_for_vehicle(v, state)
}

/// `true` when the car must wait (at entry stop line and signal is red).
pub fn should_hold_for_light(v: &Vehicle, state: TrafficLightState) -> bool {
    !may_move_through_light(v, state)
}
