🧠 Go Reloaded

📖 Overview

Go Reloaded is an automatic text processing and formatting tool written entirely in Go.
It reads an input file, detects special tags and punctuation, and produces a clean, corrected output file.
Its architecture is based on a modular pipeline, where each processing stage is independent and can be tested or extended individually.
The project is designed to collaborate with AI Agents (Codex, Copilot, Claude, ChatGPT) that follow the protocol defined in AGENTS.md.

📂 Project Structure

📁 go-reloaded/
├── main.go
│   # Program entry point
│
├── pipeline/
│   ├── readInput.go              # Reads the input file
│   ├── tokenize.go               # Splits text into tokens
│   ├── replaceHex.go             # Converts hex numbers to decimal
│   ├── replaceBin.go             # Converts binary numbers to decimal
│   ├── applyCaseTransform.go     # Applies (up), (low), (cap) transformations
│   ├── formatPunctuation.go      # Fixes spacing around punctuation marks
│   ├── fixQuotes.go              # Corrects single quotes placement
│   ├── fixArticles.go            # Replaces “a” with “an” where appropriate
│   ├── applyTransformations.go   # Combines all transformation steps
│   └── writeOutput.go            # Writes the final output file
│
├── tests/
│   ├── readInput_test.go
│   ├── tokenize_test.go
│   ├── replaceHex_test.go
│   ├── replaceBin_test.go
│   ├── applyCaseTransform_test.go
│   ├── formatPunctuation_test.go
│   ├── fixQuotes_test.go
│   ├── fixArticles_test.go
│   ├── applyTransformations_test.go
│   └── writeOutput_test.go
│   # All test files follow a TDD (Test Driven Development) pipeline
│
├── docs/
│   ├── architecture.md           # Describes the internal architecture and data flow
│   ├── coding_standards.md       # Code style and naming conventions for Go
│   ├── how_to_work.md            # Workflow guide for developers and AI agents
│   ├── blueprint-index.md        # Central task and progress tracker
│   └── glossary.md               # Glossary of project terminology
│
├── tasks/
│   ├── TASK-A1.md                # Example: replaceHex() implementation
│   ├── TASK-A2.md                # Example: replaceBin() implementation
│   └── ...                       # Additional task files for AI Agents
│
├── .github/
│   ├── workflows/
│   │   └── ci.yml                # GitHub Actions for automated testing and QA
│   └── .actrc                    # Config for local CI testing with act
│
├── AGENTS.md                     # Execution protocol for AI Agents
└── README.md                     # Main project documentation


⚙️ Functional Pipeline

readInput
↓
tokenize
↓
applyTransformations
↓
formatPunctuation
↓
fixQuotes
↓
fixArticles
↓
writeOutput


🧩 Core Features

CommandDescriptionExample(hex)Converts a hexadecimal number to decimal"1E (hex)" → "30"(bin)Converts a binary number to decimal"10 (bin)" → "2"(up)Converts the previous word to uppercase"go (up)" → "GO"(low)Converts the previous word to lowercase"STOP (low)" → "stop"(cap)Capitalizes the previous word"bridge (cap)" → "Bridge"(up, n)Converts the n previous words to uppercase"so exciting (up,2)" → "SO EXCITING"(low, n)Converts the n previous words to lowercase"WOW THAT'S COOL (low,3)" → "wow that's cool"(cap, n)Capitalizes the n previous words"brooklyn bridge park (cap,3)" → "Brooklyn Bridge Park"PunctuationFixes spacing around punctuation marks"Hello , world !" → "Hello, world!"'quotes'Places single quotes correctly" ' awesome ' " → "'awesome'"a → anReplaces “a” with “an” before a vowel or “h” sound"a apple" → "an apple"

🧱 Function Breakdown

FunctionDescriptionreadInput()Reads the input filetokenize()Splits text into words, punctuation, and tagsreplaceHex()Converts hexadecimal numbers to decimalreplaceBin()Converts binary numbers to decimalapplyCaseTransform()Applies upper/lower/capital transformationsformatPunctuation()Fixes spacing around punctuation marksfixQuotes()Adjusts the placement of quotesfixArticles()Validates and corrects “a/an” usageapplyTransformations()Combines all transformation stepswriteOutput()Writes the final output file

🧪 Testing & Quality Assurance

All functions follow the Test Driven Development (TDD) philosophy.
Tests are located in the /tests folder and can be executed with:
go test ./tests/...

For local CI execution:
act -j build -W .github/workflows/ci.yml


🤖 AI Integration (AGENTS.md)

Go Reloaded is designed to work seamlessly with AI Agents.
The tasks/ folder contains markdown task files where each AI Agent performs the following steps:


Analyze & Confirm


Generate the Tests


Generate the Code


QA & Mark Complete


Usage instructions and operational rules for Agents are detailed in AGENTS.md.
See AGENTS.md for full execution guidelines.
Additionally, the workflow that every Agent follows is defined in docs/how_to_work.md:
analyze → ask for operator confirmation → implement tests → implement code → QA


🧭 Developer Docs

If you’re a new developer or agent, start here:
FilePurposedocs/architecture.mdExplains the pipeline and data flowdocs/how_to_work.mdGuides developers & agents through the work stepsdocs/coding_standards.mdClean and consistent Go code standardsdocs/blueprint-index.mdOverview of all tasks and progressAGENTS.mdExecution protocol for AI Agents

🚀 Execution

Run the program with:
go run main.go input.txt output.txt

The result will be saved in the file:
output.txt
