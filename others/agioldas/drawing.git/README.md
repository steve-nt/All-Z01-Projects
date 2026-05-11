# Drawing

A Rust application that renders geometric shapes onto a raster image and exports the result as a PNG file. Implemented as a trait-based drawing API using the [raster](https://crates.io/crates/raster) crate.

---

## Overview

The program creates a blank canvas, draws a configurable set of shapes (points, lines, rectangles, triangles, circles), assigns each shape a color from a fixed palette, and saves the image to disk. The default demo produces a 1000×1000 pixel image as `image.png`.

## Features

- **Shape primitives:** point, line, rectangle, triangle, circle
- **Unified API:** all shapes implement the `Drawable` trait (`draw`, `color`)
- **Constructors:** `new(...)` for explicit geometry; `random(width, height)` for points, lines, rectangles, and circles
- **Rectangle normalization:** corners are stored as top-left and bottom-right regardless of argument order
- **Color palette:** red, green, blue, yellow, purple (assigned at creation time)

## Technical overview

| Aspect | Description |
|--------|-------------|
| **Drawable trait** | `fn draw(&self, im: &mut Image)` and `fn color(&self) -> Color` |
| **Image I/O** | `raster::Image` in memory; `raster::save()` for PNG output |
| **RNG** | `rand` for random coordinates and colors |

## Authors

| Name | Platform |
|------|----------|
| Aleksis Gioldaseas | [@agioldas](https://platform.zone01.gr/git/agioldas) |
| Theocharoula Tarara | [@ttarara](https://platform.zone01.gr/git/ttarara) |
| Memos Foteinopoulos | [@mfoteino](https://platform.zone01.gr/git/mfoteino) |

**Repository:** [platform.zone01.gr/git/agioldas/drawing](https://platform.zone01.gr/git/agioldas/drawing)

---

## Getting started

### Prerequisites

- [Rust](https://rustup.rs) (stable toolchain)
- Dependencies are declared in `Cargo.toml`. To add `raster` manually: `cargo add raster`

### Build and run

```bash
cargo build
cargo run
```

Output is written to `image.png` in the project root.

### Tests

```bash
cargo test
```
