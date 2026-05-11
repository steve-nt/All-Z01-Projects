# Lem-in: Ant Farm Pathfinding

## Overview

**Lem-in** is a Go-based pathfinding project that simulates ants traversing through a digital ant farm. The goal is to calculate and display the optimal movement of ants from a `##start` room to an `##end` room using a set of rules and constraints. It demonstrates key concepts in algorithmic problem solving, parsing structured input, and simulating movements.

---

## Objectives

- Simulate the movement of `n` ants across a graph-like structure of rooms and tunnels.
- Find the most efficient path(s) from `##start` to `##end`.
- Output the result of the simulation with precise format and movement rules.
- Ensure robust error handling and validation of input data.
- Implement in Go using only standard packages.

---

## How It Works

- A map of an ant colony is defined via input, consisting of:
  - A number of ants.
  - Rooms with coordinates.
  - Tunnels connecting the rooms.
- Ants begin in the `##start` room.
- The goal is to reach the `##end` room using the fewest number of moves.
- Movement rules prevent traffic jams and allow one ant per room (except for `##start` and `##end`).

---

## Input Format

The input file contains:

1. **Number of Ants** – A single integer.
2. **Rooms** – Defined as `name coord_x coord_y`.
3. **Special Rooms** – Marked with `##start` or `##end`.
4. **Tunnels** – Defined as `room1-room2`.
5. **Comments** – Lines starting with `#` (ignored unless `##start`/`##end`).

---

## Output Format

The program prints:

1. The original content (input).
2. Each move in the format:
   ```
   L<ant_number>-<room_name> L<ant_number>-<room_name> ...
   ```

Each line represents a turn where ants move through valid paths.

---

## Error Handling

The program validates the input and returns a generic error on failure:
```
ERROR: invalid data format
```
Optionally, more specific errors can be printed, such as:
```
ERROR: invalid data format, invalid number of ants
ERROR: invalid data format, no start room found
```

---

## Constraints

- Room names must not start with `L` or `#` and must contain no spaces.
- Each room has unique coordinates.
- No duplicate rooms or tunnels.
- Each tunnel connects exactly two distinct rooms.
- Each room (except start/end) holds at most one ant at a time.
- Tunnels can only be used once per turn.
- Ants must avoid collisions and congestion.

---

## Usage

### Run the program

```bash
$ go run . <input_file>
```

### Example

```bash
$ go run . test1.txt
3
##start
0 1 0
##end
1 5 0
2 9 0
3 13 0
0-2
2-3
3-1

L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3
L3-1
```

---

## Bonus Feature

You can develop a **visualizer** to animate the ant movements:

```bash
$ ./lem-in input.txt | ./visualizer
```

Room coordinates are helpful in drawing the layout.

---

## Project Goals

This project will help you gain experience in:

- Graph traversal and pathfinding algorithms (BFS, Dijkstra, etc.)
- Structured input parsing
- Efficient string manipulation
- Error handling in Go
- Writing clean, idiomatic Go code
- Testing with unit files

---

## Development Guidelines

- Use only Go standard packages.
- Follow Go best practices.
- Create a modular and maintainable codebase.
- Include unit tests where applicable.

---

## Authors

- **Theocharoula Tarara**  
- **Dionysios Pappas**  
- **Stefanos Ntentopoulos**  

---

## Special Thanks

Special thanks to **GENAI (ChatGPT and Deepseek)** for the guidance and inspiration throughout the project.

---