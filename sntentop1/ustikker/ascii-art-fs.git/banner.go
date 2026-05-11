package main

import (
	"log"
	"os"
	"strings"
)

const asciiOffset = 32

// Load the ASCII art banner into a map
func loadBanner(filename string) (map[rune][]string, int) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}
	strData := strings.ReplaceAll(string(data), "\r\n", "\n")
	// Check if the first character in the data is a newline and skip it if necessary
	if len(strData) > 0 && strData[0] == '\n' {
		strData = strData[1:]
	}
	asciiLines := strings.Split(strData, "\n\n")
	asciiMap := make(map[rune][]string)
	// Create a map for ASCII characters (starting from the space character)
	for i, art := range asciiLines {
		char := rune(asciiOffset + i) // ASCII code starts at 32 for space
		asciiMap[char] = strings.Split(art, "\n")
	}
	// Determine asciiHeight dynamically
	var asciiHeight int
	for _, art := range asciiMap {
		asciiHeight = len(art)
		break // Get the height from the first character
	}
	return asciiMap, asciiHeight
}
