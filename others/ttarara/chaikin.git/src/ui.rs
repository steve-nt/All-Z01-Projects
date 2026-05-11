use macroquad::prelude::*;

pub fn draw_small_point(point: Vec2, color: Color) {
    draw_circle_lines(point.x, point.y, 4.0, 1.5, color);
}

pub fn draw_points_as_circles(points: &[Vec2], radius: f32, thickness: f32, color: Color) {
    for &p in points {
        draw_circle_lines(p.x, p.y, radius, thickness, color);
    }
}

pub fn draw_polyline(points: &[Vec2], thickness: f32, color: Color) {
    if points.len() < 2 {
        return;
    }

    for i in 0..points.len() - 1 {
        let a = points[i];
        let b = points[i + 1];
        draw_line(a.x, a.y, b.x, b.y, thickness, color);
    }
}

pub fn draw_instructions(point_count: usize, current_step: usize, is_animating: bool) {
    let x = 20.0;
    let mut y = 30.0;
    let size = 28.0;

    draw_text("Left click: add control point", x, y, size, DARKGRAY);
    y += 30.0;
    draw_text("Enter: start / restart animation", x, y, size, DARKGRAY);
    y += 30.0;
    draw_text("C: Reset/Clear board", x, y, size, DARKGRAY);
    y += 30.0;
    draw_text("Esc: quit", x, y, size, DARKGRAY);
    y += 40.0;

    let status = format!("Control points: {}", point_count);
    draw_text(&status, x, y, size, BLACK);
    y += 30.0;

    let step_text = if is_animating {
        format!("Current step: {}", current_step)
    } else {
        "Current step: not animating".to_string()
    };
    draw_text(&step_text, x, y, size, BLACK);
}

pub fn draw_message(message: &str) {
    let font_size = 32.0;
    let dims = measure_text(message, None, font_size as u16, 1.0);

    let x = screen_width() / 2.0 - dims.width / 2.0;
    let y = screen_height() - 40.0;

    draw_text(message, x, y, font_size, RED);
}