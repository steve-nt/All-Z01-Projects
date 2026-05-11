package piscine

func Abort(a, b, c, d, e int) int {
	// Create a slice with the five integers
	nums := []int{a, b, c, d, e}

	// Simple selection sort
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 5; j++ {
			if nums[i] > nums[j] {
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
	}

	// Return the median (the third element in the sorted slice)
	return nums[2]
}
