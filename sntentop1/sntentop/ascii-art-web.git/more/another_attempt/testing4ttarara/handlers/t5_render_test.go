package handlers

import "testing"

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		status   int
	}{
		{
			name:     "Valid Input",
			input:    "Hello, World!",
			expected: "Hello, World!",
			status:   200,
		},
		{
			name:     "Input with Carriage Return and Newline",
			input:    "Line1\r\nLine2",
			expected: "Line1\\nLine2",
			status:   200,
		},
		{
			name:     "Empty Input",
			input:    "   ",
			expected: "Input cannot be empty or just whitespace.",
			status:   400,
		},
		{
			name:     "Too Long Input",
			input:    string(make([]byte, 129)), // 129 characters
			expected: "Input too long. Maximum allowed length is 128 characters.",
			status:   400,
		},
		{
			name:     "Input with Invalid Characters",
			input:    "Invalid\x01Input",
			expected: "Input contains invalid characters. Only ASCII characters, tabs, and newlines are allowed.",
			status:   400,
		},
		{
			name:     "Trim Whitespace",
			input:    "  Trimmed Input  ",
			expected: "Trimmed Input",
			status:   200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, status := render(tt.input)

			// Check the status code
			if status != tt.status {
				t.Errorf("render(%q): expected status %d, got %d", tt.input, tt.status, status)
			}

			// Check the returned string
			if result != tt.expected {
				t.Errorf("render(%q): expected %q, got %q", tt.input, tt.expected, result)
			}
		})
	}
}
