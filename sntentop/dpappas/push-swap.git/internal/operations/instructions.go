package operations

import (
	"fmt"
)

// ExecuteInstruction executes a single instruction on the stacks
func Instruction(instruction string, stackA, stackB *[]int) error {
	// fmt.Printf("Executing: %s\n", instruction)
	// fmt.Printf("Before - stackA: %v, stackB: %v\n", *stackA, *stackB)

	var err error
	switch instruction {
	case "pa":
		if len(*stackB) > 0 {
			// Move the first element from stackB to stackA
			*stackA = append([]int{(*stackB)[0]}, (*stackA)...)
			*stackB = (*stackB)[1:]
		}
	case "pb":
		if len(*stackA) > 0 {
			// Move the first element from stackA to stackB
			*stackB = append([]int{(*stackA)[0]}, (*stackB)...)
			*stackA = (*stackA)[1:]
		}
	case "sa":
		if len(*stackA) > 1 {
			// Swap the first two elements of stackA
			(*stackA)[0], (*stackA)[1] = (*stackA)[1], (*stackA)[0]
		}
	case "sb":
		if len(*stackB) > 1 {
			// Swap the first two elements of stackB
			(*stackB)[0], (*stackB)[1] = (*stackB)[1], (*stackB)[0]
		}
	case "ss":
		// Swap both stackA and stackB
		if len(*stackA) > 1 {
			(*stackA)[0], (*stackA)[1] = (*stackA)[1], (*stackA)[0]
		}
		if len(*stackB) > 1 {
			(*stackB)[0], (*stackB)[1] = (*stackB)[1], (*stackB)[0]
		}
	case "ra":
		if len(*stackA) > 1 {
			// Rotate: move the first element to the end of stackA
			*stackA = append((*stackA)[1:], (*stackA)[0])
		}
	case "rb":
		if len(*stackB) > 1 {
			// Rotate: move the first element to the end of stackB
			*stackB = append((*stackB)[1:], (*stackB)[0])
		}
	case "rr":
		// Rotate both stackA and stackB
		if len(*stackA) > 1 {
			*stackA = append((*stackA)[1:], (*stackA)[0])
		}
		if len(*stackB) > 1 {
			*stackB = append((*stackB)[1:], (*stackB)[0])
		}
	case "rra":
		if len(*stackA) > 1 {
			// Reverse rotate: move the last element to the front of stackA
			*stackA = append([]int{(*stackA)[len(*stackA)-1]}, (*stackA)[:len(*stackA)-1]...)
		}
	case "rrb":
		if len(*stackB) > 1 {
			// Reverse rotate: move the last element to the front of stackB
			*stackB = append([]int{(*stackB)[len(*stackB)-1]}, (*stackB)[:len(*stackB)-1]...)
		}
	case "rrr":
		// Reverse rotate both stackA and stackB
		if len(*stackA) > 1 {
			*stackA = append([]int{(*stackA)[len(*stackA)-1]}, (*stackA)[:len(*stackA)-1]...)
		}
		if len(*stackB) > 1 {
			*stackB = append([]int{(*stackB)[len(*stackB)-1]}, (*stackB)[:len(*stackB)-1]...)
		}
	default:
		err = fmt.Errorf("invalid instruction: %s", instruction)
	}

	//fmt.Printf("After - stackA: %v, stackB: %v\n", *stackA, *stackB)
	return err
}
