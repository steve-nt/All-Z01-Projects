//! Static road-layout renderer.
//!
//! Draws the 4-way intersection: asphalt surfaces, edge lines, lane dividers,
//! yellow centre lines, and stop lines.  No vehicle state is read here.
//!
//! # Coordinate system
//!
//! ```text
//! (0,0)──────────────────────────────── x →
//!  │   grass │  N approach  │  grass
//!  │         x=520       x=760
//!  │         ┌────────────────────────── y=240
//!  │  W app  │ intersection │  E app
//!  │         └────────────────────────── y=480
//!  │   grass │  S approach  │  grass
//!  y ↓
//! ```
//!
//! Road width = 240 px (3 lanes × 40 px × 2 directions).

#![allow(clippy::cast_possible_wrap, clippy::cast_possible_truncation)]

use sdl2::pixels::Color;
use sdl2::rect::{Point, Rect};
use sdl2::render::Canvas;
use sdl2::video::Window;

use crate::config::{
    INTER_BOTTOM, INTER_LEFT, INTER_RIGHT, INTER_TOP, LANE_W, ROAD_W, WINDOW_H, WINDOW_W,
};

// ── Palette ──────────────────────────────────────────────────────────────────

const GRASS: Color = Color::RGB(86, 125, 70);
const ASPHALT: Color = Color::RGB(55, 55, 55);
const EDGE: Color = Color::RGB(210, 210, 210);
const LANE_MARK: Color = Color::RGB(160, 160, 160);
const CENTRE: Color = Color::RGB(255, 210, 0);
const STOP_LINE: Color = Color::RGB(240, 240, 240);

// ── Geometry helpers (i32 aliases used throughout SDL2 calls) ─────────────────

const IL: i32 = INTER_LEFT as i32; // 520
const IT: i32 = INTER_TOP as i32; // 240
const IR: i32 = INTER_RIGHT as i32; // 760
const IB: i32 = INTER_BOTTOM as i32; // 480
const CX: i32 = (WINDOW_W / 2) as i32; // 640  (centre x)
const CY: i32 = (WINDOW_H / 2) as i32; // 360  (centre y)
const LW: i32 = LANE_W as i32; // 40
const WW: i32 = WINDOW_W as i32; // 1280
const WH: i32 = WINDOW_H as i32; // 720
const RW: i32 = ROAD_W as i32; // 240

// Dash pattern for lane dividers.
const DASH: i32 = 20;
const GAP: i32 = 15;

// ── Public entry point ────────────────────────────────────────────────────────

/// Draw the complete static road layout onto `canvas`.
pub fn render_map(canvas: &mut Canvas<Window>) {
    // 1. Grass background.
    canvas.set_draw_color(GRASS);
    canvas.clear();

    // 2. Road surfaces.
    draw_surfaces(canvas);

    // 3. White edge lines (road boundary).
    draw_edges(canvas);

    // 4. Yellow centre lines (divide directions).
    draw_centre_lines(canvas);

    // 5. White dashed lane dividers.
    draw_lane_dividers(canvas);

    // 6. White stop lines at each approach entrance.
    draw_stop_lines(canvas);
}

// ── Road surfaces ─────────────────────────────────────────────────────────────

fn draw_surfaces(canvas: &mut Canvas<Window>) {
    canvas.set_draw_color(ASPHALT);

    // Vertical road (N/S), full window height.
    let _ = canvas.fill_rect(Rect::new(IL, 0, RW as u32, WINDOW_H));

    // Horizontal road (E/W), full window width.
    let _ = canvas.fill_rect(Rect::new(0, IT, WINDOW_W, RW as u32));
}

// ── Road edges ────────────────────────────────────────────────────────────────

fn draw_edges(canvas: &mut Canvas<Window>) {
    canvas.set_draw_color(EDGE);

    // Vertical road — left edge (x = IL), skip the intersection box.
    hline(canvas, IL, 0, IT);
    hline(canvas, IL, IB, WH);

    // Vertical road — right edge (x = IR - 1), skip the intersection box.
    hline(canvas, IR - 1, 0, IT);
    hline(canvas, IR - 1, IB, WH);

    // Horizontal road — top edge (y = IT), skip the intersection box.
    vline(canvas, IT, 0, IL);
    vline(canvas, IT, IR, WW);

    // Horizontal road — bottom edge (y = IB - 1), skip the intersection box.
    vline(canvas, IB - 1, 0, IL);
    vline(canvas, IB - 1, IR, WW);
}

