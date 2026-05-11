✍️ Go Reloaded — Coding Standards
🎯 Purpose
This document defines the coding, formatting, and style rules for the Go Reloaded project.
It applies to both developers and AI Agents in the repository, ensuring a unified architecture, consistent naming, and code quality.

📁 General Principles


Each folder (e.g., pipeline/, tests/, docs/) must have a clear purpose and not mix functionalities.


Each function must perform only one responsibility.


Code must be readable — write for the next person who will maintain it, not for yourself.


Comments should explain why something happens, not what it does.


The Agent follows the same standards as a senior developer.



🧩 Naming Conventions
TypeFormatExampleFunctionscamelCasecleanText, replaceHex, fixQuotesVariablescamelCaseinputText, tokens, outputFileConstantsUPPER_SNAKE_CASEMAX_WORD_LENGTH, DEFAULT_FILE_PATHPackageslowercase without underscorespipeline, utils, testsTest files_test.go suffixreplaceHex_test.go, fixQuotes_test.go

🗒️ Comments & Documentation
Each .go file must start with a docstring describing its purpose:
// Package pipeline contains the core text-processing functions for Go Reloaded.

Each function must include documentation:
// replaceHex replaces all hexadecimal values (e.g. 0xFF) with their decimal equivalent.
func replaceHex(text string) string { ... }

❌ Avoid:
// loop through string
for ...
✅ Prefer:
// Scan runes to detect potential punctuation spacing issues

🧠 AI Agent Behavior
Agents operating inside the repository must strictly follow these steps:


Analyze the task and check if a related function already exists.


Ask for confirmation before modifying core files.


Implement tests first (*_test.go) and then the main function.


Document the change in blueprint-index.md.


QA & Refactor after every commit.



🧪 Testing Guidelines
Every function inside pipeline/ must have a corresponding test inside tests/.
Example naming convention:
func TestReplaceHex_ValidHexToDecimal(t *testing.T) { ... }



Always use the standard testing library.


Tests should describe behavior, not implementation.


No merges are allowed unless go test ./... passes successfully.



🧰 Code Style Rules


Indentation: Tabs (not spaces).


Maximum line length: 100 characters.


Import order:


import (
    "fmt"
    "strings"
    "go-reloaded/pipeline"
)

→ Standard libraries first, third-party next, internal packages last.
Before every commit:
go fmt ./...
go vet ./...
go test ./...



Every file must end with a newline.



🧭 Error Handling
All functions that read or write files must return errors instead of panicking:
func readInput(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("cannot read input file %s: %w", path, err)
    }
    return string(data), nil
}


💅 Style Rules per Pipeline Stage
FileDescriptionExample FunctionreadInput.goReads file contentreadInput(path string) (string, error)tokenize.goSplits text into wordstokenize(text string) []stringreplaceHex.goConverts hex → decimalreplaceHex(text string) stringreplaceBin.goConverts binary → decimalreplaceBin(text string) stringapplyCaseTransform.goHandles capitalization/lowercaseapplyCaseTransform(text string) stringformatPunctuation.goAdjusts spacing around punctuationfixPunctuation(text string) stringfixQuotes.goFixes quotation marks into proper pairsfixQuotes(text string) stringfixArticles.goCorrects “a/an” articlesfixArticles(text string) stringapplyTransformations.goIntegrates all transformationsapplyTransformations(text string) stringwriteOutput.goWrites the result to the output filewriteOutput(text, path string) error

🧩 Example Commit Flow
git add .
git commit -m "Implement fixQuotes and corresponding tests"
go fmt ./...
go test ./...
git push


🔒 Quality Checklist
Before committing:
☑ All code is formatted (go fmt ./...)
☑ All tests pass (go test ./...)
☑ No unused imports
☑ Comments are updated
☑ The agent updated blueprint-index.md
☑ QA verified consistency with how_to_work.md
