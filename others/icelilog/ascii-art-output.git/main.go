package main

import (
	"fmt"
	"os"
	"strings"

	asciiart "asciiart/Utils"
)

// this works!!

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: go run . [OPTIONS] [STRING] [BANNER] || Example: go run . \"test\" standard || Options: --output=, --align=")
		return
	}

	argStr := os.Args[1]
	var width int
	var align string
	thirdBanner := false
	var styleBanner string
	var outputFile string

	// checks if the first argument is either alignment or output
	if len(argStr) >= 8 && argStr[:2] == "--" {
		// checks if the first 8 characters are "--align="
		if argStr[:8] == "--align=" {

			// getTerminalSize() returns the terminal width, height, error
			width, _, _ = asciiart.GetTerminalSize()

			// stores the 8th character and forward to the align variable
			align = strings.ToLower(argStr[8:])
			if align == "" {
				fmt.Println("Missing align type!")
				return
			}
			// checks the value of align and specifies it's allowed values
			if align != "left" && align != "right" && align != "center" && align != "justify" {
				fmt.Println("Wrong align! (right, left, center, justify)")
				return
			}
			if len(os.Args) < 3 {
				fmt.Println("Missing string!")
				return
			}
			// assigns the third Os.arg to the argStr (since the second Os.arg was alignment)
			argStr = os.Args[2]

			// Indicate that the second argument is used as the banner style to be used on the string to be processed
			thirdBanner = true

			// Prevents the use of both --output= and --align= flags simultaneously
			if strings.HasPrefix(argStr, "--output=") {
				fmt.Println("Can't use output flag and align flag same time!")
				return
			}

			// checks if the first 8 characters are "--output="
		} else if argStr[:9] == "--output=" {
			outputFile = argStr[9:]
			if outputFile == "" {
				fmt.Println("Missing output name!")
				return
			}
			// checks if a string was given
			if len(os.Args) < 3 {
				fmt.Println("Missing string!")
				return
			}
			argStr = os.Args[2]

			// Indicate that the second argument is the file in which the string will be printed on
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

	// Determine banner style
	if len(os.Args) == 2 {
		styleBanner = "standard"
	} else if len(os.Args) == 3 {
		// Check if a flag was used
		if thirdBanner {
			styleBanner = "standard"
		} else {
			// If flag was not used, then the 2nd os.Args is the font style.
			styleBanner = strings.ToLower(os.Args[2])
		}
	} else if len(os.Args) == 4 {
		styleBanner = strings.ToLower(os.Args[3])
	} else {
		fmt.Println("Usage: go run . [OPTIONS] [STRING] [BANNER] || Example: go run . \"test\" standard || Options: --output=, --align=")
		return
	}

	// Handle empty input string and newline case
	if argStr == "" {
		return // No output for an empty string
	} else if argStr == "\\n" {
		fmt.Println() // Print just one newline for literal "\n"
		return
	}

	// Takes the string argStr and splits it wherever it finds the separator "\n" and
	// stores it in a slice of string
	sepArgs := strings.Split(argStr, "\\n")

	// stores what font file to use in file variable
	file, err := os.ReadFile("fonts/" + styleBanner + ".txt")
	if err != nil {
		fmt.Println(styleBanner + " banner does not exist.")
		return
	}

	// Normalize line endings for Windows file thinkertoy
	content := strings.ReplaceAll(string(file), "\r\n", "\n")
	// Now split the file content into lines and store it in a slice of strings
	lines := strings.Split(content, "\n")

	// Checks if alignment is specified, if yes then it sends it to the printAsciiArtAlign func
	// if you should put it into a file, then it sends it to the printAsciiArtToFile func
	// otherwise it just sends it to the printAsciiArt func
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
