package piscine

// IterativePower returns nb raised to the power of power
func RecursivePower(nb int, power int) int {
	if power < 0 {
		return 0
	} else if power == 0 {
		return 1
	} else if power == 1 {
		return nb
	} else if nb == 0 {
		return 0
	}

	return nb * RecursivePower(nb, power-1)
}
