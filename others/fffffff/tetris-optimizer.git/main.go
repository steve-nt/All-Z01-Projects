package main

import (
	"fmt"
	"os"
	tetris "tetris-optimizer/helpers"
	"time"
)

func main() {
	start := time.Now() // Read tetrominos from the file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: Unable to open the file.")
		return
	}
	defer file.Close()

	tetrominos, err := tetris.GetTetrominos(file)
	if err != nil {
		fmt.Println("ERROR: Bad format .")
		return
	}

	// Normalize tetrominos
	normalizedTetrominos := tetris.NormalizeTetrominos(tetrominos)

	// Sort tetrominos by area in descending order for greedy placement
	sortedTetrominos := tetris.SortTetrominosByArea(normalizedTetrominos)

	// Solve the Tetris puzzle and print the result
	grid := tetris.SolveTetris(normalizedTetrominos, sortedTetrominos)
	if grid != nil {
		fmt.Println("Solution found:")
		tetris.PrintGrid(grid)
	} else {
		fmt.Println("ERROR: No solution found.")
	}
	stop := time.Since(start)
	fmt.Printf("%s\n", stop)
}
