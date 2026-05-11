package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Ensure that the correct number of arguments are provided
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <inputFile> <outputFile>")
		return
	}

	// Read input and output file paths
	inputFile, outputFile := os.Args[1], os.Args[2]
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Apply transformations to the input text
	modifiedText := applyTransformations(string(inputData))

	// Remove the existing output file if it exists
	if _, err := os.Stat(outputFile); err == nil {
		if err = os.Remove(outputFile); err != nil {
			fmt.Printf("Error deleting existing output file: %v\n", err)
			return
		}
	}

	// Write the modified text to the output file
	if err = os.WriteFile(outputFile, []byte(modifiedText), 0o644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
	} else {
		fmt.Println("Text modifications completed successfully.")
	}
}

// applyTransformations applies a series of transformations to the input text
func applyTransformations(text string) string {
	transformations := []func(string) string{
		replaceHexWithDecimal,
		replaceBinWithDecimal,
		capitalizeWord,
		uppercaseWord,
		capitalizePreviousNWords,
		lowercasePreviousNWords,
		uppercasePreviousNWords,
		lowercaseWord,
		correctArticles,
		correctPunctuation,
		correctSingleQuotes,
	}

	for _, transform := range transformations {
		text = transform(text)
	}
	return text
}

// replaceHexWithDecimal replaces hexadecimal numbers followed by (hex) with their decimal equivalents
func replaceHexWithDecimal(text string) string {
	return regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`).ReplaceAllStringFunc(text, func(match string) string {
		hexNumber := regexp.MustCompile(`[0-9A-Fa-f]+`).FindString(match)
		if decimalValue, err := strconv.ParseInt(hexNumber, 16, 64); err == nil {
			return fmt.Sprintf("%d", decimalValue)
		}
		return match
	})
}

// replaceBinWithDecimal replaces binary numbers followed by (bin) with their decimal equivalents
func replaceBinWithDecimal(text string) string {
	return regexp.MustCompile(`\b([01]+)\s*\(bin\)`).ReplaceAllStringFunc(text, func(match string) string {
		binNumber := regexp.MustCompile(`[01]+`).FindString(match)
		if decimalValue, err := strconv.ParseInt(binNumber, 2, 64); err == nil {
			return fmt.Sprintf("%d", decimalValue)
		}
		return match
	})
}

// capitalizeWord capitalizes the first letter of a word followed by (cap)
func capitalizeWord(text string) string {
	return regexp.MustCompile(`\b([a-zA-Z]+)\s*\(cap\)`).ReplaceAllStringFunc(text, func(match string) string {
		return strings.Title(regexp.MustCompile(`[a-zA-Z]+`).FindString(match))
	})
}

// uppercaseWord converts a word followed by (up) to uppercase
func uppercaseWord(text string) string {
	return regexp.MustCompile(`\b(\w+)\s*\(up\)`).ReplaceAllStringFunc(text, func(match string) string {
		return strings.ToUpper(regexp.MustCompile(`\w+`).FindString(match))
	})
}

// capitalizePreviousNWords capitalizes the previous n words before (cap, n)
func capitalizePreviousNWords(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		if words[i] == "(cap," && i+1 < len(words) {
			n, err := strconv.Atoi(strings.TrimSuffix(words[i+1], ")"))
			if err != nil {
				continue
			}
			for j := i - 1; j >= 0 && n > 0; j-- {
				words[j] = strings.Title(words[j])
				n--
			}
			words = append(words[:i], words[i+2:]...)
			i--
		}
	}
	return strings.Join(words, " ")
}

// uppercasePreviousNWords converts the previous n words before (up, n) to uppercase
func uppercasePreviousNWords(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		if words[i] == "(up," && i+1 < len(words) {
			n, err := strconv.Atoi(strings.TrimSuffix(words[i+1], ")"))
			if err != nil {
				continue
			}
			for j := i - 1; j >= 0 && n > 0; j-- {
				words[j] = strings.ToUpper(words[j])
				n--
			}
			words = append(words[:i], words[i+2:]...)
			i--
		}
	}
	return strings.Join(words, " ")
}

// lowercasePreviousNWords converts the previous n words before (low, n) to lowercase
func lowercasePreviousNWords(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words); i++ {
		if words[i] == "(low," && i+1 < len(words) {
			n, err := strconv.Atoi(strings.TrimSuffix(words[i+1], ")"))
			if err != nil {
				continue
			}
			for j := i - 1; j >= 0 && n > 0; j-- {
				words[j] = strings.ToLower(words[j])
				n--
			}
			words = append(words[:i], words[i+2:]...)
			i--
		}
	}
	return strings.Join(words, " ")
}

// lowercaseWord converts a word followed by (low) to lowercase
func lowercaseWord(text string) string {
	return regexp.MustCompile(`\b(\w+)\s*\(low\)`).ReplaceAllStringFunc(text, func(match string) string {
		return strings.ToLower(regexp.MustCompile(`\w+`).FindString(match))
	})
}

// correctArticles replaces "a" with "an" if the following word starts with a vowel or 'h'
func correctArticles(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words)-1; i++ {
		if (words[i] == "a" || words[i] == "A") && strings.ContainsAny(string(words[i+1][0]), "aeiouAEIOUhH") {
			words[i] = "an"
		}
	}
	return strings.Join(words, " ")
}

// correctPunctuation removes misplaced spaces before punctuation marks and adds spaces where necessary
func correctPunctuation(text string) string {
	text = regexp.MustCompile(`(\S)\s*([.,!?;:])`).ReplaceAllString(text, `$1$2`)
	text = regexp.MustCompile(`([.,!?;:])\s*([.,!?;:])`).ReplaceAllString(text, `$1$2`)
	text = regexp.MustCompile(`([.,!?;:])([^\s.,!?;:'"\"])`).ReplaceAllString(text, `$1 $2`)
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.ReplaceAll(text, " ;", ";")
	text = strings.ReplaceAll(text, " :", ":")
	return strings.TrimSpace(text)
}

// correctSingleQuotes ensures single quotes are properly aligned with words
func correctSingleQuotes(text string) string {
	return regexp.MustCompile(`\s*'\s*(.*?)\s*'`).ReplaceAllString(text, `'${1}'`)
}
