package converts

import "strings"

func FixSingleQuoteSpacing(str string) string {
	lines := strings.Split(str, "\n")
	for k, line := range lines {
		singleQuotedTexts := extractSingleQuotedText(line)
		for _, text := range singleQuotedTexts {
			// Remove leading spaces
			textTrimmed := strings.Trim(text, "'")
			textTrimmed = strings.TrimSpace(textTrimmed)
			// Update the original line with the modified text
			line = strings.ReplaceAll(line, text, "'"+textTrimmed+"'")
		}
		lines[k] = line
	}
	return strings.Join(lines, "\n")
}
func extractSingleQuotedText(input string) []string {
	var results []string
	var result string
	quoteStarted := false
	for _, r := range input {
		if r == '\'' {
			if !quoteStarted {
				quoteStarted = true
				result += string(r)
			} else {
				quoteStarted = false
				result += string(r)
				results = append(results, result)
				result = ""
			}
		} else {
			if quoteStarted {
				result += string(r)
			}
		}
	}
	return results
}
