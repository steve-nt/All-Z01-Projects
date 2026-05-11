package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Function to generate quadA pattern
func QuadA(x, y int) {
	if x <= 0 || y <= 0 {
		return
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			if (i == 0 || i == y-1) && (j == 0 || j == x-1) {
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

// Function to generate quadB pattern
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

// Function to generate quadC pattern
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

// Function to generate quadD pattern
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

// Function to generate quadE pattern
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

func createGoMod() bool {
	_, err := os.Stat("go.mod")
	if os.IsNotExist(err) {
		cmd := exec.Command("go", "mod", "init", "quadchecker")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to create go.mod: %v", err)
		}
		return true
	}
	return false
}

func deleteGoMod() {
	if err := os.Remove("go.mod"); err != nil {
		log.Printf("Failed to delete go.mod: %v", err)
	}
}

func main() {
	if len(os.Args) > 1 {
		quadName := os.Args[1]

		// Check if the command is "build"
		if quadName == "build" {
			deleteGoMod()

			executableSuffix := ""
			if runtime.GOOS == "windows" {
				executableSuffix = ".exe"
			}

			commands := []string{
				fmt.Sprintf("go build -o quadchecker%s main.go", executableSuffix),
				fmt.Sprintf("go build -o quadA%s main.go", executableSuffix),
				fmt.Sprintf("go build -o quadB%s main.go", executableSuffix),
				fmt.Sprintf("go build -o quadC%s main.go", executableSuffix),
				fmt.Sprintf("go build -o quadD%s main.go", executableSuffix),
				fmt.Sprintf("go build -o quadE%s main.go", executableSuffix),
			}

			for _, cmd := range commands {
				if err := runCommand(cmd); err != nil {
					log.Fatalf("Command failed: %s\nError: %v", cmd, err)
				}
			}

			if runtime.GOOS != "windows" {
				chmodCommands := []string{
					"chmod +x quadA",
					"chmod +x quadB",
					"chmod +x quadC",
					"chmod +x quadD",
					"chmod +x quadE",
				}

				for _, cmd := range chmodCommands {
					if err := runCommand(cmd); err != nil {
						log.Fatalf("Command failed: %s\nError: %v", cmd, err)
					}
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

	goModCreated := createGoMod()
	if goModCreated {
		defer deleteGoMod()
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

	// Quadchecker mode to compare input with generated quads
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
			if runtime.GOOS == "windows" {
				quadCommands = []string{"quadA.exe", "quadB.exe", "quadC.exe", "quadD.exe", "quadE.exe"}
			}

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

		// Fallback for handling go run .
		if len(os.Args) == 1 {
			if inputStr != "" {
				quadCommands := []string{"./quadA", "./quadB", "./quadC", "./quadD", "./quadE"}
				if runtime.GOOS == "windows" {
					quadCommands = []string{"quadA.exe", "quadB.exe", "quadC.exe", "quadD.exe", "quadE.exe"}
				}

				width := 0
				height := 0
				for _, line := range inputLines {
					if len(line) > width {
						width = len(line)
					}
					height++
				}
				width = width - 1 // To handle the trailing newline character

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

