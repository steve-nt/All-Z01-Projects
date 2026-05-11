package piscine

// ShoppingListSort sorts a slice of strings by their lengths in ascending order without using imports
func ShoppingListSort(slice []string) []string {
	n := len(slice)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if len(slice[j]) > len(slice[j+1]) {
				slice[j], slice[j+1] = slice[j+1], slice[j] // Swap if the element found is greater than the next element
			}
		}
	}
	return slice
}
