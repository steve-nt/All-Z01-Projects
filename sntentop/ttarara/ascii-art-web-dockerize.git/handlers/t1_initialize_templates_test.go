package handlers

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create dummy template files
	templateFile1 := filepath.Join(tempDir, "test1.html")
	templateFile2 := filepath.Join(tempDir, "test2.html")
	err := os.WriteFile(templateFile1, []byte("<h1>{{.Title}}</h1>"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}
	err = os.WriteFile(templateFile2, []byte("<p>{{.Content}}</p>"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	// Create a logger with a buffer to capture logs
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "TEST: ", log.LstdFlags)

	// Call InitializeTemplates with the temporary directory
	templates, err := InitializeTemplates(tempDir, logger)
	if err != nil {
		t.Fatalf("InitializeTemplates failed: %v", err)
	}

	// Ensure the templates were parsed correctly
	if templates == nil {
		t.Fatalf("Expected templates to be non-nil")
	}

	// Test rendering with one of the templates
	var output bytes.Buffer
	err = templates.ExecuteTemplate(&output, "test1.html", map[string]string{"Title": "Hello, World!"})
	if err != nil {
		t.Fatalf("Failed to execute template: %v", err)
	}
	expectedOutput := "<h1>Hello, World!</h1>"
	if output.String() != expectedOutput {
		t.Errorf("Expected output %q, got %q", expectedOutput, output.String())
	}

	// Check the logs for successful parsing
	if !bytes.Contains(logBuffer.Bytes(), []byte("Templates successfully parsed from path")) {
		t.Errorf("Expected log message about successful parsing, but it was not found")
	}

	// Test failure case by providing an invalid path
	_, err = InitializeTemplates("/invalid/path", logger)
	if err == nil {
		t.Errorf("Expected error for invalid path, got nil")
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("Failed to parse templates")) {
		t.Errorf("Expected log message about failed parsing, but it was not found")
	}
}
