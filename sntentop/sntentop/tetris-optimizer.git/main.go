package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	X, Y int
}

type Tetromino struct {
	Blocks []Point
}

type Solver struct {
	Grid        [][]rune
	Tetrominoes []Tetromino
	PlacedCount int
	Gridsize    int
}

func main() {
	//Check if user provided a filename argument
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}

	// Read the file
	tetrominoes, err := parseFile(os.Args[1])
	if err != nil || len(tetrominoes) == 0 {
		fmt.Println("ERROR")
		return
	}

	// Calculate the grid size needed for all tetrominoes
	totalBlocks := len(tetrominoes) * 4
	gridSize := int(math.Ceil(math.Sqrt(float64(totalBlocks))))

	// Create the solver with an empty grid
	solver := &Solver{
		Grid:        createEmptyGrid(gridSize),
		Tetrominoes: tetrominoes,
		PlacedCount: 0,
		GridSize:    gridSize,
	}

	// Try to solve using backtracking
	if solver.Solve(0) {
		printGrid(solver.Grid)
	} else {
		fmt.Println("ERROR")
	}
}

// parseFile function reads the input file and extracts tetrominoes
func parseFile(filename string) ([]Tetromino, error) {
	//Open the file
}
