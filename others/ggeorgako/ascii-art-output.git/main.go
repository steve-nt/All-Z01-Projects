package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args

	var filename string
	var banner string
	var word string

	if len(args) == 2 {
		word = args[1]
		filename = "output.txt"
		banner = "standard"
	} else if len(args) == 4 {
		filename = args[1]
		word = args[2]
		banner = args[3]
	} else {
		fmt.Println("Usage go run . <OPTION] <string> <banner>")
		return
	}

	file, err := os.ReadFile("banners/" + banner + ".txt")

	if err != nil {
		fmt.Println("error reading file")
	}

	ascii := printAscii(word, file)

	fileName := strings.Split(filename, "=")

	if len(fileName) > 1 {
		os.WriteFile("writes/"+fileName[1], ascii, 0777)
	} else {
		os.WriteFile("writes/"+filename, ascii, 0777)
	}

	fmt.Printf("File wrote successfully on disk folder writes\n")
}

func printAscii(words string, file []byte) []byte {
	output := []byte{}

	cleanedFile := strings.ReplaceAll(string(file), "\r", "")
	
	lines := strings.Split(cleanedFile, "\n")

	// Split by newlines to handle multiple lines of input
	inputLines := strings.Split(words, "\\n")

	for _, line := range inputLines {
		// ASCII art has 8 lines per character
		for h := 1; h < 9; h++ {
			for _, char := range line {
				// Calculate the ASCII index in banner file
				index := int(char-32) * 9
				if index+h < len(lines) {
					output = append(output, []byte(lines[index+h])...)
				}
			}
			// New line after each line of ASCII art
			if !(h == 0 && len(output) == 0) {
				output = append(output, '\n')
			}
		}

	}

	return output
}
