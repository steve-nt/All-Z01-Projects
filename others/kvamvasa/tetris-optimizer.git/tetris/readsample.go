package tetris

import (
	"log"
	"os"
	"strings"
)

// ReadSample reads the input file and stores the tiles in a 3d slice (tiles>columns>rows)
func ReadSample(file string) [][][]string {
	f, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// make string out of slice of bytes
	fileString := string(f)

	// normalize linebreaks
	fileString = strings.ReplaceAll(fileString, "\r\n", "\n")

	//remove last linebreak if it exists
	fileString = strings.TrimRight(fileString, "\n")

	// break into tiles (1D slice)
	sliceOfTiles := strings.Split(fileString, "\n\n")

	// break into lines (2D slice)
	var sliceOfTileLines [][]string
	for _, s := range sliceOfTiles {
		sliceOfTileLines = append(sliceOfTileLines, strings.Split(s, "\n"))
	}

	// break into elements (3D slice)
	var sliceOfTileElements [][][]string
	for _, tile := range sliceOfTileLines {
		var tiles [][]string
		for _, line := range tile {
			elements := strings.Split(line, "")
			tiles = append(tiles, elements)
		}
		sliceOfTileElements = append(sliceOfTileElements, tiles)
	}

	return sliceOfTileElements
}
