<h1>📦 Tetromino Solver: Assemble the Smallest Square! 🧩</h1> <br/>  
Welcome to the Tetromino Solver, a Go program designed to assemble tetrominoes into the smallest possible square while identifying each piece with a unique letter. This project is built with good practices in mind, ensures robust error handling, and includes unit tests for reliability. Let's get started!
<br/> 
<h3>🚀 Objectives</h3><br/> 
Assemble tetrominoes from a text file to form the smallest square possible.
Identify each tetromino with uppercase Latin letters:
🅰 for the first one, 🅱 for the second, and so on.
Handle invalid file formats gracefully by printing ERROR.
Deliver clean, testable, and efficient code written in Go.<br/> 
<h3>🛠️ Instructions</h3><br/> 
<h3>📜 Input</h3><br/> 
Provide the program with a single argument: the path to a text file containing tetrominoes.<br/> 

Example:<br/> 
<br/> 
....<br/> 
..##<br/> 
..##<br/> 
....<br/> 
<br/> 
....<br/> 
..##<br/> 
..#.<br/> 
..#.<br/> 
<h3>🧩 Program Behavior</h3><br/> 
Compiles successfully and expects at least one tetromino.<br/> 
Handles errors in case of:<br/> 
Invalid file structure or formatting.<br/> 
Missing input file.<br/> 
Outputs the smallest possible square with each tetromino identified by a unique letter.<br/> 
<h3>📝 File Format</h3><br/> 
The input file must:<br/> 
<br/> 
Contain one or more tetrominoes.<br/> 
Represent tetrominoes in a 4x4 grid using the following symbols:<br/> 
# for filled blocks.<br/> 
. for empty spaces.<br/> 
Separate tetrominoes with newlines.<br/> 
<h3>✅ Valid Input Example:</h3><br/> 
<br/> 
....<br/> 
.##.<br/> 
.##.<br/> 
....<br/> 
<br/> 
....<br/> 
####<br/> 
....<br/> 
....<br/> 
<h3>❌ Invalid Input Examples:</h3><br/> 
<br/> 
Improper size:<br/> 
<br/> 
...<br/> 
##..<br/> 
#...<br/> 
Missing newline between tetrominoes:<br/> 
<br/> 
....<br/> 
##..<br/> 
#...<br/> 
####<br/> 
....<br/> 
....<br/> 
<h3>🚦 Output</h3><br/> 
The program will print the smallest square with each tetromino represented by letters.<br/> 
For example:<br/> 
<br/> 
AABB<br/> 
AABB<br/> 
CCDD<br/> 
CCDD<br/> 
If an error occurs, the program will print:<br/> 
<br/> 
ERROR<br/> 

<h3>▶️ How to Run</h3><br/> 
Ensure you have Go installed on your machine. 🐹<br/> 
Clone the repository:<br/> 
<br/> 
git clone https://github.com/yourusername/tetris-optimizer.git<br/> 
cd tetris-optimizer<br/> 
Build and run the program:<br/> 
<br/> 
go build -o tetris-optimizer<br/> 
./tetris-optimizer path/to/your/input.txt<br/> 
<h3>🧪 Testing</h3><br/> 
Unit tests are included to ensure reliability. Run them with:<br/> 
<br/>  
go test ./...<br/> 
<h3>🌟 Features</h3><br/> 
Error Handling: Prints ERROR for invalid input.<br/> 
Efficiency: Assembles the smallest possible square.<br/> 
Unit Tests: Confidence in correctness and edge cases.<br/> 
Readable Output: Identifies each tetromino clearly with letters.<br/> 
<h3>💡 Examples</h3><br/> 
Input:<br/> 
<br/> 
....<br/> 
.##.<br/> 
.##.<br/> 
....<br/> 
<br/> 
....<br/> 
####<br/> 
....<br/> 
....<br/> 
<br/> 
AABB<br/> 
AABB<br/> 
CCCC<br/> 
CCCC<br/> 
<br/> 
<h3>📂 Repository Structure</h3><br/> 
<br/> 
tetris-optimizer/<br/> 
├── main.go                 # Entry point of the program<br/> 
├── solver.go               # Logic for solving the puzzle<br/> 
├── utils.go                # Utilities for parsing and validation<br/> 
├── tetris-optimizer.go     # Unit tests<br/> 
├── README.md               # Project documentation<br/> 
<br/> 
<h3>📝 Notes</h3><br/> 
The program is robust to handle most edge cases, but ensure the input file meets the specified format.
Feel free to contribute and enhance the project! 🛠️
Happy Tetris-ing! 🎮🧩✨

