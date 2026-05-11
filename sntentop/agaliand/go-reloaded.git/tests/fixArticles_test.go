package tests // Test package

import (
	"testing"

	"go-reloaded/pipeline"
)

// TestFixArticles contains several scenarios validating that "a" is
// converted to "an" when the following word starts with a vowel or 'h'.
func TestFixArticles(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "a before vowel",
			input:    []string{"There", "it", "was.", "A", "amazing", "rock!"},
			expected: []string{"There", "it", "was.", "An", "amazing", "rock!"},
		},
		{
			name:     "a before consonant",
			input:    []string{"She", "saw", "a", "cat"},
			expected: []string{"She", "saw", "a", "cat"},
		},
		{
			name:     "a before h",
			input:    []string{"He", "waited", "for", "a", "hour"},
			expected: []string{"He", "waited", "for", "an", "hour"},
		},
		{
			name:     "capital A before vowel",
			input:    []string{"It", "was", "A", "honest", "mistake"},
			expected: []string{"It", "was", "An", "honest", "mistake"},
		},
		{
			name:     "capital A before consonant",
			input:    []string{"She", "found", "A", "dog"},
			expected: []string{"She", "found", "A", "dog"},
		},
	}

	// For each scenario, run FixArticles and assert the output tokens match expected.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pipeline.FixArticles(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Test %s failed: expected length %d, got %d", tt.name, len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Fatalf("Test %s failed at index %d: expected '%s', got '%s'", tt.name, i, tt.expected[i], result[i])
				}
			}
		})
	}
}
