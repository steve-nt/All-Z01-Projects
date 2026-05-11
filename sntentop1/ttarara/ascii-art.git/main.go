package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if exactly one argument (the string to convert) is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <string>") // Prompt the user on correct usage
		return
	}

	// Get the input string from the command-line arguments
	textFromOutside := os.Args[1]

	// Read the ASCII art from the standard.txt file and store the lines
	fileLines := ReadTxt()

	// Convert the lines read from the file into a 2D array of ASCII art templates
	asciiTemplates := return2dASCIIArray(fileLines)

	// Print the ASCII art representation of the input string
	printAllStringASCII(textFromOutside, asciiTemplates)
}

// ReadStandardTxt reads the "standard.txt" file and returns its contents as a slice of strings.
func ReadTxt() []string {
	// Open the standard.txt file
	readFile, err := os.Open("standard.txt")
	if err != nil {
		fmt.Println(err) // Handle the error if the file can't be opened
		return nil
	}
	defer readFile.Close() // Ensure the file is closed after reading

	// Scanner to read the file line by line
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines) // Read line by line (line breaks as delimiters)
	var fileLines []string

	// Append each line to the fileLines slice
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines // Return all lines read from the file
}

// return2dASCIIArray processes the lines read from the file and converts them into a 2D array of ASCII templates.
func return2dASCIIArray(fileLines []string) [][]string {
	var asciiTemplates [][]string // This will store all the ASCII art templates
	counter := 0                  // This counts lines to group them into characters (9 lines per character)
	var tempAsciArray []string    // Temporary slice to store 1 character's ASCII art

	// Loop through all lines in the file
	for _, line := range fileLines {
		counter++
		if counter != 1 { // Skip the first line (usually empty or unnecessary)
			tempAsciArray = append(tempAsciArray, line) // Collect lines for the current character
		}
		if counter == 9 { // Every 9 lines correspond to a complete ASCII character
			asciiTemplates = append(asciiTemplates, tempAsciArray) // Add the character to the main list
			counter = 0                                            // Reset the counter for the next character
			tempAsciArray = nil                                    // Reset temp array for the next character
		}
	}

	return asciiTemplates // Return the 2D array of ASCII art templates
}

// printMultipleCharacter prints the ASCII art of multiple characters in a string.
func printMultipleCharacter(s string, asciiTemplates [][]string) {
	// Convert each character in the string to its corresponding ASCII art index
	tempIntArrLetter := returnAsciiCodeInt(s)

	// Each ASCII character is represented by 8 lines, so we loop through 8 times
	for i := 0; i < 8; i++ {
		// Loop through each character in the input string
		for _, v := range tempIntArrLetter {
			// Print the corresponding line of ASCII art for each character
			fmt.Print(asciiTemplates[v][i])
		}
		fmt.Println() // Print a newline after each row of ASCII art
	}
}

// returnAsciiCodeInt converts each character in the string to its ASCII template index.
func returnAsciiCodeInt(s string) []int {
	var tempIntArrLetter []int
	// Convert each character to an integer representing its position in the ASCII table
	// The ASCII art templates start at character 32 (space), so we subtract 32
	for _, v := range s {
		tempIntArrLetter = append(tempIntArrLetter, int(v)-32)
	}
	return tempIntArrLetter // Return the list of ASCII art indices
}

// printAllStringASCII processes the input string and prints its ASCII art, including handling custom newlines.
func printAllStringASCII(text string, asciiTemplates [][]string) {
	// Split the input string based on the escaped newline sequence "\n"
	substrings := returnstring2EndlineArray(text)
	lenOfsubstrings := len(substrings)

	// Loop through each substring (splits the input based on "\n")
	for index, v := range substrings {
		if v == "\\n" { // If the substring is a newline indicator
			if index == lenOfsubstrings-1 {
				// If it's the last element, print a newline
				fmt.Println("")
			} else if substrings[index-1] == "\\n" {
				// If there are two consecutive "\n" sequences, print an empty line
				fmt.Println("")
			}
		} else {
			// Otherwise, print the ASCII art of the current substring
			printMultipleCharacter(v, asciiTemplates)
		}
	}
}

// returnstring2EndlineArray splits the input string into a slice using the custom newline delimiter "\n".
func returnstring2EndlineArray(text string) []string {
	substrings := make([]string, 0)
	escapedN := "\\n"

	// Loop to find all occurrences of the custom newline "\n"
	for {
		idx := strings.Index(text, escapedN) // Find the next occurrence of "\n"
		if idx == -1 {
			substrings = append(substrings, text) // No more occurrences, append the remaining text
			break
		}

		// Append the substring before the "\n" to the list
		substrings = append(substrings, text[:idx])
		// Append the newline "\n" itself to the list
		substrings = append(substrings, escapedN)
		// Update the text to exclude the part before and including "\n"
		text = text[idx+len(escapedN):]
	}

	// Filter out any empty strings in the final slice
	var cleanedSubstrings []string
	for _, sub := range substrings {
		if sub != "" {
			cleanedSubstrings = append(cleanedSubstrings, sub)
		}
	}

	return cleanedSubstrings // Return the cleaned slice of substrings
}
