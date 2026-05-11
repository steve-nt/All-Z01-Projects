# Tetris Optimizer

A Go program that assembles a set of tetrominoes into the smallest possible square.

## 🧩 Project Description

This program takes a text file containing one or more tetromino blocks (Tetris pieces) and attempts to arrange them in the smallest square board without overlaps.

Each tetromino is identified with capital letters (`A`, `B`, `C`...) in the solution.

If the file format is incorrect or any tetromino is invalid, the program prints:
ERROR


## 📂 File Structure

```text
.
├── main.go          # Entry point of the program
├── parser.go        # File parsing and tetromino validation
├── solver.go        # Backtracking algorithm to place tetrominoes
├── board.go         # Board creation and printing
├── tetromino.go     # Tetromino data structure
├── good_inputs/     # Valid input test files
└── bad_inputs/      # Invalid input test files



---

## 🚀 How to Run

```bash
go run . path/to/input.txt

✅ Example (valid input)
go run . good_inputs/good01_seven_blocks.txt

❌ Example (invalid input)
go run . bad_inputs/bad01_not_connected.txt

✅ Visual Example
Input: good_inputs/good01_seven_blocks.txt

...#
...#
...#
...#

....

....

....

####

.###
...#
....

....

..##
.##.
....

Output:
.CCCA
..CA.
..DDA
.DD.A
BBBB.
FFFFG
.FFH.G

📜 Tetromino Format Rules
Each tetromino is exactly 4 lines of 4 characters

Only characters allowed: # and .

Each tetromino must contain exactly 4 #

Tetromino blocks must be separated by an empty line

All # characters must be connected horizontally or vertically (not diagonally)

🧪 Testing
Test input files are organized in:

good_inputs/ — expected to succeed and print a compact square

bad_inputs/ — expected to print ERROR

go run . good_inputs/filename.txt
go run . bad_inputs/filename.txt

🔁 Test Automation Script
You can run all test files automatically using the run_tests.sh script:

chmod +x run_tests.sh
./run_tests.sh

This script will:

✅ Check that bad inputs return ERROR

✅ Display the results of good inputs as printed boards

✍️ Author
Nancy Zemperligou
Zone01 Athens — tetris-optimizer project
2025