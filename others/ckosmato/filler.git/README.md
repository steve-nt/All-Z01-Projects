# Filler

A Rust implementation of the Filler game — a two-player territory competition where robots place randomly shaped pieces on a grid, each overlapping exactly one cell of their own previous territory.

---

## Requirements

- [Rust](https://www.rust-lang.org/tools/install)
- [Docker](https://www.docker.com/products/docker-desktop/) (or a Linux VM with Docker)
- The `docker_image` package provided by your school (contains the game engine, maps, and robots)

---

## Project Structure

```
filler/
├── Cargo.toml
├── .gitignore
├── README.md
└── src/
    └── main.rs
```

---

## Build

```bash
cargo build --release
```

Binary will be at `target/release/filler`.

---

## Run Tests

```bash
cargo test
```

Expected output:
```
running 25 tests
test tests::test_best_placement_finds_valid_move ... ok
test tests::test_best_placement_no_valid_move ... ok
...
test result: ok. 25 passed; 0 failed
```

---

## Test the Binary Manually

Create a file called `test_input.txt`:
```
$$$ exec p1 : [robots/bender]
Anfield 20 15:
    01234567890123456789
000 ....................
001 ....................
002 .........@..........
003 ....................
004 ....................
005 ....................
006 ....................
007 ....................
008 ....................
009 ....................
010 ....................
011 ....................
012 .........$..........
013 ....................
014 ....................
Piece 4 1:
.OO.
```

Then run:
```bash
./target/release/filler < test_input.txt
```

You should get a coordinate back like `7 2`.

---

## Docker Setup

### 1 — Prepare

Extract the `docker_image` zip from your school. Copy your source files into it:

```bash
cd docker_image
mkdir -p solution/src
cp /path/to/filler/src/main.rs solution/src/
cp /path/to/filler/Cargo.toml solution/
```

### 2 — Build the image

```bash
docker build -t filler .
```

### 3 — Run the container

```bash
docker run -v "$(pwd)/solution":/filler/solution \
           -v "$(pwd)/linux_game_engine":/filler/game_engine \
           -v "$(pwd)/linux_robots":/filler/robots \
           -v "$(pwd)/maps":/filler/maps \
           -it filler
```

### 4 — Compile inside the container

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
source ~/.cargo/env
cd /filler/solution
cargo build --release
```

### 5 — Run the game

```bash
cd /filler
./game_engine -f maps/map00 -p1 solution/target/release/filler -p2 robots/wall_e
```

---

## Audit Commands

Run each matchup 5 times, alternating `-p1` and `-p2`. You must win at least 4/5.

```bash
cd /filler

# vs wall_e on map00
./game_engine -f maps/map00 -p1 solution/target/release/filler -p2 robots/wall_e
./game_engine -f maps/map00 -p1 solution/target/release/filler -p2 robots/wall_e
./game_engine -f maps/map00 -p1 solution/target/release/filler -p2 robots/wall_e
./game_engine -f maps/map00 -p2 solution/target/release/filler -p1 robots/wall_e
./game_engine -f maps/map00 -p2 solution/target/release/filler -p1 robots/wall_e

# vs h2_d2 on map01
./game_engine -f maps/map01 -p1 solution/target/release/filler -p2 robots/h2_d2
./game_engine -f maps/map01 -p1 solution/target/release/filler -p2 robots/h2_d2
./game_engine -f maps/map01 -p1 solution/target/release/filler -p2 robots/h2_d2
./game_engine -f maps/map01 -p2 solution/target/release/filler -p1 robots/h2_d2
./game_engine -f maps/map01 -p2 solution/target/release/filler -p1 robots/h2_d2

# vs bender on map02
./game_engine -f maps/map02 -p1 solution/target/release/filler -p2 robots/bender
./game_engine -f maps/map02 -p1 solution/target/release/filler -p2 robots/bender
./game_engine -f maps/map02 -p1 solution/target/release/filler -p2 robots/bender
./game_engine -f maps/map02 -p2 solution/target/release/filler -p1 robots/bender
./game_engine -f maps/map02 -p2 solution/target/release/filler -p1 robots/bender
```

---

## How It Works

The robot uses a greedy strategy to pick the best placement each turn:

- **Territory gain** — prefers placements that capture the most empty cells
- **Aggression** — moves toward the opponent to cut off their expansion
- **Adjacency bonus** — extra score for placements next to opponent cells

The placement validator enforces all game rules:
- Exactly **one** cell overlaps own territory
- **Zero** cells overlap opponent territory
- Piece stays fully **within bounds**
