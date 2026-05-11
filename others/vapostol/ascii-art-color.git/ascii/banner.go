package ascii

import (
	"os"
	"strings"
)

// IsBanner checks if a string is a valid banner type
func IsBanner(s string) bool {
	return s == "standard" || s == "shadow" || s == "thinkertoy"
}

// ReadFontFile reads the font file into lines
func ReadFontFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	normalizedContent := strings.ReplaceAll(string(data), "\r", "")
	lines := strings.Split(normalizedContent, "\n")
	return lines, nil
}
