package utils

import (
	"strings"
)

func (m *asciiMap) FindSubstringIndexes() {

	var allIndexes []int
	startIndex := 0

	for startIndex < len(m.input) { // Ensure startIndex is within bounds
		index := strings.Index(m.input[startIndex:], m.substring)

		if index == -1 {
			break
		}

		actualIndex := startIndex + index

		for i := 0; i < len(m.substring); i++ {
			allIndexes = append(allIndexes, actualIndex+i)
		}

		startIndex = actualIndex + 1
	}

	m.substringIndexes = allIndexes
}
