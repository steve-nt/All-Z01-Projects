package main

import (
	"log"
	"os"

	"optimizer/tetris"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No source file given")
	}

	file := os.Args[1]

	//read the input file and store in a 3d slice
	tetrominoes := tetris.ReadSample(file)
	//validate and format the input to get ready for use
	tetrominoes = tetris.PrepareTiles(tetrominoes)
	//calculate the minimum side of the grid
	minSide := tetris.CalculateMinSide(tetrominoes)
	//solve the grid
	tetris.Solve(minSide, tetrominoes)
}
