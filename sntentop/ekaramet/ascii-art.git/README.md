# ASCII Art Generator

This program generates ASCII art from a given input string using different banner styles. It reads a banner file to transform each character in the input string into corresponding ASCII art.

## Usage

```bash
go run ascii-art.go [STRING] [BANNER]

    STRING: The input string you want to convert into ASCII art.
    BANNER: (Optional) The banner style to use. If not provided, the default banner style is standard.

Example

bash

go run ascii-art.go "HELLO" standard

This command will output the ASCII representation of the word "HELLO" using the "standard" banner.
Special Characters

    Newline characters (\n) in the input string can be used to create multi-line ASCII art.

bash

go run ascii-art.go "HELLO\\nWORLD" standard

Banner Files

Banner files are plain text files containing ASCII art for characters from ASCII code 32 (space) to 126 (tilde ~). Each character in the file should be represented by 8 lines of ASCII art, followed by an empty line.

Banner files should be placed in a banners/ directory with the .txt extension. For example, a file standard.txt in the banners/ directory will be used when the standard banner is selected.
Available Banners

To see the list of available banners, provide an invalid banner name or omit the banner argument:

bash

go run ascii-art.go "HELLO" invalid-banner

How It Works

    Input Handling: The program takes an input string and an optional banner name.
    Banner Loading: It loads the corresponding banner file from the banners/ directory.
    ASCII Conversion: Each character in the input string is converted to ASCII art based on the loaded banner.
    Multi-line Input: The input string can contain newline characters (\n), which will be interpreted as line breaks.

Code Overview

    main(): Handles the program's execution, including parsing the command-line arguments and printing the ASCII art.
    loadBanner(): Reads the banner file and loads the ASCII art for each character into a map.
    getAvailableBanners(): Lists all available banner files in the banners/ directory.

Error Handling

    If the banner file does not exist, the program will print an error and list the available banners.
    If there are errors while reading the banner file (e.g., unexpected EOF), the program will terminate and show an error message.

Directory Structure

Your project should be structured like this:

bash

/path/to/project/
│
├── ascii-art.go
└── banners/
    ├── standard.txt
    ├── shadow.txt
    └── ... other banner files

Requirements

    Go programming language installed.

bash

go version

Ensure you have Go installed before running the program.

Feel free to modify the banner files or create your own by following the same structure!

vbnet


This documentation outlines how to use the program, the available banners, and the internal structure of the code. Let me know if you need any further customization!

