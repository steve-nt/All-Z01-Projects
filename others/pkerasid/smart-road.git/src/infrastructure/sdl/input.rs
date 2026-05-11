use sdl2::event::Event;
use sdl2::keyboard::Keycode;

use crate::domain::lane::Direction;

/// SDL input mapped into app-level intentions.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum InputAction {
    Quit,
    ShowStatistics,
    Spawn(Direction),
    ToggleRandomSpawning,
}

#[must_use]
pub fn translate_event(event: &Event) -> Option<InputAction> {
    match event {
        Event::Quit { .. } => Some(InputAction::Quit),
        Event::KeyDown {
            keycode: Some(Keycode::Escape),
            ..
        } => Some(InputAction::ShowStatistics),
        Event::KeyDown {
            keycode: Some(Keycode::Up),
            repeat: false,
            ..
        } => Some(InputAction::Spawn(Direction::South)),
        Event::KeyDown {
            keycode: Some(Keycode::Down),
            repeat: false,
            ..
        } => Some(InputAction::Spawn(Direction::North)),
        Event::KeyDown {
            keycode: Some(Keycode::Left),
            repeat: false,
            ..
        } => Some(InputAction::Spawn(Direction::East)),
        Event::KeyDown {
            keycode: Some(Keycode::Right),
            repeat: false,
            ..
        } => Some(InputAction::Spawn(Direction::West)),
        Event::KeyDown {
            keycode: Some(Keycode::R),
            repeat: false,
            ..
        } => Some(InputAction::ToggleRandomSpawning),
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn up_arrow_maps_to_south_origin() {
        let event = Event::KeyDown {
            timestamp: 0,
            window_id: 0,
            keycode: Some(Keycode::Up),
            scancode: None,
            keymod: sdl2::keyboard::Mod::NOMOD,
            repeat: false,
        };

        assert_eq!(
            translate_event(&event),
            Some(InputAction::Spawn(Direction::South))
        );
    }

    #[test]
    fn escape_requests_statistics_screen() {
        let event = Event::KeyDown {
            timestamp: 0,
            window_id: 0,
            keycode: Some(Keycode::Escape),
            scancode: None,
            keymod: sdl2::keyboard::Mod::NOMOD,
            repeat: false,
        };

        assert_eq!(translate_event(&event), Some(InputAction::ShowStatistics));
    }
}
