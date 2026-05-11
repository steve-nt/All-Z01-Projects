//! Top-level scene coordinator.
//!
//! Reads `World` state and delegates to individual view modules.

use sdl2::render::Canvas;
use sdl2::video::Window;

use crate::application::world::World;
use crate::domain::stats::Statistics;
use crate::presentation::{map_view, stats_view, vehicle_view};

/// Render one frame. Called every iteration of the game loop.
///
/// # Errors
///
/// Propagates SDL rendering failures from vehicle drawing.
pub fn render(canvas: &mut Canvas<Window>, world: &World) -> Result<(), String> {
    // Render order matters: map first, then moving vehicles on top.
    map_view::render_map(canvas);
    vehicle_view::render_vehicles(canvas, world)
}

pub fn render_statistics(canvas: &mut Canvas<Window>, stats: &Statistics) {
    stats_view::render_statistics(canvas, stats);
}
