package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const usageMsg = `Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> "something"
EX: go run . --color=<color> <substring to be colored> "something"
EX: go run . --color=<color1> --color=<color2> <substr1> <substr2> "something"`

// colorFlags is a custom flag.Value that collects multiple --color= values.
type colorFlags []string

func (c *colorFlags) String() string {
	return strings.Join(*c, ",")
}

func (c *colorFlags) Set(value string) error {
	*c = append(*c, value)
	return nil
}

func main() {
	// Pre-scan: reject --color or --color value (space-separated).
	// The only accepted format is --color=<value>.
	for _, arg := range os.Args[1:] {
		if arg == "--color" || arg == "-color" {
			fmt.Println("Error: flag must use = syntax: --color=<value>")
			fmt.Println(usageMsg)
			os.Exit(1)
		}
	}

	var colors colorFlags
	flag.Var(&colors, "color", "color name or value (may be repeated)")
	flag.Usage = func() { fmt.Println(usageMsg) }
	flag.Parse()

	positionals := flag.Args()
	n := len(colors) // number of --color flags supplied
	p := len(positionals)

	// Positional resolution.
	//
	// Without substrings (whole string colored, or no color at all):
	//   p==1            : string
	//   p==2 && p < n+1 : string, banner   (not enough positionals to include substrings)
	//
	// With substrings (one per color flag):
	//   p==n+1          : substr1…substrN, string
	//   p==n+2          : substr1…substrN, string, banner

	var inputStr, font string
	font = "standard"
	hasSubstrs := false

	switch {
	case p == 1:
		inputStr = positionals[0]
	case p == 2 && (n == 0 || p < n+1):
		// No substrings: just string + banner.
		inputStr = positionals[0]
		font = positionals[1]
	case n > 0 && p == n+1:
		inputStr = positionals[n]
		hasSubstrs = true
	case n > 0 && p == n+2:
		inputStr = positionals[n]
		font = positionals[n+1]
		hasSubstrs = true
	default:
		fmt.Println(usageMsg)
		os.Exit(1)
	}

	if inputStr == "" {
		os.Exit(0)
	}

	// Load banner.
	bannerMap, err := LoadBanner("banners/" + font + ".txt")
	if err != nil {
		fmt.Println("Error loading font file:", err)
		os.Exit(1)
	}

	// Build ColorPair slice.
	pairs := make([]ColorPair, 0, n)
	for i, colorName := range colors {
		ansiCode, err := ColorCode(colorName)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println(usageMsg)
			os.Exit(1)
		}
		substr := ""
		if hasSubstrs {
			substr = positionals[i]
		}
		pairs = append(pairs, ColorPair{AnsiCode: ansiCode, Substr: substr})
	}

	// Split on literal \n and render.
	lines := strings.Split(inputStr, "\\n")
	RenderWithColor(lines, bannerMap, pairs)
}
