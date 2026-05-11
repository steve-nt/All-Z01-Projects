package piscine

// Define the food struct
type food struct {
	preptime int
}

// FoodDeliveryTime takes an order string and returns the preparation time
func FoodDeliveryTime(order string) int {
	// Convert the order to lowercase to ensure case insensitivity
	order = toLower(order)

	// Initialize the preparation time variable
	var preptime int

	// Determine the preparation time based on the order item
	switch order {
	case "burger":
		preptime = 15
	case "chips":
		preptime = 10
	case "nuggets":
		preptime = 12
	default:
		// If the order item does not exist, return 404
		return 404
	}

	// Return the preparation time
	return preptime
}

// toLower function converts a string to lowercase without using strings package
func toLower(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			char += 'a' - 'A'
		}
		result += string(char)
	}
	return result
}
