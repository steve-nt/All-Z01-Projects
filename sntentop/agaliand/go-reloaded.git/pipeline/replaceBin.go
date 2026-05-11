package pipeline // This file belongs to the "pipeline" package

import (
	"fmt"     // Used for debug printing if needed
	"strconv" // Needed to parse binary strings to integers
	"strings" // Needed for string comparison
)

// ReplaceBin scans tokens for "(bin)" markers and converts the previous token
// from binary (base 2) to decimal (base 10). Invalid binaries are ignored.
func ReplaceBin(tokens []string) []string {
	// Result will accumulate tokens, replacing binary tokens as we go.
	var result []string

	// Walk through each token in the input slice.
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// If token is the case-insensitive marker "(bin)", attempt conversion.
		if strings.EqualFold(token, "(bin)") {
			// Only attempt conversion if we already have a previous token.
			if len(result) > 0 {
				binWord := result[len(result)-1] // Candidate binary string

				// Parse the candidate as a base-2 integer.
				value, err := strconv.ParseInt(binWord, 2, 64)
				if err == nil {
					// Successful conversion: replace the previous token with decimal string.
					result[len(result)-1] = fmt.Sprintf("%d", value)
				}
				// If parsing fails, keep the previous token unchanged.
			}
			// Do not append the marker itself to the result.
			continue
		}

		// Normal token: append as-is.
		result = append(result, token)
	}

	// Return transformed token slice.
	return result
}
