#![allow(
    clippy::cast_possible_truncation,
    clippy::cast_possible_wrap,
    clippy::cast_sign_loss
)]

use sdl2::pixels::Color;
use sdl2::rect::Rect;
use sdl2::render::Canvas;
use sdl2::video::Window;

use crate::config::{WINDOW_H, WINDOW_W};
use crate::domain::stats::Statistics;

const BACKGROUND: Color = Color::RGB(24, 33, 39);
const PANEL: Color = Color::RGB(236, 233, 224);
const PANEL_BORDER: Color = Color::RGB(194, 122, 58);
const TITLE: Color = Color::RGB(34, 53, 61);
const BODY: Color = Color::RGB(52, 52, 52);
const VALUE: Color = Color::RGB(168, 66, 34);

const GLYPH_W: usize = 4;
const GLYPH_H: usize = 5;

pub fn render_statistics(canvas: &mut Canvas<Window>, stats: &Statistics) {
    canvas.set_draw_color(BACKGROUND);
    canvas.clear();

    let panel = Rect::new(120, 70, WINDOW_W - 240, WINDOW_H - 140);
    canvas.set_draw_color(PANEL);
    let _ = canvas.fill_rect(panel);
    canvas.set_draw_color(PANEL_BORDER);
    let _ = canvas.draw_rect(panel);

    draw_centered_text(canvas, 120, 6, TITLE, "FINAL STATISTICS");

    let rows = [
        format!("VEHICLES COMPLETED: {}", stats.vehicles_completed),
        format!("MAX SPEED: {:.1}", stats.max_speed),
        format!("MIN SPEED: {:.1}", stats.min_speed),
        format!("MAX TIME TO PASS: {:.1}", stats.max_time_to_pass),
        format!("MIN TIME TO PASS: {:.1}", stats.min_time_to_pass),
        format!("CLOSE CALLS: {}", stats.close_calls),
    ];

    for (idx, row) in rows.iter().enumerate() {
        let y = 210 + idx as i32 * 58;
        let color = if idx == 0 || idx == rows.len() - 1 {
            VALUE
        } else {
            BODY
        };
        draw_text(canvas, 180, y, 4, color, row);
    }

    draw_centered_text(canvas, 610, 3, TITLE, "ESC TO EXIT");
}

fn draw_centered_text(canvas: &mut Canvas<Window>, y: i32, scale: i32, color: Color, text: &str) {
    let width = text_pixel_width(text, scale);
    let x = (WINDOW_W as i32 - width) / 2;
    draw_text(canvas, x, y, scale, color, text);
}

fn draw_text(canvas: &mut Canvas<Window>, x: i32, y: i32, scale: i32, color: Color, text: &str) {
    canvas.set_draw_color(color);

    let mut cursor_x = x;
    for ch in text.chars() {
        draw_glyph(canvas, cursor_x, y, scale, ch);
        cursor_x += (GLYPH_W as i32 + 1) * scale;
    }
}

fn draw_glyph(canvas: &mut Canvas<Window>, x: i32, y: i32, scale: i32, ch: char) {
    let rows = glyph_rows(ch.to_ascii_uppercase());

    for (row_idx, row) in rows.iter().enumerate() {
        for (col_idx, pixel) in row.chars().enumerate() {
            if pixel != '#' {
                continue;
            }

            let rect = Rect::new(
                x + col_idx as i32 * scale,
                y + row_idx as i32 * scale,
                scale as u32,
                scale as u32,
            );
            let _ = canvas.fill_rect(rect);
        }
    }
}

fn text_pixel_width(text: &str, scale: i32) -> i32 {
    let char_count = text.chars().count() as i32;
    if char_count == 0 {
        0
    } else {
        char_count * (GLYPH_W as i32 + 1) * scale - scale
    }
}

fn glyph_rows(ch: char) -> [&'static str; GLYPH_H] {
    match ch {
        'A' => [" ## ", "#  #", "####", "#  #", "#  #"],
        'C' => [" ###", "#   ", "#   ", "#   ", " ###"],
        'D' => ["### ", "#  #", "#  #", "#  #", "### "],
        'E' => ["####", "#   ", "### ", "#   ", "####"],
        'F' => ["####", "#   ", "### ", "#   ", "#   "],
        'H' => ["#  #", "#  #", "####", "#  #", "#  #"],
        'I' => ["####", " ## ", " ## ", " ## ", "####"],
        'K' => ["#  #", "# # ", "##  ", "# # ", "#  #"],
        'L' => ["#   ", "#   ", "#   ", "#   ", "####"],
        'M' => ["#  #", "####", "# ##", "#  #", "#  #"],
        'N' => ["#  #", "## #", "# ##", "#  #", "#  #"],
        'O' | '0' => [" ## ", "#  #", "#  #", "#  #", " ## "],
        'P' => ["### ", "#  #", "### ", "#   ", "#   "],
        'S' => [" ###", "#   ", " ## ", "   #", "### "],
        'T' => ["####", " ## ", " ## ", " ## ", " ## "],
        'U' => ["#  #", "#  #", "#  #", "#  #", " ## "],
        'V' => ["#  #", "#  #", "#  #", " ## ", " ## "],
        'X' => ["#  #", " ## ", " ## ", " ## ", "#  #"],
        '1' => ["  # ", " ## ", "  # ", "  # ", " ###"],
        '2' => [" ## ", "#  #", "  # ", " #  ", "####"],
        '3' => ["### ", "   #", " ## ", "   #", "### "],
        '4' => ["#  #", "#  #", "####", "   #", "   #"],
        '5' => ["####", "#   ", "### ", "   #", "### "],
        '6' => [" ## ", "#   ", "### ", "#  #", " ## "],
        '7' => ["####", "   #", "  # ", " #  ", "#   "],
        '8' => [" ## ", "#  #", " ## ", "#  #", " ## "],
        '9' => [" ## ", "#  #", " ###", "   #", " ## "],
        ':' => ["    ", " ## ", "    ", " ## ", "    "],
        '.' => ["    ", "    ", "    ", " ## ", " ## "],
        ' ' => ["    ", "    ", "    ", "    ", "    "],
        _ => ["####", "#  #", "#  #", "#  #", "####"],
    }
}
