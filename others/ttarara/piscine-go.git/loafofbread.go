package piscine

import "github.com/01-edu/z01"

// LoafOfBread processes the input string according to the specified rules.
func LoafOfBread(str string) string {
	if len(str) < 5 {
		return "Invalid Output\n"
	}

	result := ""
	count := 0

	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			result += string(str[i])
			count++
		}

		if count == 5 {
			if i+1 < len(str) {
				i++ // Skip the next character
			}
			count = 0
			result += " "
		}
	}

	return result + "\n"
}

// PrintString prints the string using z01.PrintRune
func PrintString(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
}
