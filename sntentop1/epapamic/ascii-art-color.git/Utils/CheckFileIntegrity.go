package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
)

// Designated path of input file.
const mapPath = "Data/standard.txt"

// Original hash of the Data/standard.txt file.
var hash = "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"

// Creates a hash string for the designated file.
// (In this case, the file is declared inside the function).
func createHash(mapPath string) (string, error) {
	// Read the file contents
	data, err := os.ReadFile(mapPath)
	if err != nil {
		return "", err
	}

	// Compute the hash of the file contents
	hash := sha256.New()
	hash.Write(data)
	actualHash := fmt.Sprintf("%x", hash.Sum(nil))

	return actualHash, nil
}

// Calls createHash to generate a new hash for the internaly designated file,
// compares the new hash with the original, returns non-nil error if they are non-matching.
func checkFileIntegrity(hashCreator func(string) (string, error)) error {
	currentFileHash, err := hashCreator(mapPath)
	if err != nil {
		return err
	}
	if currentFileHash != hash {
		return errors.New("ERRROR: malformed Data/standard.txt")
	}
	return nil
}

// Validates that the Data/standard.txt file hasn't been tempered with.
func CheckFileIntegrity() error {
	return checkFileIntegrity(func(mapPath string) (string, error) {
		return createHash(mapPath)
	})
}
