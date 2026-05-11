package piscine

func TrimAtoi(s string) int {
	sign := 1
	num := 0
	foundDigit := false

	for _, char := range s {
		if char == '-' && !foundDigit {
			sign = -1
		} else if char >= '0' && char <= '9' {
			foundDigit = true
			num = num*10 + int(char-'0')
		}
	}

	return sign * num
}
