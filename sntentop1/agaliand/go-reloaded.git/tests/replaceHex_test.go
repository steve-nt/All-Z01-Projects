package tests // This file belongs to the "tests" package, located in the /tests directory

import (
	"fmt"     // Used for printing debug messages to the console
	"reflect" // Used to compare slices (arrays) in tests
	"testing" // The official Go testing package

	"go-reloaded/pipeline" // Import the package where ReplaceHex() is defined
)

// TestReplaceHex validates ReplaceHex behavior across several scenarios.
// For each case we provide tokenized input, run ReplaceHex and compare to expected output.
func TestReplaceHex(t *testing.T) {
	// Define multiple test scenarios to ensure the function works in all cases
	tests := []struct {
		name     string   // The name/description of the test case
		input    []string // The tokenized input text
		expected []string // The expected result after running ReplaceHex
	}{
		{
			name:     "Simple hex conversion",                           // Test 1: Basic example
			input:    []string{"1E", "(hex)", "files", "were", "added"}, // "1E" in hex = 30 in decimal
			expected: []string{"30", "files", "were", "added"},          // Expected output after conversion
		},
		{
			name:     "Lowercase hex value",                         // Test 2: Lowercase hex letters should also work
			input:    []string{"ff", "(hex)", "is", "the", "value"}, // "ff" = 255 decimal
			expected: []string{"255", "is", "the", "value"},
		},
		{
			name:     "Invalid hex string should stay unchanged", // Test 3: Handle invalid hex input
			input:    []string{"xyz", "(hex)", "test"},           // "xyz" is not valid hexadecimal
			expected: []string{"xyz", "test"},                    // Should stay unchanged
		},
		{
			name:     "Multiple hex conversions in one sentence",                     // Test 4: Multiple conversions in one line
			input:    []string{"A", "(hex)", "and", "10", "(hex)", "are", "numbers"}, // "A"=10, "10"=16
			expected: []string{"10", "and", "16", "are", "numbers"},
		},
		{
			name:     "Hex marker at beginning does nothing", // Test 5: Edge case - starts with (hex)
			input:    []string{"(hex)", "FF", "test"},        // No previous word to convert
			expected: []string{"FF", "test"},                 // Nothing changes
		},
	}

	// Loop through all defined test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // Run each test separately by name
			// Debug: show the test case name and input
			fmt.Printf("\n[DEBUG] Running test: %s\n", tt.name)
			fmt.Printf("[DEBUG] Input tokens: %v\n", tt.input)

			// Execute function under test
			result := pipeline.ReplaceHex(tt.input)

			// Debug: show output and expected values for easier diagnosis
			fmt.Printf("[DEBUG] Output tokens:   %v\n", result)
			fmt.Printf("[DEBUG] Expected tokens: %v\n", tt.expected)

			// Compare the actual output with the expected result and fail on mismatch
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\n❌ Test failed: %s\nExpected: %v\nGot:      %v", tt.name, tt.expected, result)
			} else {
				fmt.Printf("✅ Test passed: %s\n", tt.name)
			}
		})
	}
}
