package pipeline

import (
	"strings"
)

// FormatPunctuation applies punctuation formatting rules:
// - groups of punctuation (e.g. "...", "!!", "?!") attach to previous token
// - single punctuation marks attach to previous token
// - handle single quotes via FixQuotes
// - trim double spaces
func FormatPunctuation(tokens []string) []string {
	// Attach any punctuation-only tokens (made of .,!?:;) to the previous token
	if len(tokens) == 0 {
		return tokens
	}

	var out []string
	for i := 0; i < len(tokens); i++ {
		// Read current token
		t := tokens[i]
		// If token consists only of punctuation and we have a previous token,
		// attach this punctuation to the end of the previous token (no space).
		if isPunctuationSequence(t) && len(out) > 0 {
			out[len(out)-1] = out[len(out)-1] + t
		} else {
			// Otherwise, keep token as-is.
			out = append(out, t)
		}
	}

	// Now run FixQuotes to collapse quoted tokens and then clean up spacing
	out = FixQuotes(out)
	for i := 0; i < len(out); i++ {
		// Collapse multiple spaces inside the token and trim edges.
		s := out[i]
		for strings.Contains(s, "  ") {
			s = strings.ReplaceAll(s, "  ", " ")
		}
		out[i] = strings.TrimSpace(s)
	}

	return out
}

// isPunctuationSequence returns true if the token consists only of the
// punctuation characters we care about: . , ! ? : ; (one or more times)
func isPunctuationSequence(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch r {
		case '.', ',', '!', '?', ':', ';':
			// allowed
		default:
			return false
		}
	}
	return true
}
