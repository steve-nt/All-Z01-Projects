# Go Vocabulary Guide 📚
*A beginner-friendly reference for the go-reloaded project*

Welcome to your Go vocabulary cheat sheet! This guide covers every function, pattern, and symbol used in the go-reloaded project. Think of it as your decoder ring for understanding Go code! 🔍

---

## 📦 Standard Functions

Functions from Go's standard library that we use throughout the project.

| Function | Package | Purpose | Example Usage |
|----------|---------|---------|---------------|
| `fmt.Println()` | fmt | Print text to console with newline | `fmt.Println("Hello!")` → prints "Hello!" |
| `fmt.Printf()` | fmt | Print formatted text (like printf in C) | `fmt.Printf("Error: %v\n", err)` → prints error message |
| `os.Args` | os | Get command-line arguments as slice | `os.Args[1]` → gets first argument after program name |
| `os.Exit()` | os | Exit program with status code | `os.Exit(1)` → exits with error code 1 |
| `os.ReadFile()` | os | Read entire file into memory | `content, err := os.ReadFile("file.txt")` |
| `os.WriteFile()` | os | Write data to file | `os.WriteFile("out.txt", []byte(text), 0644)` |
| `regexp.MustCompile()` | regexp | Create regex pattern (panics if invalid) | `re := regexp.MustCompile(`\d+`)` → matches digits |
| `regexp.Compile()` | regexp | Create regex pattern (returns error if invalid) | `re, err := regexp.Compile(`\d+`)` |
| `strconv.ParseInt()` | strconv | Convert string to integer in any base | `val, err := strconv.ParseInt("FF", 16, 64)` → hex to int |
| `strconv.FormatInt()` | strconv | Convert integer to string in any base | `strconv.FormatInt(255, 10)` → "255" |
| `strconv.Atoi()` | strconv | Convert string to int (base 10 only) | `num, err := strconv.Atoi("42")` → 42 |
| `strings.Fields()` | strings | Split string into words (by whitespace) | `strings.Fields("hello world")` → ["hello", "world"] |
| `strings.Join()` | strings | Join slice of strings with separator | `strings.Join(["a", "b"], " ")` → "a b" |
| `strings.ToUpper()` | strings | Convert string to uppercase | `strings.ToUpper("hello")` → "HELLO" |
| `strings.ToLower()` | strings | Convert string to lowercase | `strings.ToLower("HELLO")` → "hello" |
| `strings.TrimSpace()` | strings | Remove leading/trailing whitespace | `strings.TrimSpace(" hello ")` → "hello" |
| `len()` | builtin | Get length of slice, string, etc. | `len("hello")` → 5 |

> 💡 **Pro Tip**: `fmt.Printf` uses format verbs like `%v` (any value), `%s` (string), `%d` (decimal). It's like Mad Libs for programmers!

---

## 🔍 Regex Patterns

All the regular expressions used in the project, decoded for human understanding 🤖.

### Number Conversion Patterns

| Pattern | Explanation | Example Match |
|---------|-------------|---------------|
| `([0-9A-Fa-f]+)\s*\(hex\)` | **Hex Pattern**: `([0-9A-Fa-f]+)` = one or more hex digits, `\s*` = zero or more spaces, `\(hex\)` = literal "(hex)" | `"FF (hex)"`, `"1E(hex)"` |
| `([01]+)\s*\(bin\)` | **Binary Pattern**: `([01]+)` = one or more 0s or 1s, `\s*` = optional spaces, `\(bin\)` = literal "(bin)" | `"101 (bin)"`, `"1111(bin)"` |

### Case Transform Patterns

| Pattern | Explanation | Example Match |
|---------|-------------|---------------|
| `(\S+(?:\s+\S+)*)\s*\(up(?:,\s*(\d+))?\)` | **Uppercase Pattern**: `(\S+(?:\s+\S+)*)` = words, `\(up(?:,\s*(\d+))?\)` = "(up)" or "(up, number)" | `"hello (up)"`, `"go lang (up, 2)"` |
| `(\S+(?:\s+\S+)*)\s*\(low(?:,\s*(\d+))?\)` | **Lowercase Pattern**: Same structure as up, but for lowercase | `"HELLO (low)"`, `"GO LANG (low, 1)"` |
| `(\S+(?:\s+\S+)*)\s*\(cap(?:,\s*(\d+))?\)` | **Capitalize Pattern**: Same structure for capitalization | `"hello (cap)"`, `"go lang (cap, 2)"` |

