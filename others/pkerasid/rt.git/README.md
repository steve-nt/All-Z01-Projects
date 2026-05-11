# RT — Rust Ray Tracer

A ray tracer written in Rust that renders 3D scenes to PPM images. It supports
spheres, cubes, flat planes, and cylinders with a movable camera, a point light
with adjustable brightness, diffuse shading, specular highlights, and shadows.

```sh
cargo run -- --scene scenes/all_objects.ron > output.ppm
```

## Quick Start

Build and render a scene:

```sh
cargo run --release -- --scene scenes/all_objects.ron > output.ppm
```

For a fast preview at low resolution:

```sh
cargo run -- --scene scenes/all_objects.ron --width 200 --height 150 > preview.ppm
```

Open the result in any image viewer that supports PPM (GIMP, feh, Eye of GNOME):

```sh
xdg-open output.ppm
```

## Usage

```
cargo run -- [--scene <path>] [--width <pixels>] [--height <pixels>]
```

| Flag | Description |
|---|---|
| `--scene <path>` | Load a scene file (RON format). Defaults to the built-in demo scene. |
| `--width <pixels>` | Override the image width from the scene file. |
| `--height <pixels>` | Override the image height from the scene file. |

Output is written to stdout, so use shell redirection to save it:

```sh
cargo run -- --scene scenes/sphere.ron > sphere.ppm
```

## Scene Files

Scenes are written in [RON](https://github.com/ron-rs/ron) format and live under
`scenes/`. The included scenes are:

| File | Description |
|---|---|
| `scenes/sphere.ron` | A single sphere with a ground plane. |
| `scenes/plane_cube_low_brightness.ron` | A cube and ground plane at low brightness. |
| `scenes/all_objects.ron` | One of each object type at full brightness. |
| `scenes/all_objects_alt_camera.ron` | Same as above from a different camera position. |
| `scenes/demo.ron` | Developer sandbox with multiple spheres. |

## PPM Output

Output uses the `P3` PPM format — plain ASCII, no image library required:

```text
P3
800 600
255
255 0 0
0 255 0
0 0 255
...
```

- `P3` — full-color ASCII PPM.
- `800 600` — image width and height in pixels.
- `255` — maximum value per RGB channel.
- Each following line is one pixel as `red green blue`.

## Prebuilt Images

The four mandatory renders live in `deliverables/` at 800×600 and can be
regenerated from a clean checkout with:

```sh
./generate_deliverables.sh
```

Individual commands:

```sh
# Sphere scene
cargo run --release -- --scene scenes/sphere.ron > deliverables/sphere.ppm

# Plane and cube with reduced brightness
cargo run --release -- --scene scenes/plane_cube_low_brightness.ron \
    > deliverables/plane_cube_low_brightness.ppm

# One of each object (sphere, cube, cylinder, plane)
cargo run --release -- --scene scenes/all_objects.ron \
    > deliverables/all_objects.ppm

# Same scene from a different camera position
cargo run --release -- --scene scenes/all_objects_alt_camera.ron \
    > deliverables/all_objects_alt_camera.ppm
```

## Development

Run the test suite and formatter:

```sh
cargo fmt
cargo test
```

## Documentation

### How to create objects

Objects are defined in a `.ron` scene file under the `objects` list. Each entry
is a tagged variant with its fields in parentheses.

**Sphere** — defined by a center point and a radius:

```ron
Sphere((
    center: (x: 0.0, y: 0.0, z: -2.0),
    radius: 0.5,
    color: (r: 0.8, g: 0.2, b: 0.2),
))
```

**Cube** — axis-aligned, defined by a center point and a side length:

```ron
Cube((
    center: (x: 1.0, y: 0.0, z: -2.0),
    size: 0.6,
    color: (r: 0.8, g: 0.5, b: 0.2),
))
```

**Flat plane** — defined by any point on the plane and an outward normal:

```ron
Plane((
    point: (x: 0.0, y: -0.75, z: 0.0),
    normal: (x: 0.0, y: 1.0, z: 0.0),
    color: (r: 0.5, g: 0.5, b: 0.5),
))
```

**Cylinder** — finite, defined by a center point, an axis direction, a radius,
and a height:

```ron
Cylinder((
    center: (x: -1.0, y: 0.0, z: -2.0),
    axis: (x: 0.0, y: 1.0, z: 0.0),
    radius: 0.3,
    height: 1.0,
    color: (r: 0.2, g: 0.4, b: 0.8),
))
```

### How to change brightness

Brightness is controlled by two fields at the top level of the scene file:

- `brightness` inside the `light` block — controls how strongly the point light
  illuminates surfaces. `0.95` is bright, `0.25` is dim.
- `ambient` — sets the minimum light level so shadowed areas are not fully
  black. `0.18` is a natural default, `0.06` gives a darker mood.

```ron
light: (
    position: (x: -2.7, y: 4.0, z: 1.3),
    brightness: 0.25,
),
ambient: 0.06,
```

### How to move the camera

The camera is configured in the `camera` block. Change `origin` to move the
camera and `look_at` to point it at a different part of the scene:

```ron
camera: (
    origin: (x: 3.5, y: 1.2, z: 1.0),
    look_at: (x: 0.0, y: -0.3, z: -2.0),
    up: (x: 0.0, y: 1.0, z: 0.0),
    vertical_fov_degrees: 55.0,
),
```

| Field | Description |
|---|---|
| `origin` | Where the camera sits in the world. |
| `look_at` | The point the camera is aimed at. |
| `up` | Which direction is "up"; `(0, 1, 0)` keeps the camera level. |
| `vertical_fov_degrees` | Field of view angle; `55.0` gives a natural perspective. |
