use std::thread;
use std::time::Duration;
use std::time::Instant;

use crate::application::commands::Command;
use crate::application::world::World;
use crate::config;
use crate::domain::stats::Statistics;
use crate::infrastructure::random::vehicle_generator::VehicleGenerator;
use crate::infrastructure::sdl::input::{self, InputAction};
use crate::presentation::scene;

const RANDOM_SPAWN_INTERVAL_SECS: f32 = 1.0;

/// Owns the SDL2 context, canvas, and top-level simulation state.
/// Runs the fixed-timestep game loop.
pub struct App {
    canvas: sdl2::render::Canvas<sdl2::video::Window>,
    events: sdl2::EventPump,
    world: World,
    vehicle_generator: VehicleGenerator,
    random_spawning: bool,
    random_spawn_accumulator: f32,
}

impl App {
    /// Initialise SDL2, create a window and canvas.
    ///
    /// # Errors
    ///
    /// Returns an error string if SDL2 initialisation or window creation fails.
    pub fn build() -> Result<Self, String> {
        let sdl = sdl2::init()?;
        let video = sdl.video()?;

        let window = video
            .window("Kifisos App", config::WINDOW_W, config::WINDOW_H)
            .position_centered()
            .build()
            .map_err(|e| e.to_string())?;

        let canvas = window
            .into_canvas()
            .present_vsync()
            .build()
            .map_err(|e| e.to_string())?;

        let events = sdl.event_pump()?;

        Ok(Self {
            canvas,
            events,
            world: World::new(),
            vehicle_generator: VehicleGenerator::new(),
            random_spawning: false,
            random_spawn_accumulator: 0.0,
        })
    }

    /// Run the event/update/render loop until the user closes the window or
    /// presses `Esc`.
    ///
    /// # Errors
    ///
    /// Propagates any SDL2 error that surfaces during the loop.
    pub fn run(mut self) -> Result<(), String> {
        let mut last_tick = Instant::now();
        let mut accumulator = 0.0_f32;
        let show_statistics = 'main: loop {
            // ── Events ────────────────────────────────────────────────────
            for event in self.events.poll_iter() {
                match input::translate_event(&event) {
                    Some(InputAction::Quit) => break 'main false,
                    Some(InputAction::ShowStatistics) => break 'main true,
                    Some(InputAction::Spawn(origin)) => {
                        let request = self.vehicle_generator.spawn_request_for_origin(origin);
                        self.world.apply_command(Command::SpawnVehicle(request));
                    }
                    Some(InputAction::ToggleRandomSpawning) => {
                        self.random_spawning = !self.random_spawning;
                    }
                    None => {}
                }
            }

            // ── Fixed timestep ────────────────────────────────────────────
            let now = Instant::now();
            let elapsed = now.duration_since(last_tick);
            last_tick = now;
            // Clamp to avoid spiral-of-death after a pause or debug break.
            let delta = elapsed.as_secs_f32().min(0.25);
            accumulator += delta;

            while accumulator >= config::FIXED_DT {
                if self.random_spawning {
                    self.random_spawn_accumulator += config::FIXED_DT;
                    while self.random_spawn_accumulator >= RANDOM_SPAWN_INTERVAL_SECS {
                        let request = self.vehicle_generator.random_spawn_request();
                        self.world.apply_command(Command::SpawnVehicle(request));
                        self.random_spawn_accumulator -= RANDOM_SPAWN_INTERVAL_SECS;
                    }
                }
                self.world.tick(config::FIXED_DT);
                accumulator -= config::FIXED_DT;
            }

            // ── Render ────────────────────────────────────────────────────
            scene::render(&mut self.canvas, &self.world)?;
            self.canvas.present();
        };

        if show_statistics {
            let stats = self.world.statistics().clone();
            self.show_statistics(&stats)?;
        }

        Ok(())
    }

    fn show_statistics(&mut self, stats: &Statistics) -> Result<(), String> {
        self.canvas
            .window_mut()
            .set_title("Kifisos App - Statistics")
            .map_err(|e| e.to_string())?;

        loop {
            scene::render_statistics(&mut self.canvas, stats);
            self.canvas.present();

            for event in self.events.poll_iter() {
                match input::translate_event(&event) {
                    Some(InputAction::Quit | InputAction::ShowStatistics) => return Ok(()),
                    Some(InputAction::Spawn(_) | InputAction::ToggleRandomSpawning) | None => {}
                }
            }

            thread::sleep(Duration::from_millis(16));
        }
    }
}
