use crate::render::*;
use macroquad::prelude::*;

/// distance from road center to lane center
const LANE_OFFSET: f32 = 36.0 / 2.0;

/// distance from intersection to spawn point
pub const SPAWN_OFFSET: f32 = 400.0;

/// safe distance between vehicles (same pitch as `can_spawn` / `move_vehicle`)
pub const SAFETY_GAP: f32 = 40.0;

/// vehicle speed
const SPEED: f32 = 2.0;

pub const VEHICLE_WIDTH: f32 = 32.;
pub const VEHICLE_HEIGHT: f32 = 32.;

#[derive(Clone, Copy, PartialEq)]
pub enum Direction {
    North,
    South,
    East,
    West,
}

#[derive(Clone, Copy, PartialEq)]
pub enum Route {
    Straight,
    Left,
    Right,
}

#[derive(Clone, Copy, PartialEq)]
pub struct Vehicle {
    pub x: f32,
    pub y: f32,
    pub direction: Direction,
    /// Driving choice at the intersection; set to `Straight` after the turn is applied.
    pub route: Route,
    /// Spawn-time route; used for colour only (green / red / blue).
    initial_route: Route,
}

impl Vehicle {
    /// Create a new vehicle spawning from the given direction
    pub fn new(direction: Direction) -> Self {
        let route = Self::random_route();
        let (x, y) = Self::spawn_pos(direction);
        Self {
            x,
            y,
            direction,
            route,
            initial_route: route,
        }
    }

    pub fn random_route() -> Route {
        // Chosen 100% randomly
        match (get_time() * 1000.0) as u32 % 3 {
            0 => Route::Straight,
            1 => Route::Left,
            _ => Route::Right,
        }
    }

    fn same_lane(a: &Vehicle, b: &Vehicle) -> bool {
        match a.direction {
            Direction::North | Direction::South => (a.x - b.x).abs() < 0.5,
            Direction::East | Direction::West => (a.y - b.y).abs() < 0.5,
        }
    }

    /// Minimum centerline separation along the road so bumpers keep `SAFETY_GAP`.
    fn min_spawn_separation() -> f32 {
        VEHICLE_HEIGHT + SAFETY_GAP
    }

    pub fn can_spawn(&self, others: &[Vehicle]) -> bool {
        others.iter().all(|other| {
            if other.direction != self.direction || !Self::same_lane(self, other) {
                return true;
            }
            let along = match self.direction {
                Direction::North | Direction::South => (self.y - other.y).abs(),
                Direction::East | Direction::West => (self.x - other.x).abs(),
            };
            along > Self::min_spawn_separation()
        })
    }

    /// Calculate spawn position outside the intersection
    pub fn spawn_pos(dir: Direction) -> (f32, f32) {
        match dir {
            // Right lane (correct side for traffic direction)
            Direction::North => (CX + LANE_OFFSET - 16.0, CY + SPAWN_OFFSET - 32.),
            Direction::South => (CX - LANE_OFFSET - 16.0, CY - SPAWN_OFFSET),
            Direction::East => (CX - SPAWN_OFFSET, CY + LANE_OFFSET - 16.0),
            Direction::West => (CX + SPAWN_OFFSET - 32., CY - LANE_OFFSET - 16.0),
        }
    }

    pub fn color(&self) -> Color {
        match self.initial_route {
            Route::Straight => GREEN,
            Route::Left => RED,
            Route::Right => BLUE,
        }
    }
    /// Vehicle center is inside the intersection box (small margin for turn trigger).
    fn at_intersection(&self) -> bool {
        let margin = 10.0;
        let cx = self.x + VEHICLE_WIDTH / 2.0;
        let cy = self.y + VEHICLE_HEIGHT / 2.0;
        cx > CX - ROAD_HALF - margin
            && cx < CX + ROAD_HALF + margin
            && cy > CY - ROAD_HALF - margin
            && cy < CY + ROAD_HALF + margin
    }

    /// Align top-left with the lane center for the current heading (same as `spawn_pos`).
    fn snap_to_lane_for_direction(&mut self) {
        match self.direction {
            Direction::North => self.x = CX + LANE_OFFSET - 16.0,
            Direction::South => self.x = CX - LANE_OFFSET - 16.0,
            Direction::East => self.y = CY + LANE_OFFSET - 16.0,
            Direction::West => self.y = CY - LANE_OFFSET - 16.0,
        }
    }

