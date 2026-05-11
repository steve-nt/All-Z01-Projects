package piscine

func CountIf(f func(string) bool, tab []string) int {
	count := 0

	for _, r := range tab {
		if f(r) == true {
			count++
		}
	}

	return count
}
