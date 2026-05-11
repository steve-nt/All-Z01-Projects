use std::io::{self, Write};

use crate::{
    camera::Camera,
    color::{max_color_value, Color},
    ray::Ray,
    scene::{demo_scene, Scene},
};

const SHADOW_BIAS: f64 = 0.001;

pub fn render<W: Write>(output: &mut W, width: usize, height: usize) -> io::Result<()> {
    let scene = demo_scene();

    render_scene(output, &scene, width, height)
}

pub fn render_scene<W: Write>(
    output: &mut W,
    scene: &Scene,
    width: usize,
    height: usize,
) -> io::Result<()> {
    write_ppm_header(output, width, height)?;

    let camera = Camera::new(
        scene.camera.origin,
        scene.camera.look_at,
        scene.camera.up,
        scene.camera.vertical_fov_degrees,
        width as f64 / height as f64,
    );

    for y in 0..height {
        for x in 0..width {
            let ray = camera.ray_for(x, y, width, height);
            let color = ray_color(ray, scene);
            write_color(output, color)?;
        }
    }

    Ok(())
}

fn write_ppm_header<W: Write>(output: &mut W, width: usize, height: usize) -> io::Result<()> {
    writeln!(output, "P3")?;
    writeln!(output, "{} {}", width, height)?;
    writeln!(output, "{}", max_color_value())
}

fn write_color<W: Write>(output: &mut W, color: Color) -> io::Result<()> {
    let (r, g, b) = color.to_rgb();
    writeln!(output, "{r} {g} {b}")
}

fn ray_color(ray: Ray, scene: &Scene) -> Color {
    if let Some(hit) = scene.hit(ray, SHADOW_BIAS, f64::INFINITY) {
        shade(ray, hit, scene)
    } else {
        sky_color(ray)
    }
}

fn shade(ray: Ray, hit: crate::geometry::Hit, scene: &Scene) -> Color {
    let to_light = scene.light.position - hit.point;
    let light_distance = to_light.length();
    let light_direction = to_light / light_distance;
    let shadow_ray = Ray::new(hit.point + hit.normal * SHADOW_BIAS, light_direction);
    let shadowed = scene.hit(shadow_ray, SHADOW_BIAS, light_distance).is_some();
    let shadow_strength = if shadowed { 0.35 } else { 1.0 };

    let diffuse = hit.normal.dot(light_direction).max(0.0);
    let view_direction = -ray.direction;
    let halfway = (light_direction + view_direction).normalized();
    let specular = hit.normal.dot(halfway).max(0.0).powf(80.0) * 0.25;
    let intensity = scene.ambient + scene.light.brightness * diffuse * shadow_strength;

    hit.color * intensity + Color::new(1.0, 1.0, 1.0) * specular * shadow_strength
}

fn sky_color(ray: Ray) -> Color {
    let t = 0.5 * (ray.direction.y + 1.0);

    Color::lerp(Color::new(1.0, 1.0, 1.0), Color::new(0.5, 0.7, 1.0), t)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn smoke_render_scene_file_produces_valid_ppm() {
        let scene = crate::scene::Scene::from_ron_file("scenes/sphere.ron").unwrap();
        let mut output = Vec::new();

        render_scene(&mut output, &scene, 4, 3).unwrap();

        let ppm = String::from_utf8(output).unwrap();
        let lines: Vec<&str> = ppm.lines().collect();

        assert_eq!(lines[0], "P3");
        assert_eq!(lines[1], "4 3");
        assert_eq!(lines[2], "255");
        assert_eq!(lines.len(), 15);
    }

    #[test]
    fn ppm_writer_outputs_expected_header_and_pixels() {
        let mut output = Vec::new();

        render(&mut output, 2, 2).unwrap();

        let ppm = String::from_utf8(output).unwrap();
        let lines: Vec<&str> = ppm.lines().collect();

        assert_eq!(lines[0], "P3");
        assert_eq!(lines[1], "2 2");
        assert_eq!(lines[2], "255");
        assert_eq!(lines.len(), 7);
    }
}
