package main

import (
	"fmt"
	"lem-in/models"
	"lem-in/repositories"
	"lem-in/services"
	"lem-in/visualizer"

	"os"
)

func main() {
	filename, visualizerEnabled := Args()
	if filename == "" {
		return
	}
	// Create repository to handle file data operations
	fileDataRepository := repositories.NewFileDataRepository(os.Args[1])

	numOfAnts, rooms, err := fileDataRepository.FetchData()
	if err != nil || numOfAnts <= 0 || len(rooms) < 2 {
		fmt.Println("ERROR: ", err)
		return
	}

	ants := make([]*models.Ant, numOfAnts)
	for i := 0; i < numOfAnts; i++ {
		ants[i] = models.NewAnt(i+1, rooms[0])
	}
	// Create service to find paths
	pathsService := services.NewPathService(rooms[0], rooms[1])
	paths := pathsService.FindAllPaths()
	if len(paths) == 0 {
		fmt.Println("ERROR: No paths found")
		return
	}
	// Create service to find good subsets of paths
	subsetService := services.NewSubsetService(paths, numOfAnts)
	subsets := subsetService.GetGoodSubsets()
	if len(subsets) == 0 {
		fmt.Println("ERROR: No good subsets found")
		return
	}
	// Create service to optimize path selection for fastest ant movement
	optimizationService := services.NewPathOptimizationService(subsets)
	bestSubset, bestAlloc, bestTime := optimizationService.Optimize(numOfAnts)
	if len(bestSubset) == 0 || len(bestAlloc) == 0 || bestTime == 0 {
		fmt.Println("ERROR: No good subsets found")
		return
	}

	printInputFile(os.Args[1])
	fmt.Println()

	// Create service to allocate ants to paths
	antAllocationService := services.NewAntAllocationService(ants, bestSubset, bestAlloc)
	antAllocationService.Allocate()

	// Create service to simulate ant movement through rooms
	simulationService := services.NewSimulationService(ants)
	rounds := simulationService.Simulate(bestTime)

	for _, round := range rounds {
		for _, move := range round {
			fmt.Print(move)
		}
		if len(round) > 0 {
			fmt.Println()
		}
	}

	if visualizerEnabled {
		visualizer.Visualization(rooms[0], rooms[1], rooms, len(ants), rounds)
	}
}

// Args handles command line argument parsing
func Args() (string, bool) {
	var filename string
	visualizer := false
	// os.Args = []string{"", "examples/example09.txt"}

	for _, arg := range os.Args[1:] {
		if arg == "--visualizer" || arg == "-v" {
			visualizer = true
		} else if filename == "" {
			filename = arg
		} else {
			fmt.Println("Usage: go run main.go [filename] [--visualizer | -v]")
			return "", false
		}
	}

	if filename == "" {
		fmt.Println("Usage: go run main.go [filename] [--visualizer | -v]")
		return "", false
	}

	return filename, visualizer
}

// printInputFile reads and prints the content of the input file
func printInputFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	fmt.Print(string(content))
	fmt.Println()
}
