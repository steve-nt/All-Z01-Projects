package converts

func processPunctuations(result []string, word string) []string {
	for len(word) > 0 && isPunctuation(rune(word[0])) {
		// Append punctuation to the last word in the result
		result[len(result)-1] += string(word[0])
		word = word[1:]
	}
	// Append the remaining part of the word if it is not empty
	if word != "" {
		result = append(result, word)
	}
	return result
}
func isPunctuation(char rune) bool {
	punctuations := []rune{'.', ',', '!', '?', ':', ';'}
	for _, p := range punctuations {
		if p == char {
			return true
		}
	}
	return false
}
