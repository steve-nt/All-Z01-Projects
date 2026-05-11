package main

import (
	"fmt"
	"os"
	"strings" // Required for strings.HasPrefix
)

func main() {
	// Parse and validate command-line arguments
	textFromOutside, banner, colorCode, substringToColor := HandleArguments()

	// Load ASCII banner template
	asciiTemplates := LoadASCIIBanner(banner)
	if asciiTemplates == nil {
		fmt.Println("Error loading banner file. Make sure the banner name is correct.")
		return
	}

	// Process and print ASCII text with optional color
	if colorCode != "" && substringToColor != "" {
		coloredText := ColorSubstring(textFromOutside, substringToColor, colorCode)
		ProcessAndPrintASCII(coloredText, asciiTemplates)
	} else {
		ProcessAndPrintASCII(textFromOutside, asciiTemplates)
	}
}

// HandleArguments parses command-line arguments and returns the input text, banner type, color code, and substring to color.
func HandleArguments() (string, string, string, string) {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		PrintUsage()
		os.Exit(1)
	}

	// Set defaults
	textFromOutside := ""
	banner := "standard"
	colorCode, substringToColor := "", ""

	// Parse arguments based on count
	if len(os.Args) == 2 {
		// Only text is provided
		textFromOutside = os.Args[1]
	} else if len(os.Args) == 3 {
		// Text and banner provided
		textFromOutside = os.Args[1]
		banner = os.Args[2]
	} else if len(os.Args) == 4 && strings.HasPrefix(os.Args[1], "--color=") {
		// Color flag, substring, and text provided
		var valid bool
		colorCode, valid = ParseColorFlag(os.Args[1])
		if !valid {
			os.Exit(1)
		}
		substringToColor = os.Args[2]
		textFromOutside = os.Args[3]
	}

	return textFromOutside, banner, colorCode, substringToColor
}
