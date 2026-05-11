package converts

import (
	"strings"
)

func ProcessString(str string) string {
	lines := strings.Split(str, "\n")
	for k, line := range lines {
		line = removeSpacesFromBrackets(line)
		words := strings.Split(line, " ")
		var result []string
		i := 0
		for i < len(words) {
			word := words[i]
			switch {
			case checkCommandPrefix(word):
				command, num, extra := handleCommand(word)
				result = executeCommand(result, command, num, extra)
			case word == "(hex)":
				// Process hexadecimal conversion
				result = processHexCommand(result)
			case word == "(bin)":
				// Process binary conversion
				result = processBinCommand(result)
			case word == "a" || word == "A" || word == "an" || word == "An":
				// Process corrections for "a" and "an"
				result = processArticleCorrection(result, word, words, i)
			case len(word) > 0 && isPunctuation(rune(word[0])): // Check if word is non-empty
				result = processPunctuations(result, word)
			default:
				// Add word as is if no command is applicable
				result = append(result, word)
			}
			i++
		}
		lines[k] = strings.TrimSpace(strings.Join(result, " "))
	}
	return strings.Join(lines, "\n")
}

func removeSpacesFromBrackets(str string) string {
	var result string
	foundBracket := false
	for _, r := range str {
		if r == '(' {
			foundBracket = true
		} else if r == ')' {
			foundBracket = false
		}
		if !(foundBracket && r == ' ') {
			result += string(r)
		}
	}
	return result
}
