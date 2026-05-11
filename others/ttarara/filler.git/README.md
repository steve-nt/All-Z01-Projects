# Filler

Rust bot for the 01-edu `filler` project.

`filler` is an algorithmic game where two robots take turns placing random pieces on a board. A valid move must overlap exactly one cell of the player's own territory, must not overlap the opponent, and must stay inside the board. The player with the largest occupied area wins.

## Project Status

The required project is ready:

- Docker image builds successfully.
- The bot runs inside the official game engine container.
- Unit tests pass.
- Required benchmark matchups pass: `wall_e`, `h2_d2`, and `bender`.
- `terminator` is treated as optional bonus, as described in the subject.

Latest local check:

```text
wall_e  on map00: 5/5
h2_d2   on map01: 5/5
bender  on map02: 5/5
```

## Layout

```text
docker_image/
  Dockerfile
  maps/
  linux_game_engine      # official bundle, gitignored
  m1_game_engine         # official bundle, gitignored
  linux_robots/          # official bundle, gitignored
  m1_robots/             # official bundle, gitignored
  solution/
    Cargo.toml
    src/
      main.rs
      board.rs
      piece.rs
      io.rs
      strategy.rs
scripts/
  build-image.sh
  smoke.sh
  bench.sh
```

The official engine and robot binaries come from the provided `filler.zip`. They are intentionally not committed.

## Requirements

- Docker
- Rust/Cargo for running unit tests on the host
- Official 01-edu `filler.zip` contents extracted into `docker_image/`

Required untracked bundle files:

```text
docker_image/linux_game_engine
docker_image/m1_game_engine
docker_image/linux_robots/
docker_image/m1_robots/
```

## Quick Start

From the repository root:

```sh
scripts/build-image.sh
scripts/smoke.sh
```

`smoke.sh` builds the Rust bot inside the container and runs one game against `wall_e`.

Run the full benchmark sweep:

```sh
scripts/bench.sh
```

Run everything in one command:

```sh
scripts/build-image.sh && scripts/smoke.sh && scripts/bench.sh
```

## Manual Docker Run

Start the container:

```sh
docker run --rm -v "$(pwd)/docker_image/solution":/filler/solution -it filler
```

Inside the container:

```sh
cd /filler/solution
rm -f Cargo.lock
cargo build --release

cd /filler
./linux_game_engine -f maps/map00 -p1 solution/target/release/filler -p2 linux_robots/wall_e
```

For Apple Silicon or ARM Linux containers, use `m1_game_engine` and `m1_robots`.

## Unit Tests

Run tests from the host:

```sh
cargo test --manifest-path docker_image/solution/Cargo.toml
```

Current test coverage includes:

- Player detection from `$$$ exec p<number>`.
- Anfield parsing.
- Piece parsing, including `O` piece cells from the bundled engine.
- Legal placement validation.
- Boundary conditions.
- Coordinate output format: `X Y\n`.
- Strategy fallback: `0 0` when no legal move exists.

## Strategy

The bot:

1. Parses the current board and incoming piece from standard input.
2. Enumerates every legal placement.
3. Rejects placements that overlap the opponent, overlap more than one own cell, or go out of bounds.
4. Builds a Chebyshev-distance field from opponent cells.
5. Chooses the legal placement closest to the opponent.
6. Prints coordinates as `X Y\n`.

If no legal placement exists, it prints `0 0\n`, matching the subject's expected behavior.

## Useful Commands

```sh
# Build Docker image
scripts/build-image.sh

# One visible game
scripts/smoke.sh

# Audit-style sweep
scripts/bench.sh

# Unit tests
cargo test --manifest-path docker_image/solution/Cargo.toml
```

## Notes

- `Cargo.lock` is gitignored for `docker_image/solution` because newer host Cargo versions can generate a lockfile format that Rust 1.63 in the official container cannot read.
- The scripts remove `Cargo.lock` before container builds for this reason.
- Bonus visualizer is not implemented.
