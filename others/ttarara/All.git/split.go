package piscine

func Split(s, sep string) []string {
	var result []string
	sepLen := len(sep)
	start := 0

	for i := 0; i <= len(s)-sepLen; {
		if s[i:i+sepLen] == sep {
			result = append(result, s[start:i])
			i += sepLen
			start = i
		} else {
			i++
		}
	}
	result = append(result, s[start:])
	return result
}
