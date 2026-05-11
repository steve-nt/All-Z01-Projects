package tetris

// initialValidation runs a quick initial validation of the tetrominoes before proceeding
// it checks that only allowed characters are included, that they have the expected number
// of lines of expected length,  and that the # count is valid
func initialValidation(sliceOfTileElements [][][]string) bool {

	var hashCounter int
	for _, tile := range sliceOfTileElements {
		//each tile needs to have 4 rows
		if len(tile) != 4 {
			return false
		}
		//and 4 columns
		for _, line := range tile {
			if len(line) != 4 {
				return false
			}
			for _, symbol := range line {
				if symbol != "." && symbol != "#" && symbol != "\n" { // no other characters allowed
					return false
				} else if symbol == "#" {
					hashCounter++
				}
			}
		}

		if hashCounter != 4 {
			return false
		}

		hashCounter = 0

	}
	return true
}
