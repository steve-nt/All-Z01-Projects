package piscine

// StringToIntSlice converts a string to a slice of int, where each int is the ASCII value of the corresponding character.
func StringToIntSlice(str string) []int {
	// Initialize an empty slice to hold the ASCII values.
	intSlice := []int(nil)

	// Iterate over each character in the string.
	for _, char := range str {
		// Append the ASCII value of the character to the slice.
		intSlice = append(intSlice, int(char))
	}

	// Return the slice of ASCII values.
	return intSlice
}
