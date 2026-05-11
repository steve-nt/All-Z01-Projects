use crate::math::Vec3;
use serde::Deserialize;

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct PointLight {
    pub position: Vec3,
    pub brightness: f64,
}
