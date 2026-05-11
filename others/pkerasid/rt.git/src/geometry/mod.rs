mod cube;
mod cylinder;
mod plane;
mod sphere;

pub use cube::Cube;
pub use cylinder::Cylinder;
pub use plane::Plane;
pub use sphere::Sphere;

use crate::{color::Color, math::Vec3, ray::Ray};
use serde::Deserialize;

#[derive(Clone, Copy, Debug)]
pub struct Hit {
    pub t: f64,
    pub point: Vec3,
    pub normal: Vec3,
    pub color: Color,
}

pub trait Hittable {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit>;
}

#[derive(Clone, Debug, Deserialize)]
pub enum Object {
    Cube(Cube),
    Cylinder(Cylinder),
    Sphere(Sphere),
    Plane(Plane),
}

impl Hittable for Object {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        match self {
            Self::Cube(cube) => cube.hit(ray, min_t, max_t),
            Self::Cylinder(cylinder) => cylinder.hit(ray, min_t, max_t),
            Self::Sphere(sphere) => sphere.hit(ray, min_t, max_t),
            Self::Plane(plane) => plane.hit(ray, min_t, max_t),
        }
    }
}
