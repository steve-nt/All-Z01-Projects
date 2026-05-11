package ascii

import (
	"fmt"
	"os"
	"unicode"
)

// CheckValidity validates input and files
func CheckValidity(inputText, banner, colorFlag string) bool {
	if IsNotASCII(inputText) {
		fmt.Println("Error: Only ASCII characters or newline symbols (\\n) are allowed.")
		return false
	}
	if _, err := os.Stat(banner); err != nil {
		fmt.Println("Error: Banner file does not exist!")
		return false
	}
	if colorFlag != "" {
		if _, err := ColorToAnsi(colorFlag); err != nil {
			fmt.Printf("Error: Invalid color format '%s'.\n", colorFlag)
			return false
		}
	}
	return true
}

// IsNotASCII checks if a string contains non-ASCII characters
func IsNotASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}
