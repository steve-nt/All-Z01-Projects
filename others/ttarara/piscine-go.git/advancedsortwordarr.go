package piscine

func AdvancedSortWordArr(a []string, f func(a, b string) int) {
	lenarray := 0
	array := a
	for i := range array {
		lenarray = i + 1
	}

	for i := 0; i < lenarray-1; i++ {
		for j := i + 1; j < lenarray; j++ {
			if f(array[i], array[j]) > 0 {
				array[i], array[j] = array[j], array[i]
			}
		}
	}
}