    /// True if vehicle bodies (actual rectangles) overlap. Uses **no** `SAFETY_GAP` inflate:
    /// inflating every pair by 40px made adjacent lanes and crossing paths mutually exclusive,
    /// jamming the whole intersection. In-lane spacing stays enforced by `gap_to_leader`.
    fn bodies_overlap_at(x: f32, y: f32, other: &Vehicle) -> bool {
        x < other.x + VEHICLE_WIDTH
            && x + VEHICLE_WIDTH > other.x
            && y < other.y + VEHICLE_HEIGHT
            && y + VEHICLE_HEIGHT > other.y
    }

    fn conflicts_with_any_at(x: f32, y: f32, all: &[Vehicle], self_idx: usize) -> bool {
        all.iter().enumerate().any(|(i, other)| {
            i != self_idx && Self::bodies_overlap_at(x, y, other)
        })
    }

    fn conflicts_with_any(&self, all: &[Vehicle], self_idx: usize) -> bool {
        Self::conflicts_with_any_at(self.x, self.y, all, self_idx)
    }

    /// Gap along the road from this vehicle's front bumper to the leader's rear bumper.
    fn gap_to_leader(&self, all: &[Vehicle], self_idx: usize) -> Option<f32> {
        let mut best: Option<f32> = None;
        for (i, other) in all.iter().enumerate() {
            if i == self_idx || other.direction != self.direction || !Self::same_lane(self, other) {
                continue;
            }
            let gap = match self.direction {
                Direction::North => {
                    if other.y >= self.y {
                        continue;
                    }
                    self.y - (other.y + VEHICLE_HEIGHT)
                }
                Direction::South => {
                    if other.y <= self.y {
                        continue;
                    }
                    other.y - (self.y + VEHICLE_HEIGHT)
                }
                Direction::East => {
                    if other.x <= self.x {
                        continue;
                    }
                    other.x - (self.x + VEHICLE_WIDTH)
                }
                Direction::West => {
                    if other.x >= self.x {
                        continue;
                    }
                    self.x - (other.x + VEHICLE_WIDTH)
                }
            };
            best = Some(match best {
                None => gap,
                Some(g) => g.min(gap),
            });
        }
        best
    }

    pub fn move_vehicle(&mut self, all: &[Vehicle], self_idx: usize) {
        // Apply turn + lane snap only if the exit pose clears other vehicles; otherwise wait
        // (avoids snapping into traffic and keeps `route` so we retry next frame).
        if self.at_intersection() {
            let needs_turn = matches!(self.route, Route::Left | Route::Right);
            let backup = if needs_turn { Some(*self) } else { None };

            if needs_turn {
                match self.route {
                    Route::Left => {
                        self.direction = match self.direction {
                            Direction::North => Direction::West,
                            Direction::South => Direction::East,
                            Direction::East => Direction::North,
                            Direction::West => Direction::South,
                        };
                    }
                    Route::Right => {
                        self.direction = match self.direction {
                            Direction::North => Direction::East,
                            Direction::South => Direction::West,
                            Direction::East => Direction::South,
                            Direction::West => Direction::North,
                        };
                    }
                    Route::Straight => {}
                }
                self.snap_to_lane_for_direction();
                self.route = Route::Straight;
                if self.conflicts_with_any(all, self_idx) {
                    *self = backup.unwrap();
                    return;
                }
            }
        }

        // Check gap to vehicle ahead
        let can_step = match self.gap_to_leader(all, self_idx) {
            None => true,
            Some(gap) => gap > SAFETY_GAP + SPEED,
        };
        if !can_step {
            return;
        }

        let (nx, ny) = match self.direction {
            Direction::North => (self.x, self.y - SPEED),
            Direction::South => (self.x, self.y + SPEED),
            Direction::East => (self.x + SPEED, self.y),
            Direction::West => (self.x - SPEED, self.y),
        };
        if Self::conflicts_with_any_at(nx, ny, all, self_idx) {
            return;
        }
        self.x = nx;
        self.y = ny;
    }

    pub fn draw(&self) {
        let color = self.color();
        draw_rectangle(self.x, self.y, VEHICLE_WIDTH, VEHICLE_HEIGHT, color);
    }
}
