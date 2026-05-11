#![allow(clippy::cast_precision_loss)]

use crate::config::{INTER_BOTTOM, INTER_LEFT, INTER_RIGHT, INTER_TOP};
use crate::domain::lane::{Direction, Route};
use crate::domain::path::{PathProgress, PathSample, RoutePath};

/// Fixed speed presets used by the simulation phases before acceleration exists.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum SpeedTier {
    Low,
    Medium,
    High,
}

impl SpeedTier {
    pub const ALL: [Self; 3] = [Self::Low, Self::Medium, Self::High];

    #[must_use]
    pub const fn units_per_second(self) -> f32 {
        // Phase 8b: tuned for ~1280×720 layout — calmer flow, clearer tier spread (px/s).
        match self {
            Self::Low => 50.0,
            Self::Medium => 100.0,
            Self::High => 150.0,
        }
    }
}

/// High-level lifecycle states for a vehicle moving through the intersection.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum VehicleState {
    Approaching,
    WaitingReservation,
    Entering,
    Inside,
    Leaving,
    Done,
}

/// Domain vehicle entity: route choice, speed, progress, and sampled path.
#[derive(Debug, Clone)]
pub struct Vehicle {
    pub id: u64,
    pub origin: Direction,
    pub route: Route,
    pub state: VehicleState,
    pub progress: PathProgress,
    pub current_speed: f32,
    pub target_speed: f32,
    pub speed_tier: SpeedTier,
    /// Set to `true` by the application layer when the intersection manager has
    /// granted this vehicle a reservation.  `false` means the vehicle must stop
    /// at the detection line.  Right-turn vehicles are granted immediately
    /// (their corner tiles are outside the conflict zone).
    pub reservation_granted: bool,
    /// Set to `true` the first time the vehicle's geometric state becomes `Inside`.
    /// Used by the application layer to know when to release the reservation
    /// (after the vehicle has been inside and is now exiting).
    pub entered_intersection: bool,
    path: RoutePath,
}

impl Vehicle {
    #[must_use]
    pub fn new(id: u64, origin: Direction, route: Route, speed_tier: SpeedTier) -> Self {
        // Phase 3 builds the geometric route once and reuses it for every tick.
        let path = RoutePath::for_lane(origin, route);
        let current_speed = SpeedTier::Low.units_per_second();

        Self {
            id,
            origin,
            route,
            state: VehicleState::Approaching,
            progress: PathProgress::new(path.total_length),
            current_speed,
            target_speed: current_speed,
            speed_tier,
            reservation_granted: false,
            entered_intersection: false,
            path,
        }
    }

    pub fn tick(&mut self, dt: f32) {
        if self.progress.is_complete() {
            self.state = VehicleState::Done;
            return;
        }

        // For now, speed changes are immediate; later phases can add acceleration.
        self.current_speed = self.target_speed;
        self.progress.advance(self.current_speed, dt);
        if self.state == VehicleState::WaitingReservation && !self.reservation_granted {
            return;
        }
        self.state = classify_state(self.sample_pose(), self.progress.is_complete());
    }

    #[must_use]
    pub fn sample_pose(&self) -> PathSample {
        self.path.sample(self.progress.s)
    }

    #[must_use]
    pub fn is_done(&self) -> bool {
        self.state == VehicleState::Done
    }
}

/// Geometric state classification used as the base for tick updates.
/// The application layer may override this with `WaitingReservation` when the
/// vehicle is stopped at the detection line without a granted reservation.
fn classify_state(sample: PathSample, is_complete: bool) -> VehicleState {
    if is_complete {
        VehicleState::Done
    } else if is_inside_intersection(sample) {
        VehicleState::Inside
    } else if is_in_approach_zone(sample) {
        VehicleState::Entering
    } else if is_in_exit_zone(sample) {
        VehicleState::Leaving
    } else {
        VehicleState::Approaching
    }
}

// Intersection bounds match the road box used by the renderer.
fn is_inside_intersection(sample: PathSample) -> bool {
    sample.position.x >= INTER_LEFT as f32
        && sample.position.x <= INTER_RIGHT as f32
        && sample.position.y >= INTER_TOP as f32
        && sample.position.y <= INTER_BOTTOM as f32
}

fn is_in_approach_zone(sample: PathSample) -> bool {
    match heading_axis(sample.heading_deg) {
        HeadingAxis::North => {
            sample.position.y > INTER_BOTTOM as f32
                && sample.position.y <= INTER_BOTTOM as f32 + 80.0
        }
        HeadingAxis::South => {
            sample.position.y >= INTER_TOP as f32 - 80.0 && sample.position.y < INTER_TOP as f32
        }
        HeadingAxis::East => {
            sample.position.x >= INTER_LEFT as f32 - 80.0 && sample.position.x < INTER_LEFT as f32
        }
        HeadingAxis::West => {
            sample.position.x > INTER_RIGHT as f32 && sample.position.x <= INTER_RIGHT as f32 + 80.0
        }
    }
}

fn is_in_exit_zone(sample: PathSample) -> bool {
    match heading_axis(sample.heading_deg) {
        HeadingAxis::North => {
            sample.position.y >= INTER_TOP as f32 - 80.0 && sample.position.y < INTER_TOP as f32
        }
        HeadingAxis::South => {
            sample.position.y > INTER_BOTTOM as f32
                && sample.position.y <= INTER_BOTTOM as f32 + 80.0
        }
        HeadingAxis::East => {
            sample.position.x > INTER_RIGHT as f32 && sample.position.x <= INTER_RIGHT as f32 + 80.0
        }
        HeadingAxis::West => {
            sample.position.x >= INTER_LEFT as f32 - 80.0 && sample.position.x < INTER_LEFT as f32
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum HeadingAxis {
    North,
    South,
    East,
    West,
}

fn heading_axis(heading_deg: f32) -> HeadingAxis {
    if (-45.0..45.0).contains(&heading_deg) {
        HeadingAxis::East
    } else if (45.0..135.0).contains(&heading_deg) {
        HeadingAxis::South
    } else if (-135.0..-45.0).contains(&heading_deg) {
        HeadingAxis::North
    } else {
        HeadingAxis::West
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_vehicle_starts_at_path_origin() {
        let vehicle = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        let pose = vehicle.sample_pose();

        assert_eq!(vehicle.id, 1);
        assert_eq!(vehicle.state, VehicleState::Approaching);
        assert!(pose.position.y < 0.0);
    }

    #[test]
    fn tick_advances_vehicle_progress() {
        let mut vehicle = Vehicle::new(1, Direction::West, Route::Straight, SpeedTier::Medium);
        vehicle.tick(1.0);

        assert!(vehicle.progress.s > 0.0);
    }

    #[test]
    fn vehicle_reaches_done_state_at_path_end() {
        let mut vehicle = Vehicle::new(1, Direction::West, Route::Straight, SpeedTier::High);
        let dt = vehicle.progress.total_length / vehicle.current_speed + 1.0;
        vehicle.tick(dt);

        assert!(vehicle.is_done());
    }

    #[test]
    fn vehicle_enters_before_crossing_the_box() {
        let mut vehicle = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        vehicle.progress.s = 210.0;
        vehicle.tick(0.0);

        assert_eq!(vehicle.state, VehicleState::Entering);
    }

    #[test]
    fn vehicle_leaves_after_clearing_the_box() {
        let mut vehicle = Vehicle::new(1, Direction::North, Route::Straight, SpeedTier::Medium);
        vehicle.progress.s = 540.0;
        vehicle.tick(0.0);

        assert_eq!(vehicle.state, VehicleState::Leaving);
    }
}
