#![allow(clippy::cast_possible_wrap, clippy::cast_sign_loss)]

use sdl2::pixels::{Color, PixelFormatEnum};
use sdl2::rect::Rect;
use sdl2::render::{BlendMode, Texture, TextureCreator};
use sdl2::surface::Surface;
use sdl2::video::WindowContext;

pub const VEHICLE_TEXTURE_W: u32 = 36;
pub const VEHICLE_TEXTURE_H: u32 = 20;

/// Phase 8a: richer top-down AV sprite.  Long axis is +X (front / hood at high x)
/// so `heading_deg` from path tangents aligns with motion without an extra offset.
///
/// # Errors
///
/// Returns an SDL error if surface allocation or texture upload fails.
pub fn build_vehicle_texture(
    texture_creator: &TextureCreator<WindowContext>,
) -> Result<Texture<'_>, String> {
    let mut surface = Surface::new(
        VEHICLE_TEXTURE_W,
        VEHICLE_TEXTURE_H,
        PixelFormatEnum::RGBA32,
    )?;

    paint_vehicle_rgba(&mut surface, false)?;

    let mut tex = texture_creator
        .create_texture_from_surface(surface)
        .map_err(|e| e.to_string())?;
    tex.set_blend_mode(BlendMode::Blend);
    Ok(tex)
}

/// Soft drop-shadow blob (same footprint) for depth under each vehicle.
///
/// # Errors
///
/// Returns an SDL error if surface allocation or texture upload fails.
pub fn build_vehicle_shadow_texture(
    texture_creator: &TextureCreator<WindowContext>,
) -> Result<Texture<'_>, String> {
    let mut surface = Surface::new(
        VEHICLE_TEXTURE_W,
        VEHICLE_TEXTURE_H,
        PixelFormatEnum::RGBA32,
    )?;

    paint_vehicle_rgba(&mut surface, true)?;

    let mut tex = texture_creator
        .create_texture_from_surface(surface)
        .map_err(|e| e.to_string())?;
    tex.set_blend_mode(BlendMode::Blend);
    Ok(tex)
}

/// `shadow_only`: flat translucent silhouette; otherwise full paint.
fn paint_vehicle_rgba(surface: &mut Surface<'_>, shadow_only: bool) -> Result<(), String> {
    let w = VEHICLE_TEXTURE_W as i32;
    let h = VEHICLE_TEXTURE_H as i32;

    if shadow_only {
        surface.fill_rect(None, Color::RGBA(0, 0, 0, 0))?;
        // Rounded-ish footprint: inset body alpha falloff via layered rects.
        for (i, alpha) in [(0, 52), (1, 38), (2, 20)] {
            surface.fill_rect(
                Rect::new(i, i, (w - 2 * i) as u32, (h - 2 * i) as u32),
                Color::RGBA(0, 0, 0, alpha),
            )?;
        }
        return Ok(());
    }

    // Palette — cool metallic shell, readable at 1:1 pixel scale.
    let outline = Color::RGBA(18, 22, 28, 255);
    let body_top = Color::RGBA(62, 118, 168, 255);
    let body_mid = Color::RGBA(44, 92, 138, 255);
    let body_shadow = Color::RGBA(28, 56, 86, 255);
    let cabin = Color::RGBA(190, 210, 224, 255);
    let cabin_divider = Color::RGBA(120, 140, 158, 255);
    let glass_hi = Color::RGBA(210, 232, 246, 255);
    let grille = Color::RGBA(35, 38, 44, 255);
    let bumper = Color::RGBA(55, 58, 64, 255);
    let head = Color::RGBA(255, 248, 210, 255);
    let tail = Color::RGBA(255, 72, 72, 255);
    let mirror = Color::RGBA(70, 74, 82, 255);
    let wheel = Color::RGBA(24, 24, 28, 255);
    let hub = Color::RGBA(90, 94, 100, 255);

    surface.fill_rect(None, Color::RGBA(0, 0, 0, 0))?;

    // Outer silhouette (1 px border).
    surface.fill_rect(Rect::new(1, 1, (w - 2) as u32, (h - 2) as u32), outline)?;
    surface.fill_rect(Rect::new(2, 2, (w - 4) as u32, (h - 4) as u32), body_mid)?;

    // Roof / cabin band (slightly narrower than body).
    surface.fill_rect(Rect::new(8, 3, 20, 14), body_top)?;
    surface.fill_rect(Rect::new(9, 4, 18, 12), cabin)?;
    // Windshield vs rear glass split.
    surface.fill_rect(Rect::new(10, 5, 7, 10), glass_hi)?;
    surface.fill_rect(Rect::new(19, 5, 7, 10), cabin)?;
    surface.fill_rect(Rect::new(17, 4, 2, 12), cabin_divider)?;

    // Hood (front, +X) and trunk shading.
    surface.fill_rect(Rect::new(w - 9, 2, 7, (h - 4) as u32), body_top)?;
    surface.fill_rect(Rect::new(2, 2, 7, (h - 4) as u32), body_shadow)?;

    // Grille + bumper at front.
    surface.fill_rect(Rect::new(w - 6, 4, 4, 12), grille)?;
    surface.fill_rect(Rect::new(w - 2, 3, 2, (h - 6) as u32), bumper)?;

    // Headlights (front corners) and tail lights.
    surface.fill_rect(Rect::new(w - 5, 1, 3, 2), head)?;
    surface.fill_rect(Rect::new(w - 5, h - 3, 3, 2), head)?;
    surface.fill_rect(Rect::new(1, 2, 2, 3), tail)?;
    surface.fill_rect(Rect::new(1, h - 5, 2, 3), tail)?;

    // Side mirrors.
    surface.fill_rect(Rect::new(w - 11, 0, 2, 2), mirror)?;
    surface.fill_rect(Rect::new(w - 11, h - 2, 2, 2), mirror)?;

    // Wheels (four corners), compact for 20 px height.
    for (x, y) in [(5, 1), (5, h - 6), (w - 9, 1), (w - 9, h - 6)] {
        surface.fill_rect(Rect::new(x, y, 4, 4), wheel)?;
        surface.fill_rect(Rect::new(x + 1, y + 1, 2, 2), hub)?;
    }

    Ok(())
}
