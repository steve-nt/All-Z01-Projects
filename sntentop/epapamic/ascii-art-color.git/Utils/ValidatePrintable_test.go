package utils

import (
	"errors"
	"os"
	"testing"
)

func TestValidatePrintable(t *testing.T) {
	testCases := []struct {
		name        string
		inputMap    asciiMap
		userInput   string
		expected    string
		expectedErr error
	}{
		{
			name:        "Test for empty arguments",
			inputMap:    asciiMap{},
			userInput:   "",
			expected:    "",
			expectedErr: errors.New("arguments are empty"),
		},
		{
			name:        "Test for non-Printable characters",
			inputMap:    asciiMap{},
			userInput:   "🍆",
			expected:    "",
			expectedErr: errors.New("characters not within ASCII or Non-Printable"),
		},
		{
			name:        "Test for valid characters",
			inputMap:    asciiMap{},
			userInput:   "Test for valid characters",
			expected:    "Test for valid characters",
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			os.Args[1] = testCase.userInput
			err := testCase.inputMap.ValidatePrintable()
			if testCase.inputMap.input != testCase.expected {
				t.Errorf("Expected: %s, Got: %s", testCase.expected, testCase.inputMap.input)
			}
			if (err != nil) && (err.Error() != testCase.expectedErr.Error()) {
				t.Errorf("Error expected: %s, Error Got: %s", testCase.expectedErr, err)
			}
		})
	}
}

func TestIsPrintableByte(t *testing.T) {
	testCases := []struct {
		testName string
		input    rune
		expected bool
	}{
		{
			testName: "5",
			input:    '5',
			expected: true,
		},
		{
			testName: "F",
			input:    'F',
			expected: true,
		},
		{
			testName: "🤕",
			input:    '🤕',
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			result := isPrintableByte(testCase.input)
			if result != testCase.expected {
				t.Errorf("Expected: %v, Got: %v", testCase.expected, result)
			}
		})
	}
}
