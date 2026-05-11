mod algorithm;
mod animations;
mod app;
mod ui;

use app::App;
use macroquad::prelude::*;

fn window_conf() -> Conf {
    Conf {
        window_title: "Chaikin Animation".to_string(),
        window_width: 1000,
        window_height: 700,
        ..Default::default()
    }
}

#[macroquad::main(window_conf)]
async fn main() {
    let mut app = App::new();

    loop {
        clear_background(WHITE);

        if app.handle_input() {
            break;
        }

        app.update();
        app.draw();

        next_frame().await;
    }
}
