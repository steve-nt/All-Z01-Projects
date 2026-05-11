package main

import (
	"fmt"
	"os"
	"push-swap/utils"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		return
	}

	// Parse the input argument into a space-separated list of numbers
	arg := os.Args[1]
	parts := strings.Fields(arg)

	if len(parts) == 0 {
		fmt.Println("Error: Input is empty.")
		return
	}

	// parseStack input strings into integers for stackA
	stackA, err := utils.ParseStack(parts)
	if err != nil {
		fmt.Println("Error")
		return
	}

	if utils.IsSorted(stackA) {
		return
	}

	stackB := []int{}

	commands := utils.SortFinal(&stackA, &stackB)

	for _, command := range commands {
		fmt.Println(command)
	}
}
