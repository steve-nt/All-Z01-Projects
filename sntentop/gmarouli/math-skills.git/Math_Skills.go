package main

import (
	"fmt"
	"bufio"
	"os"
	"math"
	"sort"
	"strconv"
)

// rounds a float64 to the nearest integer
func round(x float64) int {
	return int(math.Round(x))
}

// calculates Average
func calculateAverage(numbers []int) float64 {
	var total int
	for _, num := range numbers {
		total += num
	}
	return float64(total) / float64(len(numbers))
}

// calculates Median 
func calculateMedian(numbers []int) float64 {
	sort.Ints(numbers)
	n := len(numbers)
	if n%2 == 0 {
		return float64(numbers[n/2-1]+numbers[n/2]) / 2
	}
	return float64(numbers[n/2])
}

// calculates Variance 
func calculateVariance(numbers []int, mean float64) float64 {
	var sumSquaredDifferences float64
	for _, num := range numbers {
		diff := float64(num) - mean
		sumSquaredDifferences += diff * diff
	}
	return sumSquaredDifferences / float64(len(numbers))
}

// calculates StandardDeviation
func calculateStandardDeviation(variance float64) float64 {
	return math.Sqrt(variance)
}

// readDataFromFile function
func readDataFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open the file: %v", err)
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("could not convert '%s' to an integer: %v", line, err)
		}
		numbers = append(numbers, number)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading the file: %v", err)
	}

	return numbers, nil
}

//  printStats function displays the statistics rounded to the nearest integer.
func printStats(average, median, variance, stdDeviation float64) {
	fmt.Printf("Average: %d\n", round(average))
	fmt.Printf("Median: %d\n", round(median))
	fmt.Printf("Variance: %d\n", round(variance))
	fmt.Printf("Standard Deviation: %d\n", round(stdDeviation))
}

func main() {
	// Check if a filename is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run your-program.go <data_file>")
		os.Exit(1)
	}
	filename := os.Args[1]

	// Read data from the file
	data, err := readDataFromFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Check if the file contains any numbers
	if len(data) == 0 {
		fmt.Println("Error: The file contains no valid numbers.")
		os.Exit(1)
	}

	// Calculate statistics
	average := calculateAverage(data)
	median := calculateMedian(data)
	variance := calculateVariance(data, average)
	stdDeviation := calculateStandardDeviation(variance)

	// Print results
	printStats(average, median, variance, stdDeviation)
}
