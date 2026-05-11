# RT Acceptance Checklist

This checklist freezes the first passing target for the 01-edu `rt` subject.
It is based on:

- Subject: <https://github.com/01-edu/public/tree/master/subjects/rt>
- Audit: <https://github.com/01-edu/public/blob/master/subjects/rt/audit/README.md>

## Mandatory Feature Scope

The first passing version must support:

- [x] `P3` ASCII PPM output.
- [x] Final audit renders at `800x600`.
- [x] Smaller preview renders for development and audit retesting.
- [x] Spheres.
- [x] Cubes.
- [x] Flat planes.
- [x] Cylinders.
- [x] Object locations configurable before rendering.
- [x] Camera position configurable before rendering.
- [x] Camera viewing angle or target configurable before rendering.
- [x] Point light brightness configurable before rendering.
- [x] Shadows visible in the mandatory rendered images.
- [x] Documentation explaining features and usage.

Bonus features are out of scope for the first passing version:

- [ ] Textures.
- [ ] Reflection.
- [ ] Refraction.
- [ ] Particles.
- [ ] Fluids.

## Scene And Render Input Contract

The intended user contract for the passing version is:

```sh
cargo run -- --scene scenes/sphere.ron --width 800 --height 600 > deliverables/sphere.ppm
```

Required command behavior:

- [x] `--scene <path>` selects a scene file.
- [x] `--width <pixels>` overrides the scene width.
- [x] `--height <pixels>` overrides the scene height.
- [x] If width or height is omitted, the scene default is used.
- [x] Output is written to standard output so shell redirection creates the `.ppm`.
- [x] Invalid scene input fails with a useful error message.

Required scene data:

- [x] Image width and height.
- [x] Camera position.
- [x] Camera target or viewing direction.
- [x] Light position.
- [x] Light brightness.
- [x] Ambient light value.
- [x] Object list.
- [x] Per-object color.
- [x] Sphere center and radius.
- [x] Cube center and size.
- [x] Plane point and normal.
- [x] Cylinder center, radius, height, and axis.

The chosen scene file format is `RON` because it maps cleanly to Rust structs
and keeps hand-written audit scenes readable.

## Audit Deliverables

The repo should keep these generated images for audit day:

- [x] `deliverables/sphere.ppm`: one scene with a sphere.
- [x] `deliverables/plane_cube_low_brightness.ppm`: one scene with a flat plane
      and cube, with lower brightness than the sphere image.
- [x] `deliverables/all_objects.ppm`: one scene with one sphere, one cube, one
      cylinder, and one flat plane.
- [x] `deliverables/all_objects_alt_camera.ppm`: the same all-objects scene
      rendered from another camera position.

The repo should also keep:

- [x] `README.md` with usage instructions.
- [x] `docs/ACCEPTANCE_CHECKLIST.md` with this checklist.
- [x] Example scene files under `scenes/`.
- [x] Reproducible commands for regenerating every deliverable.

## Audit Pass/Fail Checklist

- [x] Construct a scene containing at least one sphere, cube, cylinder, and flat
      plane.
      Evidence: `scenes/all_objects.ron` and `deliverables/all_objects.ppm`.
- [x] Rendered image corresponds to the scene definition.
      Evidence: object positions and colors in the scene are visually present.
- [x] Output resolution can be reduced.
      Evidence: `--width` and `--height` render a smaller valid PPM.
- [x] The same scene can be rendered after moving the camera.
      Evidence: `scenes/all_objects_alt_camera.ron`.
- [x] Moved-camera render shows the same scene from another perspective.
      Evidence: `deliverables/all_objects_alt_camera.ppm`.
- [x] Four `.ppm` pictures are provided.
      Evidence: all files listed in `Audit Deliverables` exist.
- [x] One provided image contains a sphere.
      Evidence: `deliverables/sphere.ppm`.
- [x] One provided image contains a flat plane and a cube with lower brightness
      than the sphere image.
      Evidence: `deliverables/plane_cube_low_brightness.ppm`.
- [x] One provided image contains one sphere, one cube, one cylinder, and one
      flat plane.
      Evidence: `deliverables/all_objects.ppm`.
- [x] One provided image contains the all-objects scene from another camera
      position.
      Evidence: `deliverables/all_objects_alt_camera.ppm`.
- [x] Shadows from objects are visible across the mandatory pictures.
      Evidence: shadow rays are enabled and visible in test renders.
- [x] Documentation explains how to create elements, change brightness, and move
      the camera.
      Evidence: `README.md` "Documentation" section covers each mandatory object,
      brightness fields, and camera fields with concrete examples.

## Current Implementation Gap

This checklist is the frozen target. Items remain unchecked until the
corresponding renderer, CLI, scene file, documentation, or deliverable exists in
the repository.
