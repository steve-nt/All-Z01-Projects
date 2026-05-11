package ascii

import (
	"fmt"
	"strconv"
	"strings"
)

// UnescapeString unescapes a string with escape sequences
func UnescapeString(s string) string {
	unescaped, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return s // Return the original string if unquoting fails
	}
	return unescaped
}

// GenerateAsciiArt generates the ASCII art representation of the input text
func GenerateAsciiArt(text string, fontLines []string, coloredIndices map[int]bool, colorCode string) []string {
	if text == "" {
		return nil
	}

	var outputLines []string
	asciiArtLines := make([]string, 8) // Holds the 8 lines of ASCII art for the current text line
	charIndex := 0

	for _, char := range text {
		if char == '\n' {
			// Append the current set of ASCII art lines to outputLines
			outputLines = append(outputLines, asciiArtLines...)
			// Reset asciiArtLines for the next line
			asciiArtLines = make([]string, 8)
			charIndex++
			continue
		}

		if char < 32 || char > 126 {
			continue // Skip non-printable characters
		}
		for row := 1; row <= 8; row++ {
			asciiIndex := (int(char)-32)*9 + row
			if asciiIndex < len(fontLines) {
				asciiCharLine := fontLines[asciiIndex]
				if len(asciiCharLine) == 0 {
					// For space character or empty ASCII art lines, replace with a fixed number of spaces
					asciiCharLine = "     " // Adjust the number of spaces as needed
				}
				if coloredIndices != nil && coloredIndices[charIndex] && colorCode != "" {
					resetCode := AnsiColors["reset"]
					asciiCharLine = colorCode + asciiCharLine + resetCode
				}
				asciiArtLines[row-1] += asciiCharLine
			}
		}
		charIndex++
	}
	// Append the last set of ASCII art lines to outputLines
	outputLines = append(outputLines, asciiArtLines...)
	return outputLines
}

// Function to strip ANSI codes without using regexp
func stripAnsiCodes(s string) string {
	var result []rune
	inEscape := false
	i := 0
	n := len(s)
	for i < n {
		if s[i] == '\x1b' {
			inEscape = true
			i++
			if i < n && s[i] == '[' {
				i++
				// Skip until 'm'
				for i < n && s[i] != 'm' {
					i++
				}
				if i < n {
					i++ // Skip 'm'
				}
				inEscape = false
				continue
			}
		}
		if !inEscape {
			result = append(result, rune(s[i]))
			i++
		} else {
			i++
		}
	}
	return string(result)
}

// PrintAsciiArtAlign handles alignments with color support
func PrintAsciiArtAlign(text string, fontLines []string, position string, terminalWidth int, coloredIndices map[int]bool, colorCode string) {
	lines := strings.Split(text, "\n")
	for _, lineText := range lines {
		if lineText == "" {
			fmt.Println() // Always print a newline for empty lines
			continue
		}

		// Split the line into parts (words and spaces)
		parts := splitTextPreservingSpaces(lineText)

		// Build the ASCII art lines for each part
		linesPerPart := make([][]string, len(parts)) // Holds ASCII art lines for each part
		widthsPerPart := make([]int, len(parts))     // Holds the maximum width of each part

		charIndex := 0

		for i, part := range parts {
			partLines := make([]string, 8) // There are 8 rows in ASCII art

			for _, char := range part {
				if char < 32 || char > 126 {
					continue // Skip non-printable characters
				}
				for row := 1; row <= 8; row++ {
					asciiIndex := (int(char)-32)*9 + row
					if asciiIndex < len(fontLines) {
						asciiCharLine := fontLines[asciiIndex]
						if len(asciiCharLine) == 0 {
							// For space character or empty ASCII art lines, replace with spaces
							asciiCharLine = "     " // Adjust the number of spaces as needed
						}
						if coloredIndices != nil && coloredIndices[charIndex] && colorCode != "" {
							resetCode := AnsiColors["reset"]
							asciiCharLine = colorCode + asciiCharLine + resetCode
						}
						partLines[row-1] += asciiCharLine
					}
				}
				charIndex++
			}

			// Calculate the maximum width of the part
			maxWidth := 0
			for _, line := range partLines {
				strippedLine := stripAnsiCodes(line)
				lineWidth := len(strippedLine)
				if lineWidth > maxWidth {
					maxWidth = lineWidth
				}
			}

			// Pad each line of the part to the maximum width
			for row := 0; row < 8; row++ {
				strippedLine := stripAnsiCodes(partLines[row])
				lineWidth := len(strippedLine)
				if lineWidth < maxWidth {
					padding := strings.Repeat(" ", maxWidth-lineWidth)
					partLines[row] += padding
				}
			}

			linesPerPart[i] = partLines
			widthsPerPart[i] = maxWidth
		}

		// Now, for each of the 8 lines, combine parts and apply alignment
		for row := 0; row < 8; row++ {
			lineParts := make([]string, len(parts))
			for i, partLines := range linesPerPart {
				lineParts[i] = partLines[row]
			}
			alignedLine := applyAlignment(lineParts, widthsPerPart, position, terminalWidth)
			fmt.Println(alignedLine)
		}
	}
}

