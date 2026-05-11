package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	"platform.zone01.gr/git/epapamic/math-skills/mathutils"
)

func main() {
	// Check if exactly 1 argument is provided
	if len(os.Args) != 2 {
		fmt.Println("Only 1 argument is needed")
		return
	}

	// Open the file
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read numbers from the file
	var numbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// convert each interger in each line
		num, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Error converting line to integer:", err)
			return
		}
		numbers = append(numbers, num)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// check for empty slice of numbers
	if len(numbers) == 0 {
		fmt.Println("No numbers found in the file.")
		return
	}

	// Perform calculations
	average := mathutils.CalculateAverage(numbers)
	median := mathutils.CalculateMedian(numbers)
	variance := mathutils.CalculateVariance(numbers, average)
	standardDeviation := mathutils.CalculateStandardDeviation(variance)

	// Print the results as rounded numbers
	fmt.Printf("Average: %d\n", int(math.Round(average)))
	fmt.Printf("Median: %d\n", int(math.Round(median)))
	fmt.Printf("Variance: %d\n", int(math.Round(variance)))
	fmt.Printf("standardDeviation: %d\n", int(math.Round(standardDeviation)))
}