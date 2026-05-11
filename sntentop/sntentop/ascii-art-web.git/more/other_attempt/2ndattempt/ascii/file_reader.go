// ascii/file_reader.go
package ascii

import (
	"bufio"
	"fmt"
	"os"
)

func ReadTxt(banner string) []string {
	// Construct the full path to the banner file
	filePath := fmt.Sprintf("./%s.txt", banner) // Or provide the full path if needed

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}
	return lines
}
