package ascii

import (
	"fmt"
	"os"
	"strings"
)

// Loads the banner file and returns its lines and characters per ASCII row
func loadBanner(bannerName string) ([]string, int, error) {
	content, err := os.ReadFile(bannerName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load banner file: %w", err)
	}

	bannerLines := strings.Split(string(content), "\n")
	const characterLines = 9
	if len(bannerLines) < characterLines {
		return nil, 0, fmt.Errorf("invalid banner file format")
	}
	return bannerLines, characterLines, nil
}

func GenerateAsciiArt(input, bannerName string) (string, error) {
	// Normalize the input by removing \r (carriage returns)
	normalizedInput := strings.ReplaceAll(input, "\r", "")

	// Load banner and lines per character
	banner, linesPerChar, err := loadBanner(bannerName)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	words := strings.Split(normalizedInput, "\n") // Split into lines

	for _, word := range words {
		if word == "" { // Handle empty lines
			result.WriteString("\n")
			continue
		}

		lines := make([]strings.Builder, linesPerChar)

		for _, char := range word {
			// Validate and fetch ASCII art index
			startIndex, err := getCharIndex(char, linesPerChar)
			if err != nil {
				return "", fmt.Errorf("unsupported character: %c", char)
			}

			for i := 0; i < linesPerChar; i++ {
				lines[i].WriteString(banner[startIndex+i])
			}
		}

		// Append the generated ASCII lines for this word
		for i := 0; i < linesPerChar; i++ {
			result.WriteString(lines[i].String())
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

// getCharIndex calculates the starting index of a character in the banner
func getCharIndex(char rune, linesPerChar int) (int, error) {
	if char == ' ' {
		return 0, nil // Handle space
	}
	if char == '\n' {
		return -1, nil // Allow newline, but return a marker to skip
	}
	if char >= 32 && char <= 126 {
		return (int(char) - 32) * linesPerChar, nil
	}
	return -1, fmt.Errorf("unsupported character: %c", char)

}
