package utils

import (
	"strings"
)

// Compares input and extracts corresponding ascii-art values from the map.
func asciiOutput(input string, m *asciiMap) []string {
	var output []string  // Slice holding the completed lines
	var line string      // Carrier for building a line
	if len(input) == 0 { // Check for empty input
		return output
	}
	for i := 0; i < 8; i++ { // Loop trought all the map values
		for index, char := range input { // For every character given
			if Exists(index, m.substringIndexes) {
				line = line + m.color + m.content[char][i] + "\033[37m"
			} else {
				line = line + m.content[char][i] // Extract output from map , build a string holding each characters corresponding line
			}
		}
		output = append(output, line) // Append each completed line to the final slice
		line = ""                     // Reset line so we can repeat
	}
	return output
}

// Converts the map into a formated string.
// Any arguments passed are strickly for testing purposes.
func (m *asciiMap) formatAsciiArt(fun func(string, *asciiMap) []string) {
	inputSlice := strings.Split(m.input, "\\n") // Seperate input at literal "\n"
	var result string
	var outputArt string
	for _, word := range inputSlice {
		output := fun(word, m)                 // Extract the map values of the word
		outputArt = strings.Join(output, "\n") // Format the map values so that can be displayed vertically
		result += outputArt + "\n"
	}
	m.printContent = result
}

// Creates a formated string and stores it inside the asciiMap.
func (m *asciiMap) FormatAsciiArt() {
	m.formatAsciiArt(func(input string, m *asciiMap) []string {
		return asciiOutput(input, m)
	})
}
