package asciiart

import (
	"fmt"
	"os"
)

// gets passed argStr, lines
func PrintAsciiArt(sentences []string, textFile []string) {
	// loops through each word in sentences slice
	for i, word := range sentences {
		if word == "" {
			if i != 0 { // checks if the empty word is NOT the first one in the slice(for \n\n)
				fmt.Println() // Print a new line for blank words
			}
			continue
		}

		/*
			this means that the loop will run 8 times
			b/c indexlines 1 through 8,
			contain the actual representation of the character.
		*/
		for h := 1; h < 9; h++ {
			// loops through each character that was stored in []word{}
			for i := 0; i < len(word); i++ {
				// loops through the whole ascii file-returning two values
				for lineIndex, line := range textFile {
					if lineIndex == (int(word[i])-32)*9+h { // Map the character to ASCII art lines
						fmt.Print(line) // Print the character line for the current height
					}
				}
			}
			fmt.Println() // New line after each line of ASCII art
		}
	}
}

// gets passed argStr, lines, createdFile
// toFile is a pointer to createdFile represented by os.File, allowing the function to use its address (pointer) to write ascii art to that file directly.
func PrintAsciiArtToFile(sentences []string, textFile []string, toFile *os.File) {
	// loops through each word in sentences slice
	for i, word := range sentences {
		if word == "" {
			if i != 0 { // checks if the empty word is NOT the first one in the slice
				_, err := toFile.WriteString("\n") // writes /n to the toFile file
				if err != nil {
					return
				}
			}
			continue
		}
		// Loop over the height of each character (assuming 8 lines)
		for h := 1; h < 9; h++ {
			// loops through each character that was stored in []word{}
			for i := 0; i < len(word); i++ {
				// loops through the whole ascii file-returning two values, finds the specific line of the ASCII art for the current character
				for lineIndex, line := range textFile {
					if lineIndex == (int(word[i])-32)*9+h {
						// writes the character line for the current height into the file
						_, err := toFile.WriteString(line)
						if err != nil {
							return
						}
					}
				}
			}
			// prints a newline so it can move on to the 2nd, 3rd etc line of the ascii art form of the character
			_, err := toFile.WriteString("\n")
			if err != nil {
				return
			}
		}
	}
	// check if useless
	_, err := toFile.WriteString("\n")
	if err != nil {
		return
	}
}
