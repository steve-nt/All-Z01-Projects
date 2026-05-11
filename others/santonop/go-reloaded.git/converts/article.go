package converts

func processArticleCorrection(result []string, word string, words []string, i int) []string {
	if i+1 < len(words) {
		nextWordStartsWithVowel := isVowel(words[i+1][0])
		switch word {
		case "a":
			if nextWordStartsWithVowel {
				word = "an"
			}
		case "an":
			if !nextWordStartsWithVowel {
				word = "a"
			}
		case "A":
			if nextWordStartsWithVowel {
				word = "An"
			}
		case "An":
			if !nextWordStartsWithVowel {
				word = "A"
			}
		}
	}
	result = append(result, word)
	return result
}
func isVowel(char byte) bool {
	return char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' ||
		char == 'h' || char == 'A' || char == 'E' || char == 'I' || char == 'O' || char == 'U' || char == 'H'
}
