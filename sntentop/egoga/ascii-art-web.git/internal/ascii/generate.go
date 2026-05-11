package ascii

import (
	"fmt"
	"strings"
)

// Map to store ASCII representations of characters
var asciiArt = make(map[rune][]string)

func GenerateTextToAscii(text string, bannerPath string) (string, error) {
	// Load banner file (assuming file name 'standard.txt')
	err := loadBanner(bannerPath)
	if err != nil {
		return "", fmt.Errorf("error loading banner: %v", err)
	}

	//Normalize input to handle different escape sequence interpretations
	text = strings.ReplaceAll(text, "\\n", "\n")
	lines := strings.Split(text, "\n")

	var output strings.Builder

	for _, line := range lines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		// Build each of the 8 ASCII lines
		asciiLines := make([]string, charHeight)
		for _, char := range line {
			art, exists := asciiArt[char]
			if !exists {
				return "", fmt.Errorf("character %q not supported in banner", char)
			}
			for i := 0; i < charHeight; i++ {
				asciiLines[i] += art[i]
			}
		}

		// Append the ASCII lines to the output
		for _, l := range asciiLines {
			output.WriteString(l + "\n")
		}
	}
	return output.String(), nil
}
