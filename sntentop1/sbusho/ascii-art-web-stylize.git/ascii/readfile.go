package ascii

import (
	"bufio"
	"os"
)

// ReadFile opens the ASCII art font file and returns each line in a slice of strings
func ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file will be closed when this function returns.

	scanner := bufio.NewScanner(file)
	var lines []string
	lineCount := 8

	for scanner.Scan() {
		if lineCount == 8 {
			lineCount = 0
			continue
		}
		lines = append(lines, scanner.Text())
		lineCount++
	}
	return lines, nil
}
