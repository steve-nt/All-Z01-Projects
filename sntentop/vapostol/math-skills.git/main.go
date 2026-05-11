package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	// Check if the user provided a file name as a command-line argument.
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run program.go data.txt")
		return
	}

	// Open the file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close() // Close the file

	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	var numbers []float64

	// Add each number to the list.
	for scanner.Scan() {
		num, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		numbers = append(numbers, num)
	}

	// Check if there were any errors reading the file.
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	// Calculate the average, median, variance, and standard deviation.
	avg := calculateAverage(numbers)
	median := calculateMedian(numbers)
	variance := calculateVariance(numbers)
	stddev := calculateStandardDeviation(variance)

	// Print the results.
	fmt.Printf("Average: %d\n", int64(math.Round(avg)))
	fmt.Printf("Median: %d\n", int64(math.Round(median)))
	fmt.Printf("Variance: %d\n", int64(variance+0.5))
	fmt.Printf("Standard Deviation: %d\n", int64(math.Round(stddev)))
}

// Calculate the average (The average value of a set of numbers)
func calculateAverage(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

// Calculate the median (The middle value of a set of numbers when they're sorted in order.)
func calculateMedian(numbers []float64) float64 {
	sort.Float64s(numbers) // Sort the list of numbers.
	length := len(numbers)
	if length%2 == 0 {
		// If the list has an even number of elements, the median is the average of the two middle elements.
		return (numbers[length/2-1] + numbers[length/2]) / 2
	} else {
		// If the list has an odd number of elements, the median is the middle element.
		return numbers[length/2]
	}
}

// Calculate the variance (A measure of how spread out a set of numbers is from the mean.)
func calculateVariance(numbers []float64) float64 {
	avg := calculateAverage(numbers)
	sum := 0.0
	for _, num := range numbers {
		diff := num - avg        // Calculate the difference between each number and the average.
		sum += math.Pow(diff, 2) // Add the squared difference to the sum.
	}
	return sum / float64(len(numbers))
}

// Calculate the standard deviation (The square root of the variance, which is a more useful measure of how spread out a set of numbers is.)
func calculateStandardDeviation(variance float64) float64 {
	return math.Sqrt(variance)
}
