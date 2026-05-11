package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"ascii-art-color/ascii"
)

var (
	colorFlag   string // Color to apply
	text2color  string // Text or substring to be colored
	outputFlag  string // Output file to save the result
	alignFlag   string // Text alignment
	reverseFlag string // Reverse option (not implemented yet)
	banner      string // Banner file to use
)

func main() {
	// Parse flags and options
	flag.StringVar(&colorFlag, "color", "", "Set the text color")
	flag.StringVar(&outputFlag, "output", "", "Output file")
	flag.StringVar(&alignFlag, "align", "left", "Text alignment (left, center, right, justify)")
	flag.StringVar(&reverseFlag, "reverse", "", "Reverse this file")
	flag.Parse()

	if !ascii.OnlyFlagsEqual() {
		return
	}

	// Validate alignFlag
	validAlignments := map[string]bool{
		"left":    true,
		"center":  true,
		"right":   true,
		"justify": true,
	}
	if !validAlignments[alignFlag] {
		fmt.Println("Error: Invalid alignment type. Allowed types are left, center, right, justify.")
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println("Example: go run . --align=justify something standard")
		return
	}

	// Get positional arguments
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Error: No input string provided.")
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println("EX: go run . --output=<fileName.txt> something standard")
		return
	}

	// Default banner
	banner = "standard.txt"

	// Initialize variables
	var inputText string
	var substringArg string

	// Check if the last argument is a valid banner
	if ascii.IsBanner(args[len(args)-1]) {
		banner = args[len(args)-1] + ".txt"
		args = args[:len(args)-1] // Remove the banner from args
	}

	// Now, inputText is the last argument
	inputText = args[len(args)-1]
	args = args[:len(args)-1] // Remove inputText from args

	inputText = ascii.UnescapeString(inputText)
	if substringArg != "" {
		substringArg = ascii.UnescapeString(substringArg)
	}

	// Process remaining args (substringArg)
	if len(args) > 1 {
		fmt.Println("Error: Too many positional arguments.")
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		return
	}

	if len(args) == 1 {
		substringArg = args[0]
	}

	// Validate input and files
	if !ascii.CheckValidity(inputText, banner, colorFlag) {
		return
	}

	// Read the font file into fontLines
	fontLines, err := ascii.ReadFontFile(banner)
	if err != nil {
		fmt.Println("Error reading font file:", err)
		return
	}

	// Prepare coloredIndices and colorCode
	var coloredIndices map[int]bool
	var colorCode string

	// If colorFlag is set
	if colorFlag != "" {
		// Determine text to color
		if substringArg != "" {
			text2color = substringArg
		} else {
			text2color = inputText
		}
		// Find indices
		substringIndices := ascii.FindSubstringIndices(inputText, text2color)
		coloredIndices = make(map[int]bool)
		for _, startIdx := range substringIndices {
			for i := startIdx; i < startIdx+len(text2color); i++ {
				coloredIndices[i] = true
			}
		}
		// Get color code
		colorCode, _ = ascii.ColorToAnsi(colorFlag)
	}

	// Get terminal width
	terminalWidth := ascii.GetTerminalWidth()

	if alignFlag != "left" {
		// Use PrintAsciiArtAlign for center, right, and justify alignments
		if outputFlag != "" {
			// If output to file is specified
			// Capture the output into a string
			var outputBuilder strings.Builder
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			ascii.PrintAsciiArtAlign(inputText, fontLines, alignFlag, terminalWidth, coloredIndices, colorCode)

			w.Close()
			os.Stdout = oldStdout
			io.Copy(&outputBuilder, r)
			output := outputBuilder.String()
			err := os.WriteFile(outputFlag, []byte(output), 0o644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			ascii.PrintAsciiArtAlign(inputText, fontLines, alignFlag, terminalWidth, coloredIndices, colorCode)
		}
	} else {
		// For left alignment, use GenerateAsciiArt
		asciiArtLines := ascii.GenerateAsciiArt(inputText, fontLines, coloredIndices, colorCode)

		// Join the lines into a single string and add an extra newline at the end
		output := strings.Join(asciiArtLines, "\n") + "\n"

		// Output the ASCII art lines to file or stdout
		if outputFlag != "" {
			err := os.WriteFile(outputFlag, []byte(output), 0o644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		} else {
			fmt.Print(output)
		}
	}
}
