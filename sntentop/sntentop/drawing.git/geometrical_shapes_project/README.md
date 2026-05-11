# Geometrical Shapes Project

A Rust program that draws geometrical shapes onto a 1000x1000 PNG image using the `raster` crate.

## Requirements

- [Rust](https://www.rust-lang.org/tools/install) (edition 2024)
- Cargo (included with Rust)

## Build & Run

```bash
# Build the project
cargo build

# Run the program
cargo run
```

The program outputs `image.png` in the project root directory.

```bash
# Build and run in one step
cargo run

# Build optimized release binary
cargo build --release
./target/release/geometrical_shapes_project
```

## What It Draws

Each run generates a fresh `image.png` (1000x1000 pixels, black background) containing:

| Shape | Count | Color | Construction |
|-------|-------|-------|--------------|
| Line | 1 | Green `rgb(0, 255, 0)` | Random endpoints |
| Point | 1 | Red `rgb(255, 0, 0)` | Random position |
| Rectangle | 1 | Orange `rgb(255, 165, 0)` | Top-left `(150, 300)`, bottom-right `(50, 60)` |
| Triangle | 1 | Blue `rgb(0, 0, 255)` | Vertices `(500,500)`, `(250,700)`, `(700,800)` |
| Circles | 50 | Magenta `rgb(255, 0, 255)` | Random centers & radii (10–200 px) |

## Project Structure

```
geometrical_shapes_project/
├── src/
│   ├── main.rs                  # Entry point — assembles and draws shapes
│   └── geometrical_shapes.rs    # Shape structs, traits, and draw logic
├── Cargo.toml
└── image.png                    # Generated output (created on first run)
```

## Functions & Methods

### `Point`

| Function | Signature | Description |
|----------|-----------|-------------|
| `new` | `(x: i32, y: i32) -> Point` | Creates a point at the given coordinates. |
| `random` | `(width: i32, height: i32) -> Point` | Creates a point at a random position within the canvas bounds. |

### `Line`

| Function | Signature | Description |
|----------|-----------|-------------|
| `new` | `(p1: &Point, p2: &Point) -> Line` | Creates a line between two existing points. Copies the coordinates out of the references. |
| `random` | `(width: i32, height: i32) -> Line` | Creates a line with two randomly positioned endpoints within the canvas. |

### `Triangle`

| Function | Signature | Description |
|----------|-----------|-------------|
| `new` | `(p1: &Point, p2: &Point, p3: &Point) -> Triangle` | Creates a triangle from three vertices. The three sides are drawn as lines connecting them in order: p1→p2, p2→p3, p3→p1. |

### `Rectangle`

| Function | Signature | Description |
|----------|-----------|-------------|
| `new` | `(p1: &Point, p2: &Point) -> Rectangle` | Creates a rectangle from two corner points (top-left and bottom-right). The other two corners are computed internally. |

### `Circle`

| Function | Signature | Description |
|----------|-----------|-------------|
| `new` | `(center: &Point, radius: i32) -> Circle` | Creates a circle at the given center point with the specified radius in pixels. |
| `random` | `(width: i32, height: i32) -> Circle` | Creates a circle with a random center and a random radius between 10 and 200 pixels. |

### Internal drawing helper

| Function | Signature | Description |
|----------|-----------|-------------|
| `draw_line` | `(image: &mut Image, p1: &Point, p2: &Point, color: Color)` | Private function used by `Line`, `Triangle`, and `Rectangle`. Implements [Bresenham's line algorithm](https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm) to rasterize a straight line between two points pixel by pixel. |

`Circle::draw` uses the [Midpoint circle algorithm](https://en.wikipedia.org/wiki/Midpoint_circle_algorithm) (Bresenham variant) — it exploits 8-fold symmetry to plot all octants of the circle with a single loop.

---

## Traits

### `Drawable`

Implemented by every shape. Requires two methods:

```rust
fn color(&self) -> Color   // returns the fixed RGB color for this shape type
fn draw(&self, image: &mut Image)  // renders the shape onto the image
```

Each shape has a hardcoded color:

| Shape | Color |
|-------|-------|
| Point | Red `rgb(255, 0, 0)` |
| Line | Green `rgb(0, 255, 0)` |
| Triangle | Blue `rgb(0, 0, 255)` |
| Rectangle | Orange `rgb(255, 165, 0)` |
| Circle | Magenta `rgb(255, 0, 255)` |

### `Displayable`

Implemented by `Image` (in `main.rs`). Wraps `raster`'s `set_pixel` with a bounds check so drawing algorithms never panic when coordinates fall outside the canvas.

```rust
fn display(&mut self, x: i32, y: i32, color: Color)
// silently skips pixels where x < 0, x >= width, y < 0, or y >= height
```

---

## Dependencies

| Crate | Version | Purpose |
|-------|---------|---------|
| `raster` | 0.2 | Image creation and PNG export |
| `rand` | 0.4 | Random number generation for shapes |
