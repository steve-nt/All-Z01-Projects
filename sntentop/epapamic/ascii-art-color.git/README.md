# ASCII Art Generator

A Go program that transforms strings into stylized ASCII art with optional coloring and substring highlighting.

---

## 🧾 Description

This tool takes an input string and renders it using ASCII art fonts from a banner file (`standard.txt`). It allows you to apply color to the entire string or a specific substring.

It also handles multi-line inputs using escape sequences like `\n`.

---

## ✅ Features

- Converts text to multi-line ASCII art
- Supports full string coloring or substring highlighting
- Handles escape sequences like `\n` for multiline output
- Validates all characters to ensure they're printable ASCII
- Verifies file integrity using a SHA-256 hash

---

## ⚙️ Requirements

- Go installed
- File: `Data/standard.txt` must exist in the project directory

---

## 📦 Installation

1. Clone the repository:
```bash
git clone https://platform.zone01.gr/git/epapamic/color.git
cd color
```

2. Ensure banner file is present in the project directory color/Data/standard.txt

3. Ensure that you give a color with at least one string 

## Usage

```bash
Usage: go run . [OPTION] [STRING]
```

Example:
```bash
EX: go run . --color=<color> <substring to be colored> "something"
```

## Project Structure

```
color/
├── App/
│   └──  main.go
├── Data/
│   └──  standard.txt
├── Utils/
│   ├── AsciiOutput.go
│   ├── AsciiOutput_test.go
│   ├── CheckFileIntegrity.go
│   ├── CheckFIleIntegrity_test.go
│   ├── CreateMap.go
│   ├── CreateMap_test.go
│   ├── Exists.go
│   ├── findSubstringIndexes.go
│   ├── getAnsiColor.go
│   ├── OpenMap.go
│   ├── parseInput.go
│   ├── PrintAsciiMapCharacters.go
│   ├── PrintAsciiMapCharacters_test.go
│   ├── ValidatePrintable.go
│   └── ValidatePrintable_test.go
└── README.md
```
