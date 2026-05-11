package utils

import (
	"bufio"
)

// Populates the map, stores it inside the asciiMap.
func (m *asciiMap) CreateMap() {
	scanner := bufio.NewScanner(m.ref)
	startOfCharacter := false // Flag to indicate the NEXT line is the start of a character
	m.content = make(map[rune][]string)
	char := []string{} // Holds multiple lines representing the character
	sliceNum := 0      // Number of currently processing line
	key := ' '         // Indicates the first writable character in ascii

	defer m.ref.Close()

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 && !startOfCharacter { // Found start of character text
			startOfCharacter = true // Set switch on
			continue
		} else {
			char = append(char, line) // Fill the slice with subsequent lines
			sliceNum++                // Mark how many lines we appended
		}
		if sliceNum == 8 { // Stop when 8 lines appended
			m.content[key] = char    // Push it to the map
			key++                    // Change to the next ascii character
			startOfCharacter = false // Set switch off
			sliceNum = 0             // Reset the counter of lines we already have
			char = []string{}        // Clear the text we have appended
		}
	}
}
