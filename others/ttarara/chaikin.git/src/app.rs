use crate::animations::ChaikinAnimator;
use crate::ui;
use macroquad::prelude::*;

pub struct App {
    control_points: Vec<Vec2>,
    animator: ChaikinAnimator,
    show_no_points_message: bool,
    no_points_message_time: f64,
}

impl App {
    pub fn new() -> Self {
        Self {
            control_points: Vec::new(),
            animator: ChaikinAnimator::new(),
            show_no_points_message: false,
            no_points_message_time: 0.0,
        }
    }

    // returns true if the app should quit
    pub fn handle_input(&mut self) -> bool {
        if is_key_pressed(KeyCode::Escape) {
            return true;
        }
//reset
        if is_key_pressed(KeyCode::C) {
            self.control_points.clear();
            self.animator.clear();
        }

        if is_mouse_button_pressed(MouseButton::Left) {
            if self.animator.is_animating() {
                self.animator.stop_and_reset_step();
            } else {
                let (mx, my) = mouse_position();
                self.control_points.push(vec2(mx, my));
            }
        }

        if is_key_pressed(KeyCode::Enter) {
            match self.control_points.len() {
                0 => {
                    self.show_no_points_message = true;
                    self.no_points_message_time = get_time();
                }
                1 | 2 => {
                    self.animator.clear();
                }
                _ => {
                    self.animator
                        .start(&self.control_points, get_time());
                }
            }
        }

        false
    }

    pub fn update(&mut self) {
        if self.show_no_points_message && get_time() - self.no_points_message_time > 1.5 {
            self.show_no_points_message = false;
        }

        self.animator.update(get_time());
    }

    pub fn draw(&self) {
        let shown_points = self.animator.curve_points(&self.control_points);

        // draw current shape
        match shown_points.len() {
            0 => {}
            1 => {
                ui::draw_small_point(shown_points[0], BLUE);
            }
            2 => {
                ui::draw_polyline(shown_points, 2.0, BLUE);
                ui::draw_points_as_circles(shown_points, 4.0, 1.5, BLUE);
            }
            _ => {
                ui::draw_polyline(shown_points, 2.0, BLUE);
            }
        }

        // show original control points 
        
            ui::draw_points_as_circles(&self.control_points, 4.0, 1.5, BLACK);
        

      

        ui::draw_instructions(
            self.control_points.len(),
            self.animator.current_step(),
            self.animator.is_animating(),
        );

        if self.show_no_points_message {
            ui::draw_message("Draw at least one point first");
        }
    }

}