# PRD — go-reloaded

Keep this short (1-2 pages). The project subject + audit cases are the source of truth.

---

## 1. Problem Statement

What does this tool do? One sentence.

---

## 2. CLI Contract

- Command: `go run main.go <inputFile> <outputFile>`
- What happens if args are missing?
- What happens if the input file doesn't exist?
- What happens if the output file can't be created?

---

## 3. Transformation Functions

Document each function in your own words with at least one example.

### 3.1 Number Conversions

- `(hex)`: convert the previous word from hexadecimal to decimal
  - Example: `1E (hex)` -> `30`
- `(bin)`: convert the previous word from binary to decimal
  - Example: `10 (bin)` -> `2`

### 3.2 Case Transformations

- `(up)`, `(low)`, `(cap)`: apply to the previous word
  - Example: `go (up)` -> `GO`
- `(up, N)`, `(low, N)`, `(cap, N)`: apply to the previous N words
  - Example: `so exciting (up, 2)` -> `SO EXCITING`

### 3.3 Punctuation

- `. , ! ? : ;` attach to the previous word, space after (unless end of text)
- Punctuation groups like `...` or `!?` stay grouped

Include at least 2 examples (one with groups).

### 3.4 Quotes

- `'` always appears in pairs, wrapping word(s)
- Remove spaces immediately inside the quotes

Include 1 example with one word and 1 with multiple words.

### 3.5 Articles

- `a` becomes `an` before a vowel (`a, e, i, o, u`) or `h`

Include 2 examples (one vowel, one `h`).

---

## 4. Non-Goals

What you are NOT building. Be explicit to avoid overbuilding.

---

## 5. Acceptance Criteria

### Audit Cases

- [ ] Audit case 1: ...
- [ ] Audit case 2: ...
- [ ] Audit case 3: ...
- [ ] Audit case 4: ...

### Extra Golden Tests (minimum 2-3)

- [ ] Extra test 1: ...
- [ ] Extra test 2: ...
- [ ] Extra test 3: ...

---

## 6. Architecture

Pick one: Pipeline or FSM

- We choose: ________
- Because: ________
- Tradeoffs we accept: ________

Sketch (ASCII is fine):

---

## 7. Milestones

3-5 milestones, each testable:

1. ________
2. ________
3. ________

---

## 8. Risks / Open Questions

- ________
- ________
