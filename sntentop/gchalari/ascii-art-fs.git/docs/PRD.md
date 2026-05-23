# PRD — ascii-art-fs

## Overview

This project extends the original ascii-art program by supporting multiple banner styles.

The program renders text into ASCII art using one of the following banner templates:

- standard
- shadow
- thinkertoy

---

## Goals

- Support multiple banner styles.
- Keep compatibility with the original ascii-art project.
- Use only standard Go packages.

--- 

## Usage

```
go run . [STRING] [BANNER]
```

Examples:

```
go run . "hello"
go run . "hello" standard
go run . "Hello There!" shadow
go run . "2 you" thinkertoy
```

Any incorrect format must print the following usage message and exit:

```
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```

---

## Arguments

- [STRING]
    - required
    - text to render

- [BANNER]
    - optional
    - defaults to standard
    - available values:
        - standard
        - shadow
        - thinkertoy

---

## Parsing Logic

- One argument:
    - treat it as [STRING]
    - use standard

- Two arguments:
    - first is [STRING]
    - second is [BANNER]

- Any other format:
    - print usage message

- Unsupported banners:
    - print usage message

---

## Supported Characters

The program supports all printable ASCII characters in the range 32 (space) to 126 (`~`). This includes lowercase and uppercase letters, digits, spaces, and special characters such as `#`, `$`, `%`, `&`, `@`, `->`, and others. Characters outside this range are silently skipped.

---

## Newline Handling

The `[STRING]` argument may contain the literal two-character sequence `\n` as a line separator. Each segment separated by `\n` is rendered independently and stacked vertically in the output.

For example:
```
go run . "Hello\nWorld" shadow
```

---

## Output Behavior

The result is printed to stdout using fmt.Print().

---

## Project Structure

---

## Architecture

- ascii.LoadBanner()
    - loads banner files
- ascii.Render()
    - renders ASCII art
- main.go
    - handles:
        - arguments,
        - banner selection,
        - output

---

## Testing

The project must include unit tests covering the following cases:

- valid arguments,
- invalid usage,
- banner rendering,
- newline handling,
- empty strings,
- full ASCII range.