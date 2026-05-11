package tests

import (
	"fmt"
	"go-reloaded/pipeline"
	"reflect"
	"testing"
)

func TestFixQuotes(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Single word quote",
			input:    []string{"I", "am", "exactly", "how", "they", "describe", "me", ":", "'", " awesome ", "'"},
			expected: []string{"I", "am", "exactly", "how", "they", "describe", "me", ":", "'awesome'"},
		},
		{
			name:     "Multiple word quote",
			input:    []string{"As", "Elton", "John", "said", ":", "'", " I", "am", "the", "most", "well-known", "homosexual", "in", "the", "world ", "'"},
			expected: []string{"As", "Elton", "John", "said", ":", "'I am the most well-known homosexual in the world'"},
		},
		{
			name:     "No quotes",
			input:    []string{"Hello", "world"},
			expected: []string{"Hello", "world"},
		},
		{
			name:     "Nested quotes ignored",
			input:    []string{"He", "said", "'", "It's", "amazing", "'"},
			expected: []string{"He", "said", "'It's amazing'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Debug: print which test case is running and its input tokens.
			fmt.Println("[DEBUG] Running test:", tt.name)
			fmt.Println("[DEBUG] Input tokens: ", tt.input)

			// Run FixQuotes on the test input.
			result := pipeline.FixQuotes(tt.input)

			// Debug: show what the function returned and the expected result.
			fmt.Println("[DEBUG] Output tokens:", result)
			fmt.Println("[DEBUG] Expected tokens:", tt.expected)

			// Compare and fail the test if results differ from expectation.
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("❌ Test failed: %s\nExpected: %v\nGot:      %v", tt.name, tt.expected, result)
			} else {
				fmt.Println("✅ Test passed:", tt.name)
			}
		})
	}
}
