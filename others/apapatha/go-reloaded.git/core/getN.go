package core

import (
	"log"
	"strconv"
)

func getN(str string) int {
	var runes []rune
	for _, r := range str {
		if r >= '0' && r <= '9' {
			runes = append(runes, r)
		}
	}
	if len(runes) <= 0 {
		log.Fatal("not valid command")
	}
	n, err := strconv.Atoi(string(runes))
	if err != nil {
		log.Fatal(err.Error())
	}
	return n
}
