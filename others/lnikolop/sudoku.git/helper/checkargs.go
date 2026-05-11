package helper

import (
	"fmt"
	"os"
	"regexp"
)

func CheckArgs(args []string) {
	if len(args) <= 0 || len(args) > 9 {
		fmt.Println("Please enter 9 arguments as many as the rows of a sudoku")
		os.Exit(1)
	}
	for index, arg := range args {
		if IsValidSudokuLine(arg) {
			continue
		} else {
			fmt.Printf("Please enter as arg a valid sudoku line (e.g. 96.4...1.) in %v argument\n", index)
			os.Exit(1)
		}
	}
}

func IsValidSudokuLine(line string) bool {
	// Define a regular expression pattern to match a valid Sudoku line
	pattern := `^[1-9\.]{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(line)
}
