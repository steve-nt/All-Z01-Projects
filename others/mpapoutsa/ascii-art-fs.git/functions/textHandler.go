package functions

import (
	"fmt"
	"os"
	"strings"
)

func PrintArt(str, banner, substring, color string, colorhandler bool) {
	// Open the banner file
	file, err := os.Open(banner)
	if err != nil {
		fmt.Println("Error Opening banner provided: ", err)
		return
	}
	defer file.Close()

	// Handle special characters like \n
	str = handleSpecialChars(str)

	// Split the string by newline
	linesOfText := strings.Split(str, "\\n")

	// For each line in the string
	for _, lineText := range linesOfText {
		// If the line is empty, print an empty line
		if lineText == "" {
			fmt.Println()
			continue
		}

		// Call the appropriate function based on colorhandler flag
		if colorhandler {
			// Call AsciiArtColor if colorhandler is true
			AsciiArtColor(lineText, file, substring, color)
		} else {
			// Otherwise, call PrintAsciiArt for normal ASCII art printing
			PrintAsciiArt(lineText, file)
		}
	}
}

// handleSpecialChars processes backslash-escaped characters
func handleSpecialChars(text string) string {
	// Replace backslash-escaped characters with their actual characters
	replacements := map[string]string{
		`\!`: "!",
		`\'`: "'",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return text
}
