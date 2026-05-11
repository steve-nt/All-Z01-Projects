package piscine

import "fmt"

func QuadA(x, y int) {
	// If either x or y is non-positive, do nothing.
	if x <= 0 || y <= 0 {
		return
	}

	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if (i == 0 || i == y-1) && (j == 0 || j == x-1) {
				// Corners
				fmt.Print("o")
			} else if i == 0 || i == y-1 {
				// Top or bottom edge (excluding corners)
				fmt.Print("-")
			} else if j == 0 || j == x-1 {
				// Left or right edge (excluding corners)
				fmt.Print("|")
			} else {
				// Inside of the rectangle
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
