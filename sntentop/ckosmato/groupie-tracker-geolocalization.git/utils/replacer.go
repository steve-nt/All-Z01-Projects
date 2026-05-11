package utils

import "strings"

func Replacer(word string) string {
	for _, v := range "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" {
		word = strings.ReplaceAll(word, string(v), " ")
	}
	return word
}
