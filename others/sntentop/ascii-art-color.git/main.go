package main

// Importing packages that provide pre-written code for specific functions
import (
	"banner/banner" // This is a custom package that contains code specific to banner generation
	"fmt"           // Provides functions for formatted I/O, such as printing messages
	"os"            // Provides access to system functions, such as reading arguments passed to the program
	"path/filepath" // Provides functions to work with file paths
	"strings"       // Provides functions for working with strings (text data)
)

// This function displays an error message if the user provides incorrect inputs.
func ShowMessageError() {
	// Printing usage instructions to the user on how to correctly run the program
	fmt.Printf("  Usage: go run . [OPTION] [STRING]\n\n")
	fmt.Println("  EX: go run . --color=<color> <letters to be colored> \"something\"")
}

// This function checks if a file name corresponds to a valid banner style file.
func IsFile(fileName string) bool {
	// Ensuring the file name has a ".txt" extension; if not, appends ".txt"
	if filepath.Ext(fileName) != ".txt" {
		fileName += ".txt"
	}
	// Returns true if the file name matches one of the predefined styles
	return fileName == "standard.txt" || fileName == "shadow.txt" || fileName == "thinkertoy.txt"
}

// Main function to generate ASCII art based on user input
func Ascii_Arts_Generator(arguments []string) {
	// Storing the number of arguments provided by the user
	argument_Size := len(arguments)
	// Default file name for banner style if not specified
	fileName := "standard"

	// Handling different cases based on the number of arguments provided
	switch argument_Size {
	case 4: // When 4 arguments are provided
		// Checks if the first argument specifies a color
		if len(arguments[0]) > 8 && arguments[0][:8] == "--color=" {
			// Sets the file name to the fourth argument
			fileName = arguments[3]
			// The text to be converted into ASCII art is stored in Record
			Record := arguments[2]
			// Extracts the color specified in the first argument
			color := arguments[0][8:]
			// Counts the occurrences of "\n" to know how many lines to create
			newLineCounter := strings.Count(fileName, "\\n")
			// Processes the text for ASCII art conversion
			words := ManipulateData(Record, fileName)
			// Calls a function to generate and display the final ASCII art
			banner.Result(words, newLineCounter, GetDataFromFile(fileName), color, arguments[1])
		} else {
			// Shows an error if the first argument is not in the correct format
			ShowMessageError()
		}
	case 3: // When 3 arguments are provided
		if len(arguments[0]) > 8 && arguments[0][:8] == "--color=" {
			// Extracts the color value
			color := arguments[0][8:]
			// Initializes empty variables for word and letters
			word := ""
			letters := ""
			// Checks if the third argument is a valid file name
			if IsFile(arguments[2]) {
				// If it's a file name, sets fileName and assigns the second argument to word
				fileName = arguments[2]
				word = arguments[1]
			} else {
				// Otherwise, treats the third argument as the word to be converted
				word = arguments[2]
				letters = arguments[1]
			}
			// Calls the Result function to generate the ASCII art
			banner.Result(ManipulateData(word, fileName), strings.Count(word, "\\n"), GetDataFromFile(fileName), color, letters)
		} else {
			// Shows an error if the first argument is not in the correct format
			ShowMessageError()
		}
	case 2: // When 2 arguments are provided
		// Checks if the second argument is a valid file name
		if IsFile(arguments[1]) {
			// Calls the Result function without color options
			banner.Result(ManipulateData(arguments[0], arguments[1]), strings.Count(arguments[0], "\\n"), GetDataFromFile(arguments[1]), "", "")
		} else if arguments[0][:8] == "--color=" {
			// If the first argument specifies a color, extracts it
			banner.Result(ManipulateData(arguments[1], fileName), strings.Count(arguments[1], "\\n"), GetDataFromFile(fileName), arguments[0][8:], "")
		} else {
			// Shows an error if the format is incorrect
			ShowMessageError()
		}
	case 1: // When only 1 argument is provided
		// Calls the Result function with the default file name and no color
		banner.Result(ManipulateData(arguments[0], fileName), strings.Count(arguments[0], "\\n"), GetDataFromFile(fileName), "", "")
	default: // If the number of arguments is not 1, 2, 3, or 4, shows an error
		ShowMessageError()
	}
}

// This function retrieves ASCII art character mappings from a file
func GetDataFromFile(fileName string) map[int][]string {
	// Adds ".txt" extension if it's missing
	if filepath.Ext(fileName) != ".txt" {
		fileName += ".txt"
	}
	// Calls a function from the banner package to read and return character mappings from the file
	return banner.ReadBannerFiles(fileName)
}

// This function processes the text to be converted into ASCII art
func ManipulateData(arg, fileName string) []string {
	// Splits the input text at each "\n" to create separate lines for ASCII art
	words := strings.Split(arg, "\\n")
	// Checks if all characters in the text can be found in the ASCII art file
	if !banner.CheckIfAllCharInFile(words) {
		// Displays an error message if there are unsupported characters
		fmt.Println("You Have a Character Not found in the file >>", fileName)
		// Returns nil to indicate an error
		return nil
	}
	// Returns the processed list of words (lines)
	return words
}

// This is the main function that runs first when the program starts
func main() {
	// Retrieves all the arguments provided when the program is run
	arg := os.Args
	// Calls the ASCII art generator function with the arguments
	Ascii_Arts_Generator(arg[1:])
}
