package funcs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	rooms            []string
	connections      map[string][]string
	tunnels          [][2]string
	totalAnts        int
	start, end, mode string
	startFound       = false
	endFound         = false
	firstLine        = true
	maxFlowPaths     [][]string
)

func ParseInput() {

	if len(os.Args) != 2 {
		fmt.Println("[ERROR]: wrong argument syntax")
		os.Exit(0)
	}

	inputFile := os.Args[1]

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("[ERROR]: failed to open input file")
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if firstLine {
			totalAnts, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("[ERROR]: number of ants must be digit")
				os.Exit(0)
			} else if totalAnts == 0 {
				fmt.Println("[ERROR]: zero ant count")
				os.Exit(0)
			}
			firstLine = false
			continue
		}

		switch line {
		case "##start":
			mode = "start"
			if !startFound {
				startFound = true
			} else {
				fmt.Println("[ERROR]: second start room found")
				os.Exit(0)
			}
			continue
		case "##end":
			mode = "end"
			if !endFound {
				endFound = true
			} else {
				fmt.Println("[ERROR]: second end room found")
				os.Exit(0)
			}
			continue
		}

		if fields := strings.Fields(line); len(fields) == 3 {

			name := fields[0]

			if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
				fmt.Println("[ERROR]: room will never start with the letter L or with #")
				os.Exit(0)
			}

			if roomExists(name, rooms) {
				fmt.Printf("[ERROR]: duplicate room name detected (%s)\n", name)
				os.Exit(0)
			}

			rooms = append(rooms, fields[0])

			if mode == "start" {
				start = fields[0]
			} else if mode == "end" {
				end = fields[0]
			}
			mode = ""
		} else if parts := strings.Split(line, "-"); len(parts) == 2 {

			a, b := parts[0], parts[1]

			if tunnelExists(a, b, tunnels) {
				fmt.Printf("[ERROR]: duplicate tunnel detected (%s-%s)\n", a, b)
				os.Exit(0)
			}

			if !roomExists(a, rooms) || !roomExists(b, rooms) {
				fmt.Printf("[ERROR]: tunnel references undefined room (%s-%s)\n", a, b)
				os.Exit(0)
			}

			tunnels = append(tunnels, [2]string{a, b})

			if a == b {
				fmt.Printf("[ERROR]: invalid tunnel - self linking (%s-%s)\n", a, b)
				os.Exit(0)
			}
		} else if strings.HasPrefix(line, "#") {
			continue
		} else {
			fmt.Println("[ERROR]: invalid text format in txt file")
			os.Exit(0)
		}

	}

	if start == "" {
		fmt.Println("[ERROR]: start room not found")
		os.Exit(0)
	}

	if end == "" {
		fmt.Println("[ERROR]: end room not found")
		os.Exit(0)
	}
}

// roomExists reports whether name is already present in the rooms slice
func roomExists(name string, rooms []string) bool {
	for _, r := range rooms {
		if r == name {
			return true
		}
	}
	return false
}

// tunnelExists reports whether a–b or b–a is already in tunnels.
func tunnelExists(a, b string, tunnels [][2]string) bool {
	for _, t := range tunnels {
		if (t[0] == a && t[1] == b) || (t[0] == b && t[1] == a) {
			return true
		}
	}
	return false
}

// prints the original txt from input
func PrintFile() {
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println()
}
