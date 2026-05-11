// File: colorutils.go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

// ANSI color codes for basic colors
var colorCodes = map[string]string{
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"reset":   "\033[0m",
}

// ParseColorFlag parses the --color=<color> flag and returns the color code
func ParseColorFlag(flag string) (string, bool) {
	colorPattern := regexp.MustCompile(`^--color=([a-zA-Z]+)$`)
	colorMatch := colorPattern.FindStringSubmatch(flag)

	if colorMatch == nil {
		return "", false
	}

	colorName := colorMatch[1]
	colorCode, ok := colorCodes[colorName]
	if !ok {
		fmt.Println("Error: Unsupported color. Supported colors are: red, green, yellow, blue, magenta, cyan.")
		return "", false
	}
	return colorCode, true
}

// ColorSubstring colors only the specified substring in the given text
func ColorSubstring(text, substring, colorCode string) string {
	if strings.Contains(text, substring) {
		return strings.Replace(text, substring, colorCode+substring+colorCodes["reset"], -1)
	}
	return text // If the substring isn't found, return the text unchanged
}

// PrintUsage prints the usage message for incorrect input
func PrintUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING]")
	fmt.Println("\nEX: go run . --color=<color> <substring to be colored> \"something\"")
}
