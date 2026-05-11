package tests

import (
	"os"
	"strings"
	"testing"

	"go-reloaded/pipeline"
)

func TestWriteOutput(t *testing.T) {
	// WriteOutput expects each slice element to be a full line.
	// Provide a single-line input to match newline-per-line behavior.
	input := []string{"Hello beautiful world!"}
	tmpFile := "test_output.txt"

	err := pipeline.WriteOutput(tmpFile, input)
	if err != nil {
		t.Fatalf("WriteOutput returned error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := strings.TrimSpace(string(data))
	expected := "Hello beautiful world!"

	if content != expected {
		t.Errorf("Expected '%s', got '%s'", expected, content)
	}

	os.Remove(tmpFile)
}

func TestWriteOutput_EmptyInput(t *testing.T) {
	tmpFile := "test_output_empty.txt"
	input := []string{}

	err := pipeline.WriteOutput(tmpFile, input)
	if err != nil {
		t.Fatalf("WriteOutput returned error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := strings.TrimSpace(string(data))

	if content != "" {
		t.Errorf("Expected empty file, got '%s'", content)
	}

	os.Remove(tmpFile)
}
