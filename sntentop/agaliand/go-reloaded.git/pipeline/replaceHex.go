package pipeline

import (
	"fmt"
	"strconv"
	"strings"
)

// ReplaceHex scans through all tokens and whenever it finds "(hex)",
// it replaces the *previous word* (which is always a hexadecimal number)
// with its decimal equivalent.
func ReplaceHex(tokens []string) []string {
	// Accumulate transformed tokens in result.
	var result []string

	// Walk input tokens sequentially.
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// If current token is the "(hex)" marker, attempt to convert the previous token.
		if strings.EqualFold(token, "(hex)") {
			if len(result) > 0 {
				hexWord := result[len(result)-1] // Candidate hex string

				// Parse the candidate as a base-16 integer.
				value, err := strconv.ParseInt(hexWord, 16, 64)
				if err == nil {
					// On success, replace the previous token with its decimal string.
					result[len(result)-1] = fmt.Sprintf("%d", value)
				}
				// On parse failure, leave previous token unchanged.
			}
			// Do not append the marker itself to the output.
			continue
		}

		// Normal token: append unchanged.
		result = append(result, token)
	}

	// Return the resulting token slice after conversions.
	return result
}
