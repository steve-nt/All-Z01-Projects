package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

const (
	leftAlign    = "left"
	centerAlign  = "center"
	rightAlign   = "right"
	justifyAlign = "justify"
)

func main() {
	var alignment string
	flag.StringVar(&alignment, "align", "left", "text alignment (left, center, right, justify)")
	errorA := "Usage: go run . [OPTION] [STRING] [BANNER] \n\nExample: go run . --align=right something standard"
	errorB := "Error: alignment must be left, center, right, or justify."
	errorC := "Error: input string must be within the range of ASCII characters."

	if er1 := checkArgs(); er1 != nil {
		fmt.Println(errorA)
		os.Exit(0)
	}

	if er2 := checkForAlign(); er2 != nil {
		fmt.Println(errorA)
		os.Exit(0)
	}

	if er3 := isValidAlignment(alignment); er3 != nil {
		fmt.Println(errorB)
		os.Exit(0)
	}

	args := flag.Args()
	userInput := args[0]
	if er4 := isASCII(userInput); er4 != nil {
		fmt.Println(errorC)
		os.Exit(0)
	}

	bannerType := getBannerType(os.Args)
	ascii := mapFont(bannerType)

	terminalWidth, err := terminalWidth()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}

	printOutput(strings.Split(userInput, "\\n"), ascii, terminalWidth, alignment)
}

// 1 - check posa args exoume sta args
func checkArgs() error {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		return errors.New("invalid number of arguments")
	}
	return nil
}

// 2 - check gia flags alignment sta args
func checkForAlign() error {
	if string(os.Args[1]) == "--align" || strings.HasPrefix(string(os.Args[1]), "-align") || strings.HasPrefix(string(os.Args[1]), "-align=") {
		return errors.New("not a valid flag")
	}
	flag.Parse()
	return nil
}

// 3 - check an einai valid to align
func isValidAlignment(align string) error {
	if align != leftAlign && align != centerAlign && align != rightAlign && align != justifyAlign {
		return errors.New("invalid value")
	}
	return nil
}

// 4 - check input string
func isASCII(s string) error {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return errors.New("no ascii character")
		}
	}
	return nil
}

// 5 - check gia stye banner sta args
func getBannerType(args []string) string {
	if len(args) > 1 && (args[len(args)-1] == "shadow" || args[len(args)-1] == "thinkertoy") {
		return args[len(args)-1]
	}
	return "standard"
}

func mapFont(fileName string) map[rune][]string {
	file, err := os.Open("banner/" + fileName + ".txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer file.Close()
	asciiArr := parseFile(file)
	var asciiStart rune = 32
	ascii := make(map[rune][]string)
	for i, char := range asciiArr {
		ascii[rune(i+int(asciiStart))] = char
	}
	return ascii
}

func parseFile(file *os.File) [][]string {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var asciiChar []string
	var asciiArr [][]string
	counter := 0
	for fileScanner.Scan() {
		if counter == 8 {
			asciiChar = append(asciiChar, fileScanner.Text())
			asciiArr = append(asciiArr, asciiChar)
			asciiChar = []string{}
			counter = 0
			continue
		}
		counter++
		asciiChar = append(asciiChar, fileScanner.Text())
	}
	if err := fileScanner.Err(); err != nil {
		fmt.Println(err)
	}
	return asciiArr
}

// 999 - before the end
func terminalWidth() (int, error) {
	var dimensions [2]uint16
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)))
	if err != 0 {
		return 0, fmt.Errorf("error getting terminal size: %v", err)
	}
	return int(dimensions[1]), nil
}

func printOutput(words []string, ascii map[rune][]string, terminalWidth int, align string) {
	var alignment string
	wordsPerLine := 0
	for index, word := range words {
		wordLength := 0
		for _, runes := range word {
			if runes == ' ' && align == justifyAlign {
				wordsPerLine++
			}
			wordLength = wordLength + len(ascii[runes][4])
		}
		if wordLength > terminalWidth {
			fmt.Println("Words don't fit in terminal.")
			os.Exit(0)
		}
		switch align {
		case centerAlign:
			alignment = strings.Repeat(" ", (terminalWidth-wordLength)/2)
		case rightAlign:
			alignment = strings.Repeat(" ", terminalWidth-wordLength)
		case justifyAlign:
			if wordsPerLine == 0 {
				align = "none"
			} else {
				alignment = strings.Repeat(" ", (terminalWidth-wordLength)/wordsPerLine)
			}
		}
		for i := 0; i <= 8; i++ {
			for j, runes := range word {
				if j == 0 && align != justifyAlign {
					fmt.Print(alignment)
				}
				if align == justifyAlign && runes == ' ' {
					fmt.Print(alignment)
				}
				fmt.Print(ascii[runes][i])
			}
			if i == 8 && index != len(words)-1 {
				continue
			}
			fmt.Println()
		}
		wordsPerLine = 0
	}
}
