// This tells Go that this file belongs to the "processor" package
package processor

// Import statements - these bring in external libraries we need
import (
	"os"      // Library for file operations
	"regexp"  // Library for pattern matching (regular expressions) "magnifying glass to find patterns"
	"strconv" // Library for converting strings to numbers and vice versa "translator between text and numbers"
	"strings" // Library for working with text strings "scissors and glue for text"
)

// Main function that processes the entire text through all transformation steps
func ProcessText(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		line = convertNumbers(line)
		line = applyCaseTransforms(line)
		line = normalizePunctuation(line)
		line = processQuotes(line)
		line = fixArticles(line)
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}

// Function to convert hexadecimal and binary numbers to decimal
func convertNumbers(text string) string {
	// Convert hex: "1E (hex)" -> "30"
	hexRe := regexp.MustCompile(`([0-9A-Fa-f]+)\s*\(hex\)`)             // Pattern to find hex numbers like "FF (hex)" - + at least one, \s matches any whitespace(space,tab,newline), * zero or more repetitions, \(hex\)-literally match parentheses)
	text = hexRe.ReplaceAllStringFunc(text, func(match string) string { // For each match found, run this function
		hexStr := hexRe.FindStringSubmatch(match)[1]                  // Extract just the hex number part (without "(hex)")
		if val, err := strconv.ParseInt(hexStr, 16, 64); err == nil { // Try to convert hex to decimal (base 16 to base 10)
			return strconv.FormatInt(val, 10) // If successful, return the decimal number as text
		}
		return match // If conversion failed, return original text unchanged
	})

	// Convert bin: "10 (bin)" -> "2"
	binRe := regexp.MustCompile(`([01]+)\s*\(bin\)`)                    // Pattern to find binary numbers like "101 (bin)"
	text = binRe.ReplaceAllStringFunc(text, func(match string) string { // For each match found, run this function
		binStr := binRe.FindStringSubmatch(match)[1]                 // Extract just the binary number part (without "(bin)")
		if val, err := strconv.ParseInt(binStr, 2, 64); err == nil { // Try to convert binary to decimal (base 2 to base 10)
			return strconv.FormatInt(val, 10) // If successful, return the decimal number as text
		}
		return match // If conversion failed, return original text unchanged
	})

	return text // Return the text with all number conversions completed
}

// Function to handle case transformations like (up), (low), (cap)
func applyCaseTransforms(text string) string {
	// Handle (up), (low), (cap) with optional word count
	// (up) → uppercase only the one previous word
	reUpSingle := regexp.MustCompile(`(\S+)\s*\(up\)`)

	// (up, n) → uppercase the previous n words
	reUpMulti := regexp.MustCompile(`((?:\S+\s+){0,}?\S+)\s*\(up,\s*(\d+)\)`)

	// Handle (up, n)
	text = reUpMulti.ReplaceAllStringFunc(text, func(match string) string {
		parts := reUpMulti.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		wordsSection := parts[1]
		n, _ := strconv.Atoi(parts[2])

		wordList := strings.Fields(wordsSection)
		if len(wordList) < n {
			n = len(wordList)
		}

		start := len(wordList) - n
		for i := start; i < len(wordList); i++ {
			wordList[i] = strings.ToUpper(wordList[i])
		}

		return strings.Join(wordList, " ")
	})

	// Handle (up)
	text = reUpSingle.ReplaceAllStringFunc(text, func(match string) string {
		parts := reUpSingle.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		word := parts[1]
		return strings.ToUpper(word)
	})

	lowRe := regexp.MustCompile(`(\S+(?:\s+\S+)*)\s*\(low(?:,\s*(\d+))?\)`) // Pattern for "word (low)" or "word (low, 2)"
	text = lowRe.ReplaceAllStringFunc(text, func(match string) string {     // For each lowercase command found
		matches := lowRe.FindStringSubmatch(match) // Split the match into parts
		words := strings.Fields(matches[1])        // Split the words before (low) into a list
		count := len(words)                        // By default, transform all words
		if matches[2] != "" {                      // If a number was specified like (low, 2)
			if c, err := strconv.Atoi(matches[2]); err == nil && c < len(words) { // Convert the number from text to integer
				count = c // Only transform this many words
			}
		}
		for i := 0; i < count; i++ { // Loop through the words to transform
			words[len(words)-1-i] = strings.ToLower(words[len(words)-1-i]) // Make each word lowercase (starting from the end)
		}
		return strings.Join(words, " ") // Put the words back together with spaces
	})

	// (cap) – capitalize only the previous ONE word
	reCapSingle := regexp.MustCompile(`(\S+)\s*\(cap\)`)

	// (cap, n) – capitalize the previous n words
	reCapMulti := regexp.MustCompile(`((?:\S+\s+){0,}?\S+)\s*\(cap,\s*(\d+)\)`)

	// Handle (cap, n)
	text = reCapMulti.ReplaceAllStringFunc(text, func(match string) string {
		parts := reCapMulti.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		wordsSection := parts[1]
		n, _ := strconv.Atoi(parts[2])

		wordList := strings.Fields(wordsSection)
		if len(wordList) < n {
			n = len(wordList)
		}

		start := len(wordList) - n
		for i := start; i < len(wordList); i++ {
			w := wordList[i]
			if len(w) > 0 {
				wordList[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
			}
		}

		return strings.Join(wordList, " ")
	})

	// Handle (cap)
	text = reCapSingle.ReplaceAllStringFunc(text, func(match string) string {
		parts := reCapSingle.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}

		w := parts[1]
		if len(w) == 0 {
			return match
		}

		return strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	})

	return text // Return the text with all case transformations completed
}

// Function to fix spacing around punctuation marks
func normalizePunctuation(text string) string {
	// Remove spaces before punctuation
	text = regexp.MustCompile(`\s+([,.!?;:])`).ReplaceAllString(text, "$1") // "hello , world" becomes "hello, world"
	// Add space after punctuation if followed by letter
	text = regexp.MustCompile(`([,.!?;:])([a-zA-Z])`).ReplaceAllString(text, "$1 $2") // "hello.world" becomes "hello. world"
	return text                                                                       // Return the text with fixed punctuation spacing
}

// Function to fix spacing inside quotes
func processQuotes(text string) string {
	// Handle single quotes: "' hello world '" -> "'hello world'"
	quoteRe := regexp.MustCompile(`'\s*([^']+?)\s*'`) // Pattern to find text inside single quotes
	text = quoteRe.ReplaceAllString(text, "'$1'")     // Remove extra spaces inside the quotes
	return text                                       // Return the text with fixed quote spacing
}

// Function to change "a" to "an" before vowels (a, e, i, o, u) or 'h'
// Only capitalizes "An" at the start of a sentence
func fixArticles(text string) string {
	re := regexp.MustCompile(`([.!?]\s*)([aA])\s+([aeiouhAEIOUH])`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)
		punct := parts[1] // punctuation + spaces before
		// always capitalize "A" at start of sentence
		return punct + "An " + parts[3]
	})
}

// ReadFile reads the content of a file and returns it as a string
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func WriteFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
