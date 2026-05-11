package utils

import "fmt"

// Prints the formated string stored inside the asciiMap.
func (m *asciiMap) PrintAsciiMapCharacters() {

	switch {
	case m.color != "" && m.substring != "":
		fmt.Print(m.printContent + "\033[37m") // resets back to white
	case m.color != "":
		fmt.Print(m.color + m.printContent + "\033[37m") // resets back to white
	default:
		fmt.Print(m.printContent)
	}
}
