package asciiart

import "fmt"

// gets passed argStr, lines,  align, width
func PrintAsciiArtAlign(sentences []string, textFile []string, position string, w int) {
	// loops through each word in sentences slice
	for _, word := range sentences {
		if word == "" {
			fmt.Println() // Always print a newline for empty words
			continue
		}
		// stores the *number* of characters needed for the ASCII representation of the word.
		wordLen := 0
		// loops through each character that was stored in []word{}
		for i := 0; i < len(word); i++ {
			// loops through the whole ascii file-returning two values
			for lineIndex, line := range textFile {
				if lineIndex == (int(word[i])-32)*9+2 {
					wordLen += len(line)
					break
				}
			}
		}
		// Total space available for later use in justify alignment
		spaces := w - wordLen

		/*
			this means that the loop will run 8 times
			b/c indexlines 1 through 8,
			contain the actual representation of the character.
		*/
		for h := 1; h < 9; h++ {
			// If alignment is "justify", calculate spaces to spread out
			if position == "justify" {
				// if characters in sentences is more than one(1)
				if len(word) > 1 {
					// Calculate spaces to insert between characters(c^a^t)
					extraSpaces := spaces / (len(word) - 1)
					// IN CASE there is an uneven ammount of terminal space/ascii space, rarely happens.
					remainingSpaces := spaces % (len(word) - 1)

					// Loop through each character in word
					for i := 0; i < len(word); i++ {
						// loops through the whole ascii file- finds and print the specific line of the ASCII art for the current character
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[i])-32)*9+h {
								fmt.Print(line)
								break
							}
						}
						// Print extra spaces between characters
						if i < len(word)-1 {
							// Insert `extraSpaces` number of spaces
							for s := 0; s < extraSpaces; s++ {
								fmt.Print(" ")
							}
							// Distribute remaining spaces one by one
							if remainingSpaces > 0 {
								fmt.Print(" ")
								remainingSpaces--
							}
						}
					}
				} else {
					// If there's only one character, no need to distribute spaces

					// loops through that one character that was stored in []word{}
					for i := 0; i < len(word); i++ {
						// loops through the whole ascii file
						for lineIndex, line := range textFile {
							if lineIndex == (int(word[i])-32)*9+h {
								fmt.Print(line)
								// prints the content of the coresponding index line
								break
							}
						}
					}
				}
			} else if position == "center" {
				// Print spaces to center-align the text (before the word)
				for i := 1; i <= spaces/2; i++ {
					fmt.Print(" ")
				}
				// Loop through each character in the word
				for i := 0; i < len(word); i++ {
					// Loops through the ascii file and finds the correct line for the current character
					for lineIndex, line := range textFile {
						if lineIndex == (int(word[i])-32)*9+h {
							fmt.Print(line)
							break
						}
					}
				}
				// Print spaces to center-align the text (after the word)
				for i := 1; i <= spaces/2; i++ {
					fmt.Print(" ")
				}
			} else if position == "right" {
				// Print spaces to right-align the text (before the word)
				for i := 1; i <= spaces; i++ {
					fmt.Print(" ")
				}
				for i := 0; i < len(word); i++ {
					for lineIndex, line := range textFile {
						if lineIndex == (int(word[i])-32)*9+h {
							fmt.Print(line)
							break
						}
					}
				}
			} else if position == "left" {
				for i := 0; i < len(word); i++ {
					for lineIndex, line := range textFile {
						if lineIndex == (int(word[i])-32)*9+h {
							fmt.Print(line)
							break
						}
					}
				}
			}
			fmt.Println() // Move to the next line after finishing the current height
		}
	}
}
