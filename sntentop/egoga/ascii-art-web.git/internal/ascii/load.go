package ascii

import (
	"bufio"
	"os"
)

const (
	spaceASCII = 32
	tidleASCII = 126
	charHeight = 8
)

// Load ASCII templates from the banner file
func loadBanner(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var char rune
	var lines []string
	charCode := spaceASCII // ASCII code for the first character (space)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			// End of a character section
			if charCode <= tidleASCII && len(lines) == charHeight { // ASCII characters range from 32 (space) to 126 (~)
				char = rune(charCode)
				asciiArt[char] = append([]string{}, lines...) // Copy lines to avoid reference issues
				lines = []string{}
				charCode++ // Move to next ASCII character
			}
		} else {
			lines = append(lines, line)
		}
	}

	// Append the last character if it wasn't followed by a newline
	if charCode <= tidleASCII && len(lines) == charHeight {
		char = rune(charCode)
		asciiArt[char] = append([]string{}, lines...)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
