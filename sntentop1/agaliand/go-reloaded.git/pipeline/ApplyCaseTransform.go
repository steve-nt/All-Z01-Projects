package pipeline // Declares that this file belongs to the "pipeline" package

import (
	"strconv" // Used for converting string numbers to integers
	"strings" // Provides utilities for splitting and formatting text
	"unicode" // Provides rune-level character transformations (upper/lower)
)

// ApplyCaseTransformations processes the token list and applies case-transform rules.
// Rules are encoded in tokens like: (up), (low), (cap), (up, N), (low, N), (cap, N)
func ApplyCaseTransformations(tokens []string) []string {
	// Create an output slice that will accumulate non-marker tokens and
	// mutated tokens produced by applying markers.
	var result []string

	// Iterate over each token in the incoming slice.
	for i := 0; i < len(tokens); i++ {
		token := tokens[i] // Current token under inspection

		// If this token looks like a parenthesized marker (e.g. "(up)"), handle it
		if strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")") {

			// Strip parentheses to get the marker content.
			content := token[1 : len(token)-1]

			// Marker may include a comma and a count, like "up, 3".
			parts := strings.Split(content, ",")
			action := strings.TrimSpace(parts[0]) // e.g. "up", "low", "cap"
			count := 1                            // default affected token count

			// If a count is provided, attempt to parse it as an integer.
			if len(parts) == 2 {
				if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					count = n
				}
			}

			// Apply the requested transformation to the previously collected tokens
			// in the result slice (we operate on the last `count` tokens).
			switch strings.ToLower(action) {
			case "up":
				applyUppercase(result, count)
			case "low":
				applyLowercase(result, count)
			case "cap":
				applyCapitalize(result, count)
			default:
				// Unknown marker: ignore it silently.
			}

			// Skip appending the marker token itself to the output.
			continue
		}

		// Not a marker: append the token to the output slice.
		result = append(result, token)
	}

	// Return the transformed tokens.
	return result
}

// applyUppercase converts the last <count> tokens of the slice to uppercase.
func applyUppercase(tokens []string, count int) {
	start := max(0, len(tokens)-count) // Determine the safe starting index
	for i := start; i < len(tokens); i++ {
		tokens[i] = strings.ToUpper(tokens[i])
	}
}

// applyLowercase converts the last <count> tokens of the slice to lowercase.
func applyLowercase(tokens []string, count int) {
	start := max(0, len(tokens)-count)
	for i := start; i < len(tokens); i++ {
		tokens[i] = strings.ToLower(tokens[i])
	}
}

// applyCapitalize capitalizes the last <count> tokens.
// Capitalizing means: first letter uppercase, remaining letters lowercase.
func applyCapitalize(tokens []string, count int) {
	start := max(0, len(tokens)-count)
	for i := start; i < len(tokens); i++ {
		tokens[i] = capitalizeWord(tokens[i])
	}
}

// capitalizeWord transforms a single token so that:
// - the first character is uppercase
// - the remaining characters are lowercase
func capitalizeWord(word string) string {
	if len(word) == 0 {
		return word // No operation on empty tokens
	}

	// Convert the string to runes for Unicode-correct operations.
	runes := []rune(word)

	// Uppercase first character and lowercase remaining characters.
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}

	// Convert back to string and return.
	return string(runes)
}

// max returns the larger of two integers.
// Used to avoid negative slice indexing.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
