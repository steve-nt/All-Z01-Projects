package tetris

//reshapeTile trims unnecessary spaces from each tile
func reshapeTile(tile [][]string, coordinates [][]int) [][]string {

	//initialize min and max x,y as first #'s coordinates
	minX := coordinates[0][0]
	maxX := coordinates[0][0]
	minY := coordinates[0][1]
	maxY := coordinates[0][1]
	//find the actual values for minX, maxX, minY, maxY
	for _, hashCoordinates := range coordinates {
		if hashCoordinates[0] < minX {
			minX = hashCoordinates[0]
		} else if hashCoordinates[0] > maxX {
			maxX = hashCoordinates[0]
		}
		if hashCoordinates[1] < minY {
			minY = hashCoordinates[1]
		} else if hashCoordinates[1] > maxY {
			maxY = hashCoordinates[1]
		}
	}
	//trim rows to remove unnecessary empty spaces
	if maxY < 3 { //make sure we don't go out of range
		tile = tile[minY : maxY+1]
	} else {
		tile = tile[minY:]
	}
	//trim columns to remove unnecessary empty spaces
	for i := range tile {
		if maxX < 3 { //make sure we don't go out of range
			tile[i] = tile[i][minX : maxX+1]
		} else {
			tile[i] = tile[i][minX:]
		}
	}

	return tile
}
