package main

import "math"

func stdDev(v int) int {
	std := math.Sqrt(float64(v))

	return int(math.Round(std))
}
