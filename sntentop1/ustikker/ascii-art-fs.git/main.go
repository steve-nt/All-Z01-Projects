package main

import (
	"flag"
	"fmt"
)

func main() {
	banner = "standard.txt"

	// Check if the argument order is valid
	if !validateArgsOrder() {
		fmt.Println("Invalid argument order: Flags must precede positional arguments.")
		return
	}

	// Parse the arguments for flags and options
	if !parseFlagsAndOptions() || !onlyflagsequal() {
		return
	}
	// Check the validity of the input
	if !checkValidity() {
		return
	}
	// Load the banner (ASCII art)
	asciiMap, asciiHeight := loadBanner(banner)
	// Process the input string
	processString(flag.Arg(iofInput), asciiMap, asciiHeight)

}