// splitTextPreservingSpaces splits text into words and spaces, preserving spaces as separate elements
func splitTextPreservingSpaces(text string) []string {
	var parts []string
	var currentPart strings.Builder
	var isSpace bool

	for _, char := range text {
		if char == ' ' {
			if !isSpace && currentPart.Len() > 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
			currentPart.WriteRune(char)
			isSpace = true
		} else {
			if isSpace && currentPart.Len() > 0 {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
			currentPart.WriteRune(char)
			isSpace = false
		}
	}
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}

// applyAlignment aligns a line based on the specified position
func applyAlignment(lineParts []string, widthsPerPart []int, position string, terminalWidth int) string {
	// Calculate the total length of the line
	lineLength := 0
	for _, width := range widthsPerPart {
		lineLength += width
	}

	if lineLength >= terminalWidth || position == "left" {
		return strings.Join(lineParts, "")
	}

	switch position {
	case "center":
		space := terminalWidth - lineLength
		leftPadding := space / 2
		return strings.Repeat(" ", leftPadding) + strings.Join(lineParts, "")
	case "right":
		space := terminalWidth - lineLength
		return strings.Repeat(" ", space) + strings.Join(lineParts, "")
	case "justify":
		// Distribute spaces between words (non-space parts)
		return justifyAsciiLine(lineParts, widthsPerPart, terminalWidth)
	default:
		return strings.Join(lineParts, "")
	}
}

// justifyAsciiLine justifies an ASCII art line to fit the terminal width
func justifyAsciiLine(lineParts []string, widthsPerPart []int, terminalWidth int) string {
	// Calculate the total length of all parts (excluding spaces between them)
	lineLen := 0
	for _, width := range widthsPerPart {
		lineLen += width
	}

	// Calculate the total spaces to distribute
	totalSpaces := terminalWidth - lineLen
	if totalSpaces <= 0 {
		return strings.Join(lineParts, "") // No justification needed
	}

	// Identify gaps (non-space parts) where extra spaces can be distributed
	gapCount := len(lineParts) - 1 // Gaps are between words
	if gapCount <= 0 {
		return strings.Join(lineParts, "") // No gaps, return the line as is
	}

	// Calculate space per gap and remaining spaces
	spacesPerGap := totalSpaces / gapCount
	extraSpaces := totalSpaces % gapCount

	// Build the justified line
	var justifiedLine strings.Builder
	for i, part := range lineParts {
		justifiedLine.WriteString(part) // Add the current part
		if i < gapCount {               // Add spaces after each part except the last one
			justifiedLine.WriteString(strings.Repeat(" ", spacesPerGap))
			if extraSpaces > 0 { // Distribute extra spaces to the first few gaps
				justifiedLine.WriteString(" ")
				extraSpaces--
			}
		}
	}

	return justifiedLine.String()
}
