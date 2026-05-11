package tests

import (
	"fmt"     // For debug printing
	"reflect" // To compare slices (arrays) in tests
	"testing" // Official Go testing package

	"go-reloaded/pipeline" // Import the pipeline package containing ReplaceBin
)

// TestReplaceBin verifies that pipeline.ReplaceBin correctly converts
// binary numbers before "(bin)" into decimal numbers across several scenarios.
func TestReplaceBin(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string   // Name of the test case
		input    []string // Input tokens (words)
		expected []string // Expected output after ReplaceBin
	}{
		{
			name:     "Simple binary conversion", // Basic test
			input:    []string{"It", "has", "been", "10", "(bin)", "years"},
			expected: []string{"It", "has", "been", "2", "years"}, // 10 binary -> 2 decimal
		},
		{
			name:     "Multiple binary conversions", // More than one (bin) in same sentence
			input:    []string{"Value", "101", "(bin)", "and", "11", "(bin)", "check"},
			expected: []string{"Value", "5", "and", "3", "check"}, // 101->5, 11->3
		},
		{
			name:     "Invalid binary string should stay unchanged", // Invalid binary
			input:    []string{"This", "is", "102", "(bin)", "test"},
			expected: []string{"This", "is", "102", "test"}, // 102 not binary, stays
		},
		{
			name:     "Binary marker at beginning does nothing", // Edge case
			input:    []string{"(bin)", "10", "test"},
			expected: []string{"10", "test"}, // No previous word
		},
		{
			name:     "Lowercase binary strings", // Binary digits are lowercase, but valid
			input:    []string{"value", "1101", "(bin)"},
			expected: []string{"value", "13"}, // 1101 -> 13
		},
	}

	// Loop through all test cases and run them individually.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // Run each test separately
			// Debug output to help trace failing cases
			fmt.Printf("\n[DEBUG] Running test: %s\n", tt.name)
			fmt.Printf("[DEBUG] Input tokens: %v\n", tt.input)

			// Call ReplaceBin and inspect its output
			result := pipeline.ReplaceBin(tt.input)
			fmt.Printf("[DEBUG] Output tokens:   %v\n", result)
			fmt.Printf("[DEBUG] Expected tokens: %v\n", tt.expected)

			// Validate result
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\n❌ Test failed: %s\nExpected: %v\nGot:      %v", tt.name, tt.expected, result)
			} else {
				fmt.Printf("✅ Test passed: %s\n", tt.name)
			}
		})
	}
}
