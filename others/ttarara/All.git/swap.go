package piscine

func Swap(a *int, b *int) {
	c := *a
	d := *b
	*b = c
	*a = d
}
