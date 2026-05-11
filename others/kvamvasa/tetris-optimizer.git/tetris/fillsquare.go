package tetris

//fillSquare is a recursive function that attempts to place a tile on the grid.
//Upon success, it calls itself with the next tile. If the last tile is placed
//successfully, it returns true and the filled grid for printing.
func fillSquare(currentIndex int, tetrominoes [][][]string, grid [][]string) (bool, [][]string) {

	if len(tetrominoes) == currentIndex {

		return true, grid
	}

	for squareY := range grid {
		for squareX := range grid[squareY] {

			if enoughSpace(squareY, squareX, tetrominoes[currentIndex], grid) {

				//place tile
				for y, line := range tetrominoes[currentIndex] {
					for x, char := range line {
						if char != "." {
							grid[squareY+y][squareX+x] = char
						}
					}
				}
				//try placing next tile
				solution, _ := fillSquare(currentIndex+1, tetrominoes, grid)
				if solution {
					return true, grid
				}

				//backtrack
				for y, line := range tetrominoes[currentIndex] {
					for x := range line {
						if grid[squareY+y][squareX+x] == string('A'+rune(currentIndex)) {
							grid[squareY+y][squareX+x] = "."

						}
					}
				}

			}

		}
	}
	return false, grid

}
