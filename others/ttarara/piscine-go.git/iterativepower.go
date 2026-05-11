package piscine

// IterativePower returns nb raised to the power of power
func IterativePower(nb int, power int) int {
	result := 1

	if power < 0 {
		return 0
	}

	for i := 0; i < power; i++ {
		result = result * nb
	}

	return result
}
