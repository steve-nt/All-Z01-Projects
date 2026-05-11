package utils

import (
	"math"
	"sort"
)

func CalculateAverage(numbers []float64) (int, float64) {
	if len(numbers) == 0 {
		panic("cannot calculate average, median, variance or standard deviation of zero elements")
	}
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	average := sum / float64(len(numbers))
	return int(math.Round(average)), average
}

func CalculateMedian(numbers []float64) int {
	sortedNumbers := make([]float64, len(numbers))
	copy(sortedNumbers, numbers)
	sort.Float64s(sortedNumbers)

	middleIndex := len(sortedNumbers) / 2
	if len(numbers)%2 == 0 { // Even length slice
		median := (sortedNumbers[middleIndex-1] + sortedNumbers[middleIndex]) / 2
		return int(math.Round(median))
	} else { // Odd length slice
		return int(math.Round((sortedNumbers[middleIndex])))
	}
}

func CalculateVariance(data []float64, average float64) (int, float64) {
	variance := 0.0
	for _, value := range data {
		variance += math.Pow(value-average, 2)
	}

	variance /= float64(len(data))
	return int(math.Round(variance)), variance
}

func CalculateStandardDeviation(variance float64) int {
	stdrDev := math.Sqrt(variance)
	return int(math.Round(stdrDev))
}
