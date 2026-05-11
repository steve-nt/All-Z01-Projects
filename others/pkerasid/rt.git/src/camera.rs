use crate::{math::Vec3, ray::Ray};

#[derive(Clone, Copy, Debug)]
pub struct Camera {
    origin: Vec3,
    lower_left_corner: Vec3,
    horizontal: Vec3,
    vertical: Vec3,
}

impl Camera {
    pub fn new(
        origin: Vec3,
        look_at: Vec3,
        up: Vec3,
        vertical_fov_degrees: f64,
        aspect_ratio: f64,
    ) -> Self {
        let theta = vertical_fov_degrees.to_radians();
        let viewport_height = 2.0 * (theta / 2.0).tan();
        let viewport_width = aspect_ratio * viewport_height;

        let w = (origin - look_at).normalized();
        let u = up.cross(w).normalized();
        let v = w.cross(u);

        let horizontal = u * viewport_width;
        let vertical = v * viewport_height;
        let lower_left_corner = origin - horizontal / 2.0 - vertical / 2.0 - w;

        Self {
            origin,
            lower_left_corner,
            horizontal,
            vertical,
        }
    }

    pub fn ray_for(self, x: usize, y: usize, width: usize, height: usize) -> Ray {
        let u = x as f64 / (width - 1) as f64;
        let v = 1.0 - y as f64 / (height - 1) as f64;
        let direction =
            self.lower_left_corner + self.horizontal * u + self.vertical * v - self.origin;

        Ray::new(self.origin, direction.normalized())
    }
}
