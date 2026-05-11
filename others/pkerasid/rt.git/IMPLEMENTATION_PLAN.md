# RT Implementation Plan

This plan targets the 01-edu `rt` subject and its audit requirements:

- Subject: <https://github.com/01-edu/public/tree/master/subjects/rt>
- Audit: <https://github.com/01-edu/public/blob/master/subjects/rt/audit/README.md>

## Current Baseline

All phases complete. The project is audit-ready.

- [x] Rust project exists and builds.
- [x] PPM `P3` output exists and is valid for all scenes.
- [x] Camera, spheres, plane, cubes, cylinders, light, shadows, and shading all work.
- [x] Scene format defined (RON) with CLI flags `--scene`, `--width`, `--height`.
- [x] Cube intersection implemented.
- [x] Cylinder intersection implemented.
- [x] All 4 mandatory audit images rendered to `deliverables/` at 800×600.
- [x] User-facing documentation written in `README.md`.
- [x] 22 unit tests passing (`cargo test`).

## Phase 1 - Freeze Requirements

- [x] Convert the subject and audit into an internal acceptance checklist.
- [x] Define the minimum mandatory feature set for the first passing version:
  - sphere
  - plane
  - cube
  - cylinder
  - movable camera
  - configurable brightness
  - shadows
  - 800x600 output
  - smaller test resolutions
- [x] Decide the input contract for scenes and renders.
- [x] Decide the exact deliverables to keep in the repo for audit day.

Exit criteria:

- [x] We can point to a written pass/fail checklist for every audit question.

## Phase 2 - Reshape the Codebase

- [x] Split the current single-file renderer into clear modules:
  - math
  - ray
  - camera
  - color
  - light
  - geometry
  - scene
  - renderer
  - cli or config
- [x] Introduce shared traits or enums for hittable objects.
- [x] Keep the renderer deterministic and easy to test.
- [x] Preserve the existing working sphere and plane behavior while refactoring.

Exit criteria:

- [x] The codebase is modular enough to add cube and cylinder without turning `main.rs` into a bottleneck.

## Phase 3 - Scene Model and CLI

- [x] Define a Rust scene model that supports:
  - camera position and target
  - image width and height
  - light position and brightness
  - object color and transform data
- [x] Add a CLI interface such as:
  - `cargo run -- --scene scenes/sphere.ron`
  - `cargo run -- --scene scenes/all_objects.ron --width 400 --height 300`
- [x] Choose a simple scene format that is easy to document and edit in Rust:
  - `ron`
  - `json`
  - `yaml`
- [x] Validate scene input with useful error messages.

Exit criteria:

- [x] An auditor can render a chosen scene and reduce resolution without editing source code.

## Phase 4 - Mandatory Geometry

- [x] Keep sphere intersection as a stable reference object.
- [x] Keep plane intersection as the flat plane implementation.
- [x] Implement cube intersection.
- [x] Implement cylinder intersection.
- [x] Ensure every object supports location changes through scene data.
- [x] Confirm normals are correct for all visible faces and curved surfaces.

Exit criteria:

- [x] A single scene can contain one sphere, one plane, one cube, and one cylinder and render correctly.

## Phase 5 - Camera and Lighting Controls

- [x] Support moving the camera position.
- [x] Support changing the camera viewing direction or look-at target.
- [x] Support adjustable light brightness.
- [x] Preserve visible shadows across all required scenes.
- [x] Keep shading simple and robust before considering any bonus effects.

Exit criteria:

- [x] The same scene can be rendered from at least two different camera positions.
- [x] Lower-brightness scenes visibly differ from brighter ones.

## Phase 6 - Renderer Output and Performance

- [x] Keep `P3` PPM output valid for all scenes.
- [x] Add resolution overrides for fast test renders.
- [x] Confirm 800x600 final renders complete in acceptable time.
- [x] Avoid obvious regressions while adding geometry and parsing.
- [x] Add a repeatable command set for generating audit deliverables.

Exit criteria:

- [x] Fast preview renders are available.
- [x] Final 800x600 renders are reproducible with documented commands.

## Phase 7 - Tests and Verification

- [x] Add unit tests for:
  - vector math
  - ray-object intersections
  - normals
  - scene parsing
  - PPM header generation
- [x] Add focused regression tests for cube and cylinder edge cases.
- [x] Add smoke tests that render tiny scenes at low resolution.
- [x] Manually verify the four required scenes before audit.

Exit criteria:

- [x] `cargo test` covers the math, parser, and mandatory object intersections.
- [x] Manual checks confirm scene content, brightness differences, camera movement, and shadows.

## Phase 8 - Audit Deliverables

- [x] Prepare 4 required `.ppm` outputs:
  - sphere-only scene
  - plane + cube with lower brightness
  - one of each object
  - same all-object scene from a different camera position
- [x] Store the scenes used to generate those outputs in the repo.
- [x] Add a simple generation script or documented command list.
- [x] Confirm one command path exists for auditors to reproduce results.

Exit criteria:

- [x] The repo contains the exact 4 mandatory `.ppm` files and the scenes that generated them.

## Phase 9 - Documentation

- [x] Write user documentation that explains:
  - what the ray tracer supports
  - how to run it
  - how to reduce resolution
  - how to create each mandatory object
  - how to change brightness
  - how to move the camera
- [x] Include concrete scene examples for sphere, plane, cube, and cylinder.
- [x] Include at least one example command per required audit action.
- [x] Keep the docs aligned with the actual CLI and scene format.

Exit criteria:

- [x] A new user can create objects, change brightness, move the camera, and render the required images without guessing.

## Phase 10 - Final Audit Rehearsal

- [x] Run the full command flow from a clean terminal session.
- [x] Re-render the mandatory outputs.
- [x] Re-check every audit question against the actual repo state.
- [x] Fix any mismatch between docs, scenes, and rendered results.

Exit criteria:

- [x] The project passes the functional audit checklist end to end.

## Suggested Execution Order

- [x] First pass: phases 1 to 5
- [x] Stabilization pass: phases 6 and 7
- [x] Delivery pass: phases 8 to 10
