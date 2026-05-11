package ascii

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// GetTerminalWidth returns the width of the terminal
func GetTerminalWidth() int {
	var terminalWidth int
	var err error
	if runtime.GOOS == "windows" {
		terminalWidth, _, err = GetTerminalSizeWin()
	} else {
		terminalWidth, _, err = GetTerminalSize()
	}
	if err != nil || terminalWidth <= 0 {
		terminalWidth = 80 // Default terminal width
	}
	return terminalWidth
}

// GetTerminalSize retrieves terminal size on Unix-like systems
func GetTerminalSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	size := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(size) != 2 {
		return 0, 0, fmt.Errorf("unexpected output from stty size")
	}

	height, err := strconv.Atoi(size[0])
	if err != nil {
		return 0, 0, err
	}
	width, err := strconv.Atoi(size[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

// GetTerminalSizeWin retrieves terminal size on Windows systems
func GetTerminalSizeWin() (int, int, error) {
	cmd := exec.Command("mode", "con")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	var width, height int
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Columns:") {
			widthStr := strings.TrimSpace(strings.TrimPrefix(line, "Columns:"))
			width, err = strconv.Atoi(widthStr)
			if err != nil {
				return 0, 0, err
			}
		} else if strings.HasPrefix(line, "Lines:") {
			heightStr := strings.TrimSpace(strings.TrimPrefix(line, "Lines:"))
			height, err = strconv.Atoi(heightStr)
			if err != nil {
				return 0, 0, err
			}
		}
	}
	if width == 0 || height == 0 {
		return 0, 0, fmt.Errorf("could not determine terminal size")
	}
	return width, height, nil
}
