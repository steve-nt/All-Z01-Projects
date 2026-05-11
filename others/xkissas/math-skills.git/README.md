# Calculations Program

This program performs statistical calculations on a '.txt' file with data that represents a statistical population and each line contains one value.

## Table of Contents
- [How to Run the Program](#how-to-run-the-program)
- [Input File Format](#input-file-format)
- [Usage Instructions](#usage-instructions)
- [Output Explanation](#output-explanation)

## How to Run the Program

To run the program, follow these steps:

1. Ensure you have Go installed on your system, as this programm is writen in golang
2. Clone this [repository](https://platform.zone01.gr/git/xkissas/math-skills.git).
3. Open a terminal and navigate to the directory containing the file.
4. Run the command: go run imathyou.go <filename.txt>

Replace `<filename>` with the path to your input '.txt' file.

## Input File Format

The program expects the input file to be a '.txt' file with data and each line contains one value.

## Usage Instructions

1. Prepare your input '.txt' file with any values needed to be calculated.
2. Run the program using the command above. (go run imathyou.go <filename.txt>)
3. The program will perform statistical calculations and display the results.

## Output Explanation

The program will output the following statistics:

- Average: The rounded average of population of values provided.
- Median: The median of population of values provided.
- Variance: The variance of population of values provided.
- Standard Deviation: The standard deviation of population of values provided.

All calculations are executed with float64 type of nambers for more accuracy but the output values are displayed as integers after rounding.

For cross-checking the output values this online calculator has been used: 
https://www.calculatorsoup.com/calculators/statistics/variance-calculator.php 

## Troubleshooting

- Ensure the file path is correct and the file exists.
- Check that the file is a valid '.txt' file.
- Verify that the program has the necessary permissions to read the file.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or issues.



