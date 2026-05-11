package piscine

func IsPrintable(s string) bool {
	for _, char := range s {
		if !(char >= ' ' && char <= '~') {
			return false
		}
	}
	return true
}
