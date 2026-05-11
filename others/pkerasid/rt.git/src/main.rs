use std::{
    env,
    io::{self, BufWriter},
    process,
};

use rt::{
    config::RenderConfig,
    render_scene,
    scene::{demo_scene, Scene},
};

fn main() {
    if let Err(err) = run() {
        eprintln!("error: {err}");
        process::exit(1);
    }
}

fn run() -> Result<(), String> {
    let args: Vec<_> = env::args_os().skip(1).collect();

    if args.iter().any(|arg| arg == "--help" || arg == "-h") {
        println!("{}", rt::config::help());
        return Ok(());
    }

    let config = RenderConfig::from_args(args)?;
    let scene = match &config.scene_path {
        Some(path) => Scene::from_ron_file(path)?,
        None => demo_scene(),
    };
    let width = config.width.unwrap_or(scene.image.width);
    let height = config.height.unwrap_or(scene.image.height);

    let stdout = io::stdout();
    let mut output = BufWriter::new(stdout.lock());

    render_scene(&mut output, &scene, width, height).map_err(|err| err.to_string())
}
