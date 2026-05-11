use smart_road::infrastructure::sdl::app::App;

fn main() {
    if let Err(e) = App::build().and_then(App::run) {
        eprintln!("fatal: {e}");
        std::process::exit(1);
    }
}
