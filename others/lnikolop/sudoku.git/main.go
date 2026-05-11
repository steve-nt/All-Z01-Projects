package main

import (
	"os"

	"sudoku/board"
	"sudoku/helper"
)

func main() {
	args := os.Args[1:]
	helper.CheckArgs(args)
	board := board.CreateBoard()
}
