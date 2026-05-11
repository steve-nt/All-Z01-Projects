package main

import (
	"fmt"
	"os"
	"strings"

	asciiart "asciiart/Utils"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: go run .  [OPTION] [STRING] [BANNER] || Example: go run . --align=right  something  standard")
		return
	}

	argStr := os.Args[1]
	var width int
	var align string
	thirdBanner := false
	var styleBanner string
	var outputFile string

	if len(argStr) >= 8 && argStr[:2] == "--" {
		if argStr[:8] == "--align=" {

			width, _, _ = asciiart.GetTerminalSize()
			if width == 0 {
				width, _, _ = asciiart.GetTerminalSizeWin()
			}

			align = strings.ToLower(argStr[8:])
			if align == "" {
				fmt.Println("Missing align type!")
				return
			}

			if align != "left" && align != "right" && align != "center" && align != "justify" {
				fmt.Println("Wrong align! (right, left, center, justify)")
				return
			}
			if len(os.Args) < 3 {
				fmt.Println("Missing string!")
				return
			}

			argStr = os.Args[2]

			thirdBanner = true

			if strings.HasPrefix(argStr, "--output=") {
				fmt.Println("Can't use output flag and align flag same time!")
				return
			}

		} else if argStr[:9] == "--output=" {
			outputFile = argStr[9:]
			if outputFile == "" {
				fmt.Println("Missing output name!")
				return
			}

			if len(os.Args) < 3 {
				fmt.Println("Missing string!")
				return
			}
			argStr = os.Args[2]

			thirdBanner = true
			if strings.HasPrefix(argStr, "--align=") {
				fmt.Println("Can't use output flag and align flag same time!")
				return
			}
		} else {
			fmt.Println("Wrong flag. (--output= || --align=)")
			return
		}
	}

	if len(os.Args) == 2 {
		styleBanner = "standard"
	} else if len(os.Args) == 3 {
		if thirdBanner {
			styleBanner = "standard"
		} else {
			styleBanner = strings.ToLower(os.Args[2])
		}
	} else if len(os.Args) == 4 {
		styleBanner = strings.ToLower(os.Args[3])
	} else {
		fmt.Println("Usage: go run .  [OPTION] [STRING] [BANNER] || Example: go run . --align=right  something  standard")
		return
	}

	if argStr == "" {
		return
	} else if argStr == "\\n" {
		fmt.Println()
		return
	}

	sepArgs := strings.Split(argStr, "\\n")

	file, err := os.ReadFile("fonts/" + styleBanner + ".txt")
	if err != nil {
		fmt.Println(styleBanner + " banner does not exist.")
		return
	}

	content := strings.ReplaceAll(string(file), "\r\n", "\n")

	lines := strings.Split(content, "\n")

	if align != "" {
		asciiart.PrintAsciiArtAlign(sepArgs, lines, align, width)
	} else if outputFile != "" {
		createdFile, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Something went wrong while creating output file.")
			return
		}
		defer createdFile.Close()
		asciiart.PrintAsciiArtToFile(sepArgs, lines, createdFile)
	} else {
		asciiart.PrintAsciiArt(sepArgs, lines)
	}
}


