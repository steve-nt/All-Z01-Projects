package main

import (
	"os"
	"testing"
)

// TestHandleArguments tests the HandleArguments function.
func TestHandleArguments(t *testing.T) {
	// Temporarily override os.Args for testing
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // Restore original args after test

	// Test with valid arguments
	os.Args = []string{"main.go", "hello", "standard"}
	text, banner := HandleArguments()
	if text != "hello" || banner != "standard" {
		t.Errorf("HandleArguments() = (%q, %q); want (\"hello\", \"standard\")", text, banner)
	}

	// Test with missing arguments
	os.Args = []string{"main.go"}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("HandleArguments() did not exit on missing arguments")
		}
	}()
	HandleArguments() // This should exit
}
