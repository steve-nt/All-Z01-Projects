package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	// Start the loop from index 1 to skip the first argument (the program name)
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]    // Get the current argument
		runes := []rune(arg) // Convert the argument to a slice of runes

		// Iterate through each rune in the current argument
		for j := 0; j < len(runes); j++ {
			z01.PrintRune(runes[j]) // Print the current rune
		}
		z01.PrintRune('\n') // Print a newline after the argument
	}
}
