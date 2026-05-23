package main

import (
	"fmt"
)

func main() {
	numbers := readNumbers()

	avg := average(numbers)
	med := median(numbers)
	varr := variance(numbers, avg)
	std := stdDev(varr)

	fmt.Println("Average:", avg)
	fmt.Println("Median:", med)
	fmt.Println("Variance:", varr)
	fmt.Println("Standard Deviation:", std)
}
