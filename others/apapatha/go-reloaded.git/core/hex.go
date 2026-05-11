package core

import (
	"log"
	"strconv"
)

func hex(str string) string {
	n, err := strconv.ParseInt(str, 16, 64)
	if err != nil {
		log.Fatal(err.Error())
	}
	return strconv.Itoa(int(n))
}
