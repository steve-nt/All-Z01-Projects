package main

import (
	"fmt"
	"os"
	"strings"
)

var outputFileWritten bool // Global flag to track if output has been written

// Save the ASCII art output to the file specified in outputFlag
func saveToOutput(outputFlag string, asciiChars [][]string, asciiHeight int) {
	// Step 1: Do nothing if outputFlag is empty
	if outputFlag == "" {
		return
	}

	// Step 2: Check if the outputFlag is a valid .txt file
	if !strings.HasSuffix(outputFlag, ".txt") {
		fmt.Println("Error: Output file must have a .txt extension.")
		return
	}

	// Step 3: Determine the file mode: create/overwrite for the first time, append for subsequent calls
	var file *os.File
	var err error
	if outputFileWritten {
		// If the file has been written to, append to it
		file, err = os.OpenFile(outputFlag, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error: Could not open the file for appending: %v\n", err)
			return
		}
	} else {
		// If this is the first time, create or overwrite the file
		file, err = os.Create(outputFlag)
		if err != nil {
			fmt.Printf("Error: Could not create or open the file: %v\n", err)
			return
		}
		outputFileWritten = true // Mark that the file has been written to
	}
	defer file.Close()

	// Step 4: Convert the asciiChars (2D array) to a string
	var asciiArtContent strings.Builder
	for i := 0; i < asciiHeight; i++ {
		for j := 0; j < len(asciiChars); j++ {
			asciiArtContent.WriteString(asciiChars[j][i]) // Add each line of the character
			if j < len(asciiChars)-1 {
				asciiArtContent.WriteString(" ") // Optional space between characters
			}
		}
		asciiArtContent.WriteString("\n") // Move to the next line of the ASCII art
	}

	// Step 5: Write the content to the file
	_, err = file.WriteString(asciiArtContent.String())
	if err != nil {
		fmt.Printf("Error: Failed to write to the file: %v\n", err)
		return
	}

	// Step 6: Notify the user that the file was saved or appended successfully
	if outputFileWritten {
		fmt.Printf("Output successfully appended to %s\n", outputFlag)
	} else {
		fmt.Printf("Output successfully saved to %s\n", outputFlag)
	}
}
