use crate::{
    color::Color,
    geometry::{Hit, Hittable},
    math::Vec3,
    ray::Ray,
};
use serde::Deserialize;

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct Cylinder {
    pub center: Vec3,
    pub axis: Vec3,
    pub radius: f64,
    pub height: f64,
    pub color: Color,
}

impl Hittable for Cylinder {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        let axis = self.axis.normalized();
        let half_height = self.height / 2.0;
        let ray_to_center = ray.origin - self.center;
        let direction_parallel = axis * ray.direction.dot(axis);
        let direction_perpendicular = ray.direction - direction_parallel;
        let origin_parallel = axis * ray_to_center.dot(axis);
        let origin_perpendicular = ray_to_center - origin_parallel;

        let a = direction_perpendicular.length_squared();
        let half_b = origin_perpendicular.dot(direction_perpendicular);
        let c = origin_perpendicular.length_squared() - self.radius * self.radius;

        let mut nearest_hit = None;
        let mut closest = max_t;

        if a > 1e-8 {
            let discriminant = half_b * half_b - a * c;

            if discriminant >= 0.0 {
                let sqrt_discriminant = discriminant.sqrt();
                let roots = [
                    (-half_b - sqrt_discriminant) / a,
                    (-half_b + sqrt_discriminant) / a,
                ];

                for root in roots {
                    if root < min_t || root > closest {
                        continue;
                    }

                    let point = ray.at(root);
                    let height_offset = (point - self.center).dot(axis);

                    if height_offset.abs() <= half_height {
                        let center_on_axis = self.center + axis * height_offset;
                        closest = root;
                        nearest_hit = Some(Hit {
                            t: root,
                            point,
                            normal: (point - center_on_axis).normalized(),
                            color: self.color,
                        });
                    }
                }
            }
        }

        for cap_sign in [-1.0, 1.0] {
            let cap_center = self.center + axis * (cap_sign * half_height);
            let cap_normal = axis * cap_sign;
            let denominator = ray.direction.dot(cap_normal);

            if denominator.abs() < 1e-8 {
                continue;
            }

            let t = (cap_center - ray.origin).dot(cap_normal) / denominator;

            if t < min_t || t > closest {
                continue;
            }

            let point = ray.at(t);

            if (point - cap_center).length_squared() <= self.radius * self.radius {
                closest = t;
                nearest_hit = Some(Hit {
                    t,
                    point,
                    normal: cap_normal,
                    color: self.color,
                });
            }
        }

        nearest_hit
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn ray_hits_cylinder_side() {
        let cylinder = Cylinder {
            center: Vec3::new(0.0, 0.0, 0.0),
            axis: Vec3::new(0.0, 1.0, 0.0),
            radius: 1.0,
            height: 2.0,
            color: Color::new(0.0, 1.0, 0.0),
        };
        let ray = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = cylinder.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 2.0);
        assert_eq!(hit.normal, Vec3::new(0.0, 0.0, 1.0));
    }

    #[test]
    fn ray_hits_cylinder_cap() {
        let cylinder = Cylinder {
            center: Vec3::new(0.0, 0.0, 0.0),
            axis: Vec3::new(0.0, 1.0, 0.0),
            radius: 1.0,
            height: 2.0,
            color: Color::new(0.0, 1.0, 0.0),
        };
        let ray = Ray::new(Vec3::new(0.0, 3.0, 0.0), Vec3::new(0.0, -1.0, 0.0));

        let hit = cylinder.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 2.0);
        assert_eq!(hit.normal, Vec3::new(0.0, 1.0, 0.0));
    }
}
