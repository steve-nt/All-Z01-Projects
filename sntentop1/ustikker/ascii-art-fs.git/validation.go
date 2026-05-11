package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"

	asciiart "asciiart/colour"
)

var (
	colorFlag   string // What color are we using?
	text2color  string // Which characters will be colored?
	outputFlag  string // Do we have an output file to save the result? Which one?
	alignFlag   string // Do we use a special alignment?
	reverseFlag string // Do we need to reverse a file? ATTENTION DO NOT DO THAT YET
	iofInput    int    // Which of the PFargs is the input
	banner      string // Which file to use as a banner
)

func onlyflagsequal() bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--color") && !strings.HasPrefix(arg, "--color=") {
			fmt.Println("Usage: go run . [OPTION] [STRING]")
			fmt.Println()
			fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
			return false
		}
		if strings.HasPrefix(arg, "--output") && !strings.HasPrefix(arg, "--output=") {
			fmt.Println("Error: --output flag must be in the form --output=value")
			return false
		}
		if strings.HasPrefix(arg, "--justify") && !strings.HasPrefix(arg, "--justify=") {
			fmt.Println("Error: --justify flag must be in the form --justify=value")
			return false
		}
	}
	return true
}

func parseFlagsAndOptions() bool {
	flag.StringVar(&colorFlag, "color", "", "Set the text color")
	flag.StringVar(&outputFlag, "output", "", "Output file")
	flag.StringVar(&alignFlag, "align", "left", "Text alignment (left, center, right)")
	flag.StringVar(&reverseFlag, "reverse", "", "Reverse this file")

	flag.Parse()

	// Case 1: No Color, No Substring, No Banner

	if (len(flag.Args()) == 1) && (colorFlag == "") {
		iofInput = 0
		banner = "standard.txt"
		return true
	}

	// Case 2: Color, No Substring, No Banner
	// To add: Color Check

	if (len(flag.Args()) == 1) && (colorFlag != "") {
		iofInput = 0
		text2color = flag.Arg(0)
		banner = "standard.txt"
		return true
	}

	// Case 3: No Color, No Substring, Banner

	if (len(flag.Args()) == 2) && (colorFlag == "") && isBanner(flag.Arg(1)) {
		iofInput = 0
		return true
	}

	// Case 4: Color, No substring, Banner

	if (len(flag.Args()) == 2) && (colorFlag != "") && isBanner(flag.Arg(1)) {
		iofInput = 0
		text2color = flag.Arg(0)
		return true
	}

	// Case 5: Color, Substring, No banner

	if (len(flag.Args()) == 2) && (colorFlag != "") && !isBanner(flag.Arg(1)) {
		iofInput = 1
		text2color = flag.Arg(0)
		return true
	}

	//Case 6: Color, Substring, Banner

	if len(flag.Args()) == 3 {
		iofInput = 1
		text2color = flag.Arg(0)
		banner = flag.Arg(2) + ".txt"
		return true
	}

	fmt.Println("\"", flag.Arg(1), "\"", "is not a valid banner, please use \"standard\", \"shadow\" or \"thinkertoy\".")
	return false
}

func checkValidity() bool {
	if isNotASCII(flag.Arg(iofInput + 1)) {
		fmt.Println("Error: Only ASCII characters or newline symbols (\\n) are allowed.")
		return false
	}

	// Check if the banner file exists
	if _, err := os.Stat(banner); err != nil {
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		fmt.Println()
		fmt.Println("EX: go run . something standard")
		return false
	}

	// Validate color format
	if colorFlag != "" {
		if _, err := asciiart.ColorToAnsi(colorFlag); err != nil {
			fmt.Printf("Error: Invalid color format '%s'.\n", colorFlag)
			return false
		}
	}

	return true
}

func isNotASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

func isBanner(s string) bool {
	if strings.ToLower(s) == "standard" || strings.ToLower(s) == "shadow" || strings.ToLower(s) == "thinkertoy" {
		banner = strings.ToLower(s) + ".txt"
		return true
	}
	return false
}

func validateArgsOrder() bool {
	hasPositional := false

	// Loop through all arguments after the program name
	for i, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			if hasPositional {
				// We found a flag after positional arguments, which is not allowed
				fmt.Printf("Error: Flag '%s' found after positional arguments at position %d.\n", arg, i+1)
				return false
			}
		} else {
			// A positional argument is encountered
			hasPositional = true
		}
	}

	return true
}
