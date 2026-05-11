package pipeline // Defines the package name

import (
	"strings" // For string operations

	"golang.org/x/text/cases"    // For proper Unicode-aware capitalization
	"golang.org/x/text/language" // For language settings
)

// FixArticles converts "a" to "an" if the next word starts with a vowel or 'h'.
// Preserves capitalization of the original article.
func FixArticles(words []string) []string {
	if len(words) == 0 { // Return empty slice if input is empty
		return words
	}

	result := make([]string, len(words)) // Prepare output slice
	copy(result, words)                  // Copy input to result

	// caser used to preserve capitalization (e.g., "A" -> "An")
	caser := cases.Title(language.English)

	// Iterate through tokens except the last, since we look ahead one token.
	for i := 0; i < len(result)-1; i++ {
		original := result[i]             // Current token as-is
		word := strings.ToLower(original) // Lowercase version for equality checks
		next := result[i+1]               // Peek at the next token

		// Only change plain "a" (case-insensitive)
		if word == "a" {
			article := "a" // default replacement

			// If the next token begins with a vowel or 'h', choose "an" instead
			if startsWithVowelOrH(next) {
				article = "an"
			}

			// Preserve capitalization: if original started with uppercase, title-case the article
			if len(original) > 0 && original[0] >= 'A' && original[0] <= 'Z' {
				article = caser.String(article)
			}

			// Store the possibly-updated article back into the result slice
			result[i] = article
		}
	}

	return result
}

// startsWithVowelOrH checks if a word starts with a vowel or 'h'
func startsWithVowelOrH(s string) bool {
	if s == "" {
		return false
	}
	// trim any leading punctuation or quotes
	s = strings.TrimLeftFunc(s, func(r rune) bool {
		// trim spaces and common punctuation/quotes
		switch r {
		case ' ', '"', '\'', '(', ')', ',', '.', '!', '?', ':', ';', '<', '>':
			return true
		}
		return false
	})
	if s == "" {
		return false
	}
	// Extract the first rune of the trimmed string and lowercase it for the vowel check.
	first := strings.ToLower(string([]rune(s)[0]))
	return strings.Contains("aeiouh", first)
}
