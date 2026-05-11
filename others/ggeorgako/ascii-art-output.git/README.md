# ASCII Art Banner Generator

This project is a command-line tool written in Go to generate ASCII art banners in various styles. It accepts a string input, applies a chosen ASCII art style, and outputs the result to a specified file.

## Features
- **ASCII Art Generation**: Supports different banner styles (e.g., `standard`, `shadow`).
- **File Output**: Saves the ASCII art output to a file specified by the `--output=<fileName.txt>` flag.
- **Usage Verification**: Displays a usage message if arguments are not correctly formatted.
- **Extendability**: Easily add more ASCII art styles by implementing additional options.

## Requirements
- Go 1.16 or higher

## Usage

To run the program, use the following command format:

```bash
go run . --output=<fileName.txt> [STRING] [BANNER]
