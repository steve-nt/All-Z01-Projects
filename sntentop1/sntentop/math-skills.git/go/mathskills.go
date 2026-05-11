package main // package name

import (
	"bufio"   // package for reading files and user input from the console
	"fmt"     // package for printing to the console
	"math"    // package for mathematical functions
	"os"      // package for reading files and user input from the console
	"sort"    // package for sorting numbers
	"strconv" // package for converting strings to numbers
)

func readData(filePath string) ([]int, error) { // function to read data from a file
	file, err := os.Open(filePath) // open the file
	if err != nil {                // if there is an error
		return nil, err // return the error
	} // end if
	defer file.Close() // close the file after the function ends

	var data []int                    // create a slice to store the data as integers
	scanner := bufio.NewScanner(file) // create a scanner to read the file line by line
	for scanner.Scan() {              // loop through each line in the file and read it as a string
		num, err := strconv.Atoi(scanner.Text()) // convert the string to an integer and store it in the num variable as an integer
		if err != nil {                          // if there is an error
			return nil, err // return the error and nil slice
		} // end if
		data = append(data, num) // append the integer to the data slice and continue the loop
	}
	if err := scanner.Err(); err != nil { // if there is an error reading the file after the loop ends
		return nil, err // return the error and nil slice
	}

	return data, nil // return the data slice and nil error and end the function
}

// Calculates the average of a list of numbers
func calculateAverage(data []int) float64 { // function to calculate the average of a list of numbers
	sum := 0                   // variable to store the sum of the numbers
	for _, num := range data { // loop through each number in the list
		sum += num // add the number to the sum
	}
	return float64(sum) / float64(len(data)) // return the average of the numbers
}

// Calculates the median of a list of numbers
func calculateMedian(data []int) float64 { // function to calculate the median of a list of numbers
	sort.Ints(data) // sort the numbers in ascending order using the sort package from the standard library (sort.Ints)
	n := len(data)  // get the length of the list of numbers (n) and store it in the variable n as an integer (int) data type
	mid := n / 2    // calculate the middle index of the list of numbers and store it in the variable mid as an integer (int) data type

	if n%2 == 0 { // if the length of the list of numbers is even (n is divisible by 2 with no remainder)
		return float64(data[mid-1]+data[mid]) / 2 // return the average of the two middle numbers as a float64 data type
	}
	return float64(data[mid]) // if the length of the list of numbers is odd, return the middle number as a float64 data type
}

// Calculates the variance of a list of numbers
func calculateVariance(data []int, average float64) float64 { // function to calculate the variance of a list of numbers
	var sumSquaredDiff float64 // variable to store the sum of the squared differences between each number and the average
	for _, num := range data { // loop through each number in the list of numbers
		diff := float64(num) - average // calculate the difference between the number and the average and store it in the variable diff as a float64 data type
		sumSquaredDiff += diff * diff  // add the square of the difference to the sum of the squared differences
	}
	return sumSquaredDiff / float64(len(data)) // return the variance of the list of numbers as a float64 data type
}

// (Check the math.Pow in go)
// Calculates the standard deviation of a list of numbers
func calculateStdDeviation(variance float64) float64 { // function to calculate the standard deviation of a list of numbers
	return math.Sqrt(variance) // return the square root of the variance as a float64 data type
}

// Custom rounding function that rounds 0.5 and higher up to the next integer
func customRound(value float64) int { // function to round a float64 number to the nearest integer
	if value-math.Floor(value) >= 0.5 { // if the decimal part of the number is greater than or equal to 0.5
		return int(math.Ceil(value)) // round the number up to the next integer using the math.Ceil function from the standard library
	}
	return int(math.Floor(value)) // round the number down to the nearest integer using the math.Floor function from the standard library
}

//(may be use math.Round in go)

func main() { // main function to run the program
	if len(os.Args) != 2 { // if the number of command-line arguments is not equal to 2
		fmt.Println("Usage: go run your_program.go data.txt") // print a usage message to the console
		return                                                // end the program and return a non-zero exit status
	}

	filePath := os.Args[1]          // get the file path from the command-line arguments
	data, err := readData(filePath) // read the data from the file using the readData function
	if err != nil {                 // if there is an error reading the file
		fmt.Printf("Error reading file: %v\n", err) // print an error message to the console
		return                                      // end the program and return a non-zero exit status and end the program
	}

	average := calculateAverage(data)               // calculate the average of the data using the calculateAverage function
	median := calculateMedian(data)                 // calculate the median of the data using the calculateMedian function
	variance := calculateVariance(data, average)    // calculate the variance of the data using the calculateVariance function
	stdDeviation := calculateStdDeviation(variance) // calculate the standard deviation of the data using the calculateStdDeviation function

	fmt.Printf("Average: %d\n", customRound(average))                 // print the average to the console using the customRound function
	fmt.Printf("Median: %d\n", customRound(median))                   // print the median to the console using the customRound function
	fmt.Printf("Variance: %d\n", customRound(variance))               // print the variance to the console using the customRound function
	fmt.Printf("Standard Deviation: %d\n", customRound(stdDeviation)) // print the standard deviation to the console using the customRound function
}
