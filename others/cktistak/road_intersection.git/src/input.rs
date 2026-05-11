//! Controls:
//! - ↑ Up: spawn from south (moving north)
//! - ↓ Down: spawn from north (moving south)
//! - → Right: spawn from west (moving east)
//! - ← Left: spawn from east (moving west)
//! - r: spawn from random direction
//! - Esc: end simulation

use crate::vehicle::Direction;
use macroquad::input::KeyCode;
use macroquad::input::is_key_pressed;
use std::time::SystemTime;

pub fn get_spawn_direction() -> Option<Direction> {
    if is_key_pressed(KeyCode::Up) {
        Some(Direction::North)
    } else if is_key_pressed(KeyCode::Down) {
        Some(Direction::South)
    } else if is_key_pressed(KeyCode::Right) {
        Some(Direction::West)
    } else if is_key_pressed(KeyCode::Left) {
        Some(Direction::East)
    } else if is_key_pressed(KeyCode::R) {
        let seed = SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap()
            .subsec_nanos();
        match seed % 4 {
            0 => Some(Direction::North),
            1 => Some(Direction::South),
            2 => Some(Direction::East),
            _ => Some(Direction::West),
        }
    } else {
        None
    }
}

/// Check if user wants to exit
pub fn should_exit() -> bool {
    is_key_pressed(KeyCode::Escape)
}
