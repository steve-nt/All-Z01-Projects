package piscine

import (
	"github.com/01-edu/z01"
)

const N = 8

func printSolution(board [N]int) {
	for _, col := range board {
		z01.PrintRune(rune(col + '1'))
	}
	z01.PrintRune('\n')
}

func isSafe(board [N]int, row, col int) bool {
	for i := 0; i < row; i++ {
		if board[i] == col || board[i]-i == col-row || board[i]+i == col+row {
			return false
		}
	}
	return true
}

func solveNQueens(board [N]int, row int) {
	if row == N {
		printSolution(board)
		return
	}
	for col := 0; col < N; col++ {
		if isSafe(board, row, col) {
			board[row] = col
			solveNQueens(board, row+1)
		}
	}
}

func EightQueens() {
	var board [N]int
	solveNQueens(board, 0)
}
