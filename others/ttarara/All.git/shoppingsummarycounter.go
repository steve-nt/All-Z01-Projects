package piscine

func ShoppingSummaryCounter(str string) map[string]int {
	summary := make(map[string]int)
	word := ""

	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			word += string(str[i])
		} else {
			summary[word]++
			word = ""
		}
	}
	// Add the last word if there is one
	summary[word]++

	return summary
}
