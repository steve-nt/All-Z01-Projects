use crate::{
    color::Color,
    geometry::{Hit, Hittable},
    math::Vec3,
    ray::Ray,
};
use serde::Deserialize;

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct Sphere {
    pub center: Vec3,
    pub radius: f64,
    pub color: Color,
}

impl Hittable for Sphere {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        let oc = ray.origin - self.center;
        let a = ray.direction.length_squared();
        let half_b = oc.dot(ray.direction);
        let c = oc.length_squared() - self.radius * self.radius;
        let discriminant = half_b * half_b - a * c;

        if discriminant < 0.0 {
            return None;
        }

        let sqrt_discriminant = discriminant.sqrt();
        let mut root = (-half_b - sqrt_discriminant) / a;

        if root < min_t || root > max_t {
            root = (-half_b + sqrt_discriminant) / a;
            if root < min_t || root > max_t {
                return None;
            }
        }

        let point = ray.at(root);
        let normal = (point - self.center) / self.radius;

        Some(Hit {
            t: root,
            point,
            normal,
            color: self.color,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn unit_sphere() -> Sphere {
        Sphere {
            center: Vec3::new(0.0, 0.0, 0.0),
            radius: 1.0,
            color: Color::new(1.0, 0.0, 0.0),
        }
    }

    #[test]
    fn ray_hits_sphere_from_front() {
        let sphere = unit_sphere();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = sphere.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 2.0);
    }

    #[test]
    fn ray_misses_sphere() {
        let sphere = unit_sphere();
        let ray = Ray::new(Vec3::new(2.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));

        assert!(sphere.hit(ray, 0.001, f64::INFINITY).is_none());
    }

    #[test]
    fn normal_points_outward_and_is_unit_length() {
        let sphere = unit_sphere();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = sphere.hit(ray, 0.001, f64::INFINITY).unwrap();
        let len = hit.normal.length();

        assert!((len - 1.0).abs() < 1e-10);
        assert!(hit.normal.z > 0.0);
    }

    #[test]
    fn ray_from_inside_hits_far_surface() {
        let sphere = unit_sphere();
        let ray = Ray::new(Vec3::new(0.0, 0.0, 0.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = sphere.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 1.0);
    }
}
