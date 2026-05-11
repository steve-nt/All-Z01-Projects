package utils

import (
	"push-swap/operations"
)

func SortFinal(stackA, stackB *[]int) (command []string) {
	switch len(*stackA) {
	case 0, 1:
		return command
	case 2:
		if (*stackA)[0] > (*stackA)[1] {
			operations.Sa(stackA)
			command = append(command, "sa")
		}
	case 3:
		SortThree(stackA, stackB, &command)
	default:
		command = SortEntire(stackA, stackB)
	}
	return command
}

// SortThree sorts a stack with only three numbers using basic operations
func SortThree(stackA, stackB *[]int, command *[]string) {
	if len(*stackA) == 2 {
		if (*stackA)[0] > (*stackA)[1] {
			// Swap the first two elements
			operations.Sa(stackA)
			*command = append(*command, "sa")
		}
		return
	}

	a, b, c := (*stackA)[0], (*stackA)[1], (*stackA)[2]

	switch {
	// Case 2 1 3
	case a > b && b < c && a < c:
		operations.Sa(stackA)
		*command = append(*command, "sa")

	// Case 3 2 1
	case a > b && b > c:
		operations.Sa(stackA)
		operations.Rra(stackA)
		*command = append(*command, "sa", "rra")

	// Case 3 1 2
	case a > b && b < c && a > c:
		operations.Ra(stackA)
		*command = append(*command, "ra")

	// Case 1 3 2
	case a < b && b > c && a < c:
		operations.Sa(stackA)
		operations.Ra(stackA)
		*command = append(*command, "sa", "ra")

	// Case 2 3 1
	case a < b && b > c && a > c:
		operations.Rra(stackA)
		*command = append(*command, "rra")
	}
}

// Sorting function to sort the entire execute
func SortEntire(stackA, stackB *[]int) (command []string) {
	// If execute r fewer elements, sort them directly
	if len(*stackA) <= 3 {
		SortThree(stackA, stackB, &command)
		return command
	}

	// Move elements to stackB until stackA has 3 elements
	for len(*stackA) > 3 {
		// Find the minimum index in stackA
		minIndex := findMinIndex(stackA)
		if minIndex == 0 {
			// Move top element to stackB
			operations.Pb(stackA, stackB)
			command = append(command, "pb")
		} else if minIndex == 1 {
			// Swap first two and then move to stackB
			operations.Sa(stackA)
			operations.Pb(stackA, stackB)
			command = append(command, "sa", "pb")
		} else if minIndex <= len(*stackA)/2 {
			// Rotate stackA to get minimum element to the top
			operations.Ra(stackA)
			command = append(command, "ra")
		} else {
			// Reverse rotate stackA to get minimum element to the top
			operations.Rra(stackA)
			command = append(command, "rra")
		}
	}

	// Now sort the remaining 3 elements in stackA
	SortThree(stackA, stackB, &command)

	// Move all elements from stackB back to stackA in sorted order
	for len(*stackB) > 0 {
		operations.Pa(stackA, stackB)
		command = append(command, "pa")
	}

	rotateToMin(stackA, &command)

	return command
}

func rotateToMin(stackA *[]int, command *[]string) {
	minIndex := findMinIndex(stackA)
	if minIndex <= len(*stackA)/2 {
		for i := 0; i < minIndex; i++ {
			operations.Ra(stackA)
			*command = append(*command, "ra")
		}
	} else {
		for i := 0; i < len(*stackA)-minIndex; i++ {
			operations.Rra(stackA)
			*command = append(*command, "rra")
		}
	}
}

// Helper to find index of minimum value
func findMinIndex(stack *[]int) int {
	minIndex := 0
	for i := 1; i < len(*stack); i++ {
		if (*stack)[i] < (*stack)[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}

// isSorted checks if stackA is sorted in ascending order
// by iterating through the list and ensures that each element is less than or equal to the next one
func IsSorted(stackA []int) bool {
	if len(stackA) < 2 {
		return true
	}

	for i := 1; i < len(stackA); i++ {
		if stackA[i-1] > stackA[i] {
			return false 
		}
	}
	return true 
}
