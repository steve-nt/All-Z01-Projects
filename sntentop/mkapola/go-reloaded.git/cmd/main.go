// This tells Go that this file is the main program (the entry point)
package main

// Import statements - these bring in libraries we need
import (
	"fmt"                       // Library for printing text to the screen
	processor "go-reloaded/pkg" // Our custom library that does the text processing
	"os"                        // Library for working with the operating system (files, command line)
)

// The main function - this is where the program starts running
func main() {
	// Check if the user provided exactly 2 arguments (input file and output file)
	if len(os.Args) != 3 { // os.Args[0] is the program name, so we need 3 total (program + 2 files)
		// If wrong number of arguments, show how to use the program correctly
		fmt.Println("🔍👀 Still looking… Try: go run main.go <input_file> <output_file>") // Ctrl + . for Linux, Windows key + . for Windows
		os.Exit(1)                                                                      // Exit the program with error code 1 (means something went wrong)
	}

	// Get the file names from the command line arguments
	inputFile := os.Args[1]  // First argument is the input file name
	outputFile := os.Args[2] // Second argument is the output file name

	// Try to read the content from the input file
	content, err := processor.ReadFile(inputFile) // Call our ReadFile function
	if err != nil {                               // err = error value                            // If there was an error reading the file
		// Print the error message and exit
		fmt.Printf("🤖💔 Robot tried. Robot failed. Cannot Read: %v\n", err) // %v is a placeholder for the error message
		os.Exit(2)                                                         // Exit with error code
	}

	// Process the text through our transformation pipeline
	correctedText := processor.ProcessText(content) // This does all the magic transformations

	// Try to write the corrected text to the output file
	err = processor.WriteFile(outputFile, correctedText) // Call our WriteFile function
	if err != nil {                                      // If there was an error writing the file
		// Print the error message and exit
		fmt.Printf("🗃️💀 Output file suffered a tragic fate... Write error: %v\n", err)
		os.Exit(3) // Exit with error code
	}

	// If we get here, everything worked successfully!
	fmt.Printf("🌟✨ Success! Output written to %s\n", outputFile)
}
