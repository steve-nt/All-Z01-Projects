package piscine

func Map(f func(int) bool, arr []int) []bool {
	boolArr := make([]bool, len(arr))
	for i, v := range arr {
		boolArr[i] = f(v)
	}
	return boolArr
}
