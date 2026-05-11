use crate::{
    color::Color,
    geometry::{Hit, Hittable},
    math::Vec3,
    ray::Ray,
};
use serde::Deserialize;

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct Plane {
    pub point: Vec3,
    pub normal: Vec3,
    pub color: Color,
}

impl Hittable for Plane {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        let denominator = self.normal.dot(ray.direction);

        if denominator.abs() < 1e-6 {
            return None;
        }

        let t = (self.point - ray.origin).dot(self.normal) / denominator;

        if t < min_t || t > max_t {
            return None;
        }

        Some(Hit {
            t,
            point: ray.at(t),
            normal: self.normal,
            color: self.color,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn floor_plane() -> Plane {
        Plane {
            point: Vec3::new(0.0, 0.0, 0.0),
            normal: Vec3::new(0.0, 1.0, 0.0),
            color: Color::new(0.5, 0.5, 0.5),
        }
    }

    #[test]
    fn ray_hits_plane() {
        let plane = floor_plane();
        let ray = Ray::new(Vec3::new(0.0, 2.0, 0.0), Vec3::new(0.0, -1.0, 0.0));

        let hit = plane.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 2.0);
    }

    #[test]
    fn parallel_ray_misses_plane() {
        let plane = floor_plane();
        let ray = Ray::new(Vec3::new(0.0, 1.0, 0.0), Vec3::new(1.0, 0.0, 0.0));

        assert!(plane.hit(ray, 0.001, f64::INFINITY).is_none());
    }

    #[test]
    fn hit_returns_plane_normal() {
        let plane = floor_plane();
        let ray = Ray::new(Vec3::new(0.0, 2.0, 0.0), Vec3::new(0.0, -1.0, 0.0));

        let hit = plane.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.normal, Vec3::new(0.0, 1.0, 0.0));
    }
}
