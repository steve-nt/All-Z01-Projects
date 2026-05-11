package main

import "fmt"

func QuadB(x, y int) {
	// Check if the dimensions are positive
	if x <= 0 || y <= 0 {
		return
	}

	// Draw the top edge
	fmt.Print("/")
	for i := 0; i < x-2; i++ {
		fmt.Print("*")
	}
	if x > 1 {
		fmt.Print("\\")
	}
	fmt.Println()

	// Draw the middle part
	for i := 0; i < y-2; i++ {
		fmt.Print("*")
		for j := 0; j < x-2; j++ {
			fmt.Print(" ")
		}
		if x > 1 {
			fmt.Print("*")
		}
		fmt.Println()
	}

	// Draw the bottom edge
	if y > 1 {
		fmt.Print("\\")
		for i := 0; i < x-2; i++ {
			fmt.Print("*")
		}
		if x > 1 {
			fmt.Print("/")
		}
		fmt.Println()
	}
}
