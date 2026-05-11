package functions

import (
	"bufio"
	"fmt"
	"os"
)

func PrintAsciiArt(text string, file *os.File) {
	//Ascii values of text (space = 32, ! = 33, etc.)
	byteOfText := []byte(text)
	lines := make([][]string, len(byteOfText))

	for idx, val := range byteOfText {
		nbrVal := int(val-32) * 9

		file.Seek(0, 0) // Reset the scanner
		scanner := bufio.NewScanner(file)

		for lineNumber := 0; scanner.Scan(); lineNumber++ {
			if lineNumber >= nbrVal && lineNumber < nbrVal+9 {
				lines[idx] = append(lines[idx], scanner.Text())
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file: ", err)
		}
	}

	for lineNbr := 1; lineNbr <= 8; lineNbr++ {
		for _, val := range lines {
			if lineNbr < len(val) {
				fmt.Print(val[lineNbr])
			}
		}
		fmt.Println()
	}
}
