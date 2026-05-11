package tests

import (
	"fmt"
	"go-reloaded/pipeline"
	"reflect"
	"testing"
)

func TestFormatPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Basic punctuation spacing",
			input:    []string{"I", "was", "sitting", "over", "there", ",", "and", "then", "BAMM", "!!"},
			expected: []string{"I", "was", "sitting", "over", "there,", "and", "then", "BAMM!!"},
		},
		{
			name:     "Ellipsis spacing",
			input:    []string{"I", "was", "thinking", "...", "You", "were", "right"},
			expected: []string{"I", "was", "thinking...", "You", "were", "right"},
		},
		{
			name:     "Combined punctuation spacing",
			input:    []string{"Wait", "!?"},
			expected: []string{"Wait!?"},
		},
		{
			name:     "Colon punctuation",
			input:    []string{"Listen", ":", "carefully"},
			expected: []string{"Listen:", "carefully"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Debug: print the test name and input to aid troubleshooting.
			fmt.Println("[DEBUG] Running test:", tt.name)
			fmt.Println("[DEBUG] Input tokens: ", tt.input)

			// Execute the function under test and capture its output.
			result := pipeline.FormatPunctuation(tt.input)

			// Debug: show output and expected tokens for visibility.
			fmt.Println("[DEBUG] Output tokens:", result)
			fmt.Println("[DEBUG] Expected tokens:", tt.expected)

			// Assert equality between expected and result; fail the test on mismatch.
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("❌ Test failed: %s\nExpected: %v\nGot:      %v", tt.name, tt.expected, result)
			} else {
				fmt.Println("✅ Test passed:", tt.name)
			}
		})
	}
}
