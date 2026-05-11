package utils

func IsSorted(nums []int) bool {
	// Iterate through the list and check if each element is less than or equal to the next one
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] > nums[i+1] {
			// If we find an element greater than the next one, the list is not sorted
			return false
		}
	}
	// If we complete the loop without finding an out-of-order pair, the list is sorted
	return true
}

// HasDuplicates checks if a slice contains duplicate values
func HasDuplicates(nums []int) bool {
	seen := make(map[int]bool)
	for _, num := range nums {
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}
