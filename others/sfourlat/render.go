package main

import (
	"fmt"
	"strings"
)

// ColorPair binds an ANSI escape code to the substring it should color.
// If Substr is empty, the entire output is colored with AnsiCode.
type ColorPair struct {
	AnsiCode string
	Substr   string
}

// Render prints each line segment as 8-row tall ASCII art.
// An empty segment prints a single blank line.
func Render(lines []string, bannerMap map[rune][]string) {
	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}
		var builder strings.Builder
		for row := 0; row < 8; row++ {
			builder.Reset()
			for _, char := range line {
				if block, ok := bannerMap[char]; ok {
					builder.WriteString(block[row])
				}
			}
			fmt.Println(builder.String())
		}
	}
}

// RenderWithColor prints ASCII art like Render, applying each ColorPair to its
// target substring. First pair wins when substrings overlap.
// Falls back to plain Render if pairs is empty.
func RenderWithColor(lines []string, bannerMap map[rune][]string, pairs []ColorPair) {
	if len(pairs) == 0 {
		Render(lines, bannerMap)
		return
	}

	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}

		runes := []rune(line)

		// Build one mask per pair.
		masks := make([][]bool, len(pairs))
		for p, pair := range pairs {
			masks[p] = BuildColorMask(runes, pair.Substr)
		}

		var builder strings.Builder
		for row := 0; row < 8; row++ {
			builder.Reset()
			currentColor := "" // ANSI code active at this point in the row

			for i, char := range runes {
				block, ok := bannerMap[char]
				if !ok {
					continue
				}

				// Find which color applies at this position (first match wins).
				activeColor := ""
				for p, pair := range pairs {
					if masks[p][i] {
						activeColor = pair.AnsiCode
						break
					}
				}

				// Emit transition only when the color actually changes.
				if activeColor != currentColor {
					if currentColor != "" {
						builder.WriteString(ansiReset)
					}
					if activeColor != "" {
						builder.WriteString(activeColor)
					}
					currentColor = activeColor
				}

				builder.WriteString(block[row])
			}

			// Close any open color at the end of the row.
			if currentColor != "" {
				builder.WriteString(ansiReset)
			}

			fmt.Println(builder.String())
		}
	}
}

// BuildColorMask returns a []bool of the same length as runes.
// true at position i means the character at i should be colored.
// If substr is empty, all positions are true.
// Matching is case-sensitive; overlapping matches are all marked.
func BuildColorMask(runes []rune, substr string) []bool {
	colored := make([]bool, len(runes))

	if substr == "" {
		for i := range colored {
			colored[i] = true
		}
		return colored
	}

	subRunes := []rune(substr)
	subLen := len(subRunes)

	for i := 0; i <= len(runes)-subLen; i++ {
		match := true
		for z := 0; z < subLen; z++ {
			if runes[i+z] != subRunes[z] {
				match = false
				break
			}
		}
		if match {
			for j := 0; j < subLen; j++ {
				colored[i+j] = true
			}
		}
	}

	return colored
}
