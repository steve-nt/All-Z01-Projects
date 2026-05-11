# README

## Description
This program is written in Go and calculates statistical metrics such as **Average**, **Median**, **Variance**, and **Standard Deviation** from a list of integers provided in a text file. The program reads the input file line by line, processes the numbers, and prints the results rounded to the nearest integer.

## Requirements
- Go (1.16 or later)

## Input Format
The input file should be a text file where each line contains a single integer. For example:
```
10
20
30
40
50
```

## How to Test the Program

### 1. Clone the Repository (if applicable)
If this program is part of a repository, clone it using:
```bash
git clone <repository-url>
cd <repository-folder>
```

### 2. Prepare the Input File
Create a text file (e.g., `data.txt`) in the same directory as the program. The file should contain integers, one per line.

Example (`data.txt`):
```
10
15
20
25
30
```

### 3. Build and Run the Program
Use the following command to run the program, providing the input file name as an argument:
```bash
go run your-program.go data.txt
```

### 4. Expected Output
The program will read the numbers from the file, calculate the statistics, and print the results in the following format:
```
Average: <value>
Median: <value>
Variance: <value>
Standard Deviation: <value>
```

For example, for the `data.txt` above, the output might be:
```
Average: 20
Median: 20
Variance: 50
Standard Deviation: 7
```

### 5. Handling Errors
The program includes error handling for the following cases:
- If no file name is provided, or more than 1 file is provided it will print:
  ```
  Only 1 argument is needed
  ```
- If the file cannot be opened or read, it will display an appropriate error message.
- If the file contains invalid data (e.g., non-integer values), it will print an error and terminate.
- If the file is empty, it will print:
  ```
  No numbers found in the file.
  ```

## Notes
- The file name is dynamic, so you can test with different files by changing the argument.
- The program automatically rounds results to the nearest integer.

## Example Test
### Input File (`example.txt`):
```
5
10
15
20
25
```
### Run Command:
```bash
go run your-program.go example.txt
```
### Output:
```
Average: 15
Median: 15
Variance: 50
Standard Deviation: 7
```

