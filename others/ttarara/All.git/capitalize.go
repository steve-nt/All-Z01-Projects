package piscine

func Capitalize(s string) string {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if (i == 0 || !((runes[i-1] >= 'A' && runes[i-1] <= 'Z') || (runes[i-1] >= 'a' && runes[i-1] <= 'z') || (runes[i-1] >= '0' && runes[i-1] <= '9'))) && ((runes[i] >= 'a' && runes[i] <= 'z') || (runes[i] >= 'A' && runes[i] <= 'Z') || (runes[i] >= '0' && runes[i] <= '9')) {
			if runes[i] >= 'a' && runes[i] <= 'z' {
				runes[i] -= 'a' - 'A'
			}
		} else if runes[i] >= 'A' && runes[i] <= 'Z' {
			runes[i] += 'a' - 'A'
		}
	}
	return string(runes)
}
