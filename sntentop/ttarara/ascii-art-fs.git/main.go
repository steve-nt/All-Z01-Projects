package main

import (
	"ascii-art-fs/extra_functions" // Import the extra_functions package
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		printUsage()
		return
	}

	// Get input string and banner type from command-line arguments
	textFromOutside := os.Args[1]
	banner := "standard"
	if len(os.Args) == 3 {
		banner = os.Args[2]
	}

	// Call ReadTxt from extra_functions with the banner argument
	fileLines := extra_functions.ReadTxt(banner)
	if fileLines == nil {
		fmt.Println("Error loading banner file. Make sure the banner name is correct.")
		return
	}

	// Use the extra_functions package functions to process and print ASCII art
	asciiTemplates := extra_functions.Return2DASCIIArray(fileLines)
	extra_functions.PrintAllStringASCII(textFromOutside, asciiTemplates)
}

func printUsage() {
	fmt.Println("Usage: go run . [STRING] [BANNER]")
	fmt.Println("EX: go run . something standard")
}
