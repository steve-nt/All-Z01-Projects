pub const IMAGE_WIDTH: usize = 800;
pub const IMAGE_HEIGHT: usize = 600;

use std::{ffi::OsString, path::PathBuf};

#[derive(Debug, Default, PartialEq)]
pub struct RenderConfig {
    pub scene_path: Option<PathBuf>,
    pub width: Option<usize>,
    pub height: Option<usize>,
}

impl RenderConfig {
    pub fn from_args<I>(args: I) -> Result<Self, String>
    where
        I: IntoIterator<Item = OsString>,
    {
        let mut args = args.into_iter();
        let mut config = Self::default();

        while let Some(arg) = args.next() {
            let arg = arg
                .into_string()
                .map_err(|_| "arguments must be valid UTF-8".to_string())?;

            match arg.as_str() {
                "--scene" => {
                    let value = next_arg(&mut args, "--scene")?;
                    config.scene_path = Some(PathBuf::from(value));
                }
                "--width" => {
                    let value = next_arg(&mut args, "--width")?;
                    config.width = Some(parse_dimension("--width", &value)?);
                }
                "--height" => {
                    let value = next_arg(&mut args, "--height")?;
                    config.height = Some(parse_dimension("--height", &value)?);
                }
                "--help" | "-h" => return Err(help().to_string()),
                _ => return Err(format!("unknown argument `{arg}`\n\n{}", help())),
            }
        }

        Ok(config)
    }
}

pub fn help() -> &'static str {
    "usage: cargo run -- [--scene <path>] [--width <pixels>] [--height <pixels>]\n\n\
     examples:\n  \
     cargo run > output.ppm\n  \
     cargo run -- --scene scenes/demo.ron --width 400 --height 300 > output.ppm"
}

fn next_arg<I>(args: &mut I, flag: &str) -> Result<String, String>
where
    I: Iterator<Item = OsString>,
{
    args.next()
        .ok_or_else(|| format!("{flag} requires a value"))?
        .into_string()
        .map_err(|_| format!("{flag} value must be valid UTF-8"))
}

fn parse_dimension(flag: &str, value: &str) -> Result<usize, String> {
    let dimension = value
        .parse::<usize>()
        .map_err(|_| format!("{flag} must be a positive integer"))?;

    if dimension < 2 {
        return Err(format!("{flag} must be at least 2 pixels"));
    }

    Ok(dimension)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn args(values: &[&str]) -> Vec<OsString> {
        values.iter().map(OsString::from).collect()
    }

    #[test]
    fn parses_scene_and_resolution_overrides() {
        let config = RenderConfig::from_args(args(&[
            "--scene",
            "scenes/demo.ron",
            "--width",
            "40",
            "--height",
            "30",
        ]))
        .unwrap();

        assert_eq!(config.scene_path, Some(PathBuf::from("scenes/demo.ron")));
        assert_eq!(config.width, Some(40));
        assert_eq!(config.height, Some(30));
    }

    #[test]
    fn rejects_tiny_dimensions() {
        let err = RenderConfig::from_args(args(&["--width", "1"])).unwrap_err();

        assert_eq!(err, "--width must be at least 2 pixels");
    }
}
