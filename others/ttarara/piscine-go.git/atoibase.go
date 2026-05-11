package piscine

// AtoiBase converts a numeric string `s` from a given base to an integer.
func AtoiBase(s string, base string) int {
	// Check if the base is valid
	if !isValidBase(base) {
		return 0
	}

	// Create a map for quick lookup of character values
	baseMap := make(map[rune]int)
	for i, char := range base {
		baseMap[char] = i
	}

	// Convert the string to an integer
	result := 0
	baseLen := len(base)
	for _, char := range s {
		// Look up the character's value in the base
		value, exists := baseMap[char]
		if !exists {
			return 0
		}
		result = result*baseLen + value
	}
	return result
}

// isValidBase checks if the base string is valid.
func IsValidBase(base string) bool {
	if len(base) < 2 {
		return false
	}

	// Create a set to check for duplicate characters
	charSet := make(map[rune]bool)
	for _, char := range base {
		if char == '+' || char == '-' {
			return false
		}
		if charSet[char] {
			return false
		}
		charSet[char] = true
	}
	return true
}
