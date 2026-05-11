use std::ops::{Add, Mul};

use serde::Deserialize;

const MAX_COLOR_VALUE: u8 = 255;

#[derive(Clone, Copy, Debug, Deserialize, PartialEq)]
pub struct Color {
    pub r: f64,
    pub g: f64,
    pub b: f64,
}

impl Color {
    pub const fn new(r: f64, g: f64, b: f64) -> Self {
        Self { r, g, b }
    }

    pub fn lerp(start: Self, end: Self, t: f64) -> Self {
        start * (1.0 - t) + end * t
    }

    pub fn to_rgb(self) -> (u8, u8, u8) {
        (
            scaled_channel(self.r),
            scaled_channel(self.g),
            scaled_channel(self.b),
        )
    }
}

impl Add for Color {
    type Output = Self;

    fn add(self, other: Self) -> Self::Output {
        Self::new(self.r + other.r, self.g + other.g, self.b + other.b)
    }
}

impl Mul<f64> for Color {
    type Output = Self;

    fn mul(self, scalar: f64) -> Self::Output {
        Self::new(self.r * scalar, self.g * scalar, self.b * scalar)
    }
}

pub fn max_color_value() -> u8 {
    MAX_COLOR_VALUE
}

fn scaled_channel(value: f64) -> u8 {
    (MAX_COLOR_VALUE as f64 * value.clamp(0.0, 1.0)).round() as u8
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn color_channels_are_clamped_and_scaled() {
        assert_eq!(Color::new(-0.2, 0.5, 1.2).to_rgb(), (0, 128, 255));
    }
}
