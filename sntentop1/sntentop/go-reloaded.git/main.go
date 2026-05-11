package main

import (
	"bufio"   // Provides buffered input/output, useful for reading lines from a file
	"fmt"     // Provides formatted I/O, such as printing to the console
	"os"      // Provides functions to interact with the operating system (e.g., open files)
	"regexp"  // Provides regular expressions (used for pattern matching in strings)
	"strconv" // Provides functions for converting between strings and numbers
	"strings" // Provides string manipulation functions (e.g., split, join, replace)
)

func main() {
	// First, check if the user provided both an input and output file name in the command-line arguments.
	// `os.Args` contains command-line arguments: os.Args[0] is the program name, os.Args[1] is the input file, and os.Args[2] is the output file.
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}

	// Assign the command-line arguments to variables.
	inputFile, outputFile := os.Args[1], os.Args[2]

	// Attempt to open the input file for reading.
	// If the file doesn't exist or there is another error, print the error and exit.
	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	defer input.Close() // Ensure the file is closed when the function exits.

	// Attempt to create the output file for writing.
	// If there's an error (e.g., permissions issue), print the error and exit.
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close() // Ensure the file is closed when the function exits.

	// Create a scanner to read the input file line by line.
	scanner := bufio.NewScanner(input)

	// Process each line in the input file.
	for scanner.Scan() {
		line := processLine(scanner.Text()) // Apply all transformations to the line.

		output.WriteString(line + "\n") // Write the transformed line to the output file, followed by a newline character.

	}
}

// This function takes a line of text and applies multiple transformations to it.
func processLine(line string) string {
	line = applyHexAndBinConversions(line) // Convert hex and binary numbers to decimal.
	line = applyCaseTransformations(line)  // Adjust text case (e.g., make uppercase or lowercase).
	line = formatPunctuation(line)         // Format punctuation marks for consistency.
	line = replaceAWithAn(line)            // Replace "a" with "an" before vowels or "h".
	return line                            // Return the transformed line.
}

// This function converts hexadecimal (hex) and binary (bin) numbers in the text to decimal numbers.
func applyHexAndBinConversions(text string) string {
	// Helper function that converts a number (in string form) from a specific base (e.g., 16 for hex) to decimal.
	convert := func(match string, base int) string {
		// Extract just the number part (before " (hex)" or " (bin)") by splitting the string.
		num := strings.Split(match, " ")[0]

		// Convert the number from the given base (hexadecimal or binary) to a decimal integer.
		val, err := strconv.ParseInt(num, base, 64)
		if err == nil {
			// If conversion is successful, return the decimal value as a string.
			return fmt.Sprintf("%d", val)
		}
		// If conversion fails, return the original text (no changes).
		return match
	}

	// Regular expression to find hexadecimal numbers in the form "123 (hex)".
	hexRe := regexp.MustCompile(`\b([0-9A-Fa-f]+) \(hex\)`)
	// Replace all matches of the hex pattern with their decimal equivalents.
	text = hexRe.ReplaceAllStringFunc(text, func(match string) string { return convert(match, 16) })

	// Regular expression to find binary numbers in the form "101 (bin)".
	binRe := regexp.MustCompile(`\b([01]+) \(bin\)`)
	// Replace all matches of the bin pattern with their decimal equivalents.
	text = binRe.ReplaceAllStringFunc(text, func(match string) string { return convert(match, 2) })

	return text // Return the transformed text.
}

// This function applies transformations to change the case of certain words (e.g., uppercase, lowercase, capitalized).
func applyCaseTransformations(text string) string {
	// Regular expression to find patterns like "word (up)" or "word (low)" or "word (cap, 2)".
	// It matches a word followed by a transformation command.
	re := regexp.MustCompile(`\b((?:\w+\s*)+?)( \((up|low|cap)(?:, (\d+))?\))`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the parts of the match: the words, the command (up/low/cap), and an optional count.
		parts := re.FindStringSubmatch(match)
		words := strings.Fields(parts[1]) // The words before the command.
		command, count := parts[3], 1     // Default count is 1 unless specified.

		// If a specific count is provided (e.g., "(up, 2)"), parse it.
		if parts[4] != "" {
			count, _ = strconv.Atoi(parts[4])
		}

		// Start applying the transformation to the last 'count' words.
		start := max(0, len(words)-count)
		for i := start; i < len(words); i++ {
			switch command {
			case "up":
				words[i] = strings.ToUpper(words[i]) // Convert the word to uppercase.
			case "low":
				words[i] = strings.ToLower(words[i]) // Convert the word to lowercase.
			case "cap":
				words[i] = strings.Title(strings.ToLower(words[i])) // Capitalize the word.
			}
		}
		return strings.Join(words, " ") // Rejoin the words into a single string.
	})
}

/*
OLD ONE WAS NOT WORKING CORRECTLY
/ This function formats punctuation, ensuring correct spacing and consistency.
func formatPunctuation(text string) string {
	// Ensure ellipses ("...") are spaced correctly, allowing for exactly three dots.
	text = regexp.MustCompile(`\s*\.{3}\s*`).ReplaceAllString(text, "... ")

	// Remove extra spaces around punctuation like ! and ?.
	text = regexp.MustCompile(`\s*([!?]+)\s*`).ReplaceAllString(text, "$1")

	// Remove extra spaces around punctuation like periods, commas, semicolons, and colons.
	text = regexp.MustCompile(`\s*([.,;:])\s*`).ReplaceAllString(text, "$1 ")

	// Ensure spaces inside single (' ') and double (" ") quotes are correct.
	text = regexp.MustCompile(`'\s*(.*?)\s*'`).ReplaceAllString(text, "'$1'")
	text = regexp.MustCompile(`"\s*(.*?)\s*"`).ReplaceAllString(text, "\"$1\"")

	return strings.TrimSpace(text) // Remove leading and trailing spaces from the entire line.
}*/
// This function formats punctuation, ensuring correct spacing and consistency.
func formatPunctuation(text string) string {
	// Handle ellipses first and ensure they are correctly formatted with no extra spaces.
	text = regexp.MustCompile(`\s*\.{3}\s*`).ReplaceAllString(text, "... ")

	// Remove extra spaces around punctuation like ! and ?.
	text = regexp.MustCompile(`\s*([.!?])\s*`).ReplaceAllString(text, "$1")

	// Remove extra spaces around periods, commas, semicolons, and colons, but ignore ellipses.
	text = regexp.MustCompile(`\s*([,;:])\s*`).ReplaceAllString(text, "$1 ")

	// Handle spacing after periods, question marks, and exclamation points.
	text = regexp.MustCompile(`([.?!])\s*([A-Z])`).ReplaceAllString(text, "$1 $2")

	// Ensure spaces inside single (' ') and double (" ") quotes are correct.
	text = regexp.MustCompile(`'\s*(.*?)\s*'`).ReplaceAllString(text, "'$1'")
	text = regexp.MustCompile(`"\s*(.*?)\s*"`).ReplaceAllString(text, "\"$1\"")

	return strings.TrimSpace(text) // Remove leading and trailing spaces.
}

// This function replaces "a" with "an" if the following word starts with a vowel or "h".
func replaceAWithAn(text string) string {
	// Regular expression to find "a" or "A" followed by a vowel or the letter "h".
	aRe := regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`)

	// Replace occurrences of "a " with "an " when appropriate.
	return aRe.ReplaceAllStringFunc(text, func(match string) string {
		return strings.Replace(match, "a ", "an ", 1)
	})
}
