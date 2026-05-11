package internal

import (
	"bufio"
	"fmt"
	"os"
	exec "push_swap/internal/operations"
	run "push_swap/internal/utils"
)

func Checker() {

	if len(os.Args) < 2 {
		return
	}

	// Extract and validate input arguments
	// Initialize stack a and b
	stackA := run.ParseArgs()

	stackB := []int{}

	// Read instructions from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		instruction := scanner.Text()
		// Skip empty instructions
		if instruction == "" {
			continue
		}

		err := exec.Instruction(instruction, &stackA, &stackB)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error")
			return
		}
	}

	// Validate final state of the stacks
	if run.IsSorted(stackA) {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
