package main

import (
	"fmt"
	"os"
)

func main() {
	arguments := os.Args[1:]

	if len(arguments) > 1 {
		fmt.Println("Too many arguments")
		return
	}
	if len(arguments) == 0 {
		fmt.Println("File name missing")
		return
	}

	filename := arguments[0]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(string(content))
}
