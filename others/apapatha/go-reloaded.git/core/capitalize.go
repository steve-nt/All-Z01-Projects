package core

func capitalize(str string) string {
	out := []rune(str)
	r := out[0]
	if r >= 'a' && r <= 'z' {
		out[0] -= alphaOffset
	}
	return string(out)
}
