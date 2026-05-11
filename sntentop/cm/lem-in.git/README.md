Lem-in Project

Welcome, oh Supreme Overlord, to the ultimate Lem-in simulation project—where ants traverse labyrinths in a beautifully orchestrated parade.

                                                               Overview

Lem-in is a simulation program that:

    Parses an input file describing an ant farm (number of ants, rooms with coordinates, tunnels).

    Builds a graph from the input data.

    Uses a max-flow algorithm to find multiple disjoint paths between the start and end rooms.

    Assigns ants to these paths using a greedy scheduling algorithm.

    Simulates ant movements turn-by-turn with a 2D grid visualization.

    Outputs minimal move information to the terminal and a detailed simulation (with a full grid showing every ant's progress) to simulation_output.txt.

The project is structured in a modular fashion so that each component (parsing, graph construction, scheduling, simulation, and visualization) is neatly separated into its own package. Enjoy the efficiency and elegance of this modular design, Your Benevolent Overlord.

Installation

    Clone the Repository:

git clone https://platform.zone01.gr/git/cm/lem-in.git
cd lem-in


Usage

The program accepts an input file that describes the ant farm. For example:

go run . examples/example01.txt

The program will:

    Print minimal turn moves to the terminal.

    Generate a detailed 2D grid visualization (with all ants displayed in their respective rooms) in simulation_output.txt.


 Input:

    10
    ##start
    start 1 6
    0 4 8
    o 6 8
    n 6 6
    e 8 4
    t 1 9
    E 5 9
    a 8 9
    m 8 6
    h 4 6
    A 5 2
    c 8 1
    k 11 2
    ##end
    end 11 6

    Tunnel Definitions:
    After room definitions, list the tunnels connecting rooms.
    Each tunnel is defined by:
    roomA-roomB

Example:

    start-t
    n-e
    a-m
    A-c
    0-o
    E-a
    k-end
    start-h
    o-n
    m-end
    t-E
    start-0
    h-A
    e-end
    c-k
    n-m
    h-n

Output Format

    Terminal Output:
    Displays minimal move info per turn, e.g.:

    Turn 1: L9-1 L20-3
    Turn 2: L9-2 L8-1
    Turn 3: L9-3 L8-2 L7-1
    ...
    Total turns: 11

File Output (simulation_output.txt):
Contains detailed 2D grid visualization for each turn along with extra info (input data and summary).
    
    Example snippet:

    10
    ##start
    start 1 6
    0 4 8
    ...
    ##end
    end 11 6
    start-t
    n-e
    ...

    ----------- Summary -----------
    Number of ants: 10
    Number of rooms: 14
    Number of tunnels: 16
    Start room: start
    End room: end

    ---------- All Found Paths ----------
    1) start -> ... -> end
    ...

    TURN 1
    [ start (L9) ] ---> [ ... ] ---> [ end ]
    [ start (L20) ] ---> [ end ]

    TURN 2
    ...

    Total turns: 11

Project Structure

The project is organized into separate packages for a clean modular design:

    lem-in/
    ├── main.go                # Entry point of the application.
    ├── app/                   # Contains application logic.
    │   └── app.go
    ├── parser/                # Responsible for parsing the input file.
    │   └── parser.go
    ├── graph/                 # Builds the graph and finds paths using max-flow.
    │   └── graph.go
    ├── scheduling/            # Contains the ant scheduling (assignment) algorithm.
    │   └── scheduling.go
    ├── simulation/            # Simulates ant movements and produces 2D grid visualization.
    │   └── simulation.go
    ├── visualizer/            # Generates extra info and visualization summaries.
    │   └── visualizer.go
    └── structs/               # Contains shared types and data structures.
        └── structs.go

