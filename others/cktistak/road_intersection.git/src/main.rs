mod input;
mod render;
mod traffic_control;
mod vehicle;

use crate::render::*;

use macroquad::color::*;
use macroquad::prelude::next_frame;
use macroquad::prelude::*;
use crate::input::*;
use crate::vehicle::*;
use macroquad::prelude::clear_background;
use macroquad::text::draw_text;

fn window_conf() -> Conf {
    let (win_w, win_h) = window_size();
    Conf {
        window_title: "Road intersection".to_owned(),
        fullscreen: false,
        window_width: win_w as i32,
        window_height: win_h as i32,
        window_resizable: false,
        ..Default::default()
    }
}

pub fn draw_legend() {
    let color = GRAY;
    let font_size = 20.;
    let x = 10.;
    let mut y = 20.;
    let strings = vec![
        "Controls:",
        "Up: spawn from south (moving north)",
        "Down: spawn from north (moving south)",
        "Right: spawn from west (moving east)",
        "Left: spawn from east (moving west)",
        "r: spawn from random direction",
        "Esc: end simulation",
        "",
        "Vehicle Colors:",
        "Green = Straight",
        "Red = Left turn",
        "Blue = Right turn"
    ];
    for text in strings {
        draw_text(text, x, y, font_size, color);
        y += 20.;
    }
}

#[macroquad::main(window_conf)]
async fn main() {

    let mut vehicles: Vec<Vehicle> = vec![];
    let mut sim_time_secs: f64 = 0.0;

    loop {
        sim_time_secs += get_frame_time() as f64;
        clear_background(Color::from_rgba(34, 40, 34, 255));

        // Handle input
        if should_exit() {
            break;
        }

        if let Some(d) = get_spawn_direction() {
            let v = Vehicle::new(d);
            if v.can_spawn(&vehicles) {
                vehicles.push(v);
            }
        }

        // Static
        draw_legend();
        draw_roads();
        draw_lane_dividers();
        draw_stop_lines();

        // Simulate
        let tl = traffic_control::traffic_light_state(sim_time_secs, &vehicles);
        draw_traffic_lights(tl);
        // Sequential poses within the frame so turns and collision checks see prior updates.
        let mut frame_snap = vehicles.clone();
        for (i, vehicle) in vehicles.iter_mut().enumerate() {
            let hold = traffic_control::should_hold_for_light(vehicle, tl);
            if !hold {
                vehicle.move_vehicle(frame_snap.as_slice(), i);
            }
            frame_snap[i] = *vehicle;
            vehicle.draw();
        }
        // vehicle cleanup
        // screen is 0-800 with 60px buffer
        vehicles.retain(|v| v.x > -50.0 && v.x < 850.0 && v.y > -50.0 && v.y < 850.0);
        next_frame().await;
    }
}
