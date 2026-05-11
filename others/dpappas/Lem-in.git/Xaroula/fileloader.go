package main

import (
	"bufio"
	"os"
)

// ReadLines reads all lines from the given file and returns them as a string slice.
// Returns an error if the file can't be opened or read.
func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when the function exits

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read file line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
