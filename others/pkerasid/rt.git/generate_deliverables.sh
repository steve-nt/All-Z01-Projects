#!/usr/bin/env bash
set -e

mkdir -p deliverables

cargo run --release -- --scene scenes/sphere.ron > deliverables/sphere.ppm
cargo run --release -- --scene scenes/plane_cube_low_brightness.ron > deliverables/plane_cube_low_brightness.ppm
cargo run --release -- --scene scenes/all_objects.ron > deliverables/all_objects.ppm
cargo run --release -- --scene scenes/all_objects_alt_camera.ron > deliverables/all_objects_alt_camera.ppm

echo "All deliverables written to deliverables/"
