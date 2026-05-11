package tetris

//enoughSpace checks that a tile can be placed at the given coordinates in the grid
//without going out of range or overwriting another tile.
func enoughSpace(row, col int, tile, grid [][]string) bool {
	//make sure placing the tile won't go out of range
	squareSide := len(grid)
	if len(tile) > squareSide-row {
		return false
	} else if len(tile[0]) > squareSide-col {
		return false
	}

	//make sure all active parts of the tile will fall on empty spaces
	for y, line := range tile {
		for x := range line {
			if grid[row+y][col+x] != "." && tile[y][x] != "." {
				return false
			}
		}
	}
	return true
}
