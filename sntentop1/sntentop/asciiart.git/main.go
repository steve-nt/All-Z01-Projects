// * version 7 kind of working
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 1; i < len(lines); i += 9 { // 8 lines of ASCII art + 1 blank line between characters
		if i+8 >= len(lines) {
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8] // Capture 8 lines for each character
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
// PrintASCIIArt prints the ASCII art for the given input string, adding a dollar sign ($) at the end of each line.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8)
	runes := []rune(input)
	hasPrinted := false // To track if any output has been printed

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\n' {
			newlineCount := 1
			for i+1 < len(runes) && runes[i+1] == '\n' {
				newlineCount++
				i++
			}
			if !allLinesEmpty(lines) {
				if hasPrinted {

				}
				for _, line := range lines {
					fmt.Println(line)
				}
				hasPrinted = true
				lines = make([]string, 8)
			}

			// Adjust this part to print `newlineCount - 1` blank lines, each ending with a `$`
			for j := 1; j < newlineCount; j++ {
				fmt.Println()
			}
		} else {
			if art, ok := banner[runes[i]]; ok {
				for j := 0; j < 8; j++ {
					lines[j] += art[j]
				}
			} else {
				for j := 0; j < 8; j++ {
					lines[j] += "      "
				}
			}
		}
	}
	// Print the final set of characters (if any)
	if !allLinesEmpty(lines) {
		if hasPrinted {

		}
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

// Helper function to check if all lines are empty (used to avoid printing empty lines)
func allLinesEmpty(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run . <string>")
	}

	input := os.Args[1]

	// Replace the literal `\n` with an actual newline character
	input = strings.ReplaceAll(input, `\n`, "\n")

	banner, err := LoadBanner("standard.txt")
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	PrintASCIIArt(input, banner)
}

