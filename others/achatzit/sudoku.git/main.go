package main

import (
	"fmt"
	"os"

	"sudoku/check"
	"sudoku/solver"
)

func main() {
	args := os.Args[1:]

	// Check if the number of arguments is correct
	if len(args) != 9 {
		fmt.Println("Error: number of arguments must be 9")
		return
	}
	// Check if the arguments are valid
	if err := check.CheckArgs(args); err != nil {
		fmt.Println(err)
		return
	}

	// Create the Sudoku board from the arguments
	board, err := solver.CreateBoard(args)
	if err != "" {
		fmt.Println(err)
		return
	}

	solvedonce := false
	solved, solvedBoard := solver.SolveSudoku(board, solvedonce)
	if solved {
		// Print the solved Sudoku board
		firstSolution := solvedBoard
		solvedonce = true
		solvedagain, solvedBoard := solver.SolveSudoku(board, solvedonce)
		if solvedagain && firstSolution == solvedBoard {
			printBoard(solvedBoard)
		} else if solvedagain && firstSolution != solvedBoard {
			fmt.Println("Error: Sudoku has more than one solutions")
		}

	} else {
		// Print an error if no solution is found
		fmt.Println("Error: Sudoku cannot be solved")
	}
}

// printBoard prints the Sudoku board
func printBoard(board [9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if j == 8 {
				fmt.Printf("%d", board[i][j])
			} else {
				fmt.Printf("%d ", board[i][j])
			}
		}
		fmt.Println()
	}
}
