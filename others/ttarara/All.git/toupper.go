package piscine

func ToUpper(s string) string {
	result := []rune(s)
	for i, char := range result {
		if char >= 'a' && char <= 'z' {
			result[i] = char - ('a' - 'A')
		}
	}
	return string(result)
}
