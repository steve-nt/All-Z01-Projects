package operations

// SortThree sorts a stack with only three numbers using basic operations
func SortThree(stackA, stackB *[]int, instr *[]string) {
	if len(*stackA) == 2 && (*stackA)[0] > (*stackA)[1] {
		// Swap the first two elements
		Instruction("sa", stackA, stackB)
		*instr = append(*instr, "sa")
	} else if len(*stackA) == 3 {
		// Handle sorting of three elements
		if (*stackA)[0] > (*stackA)[1] && (*stackA)[0] < (*stackA)[2] {
			Instruction("sa", stackA, stackB)
			*instr = append(*instr, "sa")
		} else if (*stackA)[0] < (*stackA)[1] && (*stackA)[0] > (*stackA)[2] {
			Instruction("rra", stackA, stackB)
			*instr = append(*instr, "rra")
		} else if (*stackA)[0] > (*stackA)[1] && (*stackA)[1] > (*stackA)[2] {
			Instruction("sa", stackA, stackB)
			Instruction("rra", stackA, stackB)
			*instr = append(*instr, "sa", "rra")
		} else if (*stackA)[0] > (*stackA)[1] && (*stackA)[1] < (*stackA)[2] {
			Instruction("ra", stackA, stackB)
			*instr = append(*instr, "ra")
		} else if (*stackA)[0] < (*stackA)[1] && (*stackA)[1] > (*stackA)[2] {
			Instruction("sa", stackA, stackB)
			Instruction("ra", stackA, stackB)
			*instr = append(*instr, "sa", "ra")
		}
	}
}

// Sorting function to sort the entire execute
func Sort(stackA, stackB *[]int) (instr []string) {
	// If execute r fewer elements, sort them directly
	if len(*stackA) <= 3 {
		SortThree(stackA, stackB, &instr)
		return instr
	}

	// Move elements to stackB until stackA has 3 elements
	for len(*stackA) > 3 {
		// Find the minimum index in stackA
		minIndex := findMinIndex(stackA)
		if minIndex == 0 {
			// Move top element to stackB
			Instruction("pb", stackA, stackB)
			instr = append(instr, "pb")
		} else if minIndex == 1 {
			// Swap first two and then move to stackB
			Instruction("sa", stackA, stackB)
			Instruction("pb", stackA, stackB)
			instr = append(instr, "sa", "pb")
		} else if minIndex <= len(*stackA)/2 {
			// Rotate stackA to get minimum element to the top
			Instruction("ra", stackA, stackB)
			instr = append(instr, "ra")
		} else {
			// Reverse rotate stackA to get minimum element to the top
			Instruction("rra", stackA, stackB)
			instr = append(instr, "rra")
		}
	}

	// Now sort the remaining 3 elements in stackA
	SortThree(stackA, stackB, &instr)

	// Move all elements from stackB back to stackA in sorted order
	for len(*stackB) > 0 {
		Instruction("pa", stackA, stackB)
		instr = append(instr, "pa")
	}

	return instr
}

// Find the index of the minimum element in stackA
func findMinIndex(stackA *[]int) int {
	minIndex := 0
	for i := 1; i < len(*stackA); i++ {
		if (*stackA)[i] < (*stackA)[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}
