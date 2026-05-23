# ASCII Art FS

## Description

**ASCII Art FS** is a Go program that receives a string and a banner template name, then prints the string as large ASCII-art text.

This project extends the original ASCII-art project by allowing the user to choose which banner template should be used.

The program supports:

- Letters
- Numbers
- Spaces
- Special characters
- Literal `\n` sequences for line breaks
- Multiple banner templates
- Custom user-made templates

---

## Objectives

The goal of this project is to work with the Go file system API and data manipulation by loading ASCII-art templates from files.

The program follows this usage format:

```bash
go run . [STRING] [BANNER]
```

Example:

```bash
go run . "hello" standard
```

The program can also run with only a string argument. In that case, it uses the `standard` banner by default:

```bash
go run . "hello"
```

---

## Project Structure

```text
ascii-art-fs/
│
├── go.mod
├── main.go
├── main_test.go
├── README.md
│
├── ascii/
│   ├── ascii_art.go
|   └── ascii_art_test.go
│
├── docs/
│   ├── Milestones.txt
│   └── PRD.md
|    |__ Golden_tests.md
│
└── banners/
    ├── standard.txt
    ├── shadow.txt
    ├── thinkertoy.txt
    └── lineart.txt
```

### File Responsibilities

| File | Purpose |
|---|---|
| `main.go` | Handles command-line arguments and program flow |
| `ascii/ascii_art.go` | Loads banner files and renders ASCII art |
| `main_test.go` and `ascii/ascii_art_test.go` | Test files for the two go files the project needs |
| `banners/standard.txt` | Standard banner template |
| `banners/shadow.txt` | Shadow banner template |
| `banners/thinkertoy.txt` | Thinkertoy banner template |
| `banners/lineart.txt` | Custom banner template |

---

## Banner Format

Each banner file contains graphical representations for printable ASCII characters from character `32` (space) to character `126` (`~`).

Each character is represented by:

- 8 lines of ASCII art
- 1 separator line between character blocks

This means each character block takes 9 lines in the banner file.

The program reads the requested banner file from the `banners/` directory and maps each printable ASCII character to its corresponding 8-line representation.

---

## Usage

### Default Banner

If no banner is provided, the program uses `standard`.

```bash
go run . "hello"
```

Equivalent to:

```bash
go run . "hello" standard
```

---

### Standard Banner

```bash
go run . "hello" standard
```

Output:

```text
 _              _   _
| |            | | | |
| |__     ___  | | | |   ___
|  _ \   / _ \ | | | |  / _ \
| | | | |  __/ | | | | | (_) |
|_| |_|  \___| |_| |_|  \___/

```

---

### Shadow Banner

```bash
go run . "Hello There!" shadow
```

---

### Thinkertoy Banner

```bash
go run . "Hello There!" thinkertoy
```

---

### Custom Banner

This project also includes a custom banner:

```bash
go run . "hello" lineart
```

The `lineart` banner is a custom template based on the original format, using a cleaner line-art style.

---

## Handling New Lines

The program treats literal `\n` sequences as line breaks.

Example:

```bash
go run . "Hello\nThere" standard
```

This renders `Hello`, then renders `There` underneath it.

Example with an empty line:

```bash
go run . "Hello\n\nThere" standard
```

This renders `Hello`, prints one empty line, then renders `There`.

---

## Error Handling

If the arguments do not follow the required format, the program prints:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

Invalid example:

```bash
go run . "banana" standard abc
```

This is invalid because the program expects either:

```bash
go run . [STRING]
```

or:

```bash
go run . [STRING] [BANNER]
```

---

## Examples

### Example 1

```bash
go run . "hello" standard | cat -e
```

```text
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

---

### Example 2

```bash
go run . "123" shadow | cat -e
```

---

### Example 3

```bash
go run . "nice 2 meet you" thinkertoy | cat -e
```

---

### Example 4

```bash
go run . "Hello\nThere" standard | cat -e
```

---

## How It Works

The program works in three main steps.

### 1. Parse Arguments

The program checks how many arguments were passed.

Valid forms:

```bash
go run . "hello"
go run . "hello" standard
```

Invalid forms cause the usage message to be printed.

---

### 2. Load the Banner

The selected banner name is used to build a path:

```text
banners/<banner-name>.txt
```

For example:

```bash
go run . "hello" shadow
```

loads:

```text
banners/shadow.txt
```

The banner file is then read into memory and converted into a map where each printable ASCII character points to its 8-line block.

---

### 3. Render the Text

For each line of input, the program renders the text row by row.

Instead of printing one character completely before moving to the next one, it prints:

```text
row 1 of every character
row 2 of every character
row 3 of every character
...
row 8 of every character
```

This keeps the characters aligned horizontally.

---

## Allowed Packages

Only standard Go packages are used:

- `fmt`
- `os`
- `strings`
- `errors`

---

## Testing

Run the program manually with:

```bash
go run . "hello" standard
go run . "hello world" shadow
go run . "nice 2 meet you" thinkertoy
go run . "you & me" standard
go run . "Hello\nThere" standard
```

If test files are included, run:

```bash
go test ./...
```

---

## Custom Templates

The project supports additional custom templates.

To add a new template:

1. Create a new `.txt` banner file.
2. Place it in the `banners/` directory.
3. Make sure it follows the same format as the official banners.
4. Run the program using the file name without `.txt`.

Example:

```text
banners/mybanner.txt
```

Run with:

```bash
go run . "hello" mybanner
```

This project includes the custom template:

```text
lineart.txt
```

---

## Notes

- Unsupported characters outside printable ASCII range are ignored.
- Banner files must contain all printable ASCII characters from `32` to `126`.
- Each character must have exactly 8 lines.
- The banner file structure must not be changed.
- The program expects banner files to be inside the `banners/` directory.

---

## Author

- Panagiotis Valadakis
- Georgia Chalari

---

## License

This project is for educational use.
