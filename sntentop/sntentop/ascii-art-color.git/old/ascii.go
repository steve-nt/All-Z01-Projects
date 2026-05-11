package main

import (
	"fmt"
	"strings"
)

// ProcessAndPrintASCII splits the input text by custom newline, then prints the ASCII art.
func ProcessAndPrintASCII(text string, asciiTemplates [][]string) {
	substrings := ReturnStringToEndlineArray(text)

	// Iterate over each substring and print its ASCII representation
	for _, substring := range substrings {
		if substring == "\\n" {
			fmt.Println("") // Handle custom newline
		} else {
			PrintMultipleCharacter(substring, asciiTemplates)
		}
	}
}

// PrintMultipleCharacter prints ASCII art for multiple characters in a string.
func PrintMultipleCharacter(s string, asciiTemplates [][]string) {
	tempIntArrLetter := ReturnAsciiCodeInt(s)
	for i := 0; i < 8; i++ {
		for _, v := range tempIntArrLetter {
			fmt.Print(asciiTemplates[v][i])
		}
		fmt.Println()
	}
}

// ReturnAsciiCodeInt converts each character to its ASCII template index.
func ReturnAsciiCodeInt(s string) []int {
	var tempIntArrLetter []int
	for _, v := range s {
		tempIntArrLetter = append(tempIntArrLetter, int(v)-32)
	}
	return tempIntArrLetter
}

// ReturnStringToEndlineArray splits input text into substrings based on the custom newline delimiter.
func ReturnStringToEndlineArray(text string) []string {
	var substrings []string
	escapedN := "\\n"
	for {
		idx := strings.Index(text, escapedN)
		if idx == -1 {
			substrings = append(substrings, text)
			break
		}
		substrings = append(substrings, text[:idx])
		substrings = append(substrings, escapedN)
		text = text[idx+len(escapedN):]
	}

	// Filter out any empty strings in the final slice
	var cleanedSubstrings []string
	for _, sub := range substrings {
		if sub != "" {
			cleanedSubstrings = append(cleanedSubstrings, sub)
		}
	}
	return cleanedSubstrings
}
