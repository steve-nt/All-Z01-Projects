package tests

import (
	"go-reloaded/pipeline"
	"os"
	"testing"
)

func TestReadInput(t *testing.T) {
	// Create a temporary file and write simple content to it. This ensures the test
	// does not depend on any external files and can run in CI or locally.
	tmpFile, err := os.CreateTemp("", "input_*.txt")
	if err != nil {
		t.Fatalf("❌ Failed to create temporary file: %v", err)
	}
	// Ensure the temp file is removed after test finishes to avoid leaks.
	defer os.Remove(tmpFile.Name())

	// Write a simple string and close the file before reading.
	content := "Hello, Go Reloaded!"
	tmpFile.WriteString(content)
	tmpFile.Close()

	// Read back using the pipeline helper under test.
	result, err := pipeline.ReadInput(tmpFile.Name())

	// No error expected when reading the temporary file.
	if err != nil {
		t.Fatalf("❌ Unexpected error: %v", err)
	}

	// The returned content should exactly match what we wrote.
	if result != content {
		t.Errorf("❌ Expected '%s', receive '%s'", content, result)
	}

	// Log success for manual test runs; not required for CI.
	t.Logf("✅ The readInput read the file correctly : %s", tmpFile.Name())
}

// TestReadInput_FileNotFound checks the behavior when the file does not exist
func TestReadInput_FileNotFound(t *testing.T) {
	// Call the function with a non-existent path
	_, err := pipeline.ReadInput("non_existing_file.txt")

	// If it doesn't return an error, the test fails
	if err == nil {
		t.Errorf("❌ An error was expected for a non-existent file, but none was received.")
	} else {
		t.Logf("✅ the non_existent file handling works correctly: %v", err)
	}
}
