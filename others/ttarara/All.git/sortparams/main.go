package main

import (
	"os"

	"github.com/01-edu/z01"
)

func main() {
	// Get the command-line arguments (including the program name at index 0)
	args := os.Args
	n := len(args)

	// Bubble sort the arguments in ASCII order, starting from index 1 to skip the program name
	for i := 1; i < n-1; i++ {
		for j := 1; j < n-i; j++ {
			// Compare adjacent elements using their ASCII values
			if args[j] > args[j+1] {
				// Swap args[j] and args[j+1] if they are in the wrong order
				args[j], args[j+1] = args[j+1], args[j]
			}
		}
	}

	// Print each sorted argument on a new line using z01.PrintRune
	for i := 1; i < n; i++ {
		// Iterate through each character in the argument string
		for _, r := range args[i] {
			z01.PrintRune(r) // Print each character
		}
		z01.PrintRune('\n') // Print a newline after each argument
	}
}
