package piscine

import "fmt"

func QuadE(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}

	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if i == 0 && j == 0 {
				fmt.Print("A")
			} else if i == 0 && j == x-1 {
				fmt.Print("C")
			} else if i == y-1 && j == 0 {
				fmt.Print("C")
			} else if i == y-1 && j == x-1 {
				fmt.Print("A")
			} else if i == 0 || i == y-1 {
				fmt.Print("B")
			} else if j == 0 || j == x-1 {
				fmt.Print("B")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
