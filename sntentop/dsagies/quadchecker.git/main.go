package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func QuadA(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if (i == 0 || i == y-1) && (j == 0 || j == x-1) {
				// Corners
				fmt.Print("o")
			} else if i == 0 || i == y-1 {
				fmt.Print("-")
			} else if j == 0 || j == x-1 {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func QuadB(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if (i == 0 && j == 0) || (i == y-1 && j == x-1) {
				fmt.Print("/")
			} else if (i == 0 && j == x-1) || (i == y-1 && j == 0) {
				fmt.Print("\\")
			} else if i == 0 || i == y-1 || j == 0 || j == x-1 {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func QuadC(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if i == 0 && (j == 0 || j == x-1) {
				fmt.Print("A")
			} else if i == y-1 && (j == 0 || j == x-1) {
				fmt.Print("C")
			} else if i == 0 || i == y-1 {
				fmt.Print("B")
			} else if j == 0 || j == x-1 {
				fmt.Print("B")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func QuadD(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}

	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if i == 0 && j == 0 {
				fmt.Print("A")
			} else if i == 0 && j == x-1 {
				fmt.Print("C")
			} else if i == y-1 && j == 0 {
				fmt.Print("A")
			} else if i == y-1 && j == x-1 {
				fmt.Print("C")
			} else if i == 0 || i == y-1 {
				fmt.Print("B")
			} else if j == 0 || j == x-1 {
				fmt.Print("B")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func QuadE(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if i == 0 && j == 0 {
				fmt.Print("A")
			} else if i == 0 && j == x-1 {
				fmt.Print("C")
			} else if i == y-1 && j == 0 {
				fmt.Print("C")
			} else if i == y-1 && j == x-1 {
				fmt.Print("A")
			} else if i == 0 || i == y-1 {
				fmt.Print("B")
			} else if j == 0 || j == x-1 {
				fmt.Print("B")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func runCommand(cmd string) error {
	parts := strings.Split(cmd, " ")
	command := exec.Command(parts[0], parts[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func main() {
	if len(os.Args) > 1 {
		quadName := os.Args[1]

		if quadName == "build" {
			commands := []string{
				"go build -o quadA main.go", "chmod +x quadA",
				"go build -o quadB main.go", "chmod +x quadB",
				"go build -o quadC main.go", "chmod +x quadC",
				"go build -o quadD main.go", "chmod +x quadD",
				"go build -o quadE main.go", "chmod +x quadE",
				"go build -o quadchecker main.go",
			}

			for _, cmd := range commands {
				if err := runCommand(cmd); err != nil {
					log.Fatalf("Command failed: %s\nError: %v", cmd, err)
				}
			}
			return
		}

		// Generate quads based on arguments
		if len(os.Args) == 4 {
			width, _ := strconv.Atoi(os.Args[2])
			height, _ := strconv.Atoi(os.Args[3])
			switch quadName {
			case "quadA":
				QuadA(width, height)
			case "quadB":
				QuadB(width, height)
			case "quadC":
				QuadC(width, height)
			case "quadD":
				QuadD(width, height)
			case "quadE":
				QuadE(width, height)
			default:
				fmt.Println("Unknown quad name")
			}
			return
		}
	}

	// Determine executable name
	executableName := filepath.Base(os.Args[0])
	if executableName != "quadchecker" && executableName != "main" && len(os.Args) == 3 {
		width, _ := strconv.Atoi(os.Args[1])
		height, _ := strconv.Atoi(os.Args[2])
		switch executableName {
		case "quadA":
			QuadA(width, height)
		case "quadB":
			QuadB(width, height)
		case "quadC":
			QuadC(width, height)
		case "quadD":
			QuadD(width, height)
		case "quadE":
			QuadE(width, height)
		default:
			fmt.Println("Unknown quad name")
		}
		return
	}

	// Quadchecker
	if executableName == "quadchecker" || executableName == "main" {
		reader := bufio.NewReader(os.Stdin)
		var inputLines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			inputLines = append(inputLines, line)
		}
		inputStr := strings.Join(inputLines, "")
		inputStr = strings.TrimSpace(inputStr)

		if len(os.Args) == 3 {
			width := os.Args[1]
			height := os.Args[2]

			quadCommands := []string{"./quadA", "./quadB", "./quadC", "./quadD", "./quadE"}

			matches := []string{}
			for _, quadCmd := range quadCommands {
				cmd := exec.Command(quadCmd, width, height)
				output, err := cmd.Output()
				if err != nil {
					continue
				}

				if inputStr == strings.TrimSpace(string(output)) {
					matches = append(matches, fmt.Sprintf("[%s] [%s] [%s]", filepath.Base(quadCmd), width, height))
				}
			}

			if len(matches) > 0 {
				fmt.Println(strings.Join(matches, " || "))
			} else {
				fmt.Println("Not a quad function")
			}
			return
		}

		// handling go run .
		if len(os.Args) == 1 {
			if inputStr != "" {
				quadCommands := []string{"./quadA", "./quadB", "./quadC", "./quadD", "./quadE"}

				width := 0
				height := 0
				for _, line := range inputLines {
					if len(line) > width {
						width = len(line)
					}
					height++
				}
				width = width - 1 // handle the trailing newline character

				matches := []string{}
				for _, quadCmd := range quadCommands {
					cmd := exec.Command(quadCmd, strconv.Itoa(width), strconv.Itoa(height))
					output, err := cmd.Output()
					if err != nil {
						continue
					}

					if inputStr == strings.TrimSpace(string(output)) {
						matches = append(matches, fmt.Sprintf("[%s] [%d] [%d]", filepath.Base(quadCmd), width, height))
					}
				}

				if len(matches) > 0 {
					fmt.Println(strings.Join(matches, " || "))
				} else {
					fmt.Println("Not a quad function")
				}
			} else {
				fmt.Println("Not a quad function")
			}
			return
		}
	}

	fmt.Println("Usage: ./quadchecker build || ./<quadName> <width> <height> || ./quadchecker <width> <height>")
}
