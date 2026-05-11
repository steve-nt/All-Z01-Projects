package solver

import (
	"strconv"
)

const N = 9

// isValid checks if a given number can be placed at board[row][col]
func IsValid(board [N][N]int, row, col, num int) bool {
	for x := 0; x < N; x++ {
		if board[row][x] == num || board[x][col] == num {
			return false
		}
	}
	// Check the 3x3 subgrid for the same number
	startRow, startCol := row/3*3, col/3*3
	for i := startRow; i < startRow+3; i++ {
		for j := startCol; j < startCol+3; j++ {
			if board[i][j] == num {
				return false
			}
		}
	}

	// If no conflicts, the placement is valid
	return true
}

// solveSudoku uses backtracking to solve the Sudoku puzzle
func SolveSudoku(board [N][N]int, solvedonce bool) (bool, [N][N]int) {
	var row, col int
	found := false

	// Find the first empty cell (indicated by 0)
	for i := 0; i < N && !found; i++ {
		for j := 0; j < N; j++ {
			if board[i][j] == 0 {
				row, col = i, j
				found = true
				break
			}
		}
	}

	// If there are no empty cells, the puzzle is solved
	if !found {
		return true, board
	}

	// Try placing numbers 1 through 9 in the empty cell
	if !solvedonce {
		for num := 1; num <= 9; num++ {
			if IsValid(board, row, col, num) {
				board[row][col] = num
				solved, solvedBoard := SolveSudoku(board, false)
				if solved {
					return true, solvedBoard
				}
				board[row][col] = 0
			}
		}
	} else {
		for num := 9; num >= 1; num-- {
			if IsValid(board, row, col, num) {
				board[row][col] = num
				solved, solvedBoard := SolveSudoku(board, false)
				if solved {
					return true, solvedBoard
				}
				board[row][col] = 0
			}
		}
	}

	return false, board
}

// createBoard converts input strings to a Sudoku board
func CreateBoard(args []string) ([N][N]int, string) {
	var board [N][N]int

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if args[i][j] == '.' {
				board[i][j] = 0
			} else {
				num, _ := strconv.Atoi(string(args[i][j]))
				if IsValid(board, i, j, num) {
					board[i][j] = num
				} else {
					return board, "Error: duplicate numbers in row, col or 3x3 subgrid"
				}
			}
		}
	}
	count := 0
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if board[i][j] != 0 {
				count++
			}
		}
	}
	if count < 17 {
		return board, "Error: Sudoku has less than 17 numbers imputed"
	}

	return board, ""
}
