package extra_functions

import (
	"bufio"
	"fmt"
	"os"
)

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
