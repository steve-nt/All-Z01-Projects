package utils

import (
	"errors"
)

// Validates user input to be printable,
// returns non-nil error if its not.
func (m *asciiMap) ValidatePrintable() error {
	userInput := m.input
	if len(userInput) == 0 { // Check for empty input
		err := errors.New("arguments are empty")
		return err
	}
	for _, char := range userInput {
		if !isPrintableByte(char) { // Call helper function for character validation
			err := errors.New("characters not within ASCII or Non-Printable")
			return err
		}
	}
	m.input = userInput
	return nil
}

// Returns true if argument is printable, false if its not.
func isPrintableByte(b rune) bool {
	return b >= 32 && b <= 126 // Returns false of outside of range of printable characters
}
