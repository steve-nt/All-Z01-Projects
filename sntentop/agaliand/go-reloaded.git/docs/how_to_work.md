🧱 docs/how_to_work.md

⚙️ Go Reloaded — How to Work

🎯 Purpose

This document describes the workflow, the roles of developers & AI agents, and the collaboration rules within the Go Reloaded project.
Its goal is to ensure full consistency, reproducibility, and clear communication at every stage of development.

🧠 Roles

Role	Description
Developer	Writes and tests the code, creates new modules, and maintains documentation.
AI Agent	Executes automated tasks, detects consistency issues, suggests improvements, and generates documentation.
Auditor	Reviews others’ commits, evaluates tests, provides feedback, and approves merges.
🧩 Workflow

Task Creation

Every new feature or fix is logged in docs/blueprint-index.md.

Each entry includes a description, assigned developer/agent, date, and status (Pending / Done / Verified).

Development

The developer or AI agent creates a branch:

git checkout -b feature/fixQuotes


Implements the feature in the corresponding pipeline file.

Creates tests inside the tests/ folder.

Testing

Run the tests:

go test ./...


All tests must pass before committing.

If something fails, the auditor is notified in blueprint-index.md.

Commit & Push

Format and lint before committing:

go fmt ./...
go vet ./...


Then:

git add .
git commit -m "Implement replaceHex and add tests"
git push


Review

The auditor reviews the diff, comments, and test results.

Once approved, the task is marked as Verified ✅ in blueprint-index.md.

🧱 Collaboration Rules

Never modify core files without updating documentation.

All changes must go through testing.

Each commit must have a clear, descriptive message.

Every module must include a docstring and a test.

Any change to the pipeline → update architecture.md.

🧪 Testing Workflow

Test Type	Description	Example
Unit Test	Tests individual functions	TestReplaceHex_ValidHexToDecimal
Integration Test	Tests data flow between modules	TestApplyTransformations_FullFlow
Behavioral Test	Verifies if the output matches the expected behavior	TestFixQuotes_ComplexText

📘 Documentation Update Rules

After each modification:

Change	File to Update
New feature	blueprint-index.md
Architecture change	architecture.md
Code rule	coding_standards.md
New term or abbreviation	glossary.md

🔐 Pre-Commit Checklist

 All tests pass (go test ./...)

 No unused imports (go vet ./...)

 Code is properly formatted (go fmt ./...)

 Documentation is updated

 blueprint-index.md has been updated

 Auditor has been informed

🧭 End-to-End Example

# 1. Create a new feature
git checkout -b feature/fixPunctuation

# 2. Develop and test
vim pipeline/fixPunctuation.go
vim tests/fixPunctuation_test.go
go test ./...

# 3. Format & commit
go fmt ./...
git add .
git commit -m "Implement fixPunctuation and tests"

# 4. Push & update
git push
vim docs/blueprint-index.md


🤝 Goal

To achieve perfect collaboration between human and machine.
Each commit is a step toward a cleaner, more coordinated, and more transparent project.