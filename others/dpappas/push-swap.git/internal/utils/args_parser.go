package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseArguments takes the command-line arguments and converts them to a slice of integers
func ParseArgs() []int {
	args := os.Args[1:]

	// If no arguments are provided, exit the program
	if len(args) == 0 {
		os.Exit(0)
	}
	var nums []int

	// Join all arguments into a single string to handle cases where numbers are given in multiple parts
	allArgs := strings.Join(args, " ")

	// Split the argument string by spaces to separate individual numbers
	parts := strings.Fields(allArgs)

	// Convert each part into an integer
	for _, part := range parts {
		// Convert the argument to an integer
		num, err := strconv.Atoi(part)
		if err != nil {
			fmt.Println("Error")
			os.Exit(0)
		}
		nums = append(nums, num)
	}

	// Validate that there are no duplicates
	if HasDuplicates(nums) {
		fmt.Println("Error")
		os.Exit(0)
	}
	// Validate that it is not sorted already
	if IsSorted(nums) {
		os.Exit(0)
	}
	return nums
}
