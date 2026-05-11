use crate::{
    color::Color,
    geometry::{Hit, Hittable},
    math::Vec3,
    ray::Ray,
};
use serde::Deserialize;

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct Cube {
    pub center: Vec3,
    pub size: f64,
    pub color: Color,
}

impl Hittable for Cube {
    fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        let half_size = self.size / 2.0;
        let min = self.center - Vec3::new(half_size, half_size, half_size);
        let max = self.center + Vec3::new(half_size, half_size, half_size);

        let mut near_t = f64::NEG_INFINITY;
        let mut far_t = max_t;
        let mut near_normal = Vec3::new(0.0, 0.0, 0.0);
        let mut far_normal = Vec3::new(0.0, 0.0, 0.0);

        update_slab(
            ray.origin.x,
            ray.direction.x,
            min.x,
            max.x,
            Vec3::new(-1.0, 0.0, 0.0),
            &mut near_t,
            &mut far_t,
            &mut near_normal,
            &mut far_normal,
        )?;
        update_slab(
            ray.origin.y,
            ray.direction.y,
            min.y,
            max.y,
            Vec3::new(0.0, -1.0, 0.0),
            &mut near_t,
            &mut far_t,
            &mut near_normal,
            &mut far_normal,
        )?;
        update_slab(
            ray.origin.z,
            ray.direction.z,
            min.z,
            max.z,
            Vec3::new(0.0, 0.0, -1.0),
            &mut near_t,
            &mut far_t,
            &mut near_normal,
            &mut far_normal,
        )?;

        if near_t >= min_t && near_t <= max_t {
            return Some(Hit {
                t: near_t,
                point: ray.at(near_t),
                normal: near_normal,
                color: self.color,
            });
        }

        if far_t >= min_t && far_t <= max_t {
            return Some(Hit {
                t: far_t,
                point: ray.at(far_t),
                normal: far_normal,
                color: self.color,
            });
        }

        None
    }
}

fn update_slab(
    origin: f64,
    direction: f64,
    min: f64,
    max: f64,
    min_normal: Vec3,
    near_t: &mut f64,
    far_t: &mut f64,
    near_normal: &mut Vec3,
    far_normal: &mut Vec3,
) -> Option<()> {
    const EPSILON: f64 = 1e-8;

    if direction.abs() < EPSILON {
        return (origin >= min && origin <= max).then_some(());
    }

    let mut t0 = (min - origin) / direction;
    let mut t1 = (max - origin) / direction;
    let mut t0_normal = min_normal;
    let mut t1_normal = -min_normal;

    if t0 > t1 {
        std::mem::swap(&mut t0, &mut t1);
        std::mem::swap(&mut t0_normal, &mut t1_normal);
    }

    if t0 > *near_t {
        *near_t = t0;
        *near_normal = t0_normal;
    }

    if t1 < *far_t {
        *far_t = t1;
        *far_normal = t1_normal;
    }

    (*near_t <= *far_t).then_some(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn ray_hits_front_face() {
        let cube = Cube {
            center: Vec3::new(0.0, 0.0, 0.0),
            size: 2.0,
            color: Color::new(1.0, 0.0, 0.0),
        };
        let ray = Ray::new(Vec3::new(0.0, 0.0, 3.0), Vec3::new(0.0, 0.0, -1.0));

        let hit = cube.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 2.0);
        assert_eq!(hit.normal, Vec3::new(0.0, 0.0, 1.0));
    }

    #[test]
    fn ray_from_inside_hits_exit_face() {
        let cube = Cube {
            center: Vec3::new(0.0, 0.0, 0.0),
            size: 2.0,
            color: Color::new(1.0, 0.0, 0.0),
        };
        let ray = Ray::new(Vec3::new(0.0, 0.0, 0.0), Vec3::new(1.0, 0.0, 0.0));

        let hit = cube.hit(ray, 0.001, f64::INFINITY).unwrap();

        assert_eq!(hit.t, 1.0);
        assert_eq!(hit.normal, Vec3::new(1.0, 0.0, 0.0));
    }
}
