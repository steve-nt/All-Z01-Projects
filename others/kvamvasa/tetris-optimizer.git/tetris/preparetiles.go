package tetris

import (
	"log"
)

// PrepareTiles prepares the tiles for use. It replaces # with a letter
// of the alphabet, and removes unnecessary empty spaces
// It also runs checks to make sure the tiles are valid
func PrepareTiles(tetrominoes [][][]string) [][][]string {

	if !initialValidation(tetrominoes) {
		log.Fatal("Invalid file configuration")
	}

	// log # coordinates for each tile and check they form a valid tetromino
	for i, tile := range tetrominoes {
		var coordinates [][]int

		for y, line := range tile {
			for x, symbol := range line {
				if symbol == "#" {

					coordinates = append(coordinates, []int{x, y})
				}
			}
		}
		if !checkValidCoordinates(coordinates) {
			log.Fatal("ERROR:Invalid tile configuration")
		}

		//move tile up and to the left as far as it goes
		tetrominoes[i] = reshapeTile(tile, coordinates)

	}

	tetrominoes = sortAreas(tetrominoes)

	//assign a letter of the alphabet to each tetromino
	for i, tile := range tetrominoes {

		letter := string('A' + rune(i))
		for _, line := range tile {
			for x, symbol := range line {
				if symbol == "#" {
					line[x] = letter

				}
			}
		}

	}
	return tetrominoes
}
