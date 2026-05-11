package app

import (
	"fmt"
	"os"

	"lem-in/graph"
	"lem-in/parser"
	"lem-in/scheduling"
	"lem-in/simulation"
	"lem-in/visualizer"
)

// Run executes the main application workflow.
func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}
	inputFile := os.Args[1]

	// Parse input
	antCount, rooms, tunnels, err := parser.ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Basic validation
	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}
	var startFound, endFound bool
	for _, r := range rooms {
		if r.IsStart {
			startFound = true
		}
		if r.IsEnd {
			endFound = true
		}
	}
	if !startFound || !endFound {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Build graph and find paths
	g, err := graph.BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	paths, err := graph.GetOptimalPaths(g)
	if err != nil || len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Assign ants and simulate
	assignment := scheduling.AssignAnts(antCount, paths)
	extraInfo := visualizer.PrintExtraInfo(antCount, rooms, tunnels, paths, assignment)
	simulation.SimulateMultiPath(antCount, paths, assignment, extraInfo)
}
