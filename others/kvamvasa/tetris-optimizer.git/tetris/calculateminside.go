package tetris

import "math"

//CalculateMinSide calculates the minimum side of a square than can fit all tetrominoes
//in the given slice, based on their number and dimensions
func CalculateMinSide(tetrominoes [][][]string) int {
	var squareSide int
	tilesNum := len(tetrominoes)
	// calculate minimum square based on number of #
	approxSide := math.Sqrt(float64(tilesNum * 4))

	//if result has decimal numbers, get the next integer
	if approxSide > math.Floor(approxSide) {
		squareSide = int(approxSide) + 1
	} else {
		squareSide = int(approxSide)
	}

	//find the maximum height and width of the tetrominoes
	//and make sure the side is equal or larger
	var maxHeight int
	var maxWidth int
	for _, tile := range tetrominoes {
		for _, line := range tile {
			if len(tile) > maxHeight {
				maxHeight = len(tile)
			}
			if len(line) > maxWidth {
				maxWidth = len(line)
			}
		}
	}

	if maxHeight > squareSide {
		squareSide = maxHeight
	}
	if maxWidth > squareSide {
		squareSide = maxWidth
	}

	return squareSide
}
