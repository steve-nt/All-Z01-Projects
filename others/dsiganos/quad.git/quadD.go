package piscine

// To function printRow yparxei sto quadB

func QuadD(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	// top
	printRow("A", "B", "C", x)

	// middle
	for i := 0; i < y-2; i++ {
		printRow("B", " ", "B", x)
	}

	// bottom
	if y > 1 {
		printRow("A", "B", "C", x)
	}
}
