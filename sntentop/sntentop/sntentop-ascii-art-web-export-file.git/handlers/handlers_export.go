package handlers

// The `handlers` package groups related functions for organizing
// code logically. This specific package might contain utility
// functions for managing input/output or other tasks.
import (
	"fmt"           // Provides formatted I/O functions (e.g., printing output to the console).
	"log"           // Provides logging functionality to record errors and debug messages.
	"os"            // Provides functions for interacting with the operating system (e.g., files, directories).
	"path/filepath" // Provides utility functions for manipulating file paths in a portable way.
)

// ExportToFile writes a given string to a file and returns the file path or an error.
func ExportToFile(data string) (string, error) {
	// `os.Getwd()` retrieves the current working directory of the program.
	// `rootDir` holds this directory path, and `err` captures any errors encountered.
	rootDir, err := os.Getwd()
	// Logs an error message if `os.Getwd()` fails.
	if err != nil {
		log.Println("Error getting current working directory:", err)
		// Returns an empty string and the error, ending the function early.
		return "", err
	}
	// `filepath.Join` constructs a file path by combining `rootDir`, "temp", and "output.txt".
	// This ensures cross-platform compatibility (e.g., using the correct path separators for the OS).
	filePath := filepath.Join(rootDir, "temp", "output.txt")
	// `filepath.Dir` extracts the directory portion ("temp") from the full file path.
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, 0766)
	// `os.MkdirAll` creates the "temp" directory if it doesn't exist.
	// The permissions `0766` allow the owner to read, write, and execute, while others can read and write.
	if err != nil {
		log.Println("Error creating directory:", err)
		// Logs an error if directory creation fails.
		return "", err
		// Returns an error and ends the function.
	}

	err = os.WriteFile(filePath, []byte(data), 0766)
	// `os.WriteFile` writes `data` (converted to a byte slice with `[]byte(data)`) to the file at `filePath`.
	// If the file doesn't exist, it creates it with the permissions `0766`.
	if err != nil {
		log.Println("Error writing to file:", err)
		// Logs an error if writing to the file fails.
		return "", err
		// Returns an error and ends the function.
	}

	err = os.Chmod(filePath, 0766)
	// `os.Chmod` explicitly sets the file's permissions to `0766`.
	if err != nil {
		log.Println("Error setting file permissions:", err)
		// Logs an error if changing file permissions fails.
		return "", err
		// Returns an error and ends the function.
	}

	fmt.Println("Data successfully written to", filePath)
	// Prints a success message to the console with the file path.
	return filePath, nil
	// Returns the file path as a string and `nil` to indicate success.
}
