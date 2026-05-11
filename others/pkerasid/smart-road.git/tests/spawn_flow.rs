use rand::SeedableRng;
use rand::rngs::StdRng;

use smart_road::application::commands::Command;
use smart_road::application::spawn::SpawnRequest;
use smart_road::application::world::World;
use smart_road::domain::lane::{Direction, Route};
use smart_road::domain::path::RoutePath;
use smart_road::domain::vehicle::SpeedTier;
use smart_road::infrastructure::random::vehicle_generator::VehicleGenerator;

#[test]
fn spawning_vehicle_places_it_at_lane_entry() {
    let mut world = World::new();
    let request = SpawnRequest {
        origin: Direction::North,
        route: Route::Straight,
        speed_tier: SpeedTier::Medium,
    };

    world.apply_command(Command::SpawnVehicle(request));

    let vehicle = &world.vehicles()[0];
    let actual = vehicle.sample_pose();
    let expected = RoutePath::for_lane(Direction::North, Route::Straight).sample(0.0);

    assert_eq!(world.vehicles().len(), 1);
    assert_eq!(actual.position, expected.position);
}

#[test]
fn spawn_is_blocked_when_lane_entry_is_occupied() {
    let mut world = World::new();
    let request = SpawnRequest {
        origin: Direction::North,
        route: Route::Straight,
        speed_tier: SpeedTier::Medium,
    };

    world.apply_command(Command::SpawnVehicle(request));
    world.apply_command(Command::SpawnVehicle(request));

    assert_eq!(world.vehicles().len(), 1);
    assert_eq!(world.queue_len(Direction::North), 1);
}

#[test]
fn queued_spawn_fires_when_entry_clears() {
    let mut world = World::new();
    let request = SpawnRequest {
        origin: Direction::North,
        route: Route::Straight,
        speed_tier: SpeedTier::High,
    };

    world.apply_command(Command::SpawnVehicle(request));
    world.apply_command(Command::SpawnVehicle(request));
    world.tick(4.0);

    assert_eq!(world.vehicles().len(), 2);
    assert_eq!(world.queue_len(Direction::North), 0);
}

#[test]
fn random_spawning_never_creates_overlapping_vehicles() {
    let mut world = World::new();
    let mut generator = VehicleGenerator::with_rng(StdRng::seed_from_u64(42));

    for _ in 0..60 {
        for _ in 0..4 {
            let request = generator.random_spawn_request();
            world.apply_command(Command::SpawnVehicle(request));
        }
        world.tick(1.0 / 60.0);

        let vehicles = world.vehicles();
        for (idx, vehicle) in vehicles.iter().enumerate() {
            let position = vehicle.sample_pose().position;
            for other in &vehicles[idx + 1..] {
                let other_position = other.sample_pose().position;
                assert!(
                    position.distance_to(other_position) > 1.0,
                    "vehicles overlapped at spawn: {} and {}",
                    vehicle.id,
                    other.id
                );
            }
        }
    }
}
