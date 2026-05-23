package main

import "math"

func average(numbers []int) int {
	sum := 0

	for _, n := range numbers {
		sum += n
	}

	avg := float64(sum) / float64(len(numbers))

	return int(math.Round(avg))
}
