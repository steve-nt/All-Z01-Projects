package utils

import "fmt"

func (m *asciiMap) GetAnsiColor(color string) error {
	switch color {
	case "black":
		m.color = "\033[30m"
	case "red":
		m.color = "\033[31m"
	case "green":
		m.color = "\033[32m"
	case "yellow":
		m.color = "\033[33m"
	case "blue":
		m.color = "\033[34m"
	case "magenta":
		m.color = "\033[35m"
	case "cyan":
		m.color = "\033[36m"
	case "orange":
		m.color = "\033[38;5;214m"
	case "":
		return fmt.Errorf("insert color")
	default:
		return fmt.Errorf("color %s not supported", m.color)
	}

	return nil
}
