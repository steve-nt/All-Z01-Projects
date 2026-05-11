package pipeline

// Tokenize splits the input runes into tokens used by the pipeline.
// Rules:
// - whitespace splits tokens
// - punctuation sequences (.,!?:;) are returned as single tokens (e.g. "...", "!!")
// - single and double quotes, angle brackets returned as single tokens
// - markers like "(up)" or "(hex)" are kept as a single token (no internal split)
func Tokenize(input []rune) []string {
	var tokens []string
	n := len(input)
	i := 0

	// Loop over input runes; `i` is the current read index.

	// isPunct returns true for characters treated as punctuation
	// for grouping into punctuation-only tokens.
	isPunct := func(r rune) bool {
		switch r {
		case '.', ',', '!', '?', ':', ';':
			return true
		}
		return false
	}

	for i < n {
		// Read the current rune at index i
		r := input[i]

		// If current rune is whitespace, skip it and continue scanning.
		// Tokens are delimited by whitespace, so we do not emit space tokens.
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			i++
			continue
		}

			// markers like (up) or ( low, 3 ): capture everything until the next ')'
			// If we see a '(', try to capture the marker like '(up)' or '(hex)'.
			// We scan forward until a matching ')' and treat the whole
			// parentheses sequence as a single token (preserving inner spaces).
			if r == '(' {
				j := i + 1
				valid := false
				for j < n {
					if input[j] == ')' {
						valid = true
						break
					}
					j++
				}
				if valid {
					// Append the whole parentheses token and advance past it.
					tokens = append(tokens, string(input[i:j+1]))
					i = j + 1
					continue
				}
			}

		// punctuation sequences
		// Group runs of punctuation characters into a single token.
		if isPunct(r) {
			j := i + 1
			for j < n && isPunct(input[j]) {
				j++
			}
			tokens = append(tokens, string(input[i:j]))
			i = j
			continue
		}

		// single/double quote or angle bracket as single token
		// Single/double quotes and angle brackets are emitted as single-character tokens.
		if r == '\'' || r == '"' || r == '<' || r == '>' {
			tokens = append(tokens, string(r))
			i++
			continue
		}

		// otherwise collect a word until next separator
		// Otherwise, collect a word token until the next separator.
		// Separators include whitespace, parentheses, quotes, angle brackets, and punctuation.
		j := i
		for j < n {
			rr := input[j]
			if rr == ' ' || rr == '\n' || rr == '\t' || rr == '\r' {
				break
			}
			if rr == '(' || rr == ')' || rr == '\'' || rr == '"' || rr == '<' || rr == '>' || isPunct(rr) {
				break
			}
			j++
		}
		if j > i {
			// Append the substring [i:j] as a token and advance i.
			tokens = append(tokens, string(input[i:j]))
			i = j
		} else {
			// Fallback: emit single rune as a token to avoid infinite loop.
			tokens = append(tokens, string(input[i]))
			i++
		}
	}

	return tokens
}
