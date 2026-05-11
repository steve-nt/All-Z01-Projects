package internal

import (
	"fmt"
	exec "push_swap/internal/operations"
	run "push_swap/internal/utils"
)

func PushSwap() {

	// Parse and validate input
	// Create stackA and stackB from the parsed integers
	stackA := run.ParseArgs()

	stackB := []int{}

	// Get sorting instructions
	instructions := exec.Sort(&stackA, &stackB)

	// Print the sorting instructions
	for _, instruction := range instructions {
		fmt.Println(instruction)
	}
}
