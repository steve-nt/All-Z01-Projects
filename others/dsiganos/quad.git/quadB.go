package piscine

import "fmt"

func QuadB(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}

	// top row
	printRow("/", "*", "\\", x)

	// middle row
	for i := 0; i < y-2; i++ {
		printRow("*", " ", "*", x)
	}
	// bottom row
	if y > 1 {
		printRow("\\", "*", "/", x)
	}
}

func printRow(left, middle, right string, x int) {
	// left point
	fmt.Print(string(left))
	// middle point
	for i := 0; i < x-2; i++ {
		fmt.Print(string(middle))
	}
	// right point
	if x > 1 {
		fmt.Print(string(right))
	}
	fmt.Println()
}
