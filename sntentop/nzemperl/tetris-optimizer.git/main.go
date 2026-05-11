package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}

	filePath := os.Args[1]
	tetrominoes, err := ParseFile(filePath)
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	board := Solve(tetrominoes)
	PrintBoard(board)
}
