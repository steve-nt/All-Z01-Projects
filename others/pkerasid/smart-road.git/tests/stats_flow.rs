use smart_road::application::commands::Command;
use smart_road::application::spawn::SpawnRequest;
use smart_road::application::world::World;
use smart_road::config::FIXED_DT;
use smart_road::domain::lane::{Direction, Route};
use smart_road::domain::vehicle::SpeedTier;

fn north_straight_spawn() -> SpawnRequest {
    SpawnRequest {
        origin: Direction::North,
        route: Route::Straight,
        speed_tier: SpeedTier::Medium,
    }
}

#[test]
fn completing_a_vehicle_updates_final_statistics() {
    let mut world = World::new();
    world.apply_command(Command::SpawnVehicle(north_straight_spawn()));

    while !world.vehicles().is_empty() {
        world.tick(FIXED_DT);
    }

    let stats = world.statistics();
    assert_eq!(stats.vehicles_completed, 1);
    assert!((stats.max_speed - SpeedTier::High.units_per_second()).abs() < f32::EPSILON);
    assert!((stats.min_speed - SpeedTier::Low.units_per_second()).abs() < f32::EPSILON);
    assert!(stats.max_time_to_pass > 0.0);
    assert!((stats.max_time_to_pass - stats.min_time_to_pass).abs() < f32::EPSILON);
    assert_eq!(stats.close_calls, 0);
}
