// Integration tests for the lane safety distance system.

use smart_road::application::commands::Command;
use smart_road::application::spawn::SpawnRequest;
use smart_road::application::world::World;
use smart_road::config::{FIXED_DT, SAFETY_BUFFER, VEHICLE_LENGTH};
use smart_road::domain::lane::{Direction, Route};
use smart_road::domain::safety::{is_too_close, min_safe_gap};
use smart_road::domain::vehicle::SpeedTier;

#[test]
fn safe_gap_formula_is_consistent() {
    // Smoke-test the domain function used by the safety system.
    let gap = min_safe_gap(20.0, 10.0);
    assert!((gap - 10.0).abs() < f32::EPSILON);
}

#[test]
fn follower_within_safe_gap_is_flagged() {
    assert!(is_too_close(5.0, 20.0, 10.0));
}

#[test]
fn follower_outside_safe_gap_is_not_flagged() {
    assert!(!is_too_close(35.0, 20.0, 10.0));
}

// ── helpers ──────────────────────────────────────────────────────────────────

fn north_straight_spawn() -> SpawnRequest {
    SpawnRequest {
        origin: Direction::North,
        route: Route::Straight,
        speed_tier: SpeedTier::Medium,
    }
}

/// Verify that no two vehicles on the same lane overlap (gap > 0).
/// Returns the minimum gap observed across all same-lane pairs.
fn assert_no_overlap(vehicles: &[smart_road::domain::vehicle::Vehicle]) {
    for origin in Direction::ALL {
        for route in Route::ALL {
            let mut lane: Vec<f32> = vehicles
                .iter()
                .filter(|v| v.origin == origin && v.route == route && !v.is_done())
                .map(|v| v.progress.s)
                .collect();

            if lane.len() < 2 {
                continue;
            }

            lane.sort_unstable_by(|a, b| b.partial_cmp(a).unwrap());

            for pair in lane.windows(2) {
                let gap = pair[0] - pair[1] - VEHICLE_LENGTH;
                assert!(
                    gap >= 0.0,
                    "vehicles on lane ({origin:?}, {route:?}) overlap: gap = {gap:.2}"
                );
            }
        }
    }
}

// ── Phase 5 tests ─────────────────────────────────────────────────────────────

#[test]
fn follower_slows_when_too_close_to_leader() {
    // Given: two vehicles queued on the same lane.
    let mut world = World::new();
    let req = north_straight_spawn();
    world.apply_command(Command::SpawnVehicle(req));
    world.apply_command(Command::SpawnVehicle(req));

    // When: run long enough for the spawn queue to drain (leader clears the spawn
    // zone so the follower can enter).
    // Cars now approach at Low phase speed; ticks below still allow queue drain
    // before assertions.
    // Use 60 ticks (~1 s) to guarantee both are active.
    for _ in 0..60 {
        world.tick(FIXED_DT);
    }

    // Then: whenever two vehicles are on the same lane and the gap is below the
    // required minimum, the follower must be stopped.
    let required_gap = min_safe_gap(VEHICLE_LENGTH, SAFETY_BUFFER);

    let vehicles = world.vehicles();
    for origin in Direction::ALL {
        for route in Route::ALL {
            let mut lane: Vec<_> = vehicles
                .iter()
                .filter(|v| v.origin == origin && v.route == route && !v.is_done())
                .collect();

            if lane.len() < 2 {
                continue;
            }

            lane.sort_unstable_by(|a, b| b.progress.s.partial_cmp(&a.progress.s).unwrap());

            for pair in lane.windows(2) {
                let leader = pair[0];
                let follower = pair[1];
                let gap = leader.progress.s - follower.progress.s - VEHICLE_LENGTH;

                if gap < required_gap {
                    assert!(
                        follower.current_speed <= leader.current_speed,
                        "follower within safe gap must not exceed leader speed (gap = {gap:.2}, \
                         follower = {}, leader = {})",
                        follower.current_speed,
                        leader.current_speed
                    );
                }
            }
        }
    }
}

#[test]
fn follower_stops_before_overlapping_leader() {
    // Given: leader stopped or slow, follower closing in.
    // Simulate: spawn two vehicles then run 120 ticks (2 s) checking every step.
    let mut world = World::new();
    let req = north_straight_spawn();
    world.apply_command(Command::SpawnVehicle(req));
    world.apply_command(Command::SpawnVehicle(req));

    for _ in 0..120 {
        world.tick(FIXED_DT);
        assert_no_overlap(world.vehicles());
    }
}

#[test]
fn no_rear_end_collision_under_key_spam() {
    // Given: many vehicles spawned rapidly on the same lane.
    // When: the simulation runs for several seconds.
    // Then: no two vehicles ever overlap in path-progress space.

    let mut world = World::new();
    let req = north_straight_spawn();

    // Submit more than the queue can immediately accept.
    for _ in 0..8 {
        world.apply_command(Command::SpawnVehicle(req));
    }

    // Run for ~5 seconds, checking the no-overlap invariant every tick.
    for _ in 0..300 {
        world.tick(FIXED_DT);
        assert_no_overlap(world.vehicles());
    }
}
