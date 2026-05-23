<div style="display:flex; justify-content:center;">
<pre><font color="#007D9C">                   _ _                  __                   __          
  ____ ___________(_|_)     ____ ______/ /_      _________  / /___  _____
 / __ `/ ___/ ___/ / /_____/ __ `/ ___/ __/_____/ ___/ __ \/ / __ \/ ___/
/ /_/ (__  ) /__/ / /_____/ /_/ / /  / /_/_____/ /__/ /_/ / / /_/ / /    
\__,_/____/\___/_/_/      \__,_/_/   \__/      \___/\____/_/\____/_/     </font></pre>
</div>

<div align="center">A command-line tool written in Go that renders a string as large ASCII block letters using banner font files, with optional terminal color support.</div>

## Usage

```bash
go run . "STRING"
go run . "STRING" FONT
go run . --color=<color> "STRING"
go run . --color=<color> "SUBSTRING" "STRING"
go run . --color=<color> "SUBSTRING" "STRING" FONT
go run . --color=<color1> --color=<color2> "SUBSTR1" "SUBSTR2" "STRING"
go run . --color=<color1> --color=<color2> "SUBSTR1" "SUBSTR2" "STRING" FONT
```

**FONT** is optional. Valid values: `standard`, `shadow`, `thinkertoy`. Defaults to `standard`.

The `--color` flag is optional and may be repeated. Each flag accepts a color name or value (see [Color Formats](#color-formats) below).

## Color Targeting

When `--color` is used, you can optionally specify a **substring** to color. The substring is the argument that comes directly before the main string.

- **No substring** — the entire output is colored:
  ```bash
  go run . --color=red "hello"
  ```

- **With substring** — only occurrences of that substring are colored (case-sensitive, all occurrences):
  ```bash
  go run . --color=red "ell" "hello"
  ```

- **Multiple colors** — each `--color` flag is paired with its own substring in the same order:
  ```bash
  go run . --color=red --color=blue "hel" "lo" "hello"
  ```
  `hel` is rendered in red, `lo` in blue. When substrings overlap, the first flag takes priority.

- **Multiple colors, no substrings** — all colors are applied to the whole string (first flag wins on overlap):
  ```bash
  go run . --color=red --color=blue "hello"
  ```

## Color Formats

All of the following formats are accepted as the `--color` value:

| Format | Example | Notes |
|---|---|---|
| Named color | `red`, `green`, `blue` | Case-insensitive |
| Hex | `#ff0000` | Must be exactly 6 hex digits |
| RGB | `rgb(255,0,0)` | Values 0–255, spaces around commas are tolerated |
| HSL | `hsl(0,100%,50%)` | H: 0–360, S/L: 0–100% |

**Supported named colors:**
`black`, `red`, `green`, `yellow`, `blue`, `magenta`, `purple`, `cyan`, `white`, `orange`, `pink`,
`brightblack`, `darkgray`, `brightred`, `brightgreen`, `brightyellow`, `brightblue`, `brightmagenta`, `brightcyan`, `brightwhite`

**Flag format is strict — the `=` is required:**
```bash
go run . --color=red "hello"   # correct
go run . --color red "hello"   # error: rejected
```

## Newlines

The literal `\n` inside the input string is treated as a newline, splitting the output into multiple rendered lines:

```bash
go run . "Hello\nWorld"
go run . --color=cyan "Hello\nWorld"
```

## Examples

```bash
# Plain render
go run . "Hello"
go run . "Hello" shadow

# Whole string colored
go run . --color=green "1 + 1 = 2"
go run . --color="#00bfff" "Hello" thinkertoy

# Substring colored
go run . --color=red "ell" "hello"
go run . --color=magenta "World" "Hello World" shadow

# Multiple colors
go run . --color=red --color=blue "hel" "lo" "hello"
go run . --color=rgb(255,165,0) --color=cyan "He" "llo" "Hello" thinkertoy

# Multiline with color
go run . --color=yellow "Hello\nWorld"
```

## Banner Fonts

Three fonts are available, stored in the `banners/` directory:

| Font | Description |
|---|---|
| `standard` | Default. Bold block letters |
| `shadow` | Lighter shadowed style |
| `thinkertoy` | Minimal line-drawn style |

Each font file covers all printable ASCII characters (32–126). Each character is 8 lines tall. Font files must not be modified.

## Project Structure

```
ascii-art-color/
├── main.go          # Entry point, argument parsing, flag handling
├── banner.go        # Banner file loader and parser
├── render.go        # ASCII art renderer (Render, RenderWithColor, BuildColorMask)
├── color.go         # Color string → ANSI escape code resolver
├── main_test.go     # Unit tests
├── go.mod
└── banners/
    ├── standard.txt
    ├── shadow.txt
    └── thinkertoy.txt
```

## Running Tests

```bash
go test ./...
```

## Authors

- Stergios Fourlataras
- Konstantinos Koletsis