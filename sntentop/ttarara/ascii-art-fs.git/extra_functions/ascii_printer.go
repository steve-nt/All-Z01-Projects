package extra_functions

import (
	"fmt"
	"strings"
)

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

// PrintAllStringASCII handles input string processing and ASCII art printing, including custom newlines.
func PrintAllStringASCII(text string, asciiTemplates [][]string) {
	substrings := ReturnString2EndlineArray(text)
	lenOfsubstrings := len(substrings)
	for index, v := range substrings {
		if v == "\\n" {
			if index == lenOfsubstrings-1 || substrings[index-1] == "\\n" {
				fmt.Println("")
			}
		} else {
			PrintMultipleCharacter(v, asciiTemplates)
		}
	}
}

// ReturnString2EndlineArray splits input text into substrings based on the custom newline delimiter.
func ReturnString2EndlineArray(text string) []string {
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

	var cleanedSubstrings []string
	for _, sub := range substrings {
		if sub != "" {
			cleanedSubstrings = append(cleanedSubstrings, sub)
		}
	}
	return cleanedSubstrings
}
