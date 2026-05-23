package main

import (
	"math"
	"sort"
)

func median(numbers []int) int {
	sort.Ints(numbers)

	n := len(numbers)

	if n%2 == 1 {
		return numbers[n/2]
	}

	med := float64(numbers[n/2-1]+numbers[n/2]) / 2.0

	return int(math.Round(med))
}
