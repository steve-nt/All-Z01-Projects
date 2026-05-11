package operations

import (
	"errors"
	"strings"
)

// Push top element from stackB to stackA
func Pa(stackA, stackB *[]int) {
	// If stackA is not empty, the first element of stackA is inserted at the top of stackB.
	if len(*stackB) > 0 {
		*stackA = append([]int{(*stackB)[0]}, *stackA...) 
		*stackB = (*stackB)[1:]
	}
}

// Push top element from stackA to stackB
func Pb(stackA, stackB *[]int) {
	// If stackA is not empty, the first element of stackA is inserted at the top of stackB.
	if len(*stackA) > 0 {
		*stackB = append([]int{(*stackA)[0]}, *stackB...) 
		*stackA = (*stackA)[1:]                           
	}
}

// Swap first two elements of stackA, only performed if stackA has at least two elements
func Sa(stackA *[]int) {
	if len(*stackA) > 1 {
		(*stackA)[0], (*stackA)[1] = (*stackA)[1], (*stackA)[0]
	}
}

// Swap first two elements of stackB, only performed if stackB has at least two elements.
func Sb(stackB *[]int) {
	if len(*stackB) > 1 {
		(*stackB)[0], (*stackB)[1] = (*stackB)[1], (*stackB)[0]
	}
}

// Execute sa and sb, swaps the first two elements of both stackA and stackB if possible
func Ss(stackA, stackB *[]int) {
	Sa(stackA)
	Sb(stackB)
}

// Rotate stackA (first element becomes last), only performed if stackA has at least two elements.
func Ra(stackA *[]int) {
	if len(*stackA) > 1 {
		*stackA = append((*stackA)[1:], (*stackA)[0])
	}
}

// Rotate stackB, (first element becomes last), only performed if stackB has at least two elements.
func Rb(stackB *[]int) {
	if len(*stackB) > 1 {
		*stackB = append((*stackB)[1:], (*stackB)[0])
	}
}

// This rotates both stackA and stackB if they have at least two elements each
func Rr(stackA, stackB *[]int) {
	Ra(stackA)
	Rb(stackB)
}

// Reverse rotate stackA (last element becomes first), only performed if stackA has at least two elements
func Rra(stackA *[]int) {
	if len(*stackA) > 1 {
		*stackA = append([]int{(*stackA)[len(*stackA)-1]}, (*stackA)[:len(*stackA)-1]...)
	}
}

// Reverse rotate stackB (last element becomes first), only performed if stackB has at least two elements
func Rrb(stackB *[]int) {
	if len(*stackB) > 1 {
		*stackB = append([]int{(*stackB)[len(*stackB)-1]}, (*stackB)[:len(*stackB)-1]...)
	}
}

// This moves the last element of both stackA and stackB to their respective first positions.
func Rrr(stackA, stackB *[]int) {
	Rra(stackA)
	Rrb(stackB)
}

// Executes push-swap instructions
func ExecuteCommand(command string, stackA, stackB *[]int) error {

	// Trim leading and trailing whitespaces from the instruction
	command = strings.TrimSpace(command)

	if command == "" {
		return nil 
	}

	switch command {
	case "pa":
		Pa(stackA, stackB)
	case "pb":
		Pb(stackA, stackB)
	case "sa":
		Sa(stackA)
	case "sb":
		Sb(stackB)
	case "ss":
		Ss(stackA, stackB)
	case "ra":
		Ra(stackA)
	case "rb":
		Rb(stackB)
	case "rr":
		Rr(stackA, stackB)
	case "rra":
		Rra(stackA)
	case "rrb":
		Rrb(stackB)
	case "rrr":
		Rrr(stackA, stackB)
	default:
		return errors.New("Invalid Command" + command)
	}
	return nil
}
