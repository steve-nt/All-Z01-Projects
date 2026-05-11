package piscine

import "github.com/01-edu/z01"

// PrintNbrBase prints an integer in a given base using only z01.PrintRune.
func PrintNbrBase(nbr int, base string) {
	if !isValidBase(base) {
		printString("NV")
		return
	}

	baseLen := len(base)
	if nbr == 0 {
		z01.PrintRune(rune(base[0]))
		return
	}

	if nbr == -9223372036854775808 {
		printString("-9223372036854775808")
		return
	}

	if nbr < 0 {
		z01.PrintRune('-')
		nbr = -nbr
	}

	var result []rune
	for nbr > 0 {
		remainder := nbr % baseLen
		result = append([]rune{rune(base[remainder])}, result...)
		nbr = nbr / baseLen
	}

	for _, r := range result {
		z01.PrintRune(r)
	}
}

// isValidBase checks if the base string is valid.
func isValidBase(base string) bool {
	if len(base) < 2 {
		return false
	}

	charMap := make(map[rune]bool)
	for _, char := range base {
		if char == '+' || char == '-' || charMap[char] {
			return false
		}
		charMap[char] = true
	}

	return true
}

// printString prints a string using z01.PrintRune.
func printString(s string) {
	for _, char := range s {
		z01.PrintRune(char)
	}
}