//*/
/*version 6 problematic
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 { // 8 lines of ASCII art + 1 blank line between characters
		if i+8 >= len(lines) {
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8] // Capture 8 lines for each character
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8)
	runes := []rune(input)
	hasPrinted := false // To track if any output has been printed (avoid start blank line)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\n' {
			if !allLinesEmpty(lines) {
				// Print lines collected so far (if there are any)
				if hasPrinted {
					fmt.Println() // Add one blank line after a block of text
				}
				for _, line := range lines {
					fmt.Println(line)
				}
				hasPrinted = true
				lines = make([]string, 8) // Reset lines for the next word
			}
			// Check for double newlines
			if i+1 < len(runes) && runes[i+1] == '\n' {
				i++ // Skip the next newline
				if hasPrinted {
					fmt.Println() // Print exactly one blank line for double newlines
				}
			}
		} else {
			if art, ok := banner[runes[i]]; ok {
				for j := 0; j < 8; j++ {
					lines[j] += art[j] // Append the ASCII art to each line
				}
			} else {
				for j := 0; j < 8; j++ {
					lines[j] += "      " // Handle unsupported characters with blanks
				}
			}
		}
	}

	// Print the final set of characters (if any)
	if !allLinesEmpty(lines) {
		if hasPrinted {
			fmt.Println()
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	}
}

// Helper function to check if all lines are empty (used to avoid printing empty lines)
func allLinesEmpty(lines []string) bool {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run . <string>")
	}

	input := os.Args[1]

	// Replace the literal `\n` with an actual newline character
	input = strings.ReplaceAll(input, `\n`, "\n")

	banner, err := LoadBanner("standard.txt")
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	PrintASCIIArt(input, banner)
}
*/
/* Version 5
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename) // Read the entire file
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 { // 8 lines of ASCII art + 1 blank line between characters
		// Ensure we don't exceed slice bounds if the banner is incomplete
		if i+8 >= len(lines) {
			fmt.Printf("Warning: incomplete data for character %c. Filling with blanks.\n", rune(32+charIndex))
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			// Map each character from ASCII 32 (space) onward
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8] // Capture 8 lines for each character
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8) // Prepare 8 lines to hold the result
	runes := []rune(input)     // Convert input string to rune slice
	newLineTriggered := false  // To track multiple newlines

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\n' {
			// Handle actual newlines
			if !newLineTriggered { // Only trigger once for multiple newlines
				for _, line := range lines {
					fmt.Println(line)
				}
				fmt.Println()             // Single blank line after newline
				lines = make([]string, 8) // Reset for the next set of characters
			}
			newLineTriggered = true
		} else {
			newLineTriggered = false // Reset newline trigger after any other character
			if art, ok := banner[runes[i]]; ok {
				// If character exists in the banner, append its art to the lines
				for j := 0; j < 8; j++ {
					lines[j] += art[j]
				}
			} else {
				// If character is not found, fill with blanks
				fmt.Printf("Warning: character %c not found in banner. Filling with blanks.\n", runes[i])
				for j := 0; j < 8; j++ {
					lines[j] += "      "
				}
			}
		}
	}

	// Print the remaining lines
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run . <string>")
	}

	input := os.Args[1]

	// Replace the literal string `\n` with an actual newline character
	input = strings.ReplaceAll(input, `\n`, "\n")

	banner, err := LoadBanner("standard.txt")
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	PrintASCIIArt(input, banner)
}
*/
/* Verion 4
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename) // Read the entire file
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 { // 8 lines of ASCII art + 1 blank line between characters
		// Ensure we don't exceed slice bounds if the banner is incomplete
		if i+8 >= len(lines) {
			fmt.Printf("Warning: incomplete data for character %c. Filling with blanks.\n", rune(32+charIndex))
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			// Map each character from ASCII 32 (space) onward
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8] // Capture 8 lines for each character
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8) // Prepare 8 lines to hold the result
	runes := []rune(input)     // Convert input string to rune slice

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) && runes[i+1] == 'n' {
			// Handle '\n' as a newline
			for _, line := range lines {
				fmt.Println(line)
			}
			fmt.Println()             // Blank line after newline
			lines = make([]string, 8) // Reset for the next set of characters
			i++                       // Skip over the 'n' after '\'
		} else if art, ok := banner[runes[i]]; ok {
			// If character exists in the banner, append its art to the lines
			for j := 0; j < 8; j++ {
				lines[j] += art[j]
			}
		} else {
			// If character is not found, fill with blanks
			fmt.Printf("Warning: character %c not found in banner. Filling with blanks.\n", runes[i])
			for j := 0; j < 8; j++ {
				lines[j] += "      "
			}
		}
	}

	// Print the remaining lines
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run . <string>")
	}

	input := os.Args[1]
	banner, err := LoadBanner("shadow.txt")
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	PrintASCIIArt(input, banner)
}
*/
/* Version3
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename) // Read the entire file
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 {
		// Ensure we don't exceed slice bounds if the banner is incomplete
		if i+8 >= len(lines) {
			fmt.Printf("Warning: incomplete data for character %c. Filling with blanks.\n", rune(32+charIndex))
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			// Map each character from ASCII 32 (space) onward
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8]
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8)
	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) && runes[i+1] == 'n' {
			// Handle newline by printing the collected lines and resetting the buffer
			for _, line := range lines {
				fmt.Println(line)
			}
			fmt.Println()             // Print a single empty line for newline
			lines = make([]string, 8) // Reset lines for the next set
			i++                       // Skip over the 'n' after '\'
		} else if art, ok := banner[runes[i]]; ok {
			// If the character exists in the banner, append its art to each line
			for j := 0; j < 8; j++ {
				lines[j] += art[j]
			}
		} else {
			// If character isn't found in the banner, replace with blanks
			fmt.Printf("Warning: character %c not found in banner. Filling with blanks.\n", runes[i])
			for j := 0; j < 8; j++ {
				lines[j] += "      "
			}
		}
	}

	// Print the remaining collected lines
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run . <string>")
	}

	input := os.Args[1]
	banner, err := LoadBanner("standard.txt")
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	PrintASCIIArt(input, banner)
}
*/
/* Version 2
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename) // Read the entire file
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 {
		// Ensure we don't exceed slice bounds if the banner is incomplete
		if i+8 >= len(lines) {
			fmt.Printf("Warning: incomplete data for character %c. Filling with blanks.\n", rune(32+charIndex))
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			// Map each character from ASCII 32 (space) onward
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8]
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8)
	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) && runes[i+1] == 'n' {
			// Handle newline by printing the collected lines and resetting the buffer
			for _, line := range lines {
				fmt.Println(line)
			}
			fmt.Println()             // Print empty line for newline
			lines = make([]string, 8) // Reset lines for the next set
			i++                       // Skip over the 'n' after '\'
		} else if art, ok := banner[runes[i]]; ok {
			// If the character exists in the banner, append its art to each line
			for j := 0; j < 8; j++ {
				lines[j] += art[j]
			}
		} else {
			// If character isn't found in the banner, replace with blanks
			fmt.Printf("Warning: character %c not found in banner. Filling with blanks.\n", runes[i])
			for j := 0; j < 8; j++ {
				lines[j] += "      "
			}
		}
	}

	// Print the remaining collected lines
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	// Check if the input string is provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run . [STRING]")
	}

	input := os.Args[1]

	// Load the banner file
	bannerFile := "standard.txt" // Adjust this to use other banner styles (shadow, thinkertoy)
	banner, err := LoadBanner(bannerFile)
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	// Print the ASCII art for the input string
	PrintASCIIArt(input, banner)
}
*/
/*  Version 1
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// LoadBanner loads the banner file and maps each character to its ASCII representation.
func LoadBanner(filename string) (map[rune][]string, error) {
	data, err := ioutil.ReadFile(filename) // Read the entire file
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	asciiMap := make(map[rune][]string)
	charIndex := 0

	for i := 0; i < len(lines); i += 9 {
		// Ensure we don't exceed slice bounds if the banner is incomplete
		if i+8 >= len(lines) {
			fmt.Printf("Warning: incomplete data for character %c. Filling with blanks.\n", rune(32+charIndex))
			asciiMap[rune(32+charIndex)] = []string{"      ", "      ", "      ", "      ", "      ", "      ", "      ", "      "}
		} else {
			// Map each character from ASCII 32 (space) onward
			char := rune(32 + charIndex)
			asciiMap[char] = lines[i : i+8]
		}
		charIndex++
	}

	return asciiMap, nil
}

// PrintASCIIArt prints the ASCII art for the given input string.
func PrintASCIIArt(input string, banner map[rune][]string) {
	lines := make([]string, 8)

	for _, char := range input {
		if char == '\n' {
			for _, line := range lines {
				fmt.Println(line)
			}
			fmt.Println()             // Empty line for new paragraph
			lines = make([]string, 8) // Reset lines for the next set
		} else if art, ok := banner[char]; ok {
			for i := 0; i < 8; i++ {
				lines[i] += art[i] // Append ASCII art for each line
			}
		} else {
			// If character isn't found in the banner, replace with blanks
			fmt.Printf("Warning: character %c not found in banner. Filling with blanks.\n", char)
			for i := 0; i < 8; i++ {
				lines[i] += "      "
			}
		}
	}

	// Print the remaining collected lines
	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	// Check if the input string is provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run . [STRING]")
	}

	input := os.Args[1]

	// Load the banner file
	bannerFile := "standard.txt" // Adjust this to use other banner styles (shadow, thinkertoy)
	banner, err := LoadBanner(bannerFile)
	if err != nil {
		log.Fatalf("Error loading banner: %v", err)
	}

	// Print the ASCII art for the input string
	PrintASCIIArt(input, banner)
}
*/
