package main

import (
	"io"
	"os"

	"github.com/01-edu/z01"
)

func main() {
	if len(os.Args) < 2 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			printString("ERROR: reading from stdin\n")
			os.Exit(1)
		}
		printString(string(content))
		return
	}

	for _, fileName := range os.Args[1:] {
		content, err := os.ReadFile(fileName)
		if err != nil {
			printString("ERROR: open " + fileName + ": no such file or directory\n")
			os.Exit(1)
		}
		printString(string(content))
	}
}

func printString(s string) {
	for _, r := range s {
		z01.PrintRune(r)
	}
}