### Punctuation & Quote Patterns

| Pattern | Explanation | Example Match |
|---------|-------------|---------------|
| `\s+([,.!?;:])` | **Before Punctuation**: `\s+` = one or more spaces, `([,.!?;:])` = any punctuation mark | `" ,"`, `"  !"` |
| `([.!?])([a-zA-Z])` | **After Punctuation**: `([.!?])` = sentence-ending punctuation, `([a-zA-Z])` = letter | `".A"`, `"!hello"` |
| `'\s*([^']+?)\s*'` | **Quote Content**: `'` = literal quote, `\s*` = optional spaces, `([^']+?)` = anything except quotes (non-greedy), `\s*'` = optional spaces + closing quote | `"' hello '"`, `"'  test  '"` |

### Article Correction Patterns

| Pattern | Explanation | Example Match |
|---------|-------------|---------------|
| `\ba\s+([aeiouAEIOU])` | **Vowel Articles**: `\b` = word boundary, `a\s+` = "a" + spaces, `([aeiouAEIOU])` = vowel | `"a apple"`, `"a elephant"` |
| `\ba\s+[uU][nrs]` | **U-Consonant Check**: Matches "a" before u-words that sound like consonants | `"a university"`, `"a uniform"` |

> 🎯 **Regex Decoder Ring**:
> - `+` = one or more
> - `*` = zero or more  
> - `?` = zero or one (optional)
> - `\s` = any whitespace
> - `\b` = word boundary
> - `[abc]` = any character in brackets
> - `[^abc]` = any character NOT in brackets
> - `()` = capture group (remember this part)
> - `(?:...)` = non-capturing group (don't remember)

---

## ⚡ Operators & Symbols

Go operators and symbols used throughout the project.

| Symbol | Name | Purpose | Example |
|--------|------|---------|---------|
| `:=` | Short declaration | Declare and assign variable | `name := "Go"` |
| `=` | Assignment | Assign to existing variable | `name = "Golang"` |
| `==` | Equality | Compare for equality | `if name == "Go"` |
| `!=` | Not equal | Compare for inequality | `if err != nil` |
| `&&` | Logical AND | Both conditions must be true | `if x > 0 && x < 10` |
| `<` | Less than | Numeric comparison | `if count < len(words)` |
| `<=` | Less than or equal | Numeric comparison | `if count <= len(words)` |
| `[]` | Slice/Array | Create or access slice/array | `words[0]`, `[]string{}` |
| `...` | Variadic | Variable number of arguments | `fmt.Printf(format, args...)` |
| `&` | Address of | Get memory address | `&variable` |
| `*` | Pointer/Dereference | Declare pointer or get value | `*testing.T`, `*pointer` |
| `%v` | Format verb | Print any value | `fmt.Printf("Value: %v", x)` |
| `%s` | Format verb | Print string | `fmt.Printf("Name: %s", name)` |
| `%d` | Format verb | Print decimal integer | `fmt.Printf("Count: %d", 42)` |
| `%q` | Format verb | Print quoted string | `fmt.Printf("Input: %q", text)` |

> 🤓 **Memory Trick**: `:=` is like saying "Hey Go, figure out what type this should be!" while `=` is like "I already told you the type, just change the value."

---

## 🛠️ Custom Functions

Functions we wrote specifically for this project.

### 📁 main.go

| Function | Purpose | Input → Output | Notes |
|----------|---------|----------------|-------|
| `main()` | Program entry point, handles command-line args and file I/O | Command line args → File processing | The boss function that coordinates everything 👑 |

**Example Flow:**
```
Input: go run main.go input.txt output.txt
Process: Read input.txt → Transform text → Write output.txt
Output: Success message or error
```

### 📁 pkg/processor/processor.go

| Function | Purpose | Input → Output | Notes |
|----------|---------|----------------|-------|
| `ProcessText()` | Main pipeline coordinator | `"raw text"` → `"processed text"` | The conductor of our text transformation orchestra 🎼 |
| `convertNumbers()` | Convert hex/binary to decimal | `"FF (hex) and 101 (bin)"` → `"255 and 5"` | Turns programmer numbers into human numbers |
| `applyCaseTransforms()` | Handle up/low/cap commands | `"hello (up)"` → `"HELLO"` | The case-changing chameleon 🦎 |
| `normalizePunctuation()` | Fix spacing around punctuation | `"hello , world"` → `"hello, world"` | The spacing police officer 👮 |
| `processQuotes()` | Clean up quote spacing | `"' hello '"` → `"'hello'"` | Removes awkward spaces in quotes |
| `fixArticles()` | Change "a" to "an" before vowels | `"a apple"` → `"an apple"` | Grammar nazi but in a good way 📚 |

**Pipeline Flow:**
```
Raw Text → convertNumbers() → applyCaseTransforms() → 
normalizePunctuation() → processQuotes() → fixArticles() → Clean Text
```

### 📁 pkg/processor/utils.go

| Function | Purpose | Input → Output | Notes |
|----------|---------|----------------|-------|
| `ReadFile()` | Read entire file into string | `"filename.txt"` → `("file content", nil)` or `("", error)` | Your friendly file-reading robot 🤖 |
| `WriteFile()` | Write string to file | `("filename.txt", "content")` → `nil` or `error` | The file-writing wizard ✨ |

**File Permission Note:** `0644` means owner can read/write, everyone else can only read. It's like saying "This is my diary, you can read it but don't write in it!" 📖

### 📁 tests/processor_test.go

| Function | Purpose | Input → Output | Notes |
|----------|---------|----------------|-------|
| `TestConvertNumbers()` | Test hex/binary conversion | Test cases → Pass/Fail results | Makes sure our number magic actually works ✨ |
| `TestCaseTransforms()` | Test case transformations | Test cases → Pass/Fail results | Ensures CAPS and lowercase behave properly |
| `TestPunctuationNormalization()` | Test punctuation spacing | Test cases → Pass/Fail results | Checks if punctuation plays nice with spaces |
| `TestQuoteProcessing()` | Test quote spacing fixes | Test cases → Pass/Fail results | Verifies quotes don't have social distancing issues |
| `TestArticleCorrection()` | Test a/an grammar rules | Test cases → Pass/Fail results | Grammar checker's best friend 📝 |
| `TestCompleteTransformation()` | Test full pipeline | README example → Expected output | The final boss test - everything must work! 🎯 |
| `TestEdgeCases()` | Test weird inputs | Edge cases → Graceful handling | Tests what happens when users get creative 🎨 |

> 🧪 **Testing Philosophy**: "Trust, but verify" - We trust our code works, but we verify it with tests because computers are sneaky! 

---

## 🎯 Key Concepts Explained

### Error Handling Pattern
```go
result, err := someFunction()
if err != nil {
    // Handle the error
    return err
}
// Use result safely
```
This is Go's way of saying "Things might go wrong, let's be prepared!" 🛡️

### Slice Operations
```go
words := strings.Fields("hello world")  // Split into slice
word := words[0]                        // Get first element
count := len(words)                     // Get length
```

### String Manipulation Chain
```go
text = step1(text)
text = step2(text)
text = step3(text)
```
Each function takes text, transforms it, and passes it to the next step - like an assembly line for text! 🏭

---

## 🚀 Quick Reference

**Most Common Patterns:**
- **Check for errors**: `if err != nil { return err }`
- **Loop through slice**: `for _, item := range slice { ... }`
- **String formatting**: `fmt.Printf("Value: %v", variable)`
- **Regex replace**: `re.ReplaceAllString(text, replacement)`

**Remember:** Go is like a friendly but strict teacher - it wants you to handle errors explicitly and be clear about what you're doing. No surprises allowed! 🎓

---

*Happy coding! May your regex always match and your tests always pass! 🎉*