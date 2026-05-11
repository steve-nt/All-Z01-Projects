package main

import (
	"math"
)

func Solve(tetrominoes []Tetromino) [][]rune {
	size := int(math.Ceil(math.Sqrt(float64(len(tetrominoes) * 4))))

	var board [][]rune

	for {
		board = CreateBoard(size)
		if placeTetrominoes(board, tetrominoes, 0) {
			break
		}
		size++
	}

	return board
}

func placeTetrominoes(board [][]rune, tetrominoes []Tetromino, index int) bool {
	if index >= len(tetrominoes) {
		return true
	}

	size := len(board)
	tetromino := tetrominoes[index]

	for i := 0; i <= size-tetromino.Height; i++ {
		for j := 0; j <= size-tetromino.Width; j++ {
			if canPlace(board, tetromino, i, j) {
				place(board, tetromino, i, j, tetromino.Letter)
				if placeTetrominoes(board, tetrominoes, index+1) {
					return true
				}
				place(board, tetromino, i, j, '.')
			}
		}
	}

	return false
}

func canPlace(board [][]rune, tetromino Tetromino, row, col int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if tetromino.Shape[i][j] == '#' {
				boardRow := row + i - tetromino.MinRow
				boardCol := col + j - tetromino.MinCol

				// Check if it's out of bounds
				if boardRow < 0 || boardCol < 0 || boardRow >= len(board) || boardCol >= len(board) {
					return false
				}

				// Check if cell is already occupied
				if board[boardRow][boardCol] != '.' {
					return false
				}
			}
		}
	}
	return true
}

func place(board [][]rune, tetromino Tetromino, row, col int, letter rune) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if tetromino.Shape[i][j] == '#' {
				boardRow := row + i - tetromino.MinRow
				boardCol := col + j - tetromino.MinCol
				board[boardRow][boardCol] = letter
			}
		}
	}
}
