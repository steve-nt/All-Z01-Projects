package tetris

// createSquare takes a side dimension as a parameter, and initializes
// a 2d slice of . to be filled by the tetrominoes.
func createSquare(side int) [][]string {

	square := make([][]string, side)

	for i := range side {
		for j := 0; j < side; j++ {
			square[i] = append(square[i], ".")
		}
	}
	return square
}
