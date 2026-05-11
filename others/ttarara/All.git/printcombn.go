package piscine

import "github.com/01-edu/z01"

// Η συνάρτηση print εκτυπώνει έναν ακέραιο ως χαρακτήρα.
func print(i int) {
	z01.PrintRune(rune(i) + '0')
}

// Η συνάρτηση printComma εκτυπώνει ένα κόμμα και κενό.
func printComma() {
	z01.PrintRune(',')
	z01.PrintRune(' ')
}

// Η συνάρτηση printNewline εκτυπώνει μια νέα γραμμή.
func printNewline() {
	z01.PrintRune('\n')
}

// Εκτύπωση συνδυασμών αριθμών (0 έως 9) χωρίς επανάληψη για n ψηφία
func PrintCombN(n int) {
	if n < 1 || n > 9 {
		return
	}
	tab := make([]int, n)
	for i := 0; i < n; i++ {
		tab[i] = i
	}
	max := 10 - n

	for tab[0] <= max {
		for i := 0; i < n; i++ {
			print(tab[i])
		}
		if tab[0] != max {
			printComma()
		}

		tab[n-1]++
		for i := n - 1; i > 0; i-- {
			if tab[i] >= 10-(n-i) {
				tab[i-1]++
				for j := i; j < n; j++ {
					tab[j] = tab[j-1] + 1
				}
			}
		}
	}
	printNewline()
}
