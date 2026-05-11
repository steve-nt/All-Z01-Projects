use std::array;
use std::collections::VecDeque;

use crate::config::{SAFETY_BUFFER, VEHICLE_LENGTH};
use crate::domain::lane::Direction;
use crate::domain::safety::min_spawn_center_distance;
use crate::domain::vehicle::{SpeedTier, Vehicle};

const SPAWN_COOLDOWN_SECS: f32 = 1.0;

/// Pure application data describing a vehicle that should be spawned.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct SpawnRequest {
    pub origin: Direction,
    pub route: crate::domain::lane::Route,
    pub speed_tier: SpeedTier,
}

#[derive(Debug)]
pub struct SpawnSystem {
    next_vehicle_id: u64,
    queues: [VecDeque<SpawnRequest>; 4],
    cooldowns: [f32; 4],
}

impl SpawnSystem {
    #[must_use]
    pub fn new() -> Self {
        Self {
            next_vehicle_id: 1,
            queues: array::from_fn(|_| VecDeque::new()),
            cooldowns: [0.0; 4],
        }
    }

    pub fn submit(&mut self, vehicles: &mut Vec<Vehicle>, request: SpawnRequest) {
        if self.can_spawn_now(vehicles, request) {
            self.spawn_vehicle(vehicles, request);
        } else {
            self.queues[request.origin.index()].push_back(request);
        }
    }

    pub fn tick(&mut self, dt: f32, vehicles: &mut Vec<Vehicle>) {
        for cooldown in &mut self.cooldowns {
            *cooldown = (*cooldown - dt).max(0.0);
        }

        for direction in Direction::ALL {
            let idx = direction.index();
            let Some(request) = self.queues[idx].front().copied() else {
                continue;
            };

            if self.can_spawn_now(vehicles, request) {
                let _dropped = self.queues[idx].pop_front();
                self.spawn_vehicle(vehicles, request);
            }
        }
    }

    #[must_use]
    pub fn queue_len(&self, direction: Direction) -> usize {
        self.queues[direction.index()].len()
    }

    fn can_spawn_now(&self, vehicles: &[Vehicle], request: SpawnRequest) -> bool {
        self.cooldowns[request.origin.index()] <= f32::EPSILON
            && lane_entry_is_clear(vehicles, request)
    }

    fn spawn_vehicle(&mut self, vehicles: &mut Vec<Vehicle>, request: SpawnRequest) {
        let vehicle = Vehicle::new(
            self.next_vehicle_id,
            request.origin,
            request.route,
            request.speed_tier,
        );
        self.next_vehicle_id += 1;
        self.cooldowns[request.origin.index()] = SPAWN_COOLDOWN_SECS;
        vehicles.push(vehicle);
    }
}

impl Default for SpawnSystem {
    fn default() -> Self {
        Self::new()
    }
}

fn lane_entry_is_clear(vehicles: &[Vehicle], request: SpawnRequest) -> bool {
    let spawn_position = Vehicle::new(0, request.origin, request.route, request.speed_tier)
        .sample_pose()
        .position;
    let min_gap = min_spawn_center_distance(VEHICLE_LENGTH, SAFETY_BUFFER);

    vehicles
        .iter()
        .filter(|vehicle| vehicle.origin == request.origin && !vehicle.is_done())
        .all(|vehicle| vehicle.sample_pose().position.distance_to(spawn_position) >= min_gap)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::domain::lane::Route;

    fn north_spawn() -> SpawnRequest {
        SpawnRequest {
            origin: Direction::North,
            route: Route::Straight,
            speed_tier: SpeedTier::Medium,
        }
    }

    #[test]
    fn blocked_requests_are_queued() {
        let mut system = SpawnSystem::new();
        let mut vehicles = Vec::new();

        system.submit(&mut vehicles, north_spawn());
        system.submit(&mut vehicles, north_spawn());

        assert_eq!(vehicles.len(), 1);
        assert_eq!(system.queue_len(Direction::North), 1);
    }

    #[test]
    fn queued_request_spawns_after_cooldown_and_clearance() {
        let mut system = SpawnSystem::new();
        let mut vehicles = Vec::new();

        system.submit(&mut vehicles, north_spawn());
        system.submit(&mut vehicles, north_spawn());
        vehicles[0].tick(4.0);
        system.tick(4.0, &mut vehicles);

        assert_eq!(vehicles.len(), 2);
        assert_eq!(system.queue_len(Direction::North), 0);
    }
}
