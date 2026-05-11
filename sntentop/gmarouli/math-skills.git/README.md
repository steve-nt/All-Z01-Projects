## 🧮 Math Skills Project

### Program Description
This program reads a list of integers from a file, calculates statistical metrics based on those integers, and displays the results. <br>
Specifically, it performs the following tasks:

📁 File Input:

Reads a file provided as a command-line argument.
Extracts integer values from the file, ignoring invalid lines.

📊 Statistics Calculation:

Average (Mean) <br> 
Median <br>
Variance <br> 
Standard Deviation

📤 Output:

Prints the rounded values of the average, median, variance, and standard deviation to the console in a readable format. <br>

❗️ Error Handling:

Handles errors such as missing input files, non-numeric data, or empty files gracefully, providing meaningful error messages. <br>
This program is useful for basic statistical analysis of a dataset stored in a text file, with each number on a separate line.

## The statistical measures that summarize and describe the key characteristics of a set of data.

### 1. Average (Mean)
Definition: The average (or mean) is the sum of all the numbers in a dataset divided by the count of numbers. It gives an idea of the "central" value of the data.

### Steps to Calculate the Average:

The formula for the average (mean) is $ \text{Average (Mean)} = \frac{1}{n} \sum_{i=1}^{n} x_i $.

### Example:

$$
\frac{2 + 4 + 6 + 8}{4} = 5
$$


### 2. Median
Definition: The median is the middle value of a dataset when it is sorted in ascending or descending order. If there is an even number of values, the median is the average of the two middle numbers.

### Steps to Calculate the Median

1. **Sort the Dataset**: Arrange the data in ascending order.

2. **Find the Middle Value**:
   - If the dataset has an **odd** number of values, the median is the middle value.
   - If it has an **even** number of values, the median is the average of the two middle values.


### Example:

For the dataset `[1, 3, 5, 7, 9]`, the median is:

$$
\text{Median} = 5
$$


For the dataset `[1, 3, 5, 7]`, the median is:

$$
\text{Median} = \frac{3 + 5}{2} = 4
$$


### 3. Variance
Definition: Variance measures how much the numbers in a dataset deviate from the mean. It tells you how spread out the numbers are.

### Steps to calculate the Variance

$$
\text{Variance} = \frac{1}{n} \sum_{i=1}^{n} (x_i - \text{mean})^2
$$


### Example:

For the dataset `[1, 3, 5, 7, 9]` with an average of 5:

$$
\text{Variance} = \frac{1}{5} \left[ (1 - 5)^2 + (3 - 5)^2 + (5 - 5)^2 + (7 - 5)^2 + (9 - 5)^2 \right]
$$


### 4. Standard Deviation
Definition: Standard deviation is the square root of the variance. It represents how spread out the numbers are in the original units of the data.

Standard Deviation = √Variance

### Example

If the variance of the dataset is 8, then the standard deviation is:

$$
\text{Standard Deviation} = \sqrt{8} \approx 2.83
$$


### How to Run the program

Download the file: [Here](https://platform.zone01.gr/git/root/public/src/branch/master/subjects/math-skills/audit)

Run the script with ./bin/math-skills or ./run.sh math-skills, then run the program with the created file (data.txt) by the previous command. <br>
💡Tip: Run the script and the program in different workspaces. 

```bash
go run Math_Skills.go data.txt
```

### 🔒 License
This software is provided under an educational license. It is intended solely for use in academic settings, including teaching, research, and personal learning.

### 📨 Contact
For questions or issues, please contact Georgia Marouli.