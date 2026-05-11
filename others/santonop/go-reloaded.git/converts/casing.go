package converts

import (
	"fmt"
	"strconv"
	"strings"
)

func executeCommand(results []string, command string, num int, extra string) []string {
	if len(results) >= num {
		switch command {
		case "up":
			for i := 1; i <= num; i++ {
				results[len(results)-i] = strings.ToUpper(results[len(results)-i])
			}
		case "low":
			for i := 1; i <= num; i++ {
				results[len(results)-i] = strings.ToLower(results[len(results)-i])
			}
		case "cap":
			for i := 1; i <= num; i++ {
				results[len(results)-i] = capitalizeWord(results[len(results)-i])
			}
		}
	}
	results[len(results)-1] += extra
	return results
}
func checkCommandPrefix(word string) bool {
	return strings.HasPrefix(word, "(low") || strings.HasPrefix(word, "(cap") || strings.HasPrefix(word, "(up")
}
func handleCommand(word string) (string, int, string) {
	extra := ""
	indexEndBracket := strings.Index(word, ")")
	if len(word) != indexEndBracket+1 {
		extra = word[indexEndBracket+1:]
		word = word[:indexEndBracket]
	}
	word = strings.Trim(word, "()")
	words := strings.Split(word, ",")
	var num int
	command := words[0]
	if len(words) > 1 {
		var err error
		num, err = strconv.Atoi(words[1])
		if err != nil {
			fmt.Printf("Error in %s: %v", word, err)
		}
	} else {
		num = 1
	}
	return command, num, extra
}
func capitalizeWord(word string) string {
	if word[0] == '\'' && len(word) > 1 {
		word = "'" + strings.ToUpper(string(word[1])) + word[2:]
	} else {
		word = strings.ToUpper(string(word[0])) + word[1:]
	}
	return word
}
