// ascii/ascii_printer.go
package ascii

import (
	"strings"
)

func SplitByNumber(arr []string, start int, end int) []string {
	var result []string
	for i := start; i < end; i++ {
		result = append(result, arr[i])
	}
	return result
}

func ReturnString2EndlineArray(s string) []string {
	var substrings []string
	for _, rune := range s {
		if rune == 10 {
			substrings = append(substrings, "\\n")
		} else {
			substrings = append(substrings, string(rune))
		}
	}
	return substrings
}

func Return2DASCIIArray(fileLines []string) [][]string {
	var asciiTemplates [][]string
	counter := 0
	for counter < len(fileLines) {
		tempAsciiTemplate := SplitByNumber(fileLines, counter, counter+8)
		asciiTemplates = append(asciiTemplates, tempAsciiTemplate)
		counter += 9
	}
	return asciiTemplates
}

func ReturnAllStringASCII(text string, asciiTemplates [][]string) string {
	var output strings.Builder

	substrings := ReturnString2EndlineArray(text)
	lenOfsubstrings := len(substrings)
	for index, v := range substrings {
		if v == "\\n" {
			output.WriteString("\n")
		} else {
			ReturnMultipleCharacter(&output, v, asciiTemplates)
		}
		if index != lenOfsubstrings-1 {
			output.WriteString("\n")
		}
	}
	return output.String()
}

func ReturnMultipleCharacter(output *strings.Builder, s string, asciiTemplates [][]string) {
	tempIntArrLetter := ReturnAsciiCodeInt(s)
	for i := 0; i < 8; i++ {
		for _, v := range tempIntArrLetter {
			output.WriteString(asciiTemplates[v][i])
		}
		output.WriteString("\n")
	}
}
