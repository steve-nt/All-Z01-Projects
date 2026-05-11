package piscine

func NRune(s string, n int) rune {
	runes := []rune(s) // Convert the string to a slice of runes to handle Unicode characters correctly.

	if n <= 0 || n > len(runes) {
		return 0 // Return 0 if the index is out of bounds.
	}

	return runes[n-1] // Return the nth rune (1-based index).
}
