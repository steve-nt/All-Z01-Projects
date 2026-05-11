package pipeline

import (
	"os"
	"strings"
)

func JoinTokens(tokens []string) string {
	// Join token slice into a single space-separated string (used by callers
	// that need a single-line representation of tokens).
	return strings.Join(tokens, " ")
}
func WriteOutput(filename string, input []string) error {
	// Write each element of the input slice as its own line in the output file.
	// This preserves original line boundaries from the processing pipeline.
	output := strings.Join(input, "\n") + "\n"
	return os.WriteFile(filename, []byte(output), 0o644)
}
