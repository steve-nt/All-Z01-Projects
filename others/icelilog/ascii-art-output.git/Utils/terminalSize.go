package asciiart

import (
	"os"      // Package os provides a platform-independent interface to operating system functionality.
	"os/exec" // Package exec runs external commands.
	"strconv" // Package strconv provides conversions to and from string representations of basic data types.
	"strings" // Package strings provides utility functions for manipulating UTF-8 encoded strings.
)

func GetTerminalSize() (int, int, error) {
	// Create a command to execute the `stty size` command in the shell.
	cmd := exec.Command("stty", "size") // This command retrieves the current size of the terminal.
	cmd.Stdin = os.Stdin                // Set the command's standard input to the terminal's standard input.
	out, err := cmd.Output()            // Execute the command and capture the output.
	if err != nil {
		return 0, 0, err // If there's an error, return 0 for width and height, along with the error.
	}

	// Split the output into parts; the output is expected to be two numbers: height and width.
	size := strings.Split(string(out), " ")
	width, err := strconv.Atoi(strings.TrimSpace(size[1])) // Convert the width from string to int.
	if err != nil {
		return 0, 0, err //// Return 0 for width and height, along with the error if conversion fails.
	}
	height, err := strconv.Atoi(strings.TrimSpace(size[0])) // Convert the height from string to int.
	if err != nil {
		return 0, 0, err
	}
	// Return the width, height, and a nil error.

	return width, height, nil
}
