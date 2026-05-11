package check

import (
	"fmt"
	"strings"
)

var validChars = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '.'}

func CheckArgs(args []string) error {
	if len(args) != 9 {
		return fmt.Errorf("Error: there must be exactly 9 arguments")
	}
	for _, arg := range args {
		if len(arg) != 9 {
			return fmt.Errorf("Error: each argument must be exactly 9 characters long")
		}
		if strings.ContainsAny(arg, "0123456789.") == false {
			return fmt.Errorf("Error: each argument must only contain digits and dots")
		}
	}
	return nil
}
