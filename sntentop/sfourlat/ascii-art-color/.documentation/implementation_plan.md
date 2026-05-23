# ASCII-Art Color — Implementation Plan

## Overview

Extend the existing ascii-art program with a `--color=<color>` flag that colorizes
all occurrences of a given substring inside the ASCII art output. If no substring is
provided, the entire output is colored. All other existing functionality (banner
selection, `\n` splitting, single-string mode) must continue to work unchanged.

---

## New files

| File | Purpose |
|---|---|
| `color.go` | Color string → ANSI escape code resolution |

## Modified files

| File | What changes |
|---|---|
| `main.go` | Argument parsing extended to handle `--color=` flag |
| `render.go` | New `RenderWithColor` function + `buildColorMask` helper |
| `main_test.go` | New test cases for color parsing, mask logic, and rendering |

---

## 1. Argument parsing (`main.go`)

### Valid call signatures

```
go run . [STRING]
go run . [STRING] [BANNER]
go run . --color=<color> [STRING]
go run . --color=<color> [SUBSTRING] [STRING]
go run . --color=<color> [SUBSTRING] [STRING] [BANNER]
```

### Parsing logic

Loop over `os.Args[1:]` and separate flags from positional arguments:

- If arg starts with `--color=` → extract value after `=` as `colorName`
- If arg starts with `--` but is not `--color=...` → print usage, exit 1
- Otherwise → append to positional slice

Then resolve positionals:

| `--color` present | Positional count | Meaning |
|---|---|---|
| No | 1 | `string` |
| No | 2 | `string`, `banner` |
| Yes | 1 | `string` (whole string colored) |
| Yes | 2 | `substring`, `string` |
| Yes | 3 | `substring`, `string`, `banner` |
| Any | anything else | print usage, exit 1 |

### Usage message (exact format required by spec)

```
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

---

## 2. Color resolution (`color.go`)

One exported function: `colorCode(name string) (string, error)`

Returns an ANSI foreground escape sequence, or an error if the format is not
recognized.

### Named colors (ANSI 3/4-bit)

| Input | Code |
|---|---|
| `black` | `\033[30m` |
| `red` | `\033[31m` |
| `green` | `\033[32m` |
| `yellow` | `\033[33m` |
| `blue` | `\033[34m` |
| `magenta` / `purple` | `\033[35m` |
| `cyan` | `\033[36m` |
| `white` | `\033[37m` |
| `orange` | `\033[93m` (bright yellow) |
| `pink` | `\033[95m` (bright magenta) |
| `brightred`, `brightgreen`, etc. | `\033[9Xm` variants |

Matching is case-insensitive (`strings.ToLower` before lookup).

### Hex notation — `--color=#rrggbb`

- Strip leading `#`, must be exactly 6 hex digits
- Parse into R, G, B with `strconv.ParseUint(..., 16, 8)` in pairs
- Emit `\033[38;2;R;G;Bm`

### RGB notation — `--color=rgb(r,g,b)`

- Strip `rgb(` prefix and `)` suffix
- Split on `,`, trim spaces, parse each as integer 0–255 with `strconv.Atoi`
- Emit `\033[38;2;R;G;Bm`

### HSL notation — `--color=hsl(h,s%,l%)`

- Strip `hsl(` prefix and `)` suffix
- Split on `,`, trim spaces and `%`
- Parse H (0–360 float), S (0–100 float), L (0–100 float) with `strconv.ParseFloat`
- Convert HSL → RGB (pure Go, no external packages — implement the standard
  algorithm using the hue-to-rgb helper)
- Emit `\033[38;2;R;G;Bm`

### Reset constant

```go
const ansiReset = "\033[0m"
```

### Error handling

If `colorCode` cannot match any format, return `("", error)`. `main.go` prints the
error and usage then exits 1.

---

## 3. Color mask (`render.go` — `buildColorMask`)

```
buildColorMask(runes []rune, substr string) []bool
```

- Returns a `[]bool` of the same length as `runes`
- `true` at position `i` means the character at `i` should be colored
- If `substr == ""` → all positions `true`
- Otherwise: scan `runes` with a sliding window of `len([]rune(substr))`
  - Case-sensitive comparison, rune by rune
  - On match, mark all positions in that window `true`
  - Overlapping matches are all marked (scan does not skip after a match)

---

## 4. Colored renderer (`render.go` — `RenderWithColor`)

```
RenderWithColor(lines []string, bannerMap map[rune][]string, ansiColor string, substr string)
```

For each line segment:

1. If empty → `fmt.Println()` and continue
2. Build `colored` mask via `buildColorMask`
3. For each row 0–7:
   - Iterate over runes with their index
   - Track `inColor bool`
   - At each character: if `colored[i] != inColor`, emit the appropriate
     transition (`ansiColor` to start, `ansiReset` to stop)
   - Write `bannerMap[char][row]` to the builder
   - After all characters: if `inColor` is still true, emit `ansiReset`
   - `fmt.Println(builder.String())`

This keeps the ANSI codes tightly wrapped around only the colored character
blocks, with no leaking codes between rows.

---

## 5. Tests (`main_test.go`)

### `colorCode` tests

- Named colors resolve correctly (spot-check red, green, blue, orange, pink)
- Case-insensitive: `RED`, `Red`, `rEd` all resolve
- Aliases: `purple` → magenta code, `orange` → bright yellow code
- Unknown name → returns error
- `#ff0000` → `\033[38;2;255;0;0m`
- `#FF0000` → same (case-insensitive hex)
- `rgb(255,0,0)` → `\033[38;2;255;0;0m`
- `rgb(255, 0, 0)` → same (spaces tolerated)
- `hsl(0,100%,50%)` → `\033[38;2;255;0;0m`
- Malformed inputs → error

### `buildColorMask` tests

- Empty substr → all true
- Substr not present → all false
- Single char substr, one occurrence
- Single char substr, multiple occurrences (`l` in `hello` → positions 2 and 3)
- Multi-char substr, multiple occurrences (`kit` in `a king kitten have kit`)
- Case-sensitive: `B` in `RGB()` → only position 1 colored, not `b` if present
- Substr longer than string → all false

### `RenderWithColor` tests

- Empty segment → single blank line (no color codes)
- ANSI color code appears in output when coloring is active
- ANSI reset code appears after colored section when partial coloring
- No ANSI codes in output when no characters match the substr
- Output has same number of rows (8) as plain `Render` for same input

---

## 6. Edge cases to handle explicitly

| Case | Expected behavior |
|---|---|
| `--color red "banana"` (space not `=`) | Usage message, exit 1 |
| `--color=` (empty color value) | Error: unknown color, usage, exit 1 |
| Substring not found in string | Output renders normally, no color applied |
| Substring equals full string | Entire output colored |
| String with `\n` splits and color applied per segment | Each segment colored independently |
| Banner arg combined with color flag | Both work together |
| `B` in `RGB()` | Only the `B` block colored, `(`, `R`, `G`, `)` uncolored |

---

## 7. Constraints reminder

- Standard library only — no `regexp`, no external packages
- Substring matching: manual rune-by-rune sliding window
- HSL→RGB conversion: implement from scratch using standard arithmetic
- No modification to banner files
- `Render` (no color) remains unchanged for backward compatibility