package piscine

func JumpOver(str string) string {
	runes := []rune(str)
	var str2 string

	if len(str) >= 3 {
		for i := 2; i < len(str); i = i + 3 {
			str2 += string(runes[i])
		}
	}
	return str2 + "\n"
}
