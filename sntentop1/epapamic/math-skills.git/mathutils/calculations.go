package mathutils

import (
	"math"
	"sort"
)

// Function to calculate the average
func CalculateAverage(numbers []int) float64 {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return float64(sum) / float64(len(numbers))
}

// Function to calculate the median
func CalculateMedian(numbers []int) float64 {
	// sort numbers in ascending order
	sort.Ints(numbers)
	if len(numbers) % 2 == 0 {
		// if length of numbers is an even number we need the number before average and the average number 
		return float64(numbers[len(numbers) / 2 - 1] + numbers[len(numbers) / 2]) / 2.0
	}
	return float64(numbers[len(numbers) / 2])
}

// Function to calculate the variance
func CalculateVariance(numbers []int, average float64) float64 {
	varianceSum := 0.0
	for _, num :=range numbers {
		// calculates the squared difference between the current number and the average
		varianceSum += math.Pow(float64(num) - average, 2)
	}
	//  calculates the variance by dividing by the length of the number slice
	return varianceSum / float64(len(numbers))
}

// Function to calculate the standard deviation
func CalculateStandardDeviation(variance float64) float64 {
	return math.Sqrt(variance)
}