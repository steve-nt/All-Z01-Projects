package main

import (
	"os" // Import the os package to handle command-line arguments

	"github.com/01-edu/z01" // Import the z01 package to use the PrintRune function for printing characters
)

func main() {
	args := os.Args[1:] // Get the command-line arguments, excluding the program name
	if len(args) == 0 { // If there are no arguments, exit the function
		return
	}

	letters := "abcdefghijklmnopqrstuvwxyz" // Define a string containing lowercase letters
	if args[0] == "--upper" {               // Check if the first argument is "--upper"
		letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // If true, use uppercase letters instead
		args = args[1:]                        // Remove the "--upper" argument from the list
	}

	for _, arg := range args { // Iterate over the remaining arguments
		num := 0                // Initialize the number to zero
		for _, c := range arg { // Iterate over each character in the argument
			if c < '0' || c > '9' { // Check if the character is not a digit
				num = -1 // If it's not a digit, set num to -1 and break the loop
				break
			}
			num = num*10 + int(c-'0') // Convert the character to a number and add it to num
		}
		if num >= 1 && num <= 26 { // Check if the number is between 1 and 26
			z01.PrintRune(rune(letters[num-1])) // Print the corresponding letter
		} else {
			z01.PrintRune(' ') // If the number is not valid, print a space
		}
	}
	z01.PrintRune('\n') // Print a newline character at the end
}
