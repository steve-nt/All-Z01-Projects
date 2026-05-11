use crate::{
    color::Color,
    geometry::{Cube, Cylinder, Hit, Hittable, Object, Plane, Sphere},
    light::PointLight,
    math::Vec3,
    ray::Ray,
};
use serde::Deserialize;
use std::{fs, path::Path};

#[derive(Clone, Debug, Deserialize)]
pub struct Image {
    pub width: usize,
    pub height: usize,
}

#[derive(Clone, Copy, Debug, Deserialize)]
pub struct CameraSettings {
    pub origin: Vec3,
    pub look_at: Vec3,
    pub up: Vec3,
    pub vertical_fov_degrees: f64,
}

#[derive(Clone, Debug, Deserialize)]
pub struct Scene {
    pub image: Image,
    pub camera: CameraSettings,
    pub light: PointLight,
    pub ambient: f64,
    pub objects: Vec<Object>,
}

impl Scene {
    pub fn from_ron_file(path: impl AsRef<Path>) -> Result<Self, String> {
        let path = path.as_ref();
        let contents = fs::read_to_string(path)
            .map_err(|err| format!("failed to read scene `{}`: {err}", path.display()))?;

        ron::from_str(&contents)
            .map_err(|err| format!("failed to parse scene `{}`: {err}", path.display()))
    }

    pub fn hit(&self, ray: Ray, min_t: f64, max_t: f64) -> Option<Hit> {
        let mut closest = max_t;
        let mut nearest_hit = None;

        for object in &self.objects {
            if let Some(hit) = object.hit(ray, min_t, closest) {
                closest = hit.t;
                nearest_hit = Some(hit);
            }
        }

        nearest_hit
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn parses_sphere_ron_scene() {
        let scene = Scene::from_ron_file("scenes/sphere.ron").unwrap();

        assert_eq!(scene.image.width, 800);
        assert_eq!(scene.image.height, 600);
        assert_eq!(scene.objects.len(), 2);
        assert_eq!(scene.camera.origin, Vec3::new(0.0, 0.55, 3.4));
    }

    #[test]
    fn missing_file_gives_useful_error() {
        let err = Scene::from_ron_file("scenes/nonexistent.ron").unwrap_err();

        assert!(err.contains("failed to read scene"));
    }
}

pub fn demo_scene() -> Scene {
    Scene {
        image: Image {
            width: crate::config::IMAGE_WIDTH,
            height: crate::config::IMAGE_HEIGHT,
        },
        camera: CameraSettings {
            origin: Vec3::new(0.0, 0.45, 3.0),
            look_at: Vec3::new(0.1, -0.1, -1.6),
            up: Vec3::new(0.0, 1.0, 0.0),
            vertical_fov_degrees: 55.0,
        },
        light: PointLight {
            position: Vec3::new(-2.7, 4.0, 1.3),
            brightness: 0.95,
        },
        ambient: 0.18,
        objects: vec![
            Object::Sphere(Sphere {
                center: Vec3::new(0.75, -0.15, -1.75),
                radius: 0.72,
                color: Color::new(0.45, 0.74, 0.22),
            }),
            Object::Sphere(Sphere {
                center: Vec3::new(-0.75, -0.25, -1.95),
                radius: 0.48,
                color: Color::new(0.72, 0.74, 0.72),
            }),
            Object::Sphere(Sphere {
                center: Vec3::new(-0.38, -0.22, -1.48),
                radius: 0.16,
                color: Color::new(0.38, 0.62, 0.18),
            }),
            Object::Cube(Cube {
                center: Vec3::new(1.55, -0.32, -2.25),
                size: 0.62,
                color: Color::new(0.76, 0.34, 0.22),
            }),
            Object::Cylinder(Cylinder {
                center: Vec3::new(-1.55, -0.26, -2.35),
                axis: Vec3::new(0.0, 1.0, 0.0),
                radius: 0.26,
                height: 0.96,
                color: Color::new(0.22, 0.42, 0.76),
            }),
            Object::Plane(Plane {
                point: Vec3::new(0.0, -0.75, 0.0),
                normal: Vec3::new(0.0, 1.0, 0.0),
                color: Color::new(0.47, 0.56, 0.66),
            }),
        ],
    }
}
