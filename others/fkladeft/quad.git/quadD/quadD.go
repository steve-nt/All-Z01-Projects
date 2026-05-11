package main

import "fmt"

func QuadD(x, y int) {
	// Check if the dimensions are positive
	if x <= 0 || y <= 0 {
		return
	}

	// Draw the top edge
	fmt.Print("A")
	for i := 0; i < x-2; i++ {
		fmt.Print("B")
	}
	if x > 1 {
		fmt.Print("C")
	}
	fmt.Println()

	// Draw the middle part
	for i := 0; i < y-2; i++ {
		fmt.Print("B")
		for j := 0; j < x-2; j++ {
			fmt.Print(" ")
		}
		if x > 1 {
			fmt.Print("B")
		}
		fmt.Println()
	}

	// Draw the bottom edge
	if y > 1 {
		fmt.Print("A")
		for i := 0; i < x-2; i++ {
			fmt.Print("B")
		}
		if x > 1 {
			fmt.Print("C")
		}
		fmt.Println()
	}
}
