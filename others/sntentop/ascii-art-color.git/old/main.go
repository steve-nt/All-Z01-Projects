package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse and validate command-line arguments
	textFromOutside, banner := HandleArguments()

	// Load ASCII banner template
	asciiTemplates := LoadASCIIBanner(banner)

	// If the ASCII templates are nil, exit due to file reading error
	if asciiTemplates == nil {
		fmt.Println("Error loading banner file. Make sure the banner name is correct.")
		return
	}

	// Process the input string and print the ASCII representation
	ProcessAndPrintASCII(textFromOutside, asciiTemplates)
}

// HandleArguments parses the command-line arguments and returns the input text and banner name.
func HandleArguments() (string, string) {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		fmt.Println("EX: go run . something standard")
		os.Exit(1)
	}

	// Get the input string
	textFromOutside := os.Args[1]

	// Get the banner type, default to "standard" if not provided
	banner := "standard"
	if len(os.Args) == 3 {
		banner = os.Args[2]
	}

	return textFromOutside, banner
}
