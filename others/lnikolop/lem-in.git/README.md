# 🐜 Lem-in: Ant Farm Pathfinding

# 📌 Description - Overview

Lem-in is a Go program that simulates ants finding the shortest path through a digital ant farm (colony) from a start room to an end room. The program reads colony data from a specially-formatted input file, validates it, finds the optimal paths for the ants, and displays their movements turn by turn.
The goal is to move all ants across a graph (rooms connected by tunnels) in as few turns as possible, adhering to the rule that only one ant may occupy a room or tunnel segment at any given time (excluding the start and end rooms).

This project showcases efficient graph traversal, conflict-free path selection, and turn-minimized entity scheduling.

# 🧩 Data Processing Steps - Features

Parses and validates colony data files

Models rooms and tunnels as a graph structure

Implements pathfinding algorithms to find optimal paths

Simulates ant movements while respecting constraints:

One ant per room (except start/end)

Each tunnel used only once per turn

Displays results in the required format

Comprehensive error handling

# 🔍 Pathfinding Algorithm

The project uses a Depth-First Search (DFS) approach to discover all valid paths from the start room to the end room. This provides a complete view of the graph's structure and helps in analyzing every possible route that ants can take.

# 📈 Path Optimization
Once all paths are discovered, the following optimization process occurs:

Generate subsets of paths which only includes subsets where no paths share intermediate rooms.  
Filter subsets  
Choose subsets that maximize the number of simultaneous paths (no collisions).  
Select best subset  
Among valid subsets, the one yielding the lowest number of turns for all ants to reach the end is chosen. 
 
# 🧮 Ant Allocation Strategy
With an optimal subset of non-conflicting paths found:

Ants are allocated to paths based on their lengths.  
Shorter paths may receive more ants to balance the arrival times.  
Allocation aims to equalize total travel time and minimize the number of turns.  
An internal FindBestAllocation() routine handles this by simulating ant distribution and evaluating turn counts for each valid configuration.

# 📤 Printing Results
The output is formatted step-by-step to show how ants move:

L1-roomX L2-roomY ...  
Each line represents a single turn.  
Movements are printed in the order ants are dispatched, showing the room they move into on each turn.

# Installation
Ensure you have Go installed (version 1.16 or higher recommended)

Clone the repository:  
```
git clone https://platform.zone01.gr/git/lnikolop/lem-in.git
```

Install some required packages:
```
sudo apt-get update
sudo apt-get install -y \
    libvorbis-dev \
    libopenal-dev \
    libx11-dev \
    libgl1-mesa-dev \
    xorg-dev \
    libglx-dev
```

Build the program:  
```
./go build
``` 

# 🚀 Usage
Run the program with an input file containing your map definition:
``` 
./lem-in examples/example00.txt 
```
Example example00.txt:  
```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1
```
Example output:  
```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1

L1-2 
L1-3 L2-2 
L1-1 L2-3 L3-2 
L2-1 L3-3 L4-2 
L3-1 L4-3 
L4-1 
```

