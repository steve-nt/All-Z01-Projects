package utils

import (
	"errors"
	"strconv"
)

// parseStack validates the input and initializes stackA
// by converting the input string into a slice of integers
func ParseStack(parts []string) ([]int, error) {
	if len(parts) == 0 {
		return nil, errors.New("empty input")
	}

	stackA := []int{}          
	seen := make(map[int]bool) 

	for _, part := range parts {

		// Convert the string to an integer
		num, err := strconv.Atoi(part)

		if err != nil {
			return nil, errors.New("Invalid integer: " + part)
		}

		// Check for duplicates
		if seen[num] {
			return nil, errors.New("Duplicate found: " + strconv.Itoa(num))
		}
		seen[num] = true             
		stackA = append(stackA, num) 
	}

	return stackA, nil
}

