package main

import (
	"bufio"
	"fmt"
	"os"
)

// LoadASCIIBanner reads and returns the ASCII banner file as a 2D array.
func LoadASCIIBanner(banner string) [][]string {
	fileLines := ReadTxt(banner)
	if fileLines == nil {
		return nil
	}

	// Convert file lines into a 2D array for ASCII art characters
	return Return2dASCIIArray(fileLines)
}

// ReadTxt reads the specified banner file and returns its contents as a slice of strings.
func ReadTxt(banner string) []string {
	fileName := banner + ".txt"
	readFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	return fileLines
}

// Return2dASCIIArray converts lines read from the file into a 2D array of ASCII templates.
func Return2dASCIIArray(fileLines []string) [][]string {
	var asciiTemplates [][]string
	counter := 0
	var tempAsciArray []string

	for _, line := range fileLines {
		counter++
		if counter != 1 {
			tempAsciArray = append(tempAsciArray, line)
		}
		if counter == 9 {
			asciiTemplates = append(asciiTemplates, tempAsciArray)
			counter = 0
			tempAsciArray = nil
		}
	}
	return asciiTemplates
}
