package asciiart

import (
	"fmt"
	"strings"
)

// PrintAsciiArtAlign prints ASCII art for a slice of sentences,
// aligning it according to the specified position and width.
func PrintAsciiArtAlign(sentences []string, textFile []string, position string, w int) {
	// Loops through each sentence in the sentences slice
	for _, sentence := range sentences {
		if sentence == "" {
			fmt.Println() // Always print a newline for empty sentences
			continue
		}

		// Split the sentence into words
		words := splitIntoWords(sentence)
		wordLens := make([]int, len(words))
		totalWordLen := 0

		// Calculate lengths of each word's ASCII representation
		for i, word := range words {
			if word == "" {
				continue
			}
			wordLen := 0
			for j := 0; j < len(word); j++ {
				for lineIndex, line := range textFile {
					if lineIndex == (int(word[j])-32)*9+2 {
						wordLen += len(line)
						break
					}
				}
			}
			wordLens[i] = wordLen
			totalWordLen += wordLen
		}

		// Total space available for justification
		spaces := w - totalWordLen

		// Number of gaps between words
		numGaps := len(words) - 1

		// Calculate extra spaces to distribute
		var extraSpaces int
		if numGaps > 0 {
			extraSpaces = spaces / numGaps
		}

		// Loop for each height of ASCII representation
		for h := 1; h < 9; h++ {
			if position == "justify" {
				for i, word := range words {
					// Print the ASCII representation of the current word
					for j := 0; j < len(word); j++ {
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[j])-32)*9+h {
								fmt.Print(line)
								break
							}
						}
					}
					// Print spaces between words if not the last word
					if i < numGaps {
						for s := 0; s < extraSpaces; s++ {
							fmt.Print(" ")
						}
					}
				}
			} else if position == "center" {
				// Print spaces to center-align the text (before the words)
				for i := 0; i < spaces/2; i++ {
					fmt.Print(" ")
				}
				// Print each word with a space afterwards
				for i, word := range words {
					for j := 0; j < len(word); j++ {
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[j])-32)*9+h {
								fmt.Print(line)
								break
							}
						}
					}
					if i < len(words)-1 {
						fmt.Print(" ") // Add space between words
					}
				}
				// Print spaces to center-align the text (after the words)
				for i := 0; i < spaces/2-len(words)-1; i++ {
					fmt.Print(" ")
				}
			} else if position == "right" {
				// Print spaces to right-align the text (before the words)
				for i := 0; i < spaces-len(words)-1; i++ {
					fmt.Print(" ")
				}
				// Print each word with a space afterwards
				for i, word := range words {
					for j := 0; j < len(word); j++ {
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[j])-32)*9+h {
								fmt.Print(line)
								break
							}
						}
					}
					if i < len(words)-1 {
						fmt.Print(" ") // Add space between words
					}
				}
			} else if position == "left" {
				// Print each word with a space afterwards
				for i, word := range words {
					for j := 0; j < len(word); j++ {
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[j])-32)*9+h {
								fmt.Print(line)
								break
							}
						}
					}
					if i < len(words)-1 {
						fmt.Print(" ") // Add space between words
					}
				}
			}
			fmt.Println() // Move to the next line after finishing the current height
		}
	}
}

// splitIntoWords splits a sentence into words based on spaces.
func splitIntoWords(sentence string) []string {
	return strings.Fields(sentence)
}
