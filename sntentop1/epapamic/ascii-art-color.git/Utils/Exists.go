package utils

func Exists(index int, sub []int) bool {
	for _, subIndex := range sub {
		if index == subIndex {
			return true
		}
	}
	return false
}
