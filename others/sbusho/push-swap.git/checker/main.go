package main

import (
	"bufio"
	"fmt"
	"os"
	"push-swap/operations"
	"push-swap/utils"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		return
	}

	// Parse the input argument to initialize stackA
	arg := os.Args[1]
	parts := strings.Fields(arg)

	// Validate stackA
	stackA, err := utils.ParseStack(parts)
	if err != nil {
		fmt.Println("Error")
		return
	}

	stackB := []int{}

	// Read instructions from standard input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()

		if err := operations.ExecuteCommand(command, &stackA, &stackB); err != nil {
			fmt.Println("Error")
			return
		}
	}

	// Check if stackA is sorted and stackB is empty
	if utils.IsSorted(stackA) && len(stackB) == 0 {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
