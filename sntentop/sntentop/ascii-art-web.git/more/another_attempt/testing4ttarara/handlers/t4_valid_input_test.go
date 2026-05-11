package handlers

import "testing"

func TestValidInput(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello, World!", true},                // Valid ASCII characters
		{"Valid\tInput\nHere\r", true},         // Includes valid tab, newline, and carriage return
		{"Invalid\x01Character", false},        // Contains non-printable ASCII
		{"AnotherInvalid\x7FCharacter", false}, // Contains DEL character
		{"", true},                             // Empty string (valid)
		{"Normal123!", true},                   // Valid input with numbers and special characters
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := validInput(tt.input)
			if result != tt.expected {
				t.Errorf("validInput(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
