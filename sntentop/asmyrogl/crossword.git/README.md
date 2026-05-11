# Crossword Solver

## Description
This project implements a JavaScript-based algorithm designed to solve custom crossword puzzles. It takes a formatted puzzle string and an array of words, validates the inputs, maps the available slots, and utilizes a constraint satisfaction and backtracking approach to populate the grid. If the puzzle has exactly one valid solution, it outputs the solved grid, otherwise, it returns an error.

## Core Features
* **Strict Input Validation:** Ensures the puzzle string contains only permitted characters (`.`, `0`, `1`, `2`, `\n`), verifies the grid is perfectly rectangular, and checks that the word list contains unique, valid strings.
* **Automated Grid Analysis:** Scans the grid to identify horizontal and vertical word slots, validating that the declared number of starting words in a cell matches the actual grid geometry.
* **Optimized Backtracking Solver:** Employs a depth-first search to explore word placements. It optimizes the process by filtering slots by word length, detecting character conflicts early, and preserving grid states for efficient backtracking.
* **Uniqueness Verification:** The algorithm terminates early if multiple valid solutions are detected, adhering to the rule that a valid crossword must have exactly one unique solution.

## How It Works
The execution flow is divided into three primary phases:

1.  **Validation Phase:** Checks the integrity of both the puzzle string and the words array. It catches formatting issues, duplicate words, or mismatches between the number of words and the number of available slots.
2.  **Analysis Phase:** Parses the string into a 2D character grid. It identifies all valid horizontal and vertical word slots, recording their coordinates and lengths.
3.  **Solving Phase:** Uses recursive backtracking to attempt placing words into the identified slots. If a word fits without conflicting with existing letters, it is placed, and the algorithm proceeds to the next word. If a dead-end is reached, it un-places the word and tries a different configuration.

## Input Format
The puzzle is represented as a single string with the following characters:
* `.` represents a black square (unplayable space).
* `0` represents an open cell where no new words begin.
* `1` represents an open cell where exactly one word begins (either horizontally or vertically).
* `2` represents an open cell where exactly two words begin (both horizontally and vertically).
* `\n` represents a line break, dividing the string into uniform rows.

## Running the Code via Node.js

You can easily run this solver locally using Node.js.

### Prerequisites
* Ensure you have [Node.js](https://nodejs.org/) installed on your machine. You can verify this by opening your terminal or command prompt & run: 
    ```bash
    node -v
    ```
### Execution Steps

#### Option 1: Run the code with the example provided
1.  **Ensure Example Code is Uncommented:** At the end of `crosswordSolver.js`, make sure lines similar to these are **uncommented**:
    ```javascript
    const emptyPuzzle = `2001\n0..0\n1000\n0..0`;
    const words = ["casa", "alan", "ciao", "anta"];
    crosswordSolver(emptyPuzzle, words);
    ```
4.  **Run the Script:** Open your terminal, navigate to the folder containing your file, and execute:
    ```bash
    node crosswordSolver.js
    ```

#### Option 2: Run the tests (test.js)
1.  **Ensure Example Code is commented out::** At the end of the file, **comment out** the example code (lines 247-254):
    ```javascript
    // const emptyPuzzle = `2001\n0..0\n1000\n0..0`;
    // const words = ["casa", "alan", "ciao", "anta"];
    // crosswordSolver(emptyPuzzle, words);
    ```
    And ensure the export statement are **present** and **uncommented**:
    ```javascript
    module.exports = crosswordSolver;
    ```

2.  **Run the tests:** Navigate to your project directory and execute:
    ```bash
    node test.js
    ```
