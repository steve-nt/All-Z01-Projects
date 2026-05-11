package piscine

func IsSorted(f func(a, b int) int, tab []int) bool {
	ascending := true
	descending := true

	for i := 1; i < len(tab); i++ {
		if f(tab[i-1], tab[i]) > 0 {
			ascending = false
		}
		if f(tab[i-1], tab[i]) < 0 {
			descending = false
		}
	}

	return ascending || descending
}
