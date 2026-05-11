package main

import (
	"ascii-art/functions"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var str string
	var colorhandler bool
	var color string
	var substring string

	inputArgs := os.Args[1:]
	banner := "banners/standard.txt"

	colorflag := strings.HasPrefix(inputArgs[0], "--color=")

	flag.StringVar(&color, "color", "", "coloring ascii chars")
	flag.Parse()

	flagArgs := flag.Args()

	if colorflag {
		colorhandler = true
	}

	if colorhandler {
		if color == "" {
			fmt.Println("Error: --color flag must be followed by a color value.")
			return
		}

		if !colorflag {
			printUsage()
			return

		}

		if len(flagArgs) > 1 && len(flagArgs) == 2 {
			substring = flagArgs[0]
			str = flagArgs[1]

		}
		if len(flagArgs) == 1 {

			str = flagArgs[0]
		}

		functions.PrintArt(str, banner, substring, color, colorhandler)
	}
	if !colorhandler {

		if len(inputArgs) < 1 || len(inputArgs) > 2 {
			fmt.Print("Usage: go run . [STRING] [BANNER]\nEX: go run . something standard")
			return
		}
		// Input Text
		str = inputArgs[0]

		str = strings.ReplaceAll(str, `\n`, "\n")

		if len(inputArgs) > 1 && inputArgs[1] != "" {
			banner = "banners/" + inputArgs[1] + ".txt"
		}

		functions.PrintArt(str, banner, substring, color, colorhandler)
	}
}

func printUsage() {
	fmt.Println("Usage: go run . [OPTION] [STRING]")
	fmt.Println("EX: go run . --color=<color> <substring to be colored> something")
}
