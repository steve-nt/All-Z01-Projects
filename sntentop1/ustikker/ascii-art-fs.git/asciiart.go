package main

import (
	"fmt"
	"strings"
)

// ANSI color codes map
var ansiColors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

// Process the input string and print the ASCII art with color
func processString(input string, asciiMap map[rune][]string, asciiHeight int) {

	// Replace literal "\n" with actual newlines and split into lines
	input = strings.ReplaceAll(input, `\n`, "\n")
	inputLines := strings.Split(input, "\n")

	// Process each line separately
	for _, line := range inputLines {
		if line == "" {
			fmt.Println() // Handle empty lines
			continue
		}

		asciiChars := buildAsciiArt(line, asciiMap, asciiHeight)

		highlightMask := buildHighlightMask(line, text2color)
		printAsciiArt(asciiChars, highlightMask, asciiHeight, colorFlag)
	}
}

// Build the ASCII art for a given line of input
func buildAsciiArt(line string, asciiMap map[rune][]string, asciiHeight int) [][]string {
	var asciiChars [][]string
	for _, char := range line { // Use rune to support Unicode
		if art, exists := asciiMap[char]; exists {
			//Add here things you want before every character, like Justify
			// for i := 0; charno == 0 && i < len(art); i++ {
			// 	art[i] = "--> " + art[i]
			// }
			asciiChars = append(asciiChars, art)
		} else {
			// Handle characters not present in the font data
			fmt.Printf("Warning: Character '%c' not found in font data.\n", char)
			asciiChars = append(asciiChars, make([]string, asciiHeight))
		}
	}
	saveToOutput(outputFlag, asciiChars, asciiHeight) // Call after asciiChars is generated
	return asciiChars
}

// Build a mask indicating which characters should be highlighted
func buildHighlightMask(line, substring string) []bool {
	mask := make([]bool, len(line))
	idx := 0
	if substring == "" {
		return mask // No highlighting if substring is empty
	}
	for idx < len(line) {
		if strings.HasPrefix(line[idx:], substring) {
			for i := 0; i < len(substring) && idx+i < len(mask); i++ { // Ensure boundaries are respected
				mask[idx+i] = true
			}
			idx += len(substring)
		} else {
			idx++
		}
	}
	return mask
}

// Print the ASCII art for the given characters with color
func printAsciiArt(asciiChars [][]string, highlightMask []bool, asciiHeight int, color string) {
	for i := 0; i < asciiHeight; i++ {
		for j, charLines := range asciiChars {
			if highlightMask[j] {
				fmt.Print(ansiColors[color])
			}
			fmt.Print(charLines[i])
			if highlightMask[j] {
				fmt.Print(ansiColors["reset"])
			}
		}
		fmt.Println() // Move to the next line of the ASCII art
	}
}
