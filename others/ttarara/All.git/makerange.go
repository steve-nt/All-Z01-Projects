package piscine

func MakeRange(min, max int) []int {
	// Ελέγχουμε αν το min είναι μεγαλύτερο ή ίσο με το max
	if min >= max {
		return nil
	}

	// Υπολογίζουμε το μέγεθος της φέτας
	size := max - min

	// Αρχικοποιούμε τη φέτα με το κατάλληλο μέγεθος
	ans := make([]int, size)

	// Γεμίζουμε τη φέτα με τις τιμές από το min μέχρι το max-1
	for i := 0; i < size; i++ {
		ans[i] = min + i
	}

	return ans
}
