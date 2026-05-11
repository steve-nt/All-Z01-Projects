package utils

import (
	"flag"
	"fmt"
	"os"
)

func (m *asciiMap) ParseInput(color string) error {
	userInput := os.Args[1:]
	switch {
	case len(userInput) == 1:
		m.input = userInput[0]

	case len(userInput) == 2:
		m.color = color
		m.input = userInput[1]

	case len(userInput) == 3 && len(flag.Args()) != 1:
		m.color = color
		m.substring = userInput[1]
		m.input = userInput[2]
	case len(userInput) < 1:
		return fmt.Errorf("should have at least one argument")
	case len(userInput) > 3:
		return fmt.Errorf("error, too many arguments")
	default:
		return fmt.Errorf("\nUsage: go run . [OPTION] [STRING]\nEX: go run . --color=<color> <substring to be colored> `something`")
	}
	return nil
}
