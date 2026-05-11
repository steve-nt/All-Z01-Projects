package ascii

import (
	"fmt"
	"os"
	"strings"
)

// OnlyFlagsEqual ensures flags are in correct format
func OnlyFlagsEqual() bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--color") && !strings.HasPrefix(arg, "--color=") {
			fmt.Println("Usage: go run . [OPTION] [STRING]")
			fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
			return false
		}
		if strings.HasPrefix(arg, "--output") && !strings.HasPrefix(arg, "--output=") {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
			fmt.Println("EX: go run . --output=<fileName.txt> something standard")
			return false
		}
		if strings.HasPrefix(arg, "--align") && !strings.HasPrefix(arg, "--align=") {
			fmt.Println("Error: --align flag must be in the form --align=<type>")
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
			fmt.Println("Example: go run . --align=justify something standard")
			return false
		}
	}
	return true
}

// FindSubstringIndices finds substring indices in text
func FindSubstringIndices(text, substring string) []int {
	var indices []int
	textLen := len(text)
	subLen := len(substring)
	for i := 0; i <= textLen-subLen; i++ {
		if text[i:i+subLen] == substring {
			indices = append(indices, i)
		}
	}
	return indices
}
