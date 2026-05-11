pub mod camera;
pub mod color;
pub mod config;
pub mod geometry;
pub mod light;
pub mod math;
pub mod ray;
pub mod renderer;
pub mod scene;

pub use config::{IMAGE_HEIGHT, IMAGE_WIDTH};
pub use renderer::{render, render_scene};
