# Crossword Solver

## 📖 Description

The **Crossword Solver** project solves empty crossword puzzles based on
a given list of words.\
It takes as input a puzzle (string) and a list of words, and fills the
crossword grid with a **unique solution**.\
If no solution exists or multiple solutions are possible, it prints
`Error`.

-   Uses **backtracking** to place words.
-   Validates input for:
    -   Invalid puzzle format.
    -   Duplicate words.
    -   Mismatch between number of slots and words.
    -   More than one or zero solutions.

------------------------------------------------------------------------

## 🚀 Usage

### 1. Run directly

To run the solver with a demo puzzle:

``` bash
node crosswordSolver.js
```

### 2. Run unit tests

The `test.js` file covers both success and error scenarios:

``` bash
node test.js
```

You will see ✓ or ✗ depending on the result.

### 3. Run examples with expected output

The `examples.js` file runs all examples and compares them with saved
"expected" outputs.

-   Generate/update expected snapshots:

    ``` bash
    node examples.js --update
    ```

-   Verify outputs against expected:

    ``` bash
    node examples.js
    ```

------------------------------------------------------------------------

## 📂 Project Structure

    crossword/
    │
    ├── crosswordSolver.js   # Main crosswordSolver function
    ├── test.js              # Unit tests for all cases
    ├── examples.js          # Examples + snapshot testing
    └── expected/            # Folder with expected output files

------------------------------------------------------------------------

## ✅ Example

**Input:**

``` js
const puzzle = '2001\n0..0\n1000\n0..0'
const words = ['casa', 'alan', 'ciao', 'anta']
crosswordSolver(puzzle, words)
```

**Output:**

    casa
    i..l
    anta
    o..n

------------------------------------------------------------------------

## 👩‍💻 Technologies

-   Node.js (no external packages required)
-   JavaScript (ES6)

------------------------------------------------------------------------

## 📝 Notes

-   If the puzzle or word list is invalid, the program prints `Error`.
-   All tests and examples run without needing a `package.json`.