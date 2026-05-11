package pipeline // This file belongs to the 'pipeline' package

import (
	"os" // Import the 'os' package to interact with the operating system (for reading files)
)

// ReadInput reads the entire content of a file and returns it as a string.
// If an error occurs (e.g., file not found), it returns an empty string and the error.
func ReadInput(path string) (string, error) {
	// Read the whole file into memory as a byte slice.
	data, err := os.ReadFile(path)
	if err != nil {
		// Propagate any read error to the caller along with an empty string.
		return "", err
	}
	// Convert bytes to string and return success (nil error).
	return string(data), nil
}
