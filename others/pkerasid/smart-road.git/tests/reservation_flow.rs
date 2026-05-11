// Integration tests for the reservation-based intersection manager.

use smart_road::application::commands::Command;
use smart_road::application::spawn::SpawnRequest;
use smart_road::application::world::World;
use smart_road::config::FIXED_DT;
use smart_road::domain::intersection::{EntryDir, TurnDir, build_conflict_map};
use smart_road::domain::lane::{Direction, Route};
use smart_road::domain::reservation::ReservationManager;
use smart_road::domain::vehicle::SpeedTier;

// ── Conflict-map smoke tests (no manager needed) ──────────────────────────────

#[test]
fn conflicting_routes_are_known_to_the_map() {
    let map = build_conflict_map();
    let ns = (EntryDir::North, TurnDir::Straight);
    let ew = (EntryDir::East, TurnDir::Straight);

    assert!(map[&ns].iter().any(|e| e.other_route == ew));
    assert!(map[&ew].iter().any(|e| e.other_route == ns));
}

#[test]
fn opposite_straights_are_not_in_each_others_conflict_list() {
    let map = build_conflict_map();
    let ns = (EntryDir::North, TurnDir::Straight);
    let sn = (EntryDir::South, TurnDir::Straight);

    assert!(!map[&ns].iter().any(|e| e.other_route == sn));
    assert!(!map[&sn].iter().any(|e| e.other_route == ns));
}

// ── ReservationManager tests ──────────────────────────────────────────────────

fn ns() -> (EntryDir, TurnDir) {
    (EntryDir::North, TurnDir::Straight)
}
fn sn() -> (EntryDir, TurnDir) {
    (EntryDir::South, TurnDir::Straight)
}
fn ew() -> (EntryDir, TurnDir) {
    (EntryDir::East, TurnDir::Straight)
}

#[test]
fn granting_reservation_blocks_conflicting_route() {
    // Given: vehicle A holds NS.
    // When:  vehicle B requests EW.
    // Then:  manager denies B because NS and EW share tile (0,0).
    let mut manager = ReservationManager::new();

    assert!(
        manager.request(1, ns()),
        "vehicle 1 should receive NS reservation"
    );
    assert!(
        !manager.request(2, ew()),
        "vehicle 2 must be denied while NS is active"
    );
}

#[test]
fn non_conflicting_routes_can_hold_simultaneous_reservations() {
    // Given: vehicle A holds NS (inner col 0).
    // When:  vehicle B requests SN (inner col 3).
    // Then:  manager grants B because NS and SN share no tiles.
    let mut manager = ReservationManager::new();

    assert!(manager.request(1, ns()));
    assert!(
        manager.request(2, sn()),
        "SN must be granted concurrently with NS"
    );
}

#[test]
fn reservation_released_on_vehicle_exit() {
    // Given: vehicle A holds NS and blocks EW.
    // When:  vehicle A exits and releases its reservation.
    // Then:  vehicle B's subsequent EW request is granted.
    let mut manager = ReservationManager::new();

    manager.request(1, ns());
    assert!(!manager.request(2, ew()));

    manager.release(1, ns()); // vehicle 1 leaves the conflict zone

    assert!(
        manager.request(2, ew()),
        "EW must be grantable after NS reservation is released"
    );
}

#[test]
fn earliest_detection_time_wins_among_conflicting_requests() {
    // `ReservationManager` is call-order sensitive.  In the full sim, `World`
    // reorders requests by detection time + aging (Phase 8c) before calling
    // `request`.  Here we model arrival order by call order: vehicle 1 first.
    let mut manager = ReservationManager::new();

    // Tick 0: both arrive at detection zone; vehicle 1 is processed first.
    assert!(manager.request(1, ns()), "first arrival wins NS");
    assert!(!manager.request(2, ew()), "second arrival blocked");

    // Vehicle 1 crosses and exits.
    manager.release(1, ns());

    // Now vehicle 2 can proceed.
    assert!(
        manager.request(2, ew()),
        "vehicle 2 proceeds after vehicle 1 clears"
    );
}

#[test]
fn every_two_approach_route_pair_completes_without_close_calls() {
    for origin_a in Direction::ALL {
        for route_a in Route::ALL {
            for origin_b in Direction::ALL {
                if origin_a == origin_b {
                    continue;
                }

                for route_b in Route::ALL {
                    let mut world = World::new();

                    world.apply_command(Command::SpawnVehicle(SpawnRequest {
                        origin: origin_a,
                        route: route_a,
                        speed_tier: SpeedTier::High,
                    }));
                    world.apply_command(Command::SpawnVehicle(SpawnRequest {
                        origin: origin_b,
                        route: route_b,
                        speed_tier: SpeedTier::High,
                    }));

                    for _ in 0..3600 {
                        if world.statistics().vehicles_completed == 2 {
                            break;
                        }
                        world.tick(FIXED_DT);
                    }

                    assert_eq!(
                        world.statistics().vehicles_completed,
                        2,
                        "pair did not complete: {origin_a:?}/{route_a:?} with {origin_b:?}/{route_b:?}"
                    );
                    assert_eq!(
                        world.statistics().close_calls,
                        0,
                        "close call for pair: {origin_a:?}/{route_a:?} with {origin_b:?}/{route_b:?}"
                    );
                }
            }
        }
    }
}
