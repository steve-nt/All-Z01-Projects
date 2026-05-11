package main

import (
	utils "ASCII-ART/Utils"
	"flag"
	"log"
)

func main() {

	// Initialize Map
	asciiMap := utils.CreateMap()

	// Initialize Flag
	flagColor := flag.String("color", "none", "Color the output")

	flag.Parse()

	color := *flagColor

	// Parse user input
	err := asciiMap.ParseInput(color)
	if err != nil {
		log.Fatal(err)
	}

	// Load ANSI Color
	err = asciiMap.GetAnsiColor(color)
	if err != nil {
		log.Fatal(err)
	}

	// Get substring Indexes
	asciiMap.FindSubstringIndexes()

	// // Check for malformed data file
	err = utils.CheckFileIntegrity()
	if err != nil {
		log.Fatal("Data/standard.txt malformed")
	}

	// Get reference to the map
	err = asciiMap.OpenMap()
	if err != nil {
		log.Fatal("Failed to open file")
	}

	// Validate that the input is not empty and every character is printable and within Ascii
	err = asciiMap.ValidatePrintable()
	if err != nil {
		log.Fatal("Contains non-Printable characters")
	}

	// Create the map
	asciiMap.CreateMap()

	// Replace each character provided by the user to the corresponding ascii art representation
	asciiMap.FormatAsciiArt()

	// Print the map
	asciiMap.PrintAsciiMapCharacters()
}
