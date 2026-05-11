# Go-Reloaded: Text Correction Tool

A **professional-grade** Go program that intelligently processes and corrects text files using a **pipeline architecture**. Built with **TDD principles** and designed for reliability, performance, and maintainability. The code is ai generated and some parts of the text from the rest of the files.
Created for educational purposes, so expect plenty of comments inside the code!

## 🎯 Core Features

### **Number Base Conversion**
- **Hexadecimal to Decimal**: `1E (hex)` → `30`
- **Binary to Decimal**: `10 (bin)` → `2`
- Handles edge cases: `0 (hex)`, `FF (hex)`, `1111 (bin)`

### **Case Transformations**
- **Single Word**: `hello (up)` → `HELLO`
- **Multi-Word**: `go reloaded (cap, 2)` → `Go Reloaded`
- **Supported Tags**: `(up)`, `(low)`, `(cap)` with optional word count

### **Punctuation Normalization**
- **Spacing Correction**: `hello , world` → `hello, world`
- **Grouped Punctuation**: Preserves `...`, `!?`, `!!!`
- **Supported Marks**: `. , ! ? : ;`

### **Quote Processing**
- **Single Quotes**: `' hello world '` → `'hello world'`
- **Nested Handling**: Complex quote scenarios with proper spacing

### **Grammar Correction**
- **Article Adjustment**: `a amazing` → `an amazing`
- **Vowel Detection**: `a hour` → `an hour`
- **Smart Logic**: Preserves `a university` (consonant sound)

## 🚀 Quick Start

### **Installation**
```bash
git clone <repository-url>
cd go-reloaded
```

### **Usage**
```bash
go run cmd/main.go <input_file> <output_file>
```

### **Example**
```bash
go run cmd/main.go tests/test_input.txt output.txt
```

## 📁 Project Structure

```
go-reloaded/
├── cmd/go-reloaded/          # Main application entry point
│   └── main.go
├── pkg/processor/            # Core processing logic
│   ├── processor.go          # Main text processing pipeline
│   └── utils.go             # Utility functions
├── tests/                   # Test files and test data
│   ├── processor_test.go    # Unit tests
│   ├── test_input.txt       # Sample input
│   └── expected_output.txt  # Expected results
├── docs/                    # Documentation
│   ├── analysis/            # Technical analysis
│   └── GoldenTestSet/       # Comprehensive test cases
└── go.mod                   # Go module definition
```

## 🧪 Testing

### **Run Tests**
```bash
go test ./tests/
```

### **Test Coverage**
```bash
go test -cover ./tests/
```

### **Benchmark Tests**
```bash
go test -bench=. ./tests/
```

## 🔧 Architecture

### **Pipeline Design**
The program uses a **sequential pipeline** approach where each transformation step:
1. **Receives** input text
2. **Processes** specific patterns
3. **Passes** result to next step

RAW TEXT INPUT
      |
      v
+--------------------+
| 1️⃣ convertNumbers  |
| (hex/bin → decimal)|
+--------------------+
      |
      v
+--------------------------+
| 2️⃣ applyCaseTransforms   |
| (UP / LOW / CAP)         |
+--------------------------+
      |
      v
+--------------------------+
| 3️⃣ normalizePunctuation  |
| (Fix spacing around      |
|  punctuation marks)      |
+--------------------------+
      |
      v
+--------------------------+
| 4️⃣ processQuotes          |
| (Remove extra spaces      |
| inside quotes)            |
+--------------------------+
      |
      v
+--------------------------+
| 5️⃣ fixArticles            |
| (Change "a" → "an"       |
| before vowel sounds)      |
+--------------------------+
      |
      v
TRIM SPACES & FINAL CLEAN TEXT


### **Processing Order**
1. **Number Conversions** (hex/bin → decimal)
2. **Case Transformations** (up/low/cap)
3. **Punctuation Normalization**
4. **Quote Processing**
5. **Article Correction** (a → an)

### **Key Benefits**
- **Isolated Functions**: Each step is independently testable
- **TDD Compatible**: Perfect for test-driven development
- **Maintainable**: Easy to add new transformations
- **Debuggable**: Clear separation of concerns

## 📊 Performance

- **Time Complexity**: O(n) where n is input length
- **Memory Usage**: Minimal overhead with string processing
- **File Size Limit**: Handles large files efficiently
- **Execution Time**: < 5 minutes for any reasonable input

## 🎯 Example Transformations

### **Input**
```
I have 1E (hex) apples and 10 (bin) oranges . it (cap) is a amazing day ! ' hello world '
```

### **Output**
```
I have 30 apples and 2 oranges. It is an amazing day! 'hello world'
```

### **Complex Example**
```
Input:  the (cap, 3) quick brown fox said ' FF (hex) is 255 ' !
Output: The Quick Brown fox said 'FF is 255'!
```

## 🛠️ Development

### **Code Quality Standards**
- **Go Best Practices**: Following effective Go guidelines
- **DRY Principle**: No code duplication
- **KISS Principle**: Keep implementations simple
- **SOC**: Clear separation of concerns

### **Testing Strategy**
- **Unit Tests**: Individual function testing
- **Integration Tests**: Full pipeline testing
- **Edge Cases**: Boundary condition handling
- **TDD Approach**: Tests written before implementation

### **Error Handling**
- **File Operations**: Graceful handling of missing/invalid files
- **Input Validation**: Robust parsing of malformed input
- **Recovery**: Continues processing despite individual transformation errors

## 📚 Documentation

- **Technical Analysis**: `/docs/analysis/`
- **Test Cases**: `/docs/GoldenTestSet/`
- **Development Guide**: `/docs/agents.md`
- **Best Practices**: `/docs/good_practices.md`

## 🤝 Contributing

1. **Fork** the repository
2. **Create** feature branch: `git checkout -b feature-name`
3. **Write** tests for new functionality
4. **Implement** changes following TDD
5. **Ensure** all tests pass: `go test ./...`
6. **Submit** pull request

### **Development Workflow**
- Follow **TDD**: Write failing tests first
- Maintain **test coverage** > 90%
- Use **gofmt** and **goimports** for formatting
- Follow **Go conventions** for naming and structure

## 📄 License

This project is open source and available under the **MIT License**.

---

**Built with ❤️ using Go and TDD principles**