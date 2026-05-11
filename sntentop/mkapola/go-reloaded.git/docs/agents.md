## AI Agent Instructions for Text Correction Tool

This file provides structured instructions for AI agents working with the Go text correction project. Follow these to ensure proper behavior, safe file editing, and consistent output.

1. Project Overview

The project is a Go text processing tool that modifies input files according to tags:

(hex) → convert previous word to decimal

(bin) → convert previous word to decimal

(up), (low), (cap) → change case of previous word

(up, n), (low, n), (cap, n) → multi-word case transformation

Punctuation correction: . , ! ? : ;

Single quote spacing: ' '

Article correction: a → an before vowels or 'h'

The main entry point: main.go

Helpers: processor.go, utils.go

Tests: test/processor_test.go

2. Development Workflow

Use Development Mode

# AI Agent Development Tasks

## Task 1: Basic File I/O and Project Structure
**Functionality**: Set up project structure and implement basic file reading/writing
**TDD Approach**: 
- Write tests for file operations (reading input, writing output)
- Test error handling for missing files
**Implementation Goal**: Create main.go with file I/O functionality
**Validation**: Tests pass for file operations and error cases

## Task 2: Hexadecimal Conversion (hex)
**Functionality**: Convert hexadecimal numbers to decimal when followed by "(hex)"
**TDD Approach**:
- Write tests for hex conversion: "1E (hex)" → "30", "0x23 (hex)" → "35"
- Test edge cases: invalid hex, missing spaces
**Implementation Goal**: Implement hexToDecimal function and integration
**Validation**: All hex conversion tests pass

## Task 3: Binary Conversion (bin)
**Functionality**: Convert binary numbers to decimal when followed by "(bin)"
**TDD Approach**:
- Write tests for binary conversion: "10 (bin)" → "2", "100011 (bin)" → "35"
- Test edge cases: invalid binary, mixed formats
**Implementation Goal**: Implement binToDecimal function and integration
**Validation**: All binary conversion tests pass

## Task 4: Case Transformations - Single Word
**Functionality**: Handle (up), (low), (cap) for single words
**TDD Approach**:
- Write tests for "go (up)" → "GO", "SHOUTING (low)" → "shouting"
- Test "bridge (cap)" → "Bridge"
**Implementation Goal**: Implement case transformation functions
**Validation**: Single word case transformation tests pass

## Task 5: Case Transformations - Multiple Words
**Functionality**: Handle (up, n), (low, n), (cap, n) for multiple words
**TDD Approach**:
- Write tests for "so exciting (up, 2)" → "SO EXCITING"
- Test "sucker punch productions (cap, 3)" → "Sucker Punch Productions"
**Implementation Goal**: Extend case functions to handle word count parameter
**Validation**: Multi-word case transformation tests pass

## Task 6: Punctuation Spacing
**Functionality**: Fix spacing around punctuation marks (. , ! ? : ;)
**TDD Approach**:
- Write tests for "there ,and" → "there, and"
- Test "BAMM !!" → "BAMM!!"
**Implementation Goal**: Implement punctuation spacing correction
**Validation**: Punctuation spacing tests pass

## Task 7: Quote Handling
**Functionality**: Remove spaces inside single quotes
**TDD Approach**:
- Write tests for "' awesome '" → "'awesome'"
- Test multi-word quotes: "' I am happy '" → "'I am happy'"
**Implementation Goal**: Implement quote spacing correction
**Validation**: Quote handling tests pass

## Task 8: Article Correction (a/an)
**Functionality**: Change "a" to "an" before vowels and "h"
**TDD Approach**:
- Write tests for "a amazing" → "an amazing"
- Test "a hour" → "an hour", but keep "a university"
**Implementation Goal**: Implement article correction logic
**Validation**: Article correction tests pass

## Task 9: Integration and Command Processing
**Functionality**: Integrate all transformations in correct order
**TDD Approach**:
- Write comprehensive integration tests with multiple transformations
- Test command line argument handling
**Implementation Goal**: Complete main processing pipeline
**Validation**: Full integration tests pass

## Task 10: Edge Cases and Error Handling
**Functionality**: Handle malformed input and edge cases
**TDD Approach**:
- Write tests for incomplete parentheses, invalid commands
- Test empty files, large files
**Implementation Goal**: Robust error handling and edge case management
**Validation**: All edge case tests pass and program handles errors gracefully

## Development Guidelines
- Follow TDD: write failing tests first, implement minimal code to pass tests, then refactor
- Tasks build incrementally toward complete text processing tool
- Each task must pass all tests before proceeding to next task
- Use Go best practices and maintain clean, readable code

## Task 11: Logging and Debugging

**Functionality**: Add logging to track text transformations step by step
**TDD Approach**:
- Write tests to verify that each transformation logs expected messages

- Test logging for multi-word transformations and punctuation fixes
**Implementation Goal**: Integrate logging (e.g., log.Printf) for transformations without changing output
**Validation**: Logs match expected messages during transformations, tests pass

**AI Tips**: Ask AI to generate logging statements for each function and provide examples for multi-word cases.

## Task 12: Unit Test Coverage Expansion

**Functionality**: Ensure all edge cases are covered for each transformation
**TDD Approach**:

- Write additional tests for:

    - (hex) with lowercase letters (1f (hex) → 31)

    - (bin) with leading zeros (0010 (bin) → 2)

    - Quotes around punctuation ('hello!')

- Add tests for text with multiple transformations in one line
**Implementation Goal**: Increase test coverage for all functions
**Validation**: All edge case tests pass

**AI Tips**: Use AI to **auto-generate multiple edge case tests** from examples.

## Task 13: Performance Checks

**Functionality**: Ensure program handles large text files efficiently
**TDD Approach**:

- Write tests for large input files (>1000 lines)

- Test memory usage and execution time for combined transformations
**Implementation Goal**: Optimize processing pipeline if needed
**Validation**: Program executes correctly and within reasonable time for large files

**AI Tips**: Ask AI to suggest **streaming or buffered reading** in Go to handle large files efficiently.

## Task 14: Command-Line Flags / Options

**Functionality**: Allow optional flags like --verbose or --dry-run
**TDD Approach**:

- Write tests for flag parsing using flag package

- Test verbose mode prints logs without altering output

- Test dry-run mode shows changes without writing output file
**Implementation Goal**: Add CLI options to main.go
**Validation**: Flags work correctly; tests pass

**AI Tips**: Ask AI to generate flag parsing templates and test examples.

## Task 15: Documentation & README Completion

**Functionality**: Update README and inline documentation
**TDD Approach**:

- Write tests that verify README examples match actual program behavior (optional, can be manual)
**Implementation Goal**: Ensure README accurately describes usage, features, and input/output examples
**Validation**: All examples in README produce the expected output when run

**AI Tips**: Use AI to rewrite README sections or generate example outputs from tests.

## Task 16: Refactoring & Code Quality

**Functionality**: Ensure code is clean, readable, and idiomatic Go
**TDD Approach**:

- Run **golint**, **go vet**, and **gofmt**

- Verify refactored functions pass all previous tests
**Implementation Goal**: Maintain all functionality while improving readability and maintainability
**Validation**: All tests pass; code passes linting and formatting checks

**AI Tips**: Ask AI to refactor functions for readability, modularity, or efficiency.

## Task 17: Optional Features / Extensions

- Add new transformations (e.g., custom text replacements, trimming whitespace at line ends)

- Write tests for these new features first

- Implement minimal code to pass tests

**AI Tips**: Let AI suggest additional useful transformations and generate tests automatically.