// ── Centre lines (yellow, separate inbound from outbound) ─────────────────────

fn draw_centre_lines(canvas: &mut Canvas<Window>) {
    canvas.set_draw_color(CENTRE);

    // Vertical road centre (x = CX).
    hline(canvas, CX, 0, IT);
    hline(canvas, CX, IB, WH);

    // Horizontal road centre (y = CY).
    vline(canvas, CY, 0, IL);
    vline(canvas, CY, IR, WW);
}

// ── Lane dividers (white dashed, within same-direction bands) ─────────────────

fn draw_lane_dividers(canvas: &mut Canvas<Window>) {
    canvas.set_draw_color(LANE_MARK);

    // ── Vertical road (N/S) ──────────────────────────────────────────────────
    //
    // Southbound (left half, x = IL..CX):  dividers at IL+LW, IL+2*LW
    // Northbound (right half, x = CX..IR): dividers at CX+LW, CX+2*LW
    //
    // Draw only in the N approach (y=0..IT) and S approach (y=IB..WH).
    for offset in [LW, 2 * LW] {
        dashed_hline(canvas, IL + offset, 0, IT);
        dashed_hline(canvas, IL + offset, IB, WH);
        dashed_hline(canvas, CX + offset, 0, IT);
        dashed_hline(canvas, CX + offset, IB, WH);
    }

    // ── Horizontal road (E/W) ─────────────────────────────────────────────────
    //
    // Westbound (top half, y = IT..CY):  dividers at IT+LW, IT+2*LW
    // Eastbound (bottom half, y = CY..IB): dividers at CY+LW, CY+2*LW
    //
    // Draw only in the W approach (x=0..IL) and E approach (x=IR..WW).
    for offset in [LW, 2 * LW] {
        dashed_vline(canvas, IT + offset, 0, IL);
        dashed_vline(canvas, IT + offset, IR, WW);
        dashed_vline(canvas, CY + offset, 0, IL);
        dashed_vline(canvas, CY + offset, IR, WW);
    }
}

// ── Stop lines ────────────────────────────────────────────────────────────────

fn draw_stop_lines(canvas: &mut Canvas<Window>) {
    canvas.set_draw_color(STOP_LINE);

    // North approach: southbound lanes (x = IL..CX) stop at y = IT.
    let _ = canvas.fill_rect(Rect::new(IL, IT - 3, (CX - IL) as u32, 3));

    // South approach: northbound lanes (x = CX..IR) stop at y = IB.
    let _ = canvas.fill_rect(Rect::new(CX, IB, (IR - CX) as u32, 3));

    // East approach: westbound lanes (y = IT..CY) stop at x = IR.
    let _ = canvas.fill_rect(Rect::new(IR, IT, 3, (CY - IT) as u32));

    // West approach: eastbound lanes (y = CY..IB) stop at x = IL.
    let _ = canvas.fill_rect(Rect::new(IL - 3, CY, 3, (IB - CY) as u32));
}

// ── Drawing primitives ────────────────────────────────────────────────────────

/// Solid vertical line of 1 px width at column `x`, from row `y0` to `y1`.
fn hline(canvas: &mut Canvas<Window>, x: i32, y0: i32, y1: i32) {
    let _ = canvas.draw_line(Point::new(x, y0), Point::new(x, y1));
}

/// Solid horizontal line of 1 px height at row `y`, from column `x0` to `x1`.
fn vline(canvas: &mut Canvas<Window>, y: i32, x0: i32, x1: i32) {
    let _ = canvas.draw_line(Point::new(x0, y), Point::new(x1, y));
}

/// Dashed vertical line at column `x`, from row `y0` to `y1`.
fn dashed_hline(canvas: &mut Canvas<Window>, x: i32, y0: i32, y1: i32) {
    let mut y = y0;
    while y < y1 {
        let end = (y + DASH).min(y1);
        let _ = canvas.draw_line(Point::new(x, y), Point::new(x, end));
        y += DASH + GAP;
    }
}

/// Dashed horizontal line at row `y`, from column `x0` to `x1`.
fn dashed_vline(canvas: &mut Canvas<Window>, y: i32, x0: i32, x1: i32) {
    let mut x = x0;
    while x < x1 {
        let end = (x + DASH).min(x1);
        let _ = canvas.draw_line(Point::new(x, y), Point::new(end, y));
        x += DASH + GAP;
    }
}
