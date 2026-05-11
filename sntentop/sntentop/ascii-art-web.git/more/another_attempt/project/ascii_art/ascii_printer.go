package ascii

// import (
// 	"strings"
// )

// func SplitByNumber(arr []string, start int, end int) []string {
// 	var result []string
// 	for i := start; i < end; i++ {
// 		result = append(result, arr[i])
// 	}
// 	return result
// }

// func ReturnString2EndlineArray(s string) []string {
// 	var substrings []string
// 	for _, rune := range s {
// 		if rune == 10 {
// 			substrings = append(substrings, "\\n")
// 		} else {
// 			substrings = append(substrings, string(rune))
// 		}
// 	}
// 	return substrings
// }

// func Return2DASCIIArray(fileLines []string) [][]string {
// 	var asciiTemplates [][]string
// 	counter := 0
// 	for counter < len(fileLines) {
// 		tempAsciiTemplate := SplitByNumber(fileLines, counter, counter+8)
// 		asciiTemplates = append(asciiTemplates, tempAsciiTemplate)
// 		counter += 9
// 	}
// 	return asciiTemplates
// }

// func ReturnAllStringASCII(text string, asciiTemplates [][]string) string {
// 	var output strings.Builder

// 	substrings := ReturnString2EndlineArray(text)
// 	for _, v := range substrings {
// 		if v == "\\n" {
// 			// Don't add newline here if you want everything in one line.
// 			// Instead, just continue to the next character.
// 			continue
// 		} else {
// 			ReturnMultipleCharacter(&output, v, asciiTemplates)
// 		}
// 	}
// 	return output.String()
// }

// func ReturnMultipleCharacter(output *strings.Builder, s string, asciiTemplates [][]string) {
// 	tempIntArrLetter := ReturnAsciiCodeInt(s)

// 	// Create an array to store each row of the ASCII art for the entire input
// 	rows := make([]string, 8) // 8 rows per character template

// 	// Loop through each character and append its rows to the appropriate index in `rows`
// 	for _, charIndex := range tempIntArrLetter {
// 		for i := 0; i < 8; i++ {
// 			rows[i] += asciiTemplates[charIndex][i] // Append characters horizontally (no newlines here)
// 		}
// 	}

// 	// Append the final rows to the output, each row on the same line
// 	for _, row := range rows {
// 		output.WriteString(row) // Append without newline
// 	}
// }
