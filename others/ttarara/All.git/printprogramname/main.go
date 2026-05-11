package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	// Get the program name (the base name of the executable)
	programName := os.Args[0]

	// Find the last '/' or '\' to get the base name
	for i := len(programName) - 1; i >= 0; i-- {
		if programName[i] == '/' || programName[i] == '\\' {
			programName = programName[i+1:]
			break
		}
	}

	// Print each character of the program name
	for _, r := range programName {
		z01.PrintRune(r)
	}
	z01.PrintRune('\n')
}
