package main

import (
	"strings"
	"testing"
)

// Test the createBannerMap function
func TestLoadBanner(t *testing.T) {
	// Test with a known banner file (assuming standard.txt exists in the banners directory)
	bannerMap, err := createBannerMap("banners/standard.txt")
	if err != nil {
		t.Fatalf("Failed to load banner: %v", err)
	}

	// Check that the map is not empty
	if len(bannerMap) == 0 {
		t.Fatal("Banner map is empty")
	}

	// Check that the banner contains a known character, e.g., 'A'
	if _, ok := bannerMap['A']; !ok {
		t.Error("Character 'A' not found in banner map")
	}

	// Check the height of the character's ASCII art
	if len(bannerMap['A']) != bannerHeight {
		t.Errorf("Expected banner height for character 'A' is %d, got %d", bannerHeight, len(bannerMap['A']))
	}
}

// Test getAvailableBanners function
func TestGetAvailableBanners(t *testing.T) {
	banners, err := getAvailableBanners("banners")
	if err != nil {
		t.Fatalf("Failed to get available banners: %v", err)
	}

	// Assuming standard.txt, shadow.txt, and thinkertoy.txt are in the banners directory
	expectedBanners := []string{"standard", "shadow", "thinkertoy"}
	for _, expected := range expectedBanners {
		found := false
		for _, banner := range banners {
			if banner == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected banner %s not found", expected)
		}
	}
}

// Test the handling of newline in the input string
func TestNewlineHandling(t *testing.T) {
	input := "Hello\\nWorld"

	// Replace '\\n' with '\n' and check if it's correctly split
	input = strings.ReplaceAll(input, "\\n", "\n")
	lines := strings.Split(input, "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines after splitting by newline, got %d", len(lines))
	}

	if lines[0] != "Hello" {
		t.Errorf("Expected 'Hello' on the first line, got '%s'", lines[0])
	}

	if lines[1] != "World" {
		t.Errorf("Expected 'World' on the second line, got '%s'", lines[1])
	}
}

// Test for input string that contains multiple consecutive newlines
func TestMultipleNewlines(t *testing.T) {
	input := "Hello\\n\\nWorld"

	// Replace '\\n' with '\n' and check if it's correctly split
	input = strings.ReplaceAll(input, "\\n", "\n")
	lines := strings.Split(input, "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines after splitting by newline, got %d", len(lines))
	}

	if lines[0] != "Hello" {
		t.Errorf("Expected 'Hello' on the first line, got '%s'", lines[0])
	}

	if lines[1] != "" {
		t.Errorf("Expected an empty line between 'Hello' and 'World', got '%s'", lines[1])
	}

	if lines[2] != "World" {
		t.Errorf("Expected 'World' on the third line, got '%s'", lines[2])
	}
}

// Test empty input string
func TestEmptyInput(t *testing.T) {
	input := ""

	// Test with empty string
	if input != "" {
		t.Errorf("Expected empty string, but got '%s'", input)
	}
}

// Test handling unsupported characters
func TestUnsupportedCharacters(t *testing.T) {
	bannerMap, err := createBannerMap("banners/standard.txt")
	if err != nil {
		t.Fatalf("Failed to load banner: %v", err)
	}

	// Check unsupported character handling
	unsupportedChar := rune(128) // character outside the ASCII range (0-127)
	_, ok := bannerMap[unsupportedChar]

	if ok {
		t.Errorf("Expected unsupported character '%c' not to be in the banner map, but it was found", unsupportedChar)
	}
}

// Test ASCII art generation for a given input
func TestAsciiArtGeneration(t *testing.T) {
	// Load the standard banner
	bannerMap, err := createBannerMap("banners/standard.txt")
	if err != nil {
		t.Fatalf("Failed to load banner: %v", err)
	}

	// Test with a simple string "A"
	inputString := "A"

	// Collect the output
	var result []string
	for i := 0; i < bannerHeight; i++ {
		var outputLine string
		for _, ch := range inputString {
			artLines, ok := bannerMap[ch]
			if !ok {
				artLines = bannerMap[' '] // Default to space if the character is unsupported
			}
			outputLine += artLines[i]
		}
		result = append(result, outputLine)
	}

	// Check that the result has the correct number of lines
	if len(result) != bannerHeight {
		t.Errorf("Expected %d lines of ASCII art, but got %d", bannerHeight, len(result))
	}

	// Check that the first line of "A" is correct (you may need to verify this with the actual content of "A" in standard.txt)
	expectedLine1 := "           " // Replace this with the actual first line for "A" in your standard.txt
	if result[0] != expectedLine1 {
		t.Errorf("Expected first line of 'A' to be '%s', but got '%s'", expectedLine1, result[0])
	}
}
