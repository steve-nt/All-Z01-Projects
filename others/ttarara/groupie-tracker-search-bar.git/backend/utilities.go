package backend

import (
	"strings"
	"unicode"
)

// FormatLocation transforms locations like "west_melbourne-australia"
// into "West Melbourne, Australia"
func FormatLocation(location string) string {
	// Replace underscores with spaces
	location = strings.ReplaceAll(location, "_", " ")
	// Replace hyphens with commas + space
	location = strings.ReplaceAll(location, "-", ", ")

	// Capitalize each word
	words := strings.Fields(location)
	for i, word := range words {
		words[i] = capitalize(word)
	}
	return strings.Join(words, " ")
}

// capitalize returns a string with the first letter uppercase
func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
