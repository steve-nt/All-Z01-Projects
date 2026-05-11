package main

import (
	"fmt"
	"go-reloaded/converts"
	"os"
)

func modify(str string) string {
	str = converts.FixSingleQuoteSpacing(str)
	str = converts.ProcessString(str)
	str = converts.FixSingleQuoteSpacing(str)
	return str
}
func createResult(result, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	// Write some text to the file
	_, err = file.WriteString(result)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	fmt.Println("File written successfully")
}
func main() {
	args := os.Args[1:]
	filename := args[0]
	outputname := args[1]
	// Read the entire file as a byte slice
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	result := modify(string(content))
	createResult(result, outputname)
}
