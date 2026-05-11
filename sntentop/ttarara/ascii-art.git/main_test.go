package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestReturn2dASCIIArray checks if the return2dASCIIArray function works as expected
func TestReturn2dASCIIArray(t *testing.T) {
	// Pretend we have lines from a file that show two characters in ASCII art form.
	mockFileLines := []string{
		"",       // The first line is ignored (maybe it's empty)
		"Line 1", // The first part of an ASCII character (like a drawing)
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
		"Line 6",
		"Line 7",
		"Line 8", // Now we have 8 lines for one character
		"",       // Another empty line (ignored)
		"Line 9", // The first line of a second character's ASCII drawing
		"Line 10",
		"Line 11",
		"Line 12",
		"Line 13",
		"Line 14",
		"Line 15",
		"Line 16", // Now we have 8 lines for the second character
	}

	// This is what we expect to get back from the function
	expectedOutput := [][]string{
		{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5", "Line 6", "Line 7", "Line 8"},        // First character
		{"Line 9", "Line 10", "Line 11", "Line 12", "Line 13", "Line 14", "Line 15", "Line 16"}, // Second character
	}

	// Run the function we are testing
	result := return2dASCIIArray(mockFileLines)

	// Check if the number of characters we got matches what we expected
	if len(result) != len(expectedOutput) {
		t.Errorf("Expected %d characters, got %d", len(expectedOutput), len(result))
	}

	// Now we check each character's lines to make sure they are correct
	for i, char := range result {
		for j, line := range char {
			if line != expectedOutput[i][j] {
				t.Errorf("Expected line %s, got %s", expectedOutput[i][j], line)
			}
		}
	}
}

// TestReturnAsciiCodeInt checks if the returnAsciiCodeInt function converts characters to the right ASCII numbers
func TestReturnAsciiCodeInt(t *testing.T) {
	// We're going to test the function with the string "AB"
	input := "AB"
	// 'A' is 65 in ASCII, so subtracting 32 makes it 33. 'B' is 66, so it becomes 34.
	expectedOutput := []int{33, 34}

	// Run the function
	result := returnAsciiCodeInt(input)

	// Check if each number matches the expected output
	for i, v := range result {
		if v != expectedOutput[i] {
			t.Errorf("Expected %d, got %d", expectedOutput[i], v)
		}
	}
}

// TestReturnString2EndlineArray checks if returnstring2EndlineArray correctly splits strings with "\n"
func TestReturnString2EndlineArray(t *testing.T) {
	// We want to test splitting "Hello\\nWorld\\n" into separate parts
	input := "Hello\\nWorld\\n"
	// We expect it to return "Hello", "\n", "World", "\n"
	expectedOutput := []string{"Hello", "\\n", "World", "\\n"}

	// Run the function
	result := returnstring2EndlineArray(input)

	// Check if the lengths match
	if len(result) != len(expectedOutput) {
		t.Errorf("Expected length %d, got %d", len(expectedOutput), len(result))
	}

	// Check if each part matches what we expected
	for i, v := range result {
		if v != expectedOutput[i] {
			t.Errorf("Expected %s, got %s", expectedOutput[i], v)
		}
	}
}

// TestPrintMultipleCharacter checks if printMultipleCharacter prints the right ASCII art for a string
func TestPrintMultipleCharacter(t *testing.T) {
	// We're mocking (pretending) that we have templates for ASCII art characters.
	mockAsciiTemplates := make([][]string, 35)
	mockAsciiTemplates[33] = []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"} // The drawing for 'A'
	mockAsciiTemplates[34] = []string{"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"} // The drawing for 'B'

	// We expect the function to print something like:
	// A1B1
	// A2B2
	// A3B3
	// A4B4
	// A5B5
	// A6B6
	// A7B7
	// A8B8
	expectedOutput := "A1B1\nA2B2\nA3B3\nA4B4\nA5B5\nA6B6\nA7B7\nA8B8\n"

	// Capture the output of the function
	output := captureOutput(func() {
		printMultipleCharacter("AB", mockAsciiTemplates)
	})

	// Check if the output matches what we expected
	if output != expectedOutput {
		t.Errorf("Expected %s, got %s", expectedOutput, output)
	}
}

// captureOutput is a helper function that "captures" anything printed to the screen
func captureOutput(f func()) string {
	r, w, _ := os.Pipe() // Create a pipe (a way to capture text)
	defer r.Close()      // Make sure to close the pipe later

	stdout := os.Stdout                   // Save the current standard output (screen)
	defer func() { os.Stdout = stdout }() // Make sure to put it back later

	os.Stdout = w // Change where the program prints to the pipe

	f() // Run the function

	w.Close() // Close the pipe when done

	var buf bytes.Buffer // Create a buffer (storage for text)
	io.Copy(&buf, r)     // Copy the captured output into the buffer

	return buf.String() // Return everything that was captured
}
