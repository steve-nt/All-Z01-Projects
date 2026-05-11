package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
