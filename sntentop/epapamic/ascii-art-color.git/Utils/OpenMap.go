package utils

import (
	"fmt"
	"os"
)

// Core struct, holds all information relative to the operation.
type asciiMap struct {
	ref              *os.File
	input            string
	content          map[rune][]string
	printContent     string
	color            string
	substring        string
	substringIndexes []int
}

func CreateMap() *asciiMap {
	return &asciiMap{}

}

// Gives a reference to the internally designated file,
// returns non-nil error if unable to open the file.
func (m *asciiMap) OpenMap() error {
	return m.openMap(func(input string) (*os.File, error) {
		return os.Open(mapPath)
	})
}

// Any arguments passed are strickly for testing purposes.
func (m *asciiMap) openMap(opener func(string) (*os.File, error)) error {
	osFile, err := opener(mapPath)
	if err != nil {
		return fmt.Errorf("failed to open the file at: %s", mapPath)
	}
	m.ref = osFile
	return nil
}
