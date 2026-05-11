# AI Changelog

This file records AI-assisted engineering work so reviewers can inspect the
intent, constraints, and quality gates behind generated or refactored code.

## Workflow Rules

- Keep specifications and implementation plans in the repo.
- Work in small phases that map to `IMPLEMENTATION_PLAN.md`.
- Commit each phase separately when practical.
- Record the prompt intent, changed scope, and verification commands.
- Reviewers should inspect both the plan and the resulting code.

## 2026-05-03 - Phase 1: Freeze Requirements

Prompt intent:

- Implement Phase 1 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Converted the 01-edu `rt` subject and audit into an internal acceptance
  checklist.
- Defined mandatory features, planned input contract, and audit deliverables.
- Marked Phase 1 complete in `IMPLEMENTATION_PLAN.md`.

Quality gate:

- Documentation-only change.
- No code tests required.

Review focus:

- Confirm the checklist maps directly to the subject and audit requirements.
- Confirm no renderer behavior changed.

## 2026-05-03 - Phase 2: Reshape Codebase

Prompt intent:

- Implement Phase 2 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Split the single-file renderer into modules for math, rays, camera, color,
  light, geometry, scene, renderer, and config.
- Introduced a shared `Hittable` trait for geometry.
- Kept `main.rs` as a thin executable wrapper.
- Preserved sphere, plane, lighting, shadows, and PPM output behavior.

Quality gate:

```sh
cargo fmt
cargo test
```

Review focus:

- Confirm behavior stayed stable while the code moved into modules.
- Confirm the geometry module is ready for cube and cylinder additions.

## 2026-05-03 - Phase 3: Scene Model And CLI

Prompt intent:

- Implement Phase 3 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Added a Rust scene model with image, camera, light, ambient, and object data.
- Added RON scene parsing through `serde` and `ron`.
- Added CLI options for `--scene`, `--width`, and `--height`.
- Added `scenes/demo.ron`.
- Updated `README.md` with render commands.
- Marked Phase 3 complete in `IMPLEMENTATION_PLAN.md`.

Quality gate:

```sh
cargo fmt
cargo test
cargo run -- --scene scenes/demo.ron --width 4 --height 3
```

Review focus:

- Confirm CLI errors go to stderr and PPM image data stays on stdout.
- Confirm scene parsing has useful error messages.
- Confirm resolution overrides work without editing source code.

## 2026-05-03 - Phase 4: Mandatory Geometry

Prompt intent:

- Implement Phase 4 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Added axis-aligned cube intersection with face normals.
- Added finite cylinder intersection with side and cap normals.
- Added `Cube` and `Cylinder` scene object variants.
- Updated the built-in demo scene and `scenes/demo.ron` to include all object
  types.
- Added `scenes/all_objects.ron` as the first mandatory all-objects scene.
- Marked Phase 4 complete in `IMPLEMENTATION_PLAN.md`.

Quality gate:

```sh
cargo fmt
cargo test
cargo run -- --scene scenes/all_objects.ron --width 4 --height 3
```

Review focus:

- Confirm cube slab intersections handle outside and inside rays correctly.
- Confirm cylinder side and cap normals are correct.
- Confirm the all-objects scene visibly includes a sphere, plane, cube, and
  cylinder.

## 2026-05-04 - Phase 5: Camera and Lighting Controls

Prompt intent:

- Implement Phase 5 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Verified that camera `origin`, `look_at`, and `up` are fully wired through
  the scene model and renderer — no code changes needed.
- Verified that `PointLight.brightness` and scene `ambient` are already
  honoured in the shading pass.
- Added `scenes/all_objects_alt_camera.ron`: same geometry as the main scene
  but shot from the right side (`origin: 3.5, 1.2, 1.0`) to prove camera
  movement works without source edits.
- Added `scenes/all_objects_dim.ron`: same geometry with brightness `0.25` and
  ambient `0.06` to give a visibly darker result.
- Marked Phase 5 complete in `IMPLEMENTATION_PLAN.md`.

Quality gate:

```sh
cargo test
cargo run -- --scene scenes/all_objects_alt_camera.ron --width 80 --height 60 > /tmp/alt.ppm
cargo run -- --scene scenes/all_objects_dim.ron --width 80 --height 60 > /tmp/dim.ppm
```

Review focus:

- Confirm the alternate-camera render shows the scene from a clearly different
  angle with shadows still present.
- Confirm the dim render is noticeably darker than `all_objects.ron` at the
  same resolution.
- Confirm no new code was introduced — controls were already implemented in
  prior phases.

## 2026-05-04 - Phase 6: Renderer Output and Performance

Prompt intent:

- Implement Phase 6 from `IMPLEMENTATION_PLAN.md`.

Scope:

- Added `scenes/sphere.ron`: sphere-only scene at 800×600 with bright lighting.
- Added `scenes/plane_cube_low_brightness.ron`: plane and cube scene with
  brightness `0.25` and ambient `0.06` for a visibly darker result.
- Created `deliverables/` folder containing all four required 800×600 PPM
  files: `sphere.ppm`, `plane_cube_low_brightness.ppm`, `all_objects.ppm`,
  and `all_objects_alt_camera.ppm`.
- Added `generate_deliverables.sh`: a reproducible shell script that regenerates
  all four deliverables with `cargo run --release`.
- Confirmed 800×600 renders complete in ~9 seconds in release mode.
- Confirmed P3 PPM headers are valid for every deliverable.
- Ticked all Phase 6 items in `IMPLEMENTATION_PLAN.md` and `IMPLEMENTATION_PLAN_VISUAL.md`.
- Ticked all deliverable items in `docs/ACCEPTANCE_CHECKLIST.md`.

Quality gate:

```sh
bash generate_deliverables.sh
head -3 deliverables/sphere.ppm
head -3 deliverables/all_objects.ppm
```

Review focus:

- Confirm all four PPM files are present in `deliverables/` and have valid P3
  headers.
- Confirm `generate_deliverables.sh` reproduces them from a clean state.
- Confirm the plane-cube image is visibly darker than the sphere image.

## 2026-05-06 - Phase 7: Tests and Verification

Prompt intent:

- Implement Phase 7 from `IMPLEMENTATION_PLAN.md` one step at a time.

Scope:

- Added sphere intersection tests: ray hit from outside, ray miss, outward unit
  normal, ray hit from inside.
- Added plane intersection tests: ray hit, parallel ray miss, normal returned.
- Expanded math tests: cross product perpendicularity, normalized unit length.
- Added scene parsing tests: valid RON file loads correctly, missing file gives
  a useful error message.
- Added smoke render test: renders `scenes/sphere.ron` at 4×3 via the full
  pipeline and asserts a valid PPM header and pixel count.
- Ticked all Phase 7 items in `IMPLEMENTATION_PLAN.md`.
- Updated current state in `IMPLEMENTATION_PLAN_VISUAL.md`.
- Total test count: 22.

Quality gate:

```sh
cargo test
```

Review focus:

- Confirm `cargo test` passes with 22 tests.
- Confirm every Phase 7 exit criterion is met.
- Confirm Phase 8 and Phase 9 items accurately reflect work done in prior phases.
