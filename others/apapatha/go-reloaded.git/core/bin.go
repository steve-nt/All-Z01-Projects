package core

import (
	"log"
	"strconv"
)

func bin(str string) string {
	n, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		log.Fatal(err.Error())
	}
	return strconv.Itoa(int(n))
}
