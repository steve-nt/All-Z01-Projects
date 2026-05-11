package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const CharHeight = 8

// Function to load the ASCII art from the banner file into a map
func loadAsciiMap(filename string) (map[rune][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	asciiMap := make(map[rune][]string)
	scanner := bufio.NewScanner(file)

	var currentChar rune = 32 // ' ' starts at ASCII code 32
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if len(lines) == CharHeight {
			asciiMap[currentChar] = lines
			currentChar++
			lines = []string{} // Reset for next character
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return asciiMap, nil
}

// Function to convert a string to ASCII art using the loaded map
func convertToAsciiArt(input string, asciiMap map[rune][]string) string {
	var outputLines [CharHeight]string
	var output string

	for _, char := range input {
		if char == '\n' {
			output += strings.Join(outputLines[:], "\n") + "\n"
			outputLines = [CharHeight]string{}
			continue
		}

		if art, ok := asciiMap[char]; ok {
			for i := 0; i < CharHeight; i++ {
				outputLines[i] += art[i]
			}
		} else {
			for i := 0; i < CharHeight; i++ {
				outputLines[i] += " "
			}
		}
	}

	output += strings.Join(outputLines[:], "\n")
	return output
}

func main() {
	// Get command-line arguments
	args := os.Args

	// Check if at least the text argument is provided
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <\"text_to_print\"> [banner_file]")
		fmt.Println("Example: go run main.go \"HELLO WORLD\" [standard.txt]")
		os.Exit(1)
	}

	// Extract input text from command-line arguments
	inputText := args[1]

	// Default to 'standard.txt' if no banner file is specified
	bannerFile := "standard.txt"
	if len(args) > 2 {
		bannerFile = args[2]
	}

	// Load the ASCII art from the user-specified file or default to 'standard.txt'
	asciiMap, err := loadAsciiMap(bannerFile)
	if err != nil {
		fmt.Printf("Error loading ASCII art from file '%s': %v\n", bannerFile, err)
		fmt.Println("Falling back to default banner file 'standard.txt'.")

		// Attempt to load the default 'standard.txt' banner file
		asciiMap, err = loadAsciiMap("standard.txt")
		if err != nil {
			fmt.Println("Error loading default banner file 'standard.txt'. Exiting.")
			os.Exit(1)
		}
	}

	// Convert the input text to ASCII art
	asciiArt := convertToAsciiArt(inputText, asciiMap)
	fmt.Println(asciiArt)
}
