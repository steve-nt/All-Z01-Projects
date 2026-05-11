package pipeline

import "strings"

// ApplyTransformations applies the main sequence of transformations.
// Note: we intentionally do not call FormatPunctuation here so callers can
// decide when to format punctuation (some tests expect punctuation as separate tokens).
func ApplyTransformations(tokens []string) []string {
	// If there are no tokens, return an empty slice immediately.
	if len(tokens) == 0 {
		return []string{}
	}

	// First: fix articles like "a" -> "an" when appropriate.
	// We run this early so that articles inside quoted regions are corrected
	// before those regions might be merged into single tokens.
	tokens = FixArticles(tokens)

	// Second: collapse single-quote quoted sequences into a single token.
	// This merges tokens between matching single quotes into one quoted token.
	tokens = FixQuotes(tokens)

	// Third: perform numeric marker replacements. These expect stable token
	// boundaries (e.g., "101 (bin)" -> "5").
	tokens = ReplaceBin(tokens)
	tokens = ReplaceHex(tokens)

	// Fourth: apply case transformation markers such as (up), (low), (cap).
	// These helpers look for marker tokens and change surrounding tokens.
	tokens = ApplyCaseTransformations(tokens)

	// Fifth: normalize double-quoted token casing according to simple heuristics
	// (short content -> all caps, longer -> capitalize first rune, multi-word -> title-case).
	tokens = fixDoubleQuotedCase(tokens)

	// Sixth: squeeze comma + word + punctuation patterns into a tighter form,
	// e.g., [A "," B "!"] -> [A, B!]
	tokens = squeezeCommaBeforePunct(tokens)

	// Return the transformed token slice.
	return tokens
}

// fixDoubleQuotedCase transforms tokens like "hello" -> "Hello" or "hi" -> "HI".
// Heuristic: if inner content length <=2 -> uppercase entire content, else capitalize first rune.
func fixDoubleQuotedCase(tokens []string) []string {
	// Iterate over tokens and adjust tokens that are double-quoted strings.
	for i, t := range tokens {
		// Check token looks like a double-quoted value: starts and ends with '"'.
		if len(t) >= 2 && t[0] == '"' && t[len(t)-1] == '"' {
			// Extract inner content between the quotes.
			inner := t[1 : len(t)-1]

			// If inner is empty, nothing to do for this token.
			if inner == "" {
				continue
			}

			// If the quoted content contains spaces, title-case each inner word.
			if strings.Contains(inner, " ") {
				parts := strings.Fields(inner)
				for j, p := range parts {
					if len(p) == 0 {
						continue
					}
					runes := []rune(p)
					// Uppercase the first rune of the word.
					runes[0] = rune(strings.ToUpper(string(runes[0]))[0])
					// Lowercase the rest of the runes in the word.
					for k := 1; k < len(runes); k++ {
						runes[k] = rune(strings.ToLower(string(runes[k]))[0])
					}
					parts[j] = string(runes)
				}
				// Rebuild the token with quotes around the joined parts.
				tokens[i] = "\"" + strings.Join(parts, " ") + "\""
				continue
			}

			// For single-word quoted content: if length <= 2 runes, uppercase fully.
			if len([]rune(inner)) <= 2 {
				tokens[i] = "\"" + strings.ToUpper(inner) + "\""
			} else {
				// Otherwise capitalize the first rune and lowercase the rest.
				r := []rune(inner)
				first := string(r[0])
				rest := string(r[1:])
				tokens[i] = "\"" + strings.ToUpper(first) + strings.ToLower(rest) + "\""
			}
		}
	}
	return tokens
}

// squeezeCommaBeforePunct turns sequences like [A "," B "!"] into [A, B!]
func squeezeCommaBeforePunct(tokens []string) []string {
	if len(tokens) == 0 {
		return tokens
	}
	out := []string{}
	i := 0
	for i < len(tokens) {
		// need at least 4 tokens to match pattern
		if i+3 < len(tokens) && tokens[i+1] == "," && isPunctuationSequence(tokens[i+3]) {
			// keep tokens[i] as-is
			out = append(out, tokens[i])
			// combine tokens[i+2] and tokens[i+3]
			out = append(out, tokens[i+2]+tokens[i+3])
			i += 4
			continue
		}
		out = append(out, tokens[i])
		i++
	}
	return out
}
