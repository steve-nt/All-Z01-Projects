package tetris

import (
	"fmt"
)

// Solve initiates the creation of the grid with the given side, and begins the process to find a solution.
// If no solution is found it calls itself recursively and retries with a larger grid.
// Once a solution is found it prints it along with the number of empty spaces.
func Solve(side int, tetrominoes [][][]string) {
	grid := createSquare(side)
	solution, filledGrid := fillSquare(0, tetrominoes, grid)
	if solution {

		//count .
		var counter int
		for _, line := range filledGrid {
			for _, char := range line {
				if char == "." {
					counter++
				}
				fmt.Print(char + " ")

			}
			fmt.Println()
		}
		fmt.Println("Empty spaces: ", counter)

		return

	} else {
		Solve(side+1, tetrominoes)
	}

}
