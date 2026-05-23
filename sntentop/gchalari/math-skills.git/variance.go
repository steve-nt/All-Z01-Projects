package main

import "math"

func variance(numbers []int, avg int) int {
	sum := 0

	for _, n := range numbers {
		diff := n - avg
		sum += diff * diff
	}

	variance := float64(sum) / float64(len(numbers))

	return int(math.Round(variance))
}
