package main

import (
	"strings"
	"testing"
)

// Helper function to compare expected and actual results in tests
func assertEqual(t *testing.T, got, want string) {
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		t.Errorf("Got:\n%s\nWant:\n%s", got, want)
	}
}

// Test the loading of ASCII art from the standard.txt file
func TestLoadAsciiArt(t *testing.T) {
	asciiMap, err := loadAsciiMap("standard.txt")
	if err != nil {
		t.Fatalf("Error loading ASCII art: %v", err)
	}

	// Check if the ASCII art for character 'A' is correct
	if len(asciiMap['A']) != CharHeight {
		t.Errorf("Expected 8 lines for character 'A', got %d", len(asciiMap['A']))
	}
}

// Test basic ASCII art conversion
func TestConvertToAsciiArt(t *testing.T) {
	asciiMap, _ := loadAsciiMap("standard.txt")

	// Test a simple word
	result := convertToAsciiArt("HI", asciiMap)
	expected := "EXPECTED ASCII ART FOR 'HI'"
	assertEqual(t, result, expected)
}

// Test handling of newline characters
func TestConvertToAsciiArtWithNewline(t *testing.T) {
	asciiMap, _ := loadAsciiMap("standard.txt")

	// Test input with newlines
	result := convertToAsciiArt("HELLO\nWORLD", asciiMap)
	expected := "EXPECTED ASCII ART FOR 'HELLO'\nEXPECTED ASCII ART FOR 'WORLD'"
	assertEqual(t, result, expected)
}

// Test handling of unsupported characters
func TestConvertToAsciiArtUnsupportedCharacters(t *testing.T) {
	asciiMap, _ := loadAsciiMap("standard.txt")

	// Test with an unsupported character (should be blank spaces in ASCII)
	result := convertToAsciiArt("HELLO@", asciiMap)
	expected := "EXPECTED ASCII ART FOR 'HELLO' WITH BLANK FOR '@'"
	assertEqual(t, result, expected)
}
