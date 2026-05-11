// This tells Go that this file contains tests for the processor package
// The "_test" suffix tells Go this is a test file
package processor_test

// Import statements - these bring in libraries we need for testing
import (
	"testing"                   // Go's built-in testing library
	processor "go-reloaded/pkg" // The code we want to test
)

// Test function for number conversions (hex and binary to decimal)
// The "t *testing.T" parameter lets us report test results
func TestConvertNumbers(t *testing.T) {
	// Create a list of test cases - each has input text and expected output
	tests := []struct {
		input    string // The text we want to transform
		expected string // What we expect to get back
	}{
		{"1E (hex)", "30"},   // Simple hex conversion
		{"10 (bin)", "2"},    // Simple binary conversion
		{"FF (hex)", "255"},  // Larger hex number
		{"1111 (bin)", "15"}, // Larger binary number
		{"0 (hex)", "0"},     // Edge case: zero in hex
		{"0 (bin)", "0"},     // Edge case: zero in binary
		{"I have 1E (hex) apples", "I have 30 apples"},         // Hex in sentence
		{"Binary 101 (bin) is decimal", "Binary 5 is decimal"}, // Binary in sentence
	}

	// Loop through each test case and check if it works correctly
	for _, test := range tests { // "range" goes through each item in the list
		result := processor.ProcessText(test.input) // Run our function on the test input
		if result != test.expected {                // If the result doesn't match what we expected
			// Report the test failure with details about what went wrong
			t.Errorf("Number conversion failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}

// Test function for case transformations (uppercase, lowercase, capitalize)
func TestCaseTransforms(t *testing.T) {
	// Create test cases for different case transformation patterns
	tests := []struct {
		input    string // Input text with case commands
		expected string // Expected output after transformation
	}{
		{"hello (up)", "HELLO"},                        // Make one word uppercase
		{"WORLD (low)", "world"},                       // Make one word lowercase
		{"hello world (cap)", "hello World"},           // Capitalize one word (the last one)
		{"go reloaded (cap, 2)", "Go Reloaded"},        // Capitalize 2 words
		{"this is amazing (up, 2)", "this IS AMAZING"}, // Uppercase 2 words
		{"MAKE IT lower (low, 1)", "MAKE IT lower"},    // Lowercase 1 word
	}

	// Test each case transformation
	for _, test := range tests {
		result := processor.ProcessText(test.input)
		if result != test.expected {
			t.Errorf("Case transform failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}

// Test function for punctuation spacing fixes
func TestPunctuationNormalization(t *testing.T) {
	// Test cases for fixing spaces around punctuation marks
	tests := []struct {
		input    string
		expected string
	}{
		{"hello , world", "hello, world"},     // Remove space before comma
		{"what ? is this", "what? is this"},   // Remove space before question mark
		{"amazing ! right", "amazing! right"}, // Remove space before exclamation
		{"hello.world", "hello. world"},       // Add space after period
		{"end.Start", "end. Start"},           // Add space after period before capital
	}

	// Test each punctuation fix
	for _, test := range tests {
		result := processor.ProcessText(test.input)
		if result != test.expected {
			t.Errorf("Punctuation normalization failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}

// Test function for quote spacing fixes
func TestQuoteProcessing(t *testing.T) {
	// Test cases for fixing spaces inside quotes
	tests := []struct {
		input    string
		expected string
	}{
		{"' hello world '", "'hello world'"},                 // Remove extra spaces inside quotes
		{"' test '", "'test'"},                               // Simple quote fix
		{"' multiple words here '", "'multiple words here'"}, // Multiple words in quotes
	}

	// Test each quote processing case
	for _, test := range tests {
		result := processor.ProcessText(test.input)
		if result != test.expected {
			t.Errorf("Quote processing failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}

// Test function for article correction (a vs an)
func TestArticleCorrection(t *testing.T) {
	// Test cases for changing "a" to "an" before vowel sounds
	tests := []struct {
		input    string
		expected string
	}{
		{"a apple", "an apple"},             // Change "a" to "an" before vowel
		{"a amazing day", "an amazing day"}, // Another vowel case
		{"a elephant", "an elephant"},       // Another vowel case
		{"a university", "a university"},    // Keep "a" for consonant sound (even though it starts with 'u')
	}

	// Test each article correction
	for _, test := range tests {
		result := processor.ProcessText(test.input)
		if result != test.expected {
			t.Errorf("Article correction failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}

// Test function for complete transformation using the README example
func TestCompleteTransformation(t *testing.T) {
	// This tests the entire pipeline with a complex example from the README
	input := "I have 1E (hex) apples and 10 (bin) oranges. it (cap) is a amazing day! ' hello world '"
	expected := "I have 30 apples and 2 oranges. It is an amazing day! 'hello world'"

	// Run the complete transformation
	result := processor.ProcessText(input)
	if result != expected {
		t.Errorf("Complete transformation failed.\nInput: %q\nExpected: %q\nGot: %q", input, expected, result)
	}
}

// Test function for edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	// Test unusual inputs that might cause problems
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},    // Empty string should return empty
		{"   ", ""}, // Only spaces should return empty
		{"no transformations needed", "no transformations needed"}, // Text with no special patterns
		{"invalid (xyz) tag", "invalid (xyz) tag"},                 // Invalid transformation tag should be ignored
	}

	// Test each edge case
	for _, test := range tests {
		result := processor.ProcessText(test.input)
		if result != test.expected {
			t.Errorf("Edge case failed.\nInput: %q\nExpected: %q\nGot: %q", test.input, test.expected, result)
		}
	}
}
