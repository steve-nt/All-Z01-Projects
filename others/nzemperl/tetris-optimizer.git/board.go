package main

import "fmt"

func CreateBoard(size int) [][]rune {
	board := make([][]rune, size)
	for i := range board {
		board[i] = make([]rune, size)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}
	return board
}

func PrintBoard(board [][]rune) {
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}
