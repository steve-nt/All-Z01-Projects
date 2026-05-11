package handlers

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TestLogEvent tests the behavior of the logEvent function.
func TestLogEvent(t *testing.T) {
	// Reset the global history for testing
	history = []string{}

	// Capture stdout
	var outputBuffer bytes.Buffer
	stdout := os.Stdout                   // Save the original stdout
	r, w, _ := os.Pipe()                  // Create a pipe for stdout
	os.Stdout = w                         // Redirect stdout to the pipe
	defer func() { os.Stdout = stdout }() // Restore original stdout after the test

	// Define test inputs
	eventType := "INFO"
	message := "Test event occurred"

	// Call logEvent
	logEvent(eventType, message)

	// Close the pipe and read the captured output
	w.Close()
	outputBuffer.ReadFrom(r)

	// Verify history entry
	if len(history) != 1 {
		t.Fatalf("Expected 1 entry in history, got %d", len(history))
	}

	// Verify the format of the log entry
	expectedPrefix := fmt.Sprintf("[%s] %s: %s",
		time.Now().Format("2006-01-02 15:04:05"), eventType, message,
	)
	if !strings.HasPrefix(history[0], expectedPrefix[:len(expectedPrefix)-1]) {
		t.Errorf("Expected log entry to start with %q, got %q", expectedPrefix, history[0])
	}

	// Check printed history output
	expectedOutput := fmt.Sprintf("\n===== Submission and Error History =====\n%s\n", history[0])
	actualOutput := outputBuffer.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected printed output %q, got %q", expectedOutput, actualOutput)
	}
}
