package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	// Start the loop from the last argument to the first argument
	for i := len(os.Args) - 1; i > 0; i-- {
		arg := os.Args[i]    // Get the current argument
		runes := []rune(arg) // Convert the argument to a slice of runes

		// Iterate through each rune in the current argument
		for j := 0; j < len(runes); j++ {
			z01.PrintRune(runes[j]) // Print the current rune
		}
		z01.PrintRune('\n') // Print a newline after the argument
	}
}
