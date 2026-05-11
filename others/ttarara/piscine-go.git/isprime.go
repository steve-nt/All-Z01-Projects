package piscine

func IsPrime(nb int) bool {
	if nb <= 1 {
		return false
	}
	var count int
	for i := 1; i <= nb/2; i++ {
		if nb%i == 0 {
			count = count + 1
		}
	}
	if count == 1 {
		return true
	} else if count != 1 {
		return false
	}
	return false
}
