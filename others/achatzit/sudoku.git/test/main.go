package main

import (
	"fmt"
	"os"

	"sudoku/check"
	"sudoku/sudokuboard"
)

func main() {
	var sudoku sudokuboard.Sudoku
	args := os.Args[1:]
	err := check.CheckArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	errorCreation := sudoku.CreateBoard(args)
	fmt.Printf("Initial Sudoku\n\n")
	sudoku.PrintBoard()
	if errorCreation != "" {
		fmt.Println(errorCreation)
	}

	sudoku.SolveSudoku()

	if sudoku.Solved {
		fmt.Printf("Solved Sudoku\n\n")
		sudoku.PrintBoard()
		return
	} else {
		fmt.Println("Error: Sudoku cannot be solved")
		return
	}
}
