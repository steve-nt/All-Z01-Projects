package piscine

// Atoi converts a string to an integer, mimicking the behavior of the C function atoi.
func Atoi(str string) int {
	result := 0 // Initialize result to 0, which will store the final integer value.
	sign := 1   // Initialize sign to 1, assuming the number is positive.

	for i, char := range str {
		if i == 0 && char == '-' { // Check if the first character is '-'.
			sign = -1 // Set sign to -1 for negative numbers.
			continue  // Skip to the next iteration to avoid processing the sign character as a digit.
		}
		if i == 0 && char == '+' { // Check if the first character is '+'.
			continue // Skip to the next iteration to avoid processing the sign character as a digit.
		}

		if char < '0' || char > '9' { // If the character is not a valid digit (0-9), return 0.
			return 0
		}

		result = result*10 + int(char-'0') // Update result by adding the new digit.
	}

	return result * sign // Return the final result multiplied by the sign.
}
