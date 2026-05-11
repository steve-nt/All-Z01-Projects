package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go-reloaded/pipeline"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input.txt> <output.txt>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read the entire input file into a string
	text, err := pipeline.ReadInput(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Split the input into lines
	lines := strings.Split(text, "\n")

	var results []string

	for _, line := range lines {
		// If the input line is empty, preserve an empty output line.
		if line == "" {
			results = append(results, "")
			continue
		}

		// Tokenize the line into pipeline tokens (words, punctuation, markers).
		tokens := pipeline.Tokenize([]rune(line))

		// Apply the standard sequence of transformations (articles, quotes,
		// replacements, case transforms, punctuation fixes, etc.).
		tokens = pipeline.ApplyTransformations(tokens)

		// Join tokens back into a single line string and append to results.
		result := pipeline.JoinTokens(tokens)
		results = append(results, result)
	}

	// Write each processed line to the output file
	err = pipeline.WriteOutput(outputFile, results)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Println("✅ File processed successfully:", outputFile)
}
