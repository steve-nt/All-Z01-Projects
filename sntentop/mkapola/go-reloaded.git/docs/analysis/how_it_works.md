# How It Works: Go-Reloaded User Guide

Welcome to **Go-Reloaded** — the text correction tool that makes messy text behave like a well-trained pet.

## 🎯 For Users: Just Want to Fix Text?

### **Quick Setup**
```bash
# Clone the project
git clone <your-repo-url>
cd go-reloaded

# Run it (that's literally it)
go run cmd/main.go input.txt output.txt
```

### **What Does It Actually Do?**

Think of it as a **smart text editor** that follows your instructions embedded in the text itself.

#### **Number Conversions**
```
Input:  I have 1E (hex) apples and 10 (bin) oranges
Output: I have 30 apples and 2 oranges
```
*Because who doesn't want their fruit inventory in decimal?*

#### **Case Transformations**
```
Input:  make this (up) and this (low) and This (cap)
Output: make THIS and this and This

Input:  transform these words (cap, 3) please
Output: Transform These Words please
```

#### **Punctuation Cleanup**
```
Input:  Hello , world ! What's up ?
Output: Hello, world! What's up?
```
*Finally, punctuation that doesn't look like it was sneezed onto the page.*

#### **Quote Handling**
```
Input:  He said ' hello there ' to me
Output: He said 'hello there' to me
```

#### **Grammar Fixes**
```
Input:  It was a amazing day and a hour well spent
Output: It was an amazing day and an hour well spent
```

### **File Requirements**
- **Input**: Any `.txt` file with your messy text
- **Output**: A clean, corrected version
- **Size Limit**: Reasonable files (the program won't crash on War and Peace, but don't test it)

---

## 🛠️ For Contributors: Want to Make It Better?

### **Project Architecture**

This isn't your typical "throw everything in main.go" project. We use a **pipeline architecture** because:

1. **Each step does one thing well** (like a good Unix tool)
2. **Easy to test** (TDD lovers, this is for you)
3. **Easy to debug** (when things go wrong, you know exactly where)

### **Code Structure**
```
cmd/main.go     # Entry point (handles files, calls processor)
pkg/processor.go  # The main pipeline logic
pkg/utils.go      # Helper functions
tests/                      # All the tests (yes, we actually test things)
```

### **How the Pipeline Works**

Text flows through these steps **in order**:

1. **Number Conversion** → `1E (hex)` becomes `30`
2. **Case Transformation** → `hello (up)` becomes `HELLO`
3. **Punctuation Cleanup** → `hello , world` becomes `hello, world`
4. **Quote Processing** → `' text '` becomes `'text'`
5. **Article Correction** → `a amazing` becomes `an amazing`

*Each step gets the output from the previous step. Simple, predictable, debuggable.*

### **Development Workflow**

#### **1. Setting Up**
```bash
# Get the code
git clone <repo-url>
cd go-reloaded

# Run tests (they should all pass)
go test ./tests/

# Run the program
go run cmd/main.go tests/test_input.txt output.txt
```

#### **2. Making Changes**

We follow **TDD** (Test-Driven Development):

1. **Write a failing test** for your new feature
2. **Write minimal code** to make it pass
3. **Refactor** if needed
4. **Repeat**

Example: Adding a new transformation
```bash
# 1. Add test to tests/processor_test.go
# 2. Run tests (should fail)
go test ./tests/

# 3. Implement feature in pkg/processor/
# 4. Run tests again (should pass)
go test ./tests/
```

#### **3. Code Style**

- **Go conventions**: Use `gofmt` and `goimports`
- **DRY principle**: Don't repeat yourself
- **KISS principle**: Keep it simple
- **Clear naming**: `convertHexToDecimal()` not `doStuff()`

### **Common Contribution Areas**

#### **Easy Wins**
- **Add more test cases** (especially edge cases)
- **Improve error messages** (make them helpful, not cryptic)
- **Add input validation** (handle malformed tags gracefully)

#### **Medium Complexity**
- **New transformations** (e.g., `(reverse)`, `(title)`)
- **Performance optimizations** (for large files)
- **Better quote handling** (nested quotes, mixed quote types)

#### **Advanced Features**
- **Configuration files** (custom transformation rules)
- **Streaming processing** (for huge files)
- **Plugin system** (custom transformations)

### **Testing Strategy**

We have **three types of tests**:

1. **Unit Tests** → Test individual functions
2. **Integration Tests** → Test the full pipeline
3. **Golden Tests** → Compare output with expected results

```bash
# Run all tests
go test ./tests/

# Run with coverage
go test -cover ./tests/

# Run benchmarks
go test -bench=. ./tests/
```

### **Debugging Tips**

- **Add logging** to see what each pipeline step produces
- **Test individual functions** before testing the full pipeline
- **Use the test files** in `/tests/` as examples
- **Check the docs** in `/docs/analysis/` for design decisions

---

## 🚀 Advanced Usage

### **Command Line Options**
```bash
# Basic usage
go run cmd/main.go input.txt output.txt

# With verbose output (if implemented)
go run cmd/main.go -v input.txt output.txt
```

### **Batch Processing**
```bash
# Process multiple files
for file in *.txt; do
    go run cmd/main.go "$file" "corrected_$file"
done
```

### **Integration with Other Tools**
```bash
# Use with pipes (if stdin/stdout support is added)
cat messy.txt | go run cmd/main.go - - > clean.txt
```

---

## 🤔 FAQ

**Q: What if my input file has weird characters?**  
A: The program handles UTF-8 text. If you have encoding issues, convert to UTF-8 first.

**Q: Can I add custom transformations?**  
A: Currently, no. But it's a great contribution opportunity! Check the `/docs/agents.md` for implementation tasks.

**Q: What's the performance like?**  
A: O(n) time complexity. It should handle reasonable files quickly. If you're processing novels, grab a coffee.

**Q: Why a pipeline instead of a state machine?**  
A: Check `/docs/analysis/PIPELINEvsFSM.md` for the full analysis. TL;DR: Simpler to test and debug.

---

## 🎉 Contributing

1. **Fork** the repo
2. **Create a branch**: `git checkout -b my-awesome-feature`
3. **Write tests** for your feature
4. **Implement** the feature
5. **Make sure tests pass**: `go test ./tests/`
6. **Submit a PR**

### **Good First Issues**
- Add more test cases to `/docs/GoldenTestSet/`
- Improve error handling for malformed input
- Add support for different quote types (`"` vs `'`)
- Optimize string processing for large files

---

**Remember**: This tool is meant to be **simple, reliable, and easy to understand**. When in doubt, choose the simpler solution.

*Happy text correcting! 🎯*