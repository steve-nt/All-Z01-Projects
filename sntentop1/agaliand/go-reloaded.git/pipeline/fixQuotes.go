package pipeline

import (
	"strings"
)

// FixQuotes handles single quotes around words or phrases
// Example:
//
//	Input:  [I am ' awesome ']
//	Output: [I am 'awesome']
func FixQuotes(tokens []string) []string {
	openIndex := -1

	for i := 0; i < len(tokens); i++ {
		// When we encounter a single-quote token, either mark its position
		// as the opening quote or, if already open, treat as the closing quote.
		if tokens[i] == "'" {
			if openIndex == -1 {
				// Opening quote found: remember its index and continue scanning.
				openIndex = i
			} else {
				// Closing quote found at index i. Collect tokens between openIndex and i.
				var inner []string
				for j := openIndex + 1; j < i; j++ {
					// Trim any accidental spaces around inner tokens
					inner = append(inner, strings.TrimSpace(tokens[j]))
				}
				// Join inner tokens into a single space-separated string
				content := strings.Join(inner, " ")
				// Rebuild a single quoted token like 'content'
				quoted := "'" + content + "'"

				// Reconstruct token slice: tokens before openIndex, the quoted token,
				// and tokens after the closing quote.
				newTokens := make([]string, 0, len(tokens)-(i-openIndex))
				newTokens = append(newTokens, tokens[:openIndex]...)
				newTokens = append(newTokens, quoted)
				if i+1 < len(tokens) {
					newTokens = append(newTokens, tokens[i+1:]...)
				}

				tokens = newTokens

				// Reset scanning index and openIndex so scanning continues after the inserted quoted token.
				i = openIndex
				openIndex = -1
			}
		}
	}

	return tokens
}
