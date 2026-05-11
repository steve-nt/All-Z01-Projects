#![allow(clippy::cast_possible_truncation, clippy::cast_sign_loss)]

use sdl2::rect::{Point, Rect};
use sdl2::render::{BlendMode, Canvas};
use sdl2::video::Window;

use crate::application::world::World;
use crate::infrastructure::sdl::texture_store::{self, VEHICLE_TEXTURE_H, VEHICLE_TEXTURE_W};

/// Screen-space offset for the soft shadow (light from upper-left).
const SHADOW_OFFSET_X: i32 = 3;
const SHADOW_OFFSET_Y: i32 = 3;

/// Draw all active vehicles for the current frame.
///
/// # Errors
///
/// Propagates SDL failures while creating or drawing the vehicle texture.
pub fn render_vehicles(canvas: &mut Canvas<Window>, world: &World) -> Result<(), String> {
    let texture_creator = canvas.texture_creator();
    let vehicle_tex = texture_store::build_vehicle_texture(&texture_creator)?;
    let shadow_tex = texture_store::build_vehicle_shadow_texture(&texture_creator)?;

    canvas.set_blend_mode(BlendMode::Blend);

    for vehicle in world.vehicles() {
        if vehicle.is_done() {
            continue;
        }

        // Pose comes from the domain path sample; presentation only draws it.
        let pose = vehicle.sample_pose();
        let cx = pose.position.x as i32;
        let cy = pose.position.y as i32;
        let heading = f64::from(pose.heading_deg);

        let dst_shadow = Rect::from_center(
            Point::new(cx + SHADOW_OFFSET_X, cy + SHADOW_OFFSET_Y),
            VEHICLE_TEXTURE_W,
            VEHICLE_TEXTURE_H,
        );
        canvas.copy_ex(&shadow_tex, None, dst_shadow, heading, None, false, false)?;

        let dst = Rect::from_center(Point::new(cx, cy), VEHICLE_TEXTURE_W, VEHICLE_TEXTURE_H);
        // Rotation follows the path tangent so turns animate continuously.
        canvas.copy_ex(&vehicle_tex, None, dst, heading, None, false, false)?;
    }

    canvas.set_blend_mode(BlendMode::None);
    Ok(())
}
