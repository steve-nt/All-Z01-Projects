package functions

import (
	"bufio"
	"fmt"
	"os"
)

// AsciiArtColor applies colors to the specified substring or the entire string if no substring is provided.
func AsciiArtColor(str string, file *os.File, substring, color string) {
	// ANSI color codes
	colors := map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"reset":   "\033[0m",
	}

	// Get the color code, return error if invalid
	colorCode, exists := colors[color]
	if !exists {
		fmt.Println("Error: Invalid color specified.")
		return
	}

	// Process each character in the string to print the ASCII art with color
	byteOfText := []byte(str)
	lines := make([][]string, len(byteOfText))

	// For each character in the string, get its ASCII art representation
	for idx, val := range byteOfText {
		// Locate the ASCII art block for the character
		nbrVal := int(val-32) * 9
		lineArray := []string{}
		file.Seek(0, 0) // Reset file position
		scanner := bufio.NewScanner(file)

		// Collect the ASCII representation of the character
		for lineNum := 0; scanner.Scan(); lineNum++ {
			if lineNum >= nbrVal && lineNum < nbrVal+9 {
				lineArray = append(lineArray, scanner.Text())
			}
		}
		lines[idx] = lineArray

		// Check for scanning errors
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
	}

	// Build the ASCII art row by row
	for lineNbr := 0; lineNbr < 9; lineNbr++ {
		for idx, charLines := range lines {
			if lineNbr < len(charLines) {
				asciiLine := charLines[lineNbr]

				// If no substring is specified, color the entire line
				if substring == "" {
					asciiLine = colorCode + asciiLine + colors["reset"]
				} else {
					// Otherwise, color only the parts matching the substring
					inMatchedSubstring := false
					for i := 0; i <= len(str)-len(substring); i++ {
						if str[i:i+len(substring)] == substring {
							// Mark the part of the line that matches the substring
							if idx >= i && idx < i+len(substring) {
								inMatchedSubstring = true
								break
							}
						}
					}

					// Apply color if part of the substring
					if inMatchedSubstring {
						asciiLine = colorCode + asciiLine + colors["reset"]
					}
				}

				// Print the ASCII art line with color applied
				fmt.Print(asciiLine)
			}
		}
		fmt.Println() // Newline after each ASCII row
	}
}
