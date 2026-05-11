# lem-in

A Go implementation of the classic Lem-in problem: finding optimal vertex-disjoint paths through a network of rooms and simulating the movement of ants.



## ✨ Key Features

- **Input Parsing & Validation**: Reads an input file defining the number of ants, rooms, and tunnels. Checks for common errors such as:
  - Invalid or missing ant count
  - Duplicate room names or rooms starting with `L`/`#`
  - Self-linking, duplicate, or reversed tunnels
  - Tunnels referencing undefined rooms
  - Multiple `##start`/`##end` declarations
- **Graph Construction**: Builds an undirected adjacency list of rooms and tunnels.
- **Connectivity Check**: Uses BFS to ensure at least one path exists between the start and end rooms.
- **Vertex-Disjoint Path Extraction**: Implements vertex splitting + Edmonds–Karp max-flow to compute all vertex-disjoint paths from start to end.
- **Ant Distribution & Simulation**: Greedy assignment of ants to minimize total turns, followed by turn-by-turn simulation printing each move.

## Project Structure

```
lem-in/
├── funcs/
│   ├── input.go               # Input parsing and validation
│   ├── buildConnections.go    # Builds adjacency list
│   ├── startEndConnection.go  # Connectivity BFS check
│   ├── maxFlow.go             # Vertex-splitting + Edmonds–Karp algorithm
│   ├── optimalDistribution.go # Path sorting, ant assignment, and simulation
├── main.go                    # Orchestrates parsing, computation, and output
├── README.md                  # Project overview and instructions
```

## 🧠 Algorithm Overview
1. **Parsing** – strict, fail‑fast scanner; any error aborts with an explanatory message.  
2. **Max‑flow** – each vertex (except *start*/*end*) is split into **v_in** and **v_out** (capacity = 1) to enforce vertex‑disjointness.  
3. **Path extraction** – depth‑first pull of flow‑saturated edges translates the residual network back to human paths.  
4. **Greedy distribution** – ants are queued onto the shortest finishing path each turn (classic *load balancing* heuristic).  
5. **Simulation** – we iterate over turns, printing only active moves to keep output concise.


## 🚀 Getting Started
### Prerequisites
* Go 1.22+  
* A Unix‑like shell (only for running examples)

### Run
```bash
go run . <path/to/input_file>
```

The programme prints:
1. The validated map (echo),
2. All ants’ moves grouped by turn.


## 👥 Contributors
- Kostas Apostolou  
- Yana Kopilova  
- Vicky Apostolou  