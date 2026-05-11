package tetris

// checkValidCoordinates get the coordinates of the four # in a tile and returns true if they form a valid shape
func checkValidCoordinates(coordinates [][]int) bool {
	// map of connections per #
	matched := make(map[int]int)

	for i := 0; i < len(coordinates); i++ {
		for j := 0; j < len(coordinates); j++ {
			if coordinates[i][0] == coordinates[j][0] && (coordinates[j][1]-coordinates[i][1] == 1 || coordinates[j][1]-coordinates[i][1] == -1) {
				matched[i]++
			} else if coordinates[i][1] == coordinates[j][1] && (coordinates[j][0]-coordinates[i][0] == 1 || coordinates[j][0]-coordinates[i][0] == -1) {
				matched[i]++
			}
		}
	}

	// at least two # need to be connected to two # each for the shape to be continuous
	var count int
	var continuous bool
	for _, value := range matched {
		if value == 3 {
			continuous = true
			break
		}
		if value == 2 {
			count++
			// If we find at least two, we can break early
			if count >= 2 {
				continuous = true
				break
			}

		}
	}
	// all # need to have at least one connection, and at least 2 # need to have at least 2.
	if len(matched) == 4 && continuous {
		return true
	}

	return false
}
