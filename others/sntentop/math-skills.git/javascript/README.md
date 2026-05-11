
# 📊 Simple Statistics Calculator

Welcome to the **Simple Statistics Calculator**! This project allows you to upload a text file containing numbers and calculates the following statistical measures:

- **Average** 🧮
- **Median** 📈
- **Variance** 📊
- **Standard Deviation** 📐

All results are rounded to the nearest integer for simplicity.

---

## 🌟 Features

- **User-Friendly Interface**: A clean and minimalist design for easy navigation.
- **Real-Time Calculation**: Instant results upon uploading and processing the file.
- **Modular Code Structure**: Organized into separate files for better readability and maintenance.

---

## 🚀 Getting Started

### Prerequisites

- A modern web browser (Chrome, Firefox, Safari, Edge)

### Project Structure

```
statistics-calculator/
├── index.html          # Main HTML file
├── styles/
│   └── style.css       # CSS styling
└── scripts/
    ├── main.js         # Main script orchestrating the app
    ├── fileReader.js   # Handles file reading and number extraction
    ├── calculations.js # Performs statistical calculations
    └── displayResult.js# Displays results or errors
```

---

## 🔧 How to Use

1. **Clone or Download the Repository**:

   ```bash
   git clone https://github.com/your-username/statistics-calculator.git
   ```

2. **Open `index.html`**:

   - Navigate to the project folder.
   - Open `index.html` in your preferred web browser.

3. **Prepare Your Data File**:

   - Create a `.txt` file containing numbers.
   - **Format**: One number per line.
   - **Example**:

     ```
     @@
     @@
     @@
     @@
     @@
     ```

4. **Upload the File**:

   - Click on the **"Choose File"** button.
   - Select your `.txt` file.

5. **Calculate Statistics**:

   - Click on the **"Calculate"** button.
   - View the results displayed on the page.

---

## 📝 Code Explanation

### 1. `index.html`

- **Purpose**: Provides the structure of the web page.
- **Contains**:
  - File input for uploading `.txt` files.
  - A button to trigger calculations.
  - A div to display results.

### 2. `style.css`

- **Purpose**: Styles the web page for a professional and clean look.
- **Features**:
  - Responsive design centered on the page.
  - Styled buttons and input fields.
  - Result display area with clear formatting.

### 3. `main.js`

- **Purpose**: Orchestrates the application flow.
- **Function**: `calculateStatistics()`
  - Invoked when the user clicks the "Calculate" button.
  - Calls `readFileAndExtractNumbers` to read the file.
  - Handles errors or proceeds to calculations.
  - Calls calculation functions and displays results.

### 4. `fileReader.js`

- **Purpose**: Reads the uploaded file and extracts numbers.
- **Function**: `readFileAndExtractNumbers(callback)`
  - Reads the file as text.
  - Splits content into lines.
  - Parses each line to a number.
  - Filters out invalid entries.
  - Returns an array of numbers or an error message.

### 5. `calculations.js`

- **Purpose**: Performs statistical calculations.
- **Functions**:
  - `calculateAverage(numbers)`
  - `calculateMedian(numbers)`
  - `calculateVariance(numbers, mean)`
  - `calculateStandardDeviation(variance)`
- **Logic**:
  - Uses standard mathematical formulas.
  - Each function handles one specific calculation.

### 6. `displayResult.js`

- **Purpose**: Displays the results or errors to the user.
- **Functions**:
  - `displayResult(average, median, variance, stdDev)`
    - Shows the calculated statistics.
  - `displayError(message)`
    - Displays error messages.

---

## 💡 Understanding the Code

- **Asynchronous Operations**: File reading is asynchronous. The use of callbacks ensures that calculations happen only after the file is fully read.
- **Modular Design**: Separating code into different files makes it easier to manage and understand.
- **Error Handling**: The app gracefully handles cases where no file is selected or the file contains invalid data.
- **Functional Programming**: Use of array methods like `map`, `filter`, and `reduce` for data processing.

---

## 📚 Learning Resources

To understand and write similar code, you can explore the following resources:

- **JavaScript Basics**:
  - [JavaScript Guide - MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide)
  - [Eloquent JavaScript](https://eloquentjavascript.net/)

- **Asynchronous JavaScript**:
  - [Understanding Asynchronous JavaScript](https://blog.bitsrc.io/understanding-asynchronous-javascript-the-event-loop-74cd408419ff)
  - [JavaScript Promises](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Using_promises)

- **File Handling in JavaScript**:
  - [FileReader API - MDN](https://developer.mozilla.org/en-US/docs/Web/API/FileReader)
  - [Using Files from Web Applications](https://developer.mozilla.org/en-US/docs/Web/API/File_API/Using_files_from_web_applications)

- **Array Methods**:
  - [Array.prototype.map()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/map)
  - [Array.prototype.filter()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/filter)
  - [Array.prototype.reduce()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/reduce)

- **DOM Manipulation**:
  - [Document Object Model (DOM) - MDN](https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model)
  - [Introduction to the DOM](https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model/Introduction)

- **CSS Styling**:
  - [CSS Basics - MDN](https://developer.mozilla.org/en-US/docs/Learn/Getting_started_with_the_web/CSS_basics)
  - [CSS Flexbox Guide](https://css-tricks.com/snippets/css/a-guide-to-flexbox/)

- **Statistical Concepts**:
  - [Mean, Median, and Mode](https://www.mathsisfun.com/mean.html)
  - [Variance and Standard Deviation](https://www.mathsisfun.com/data/standard-deviation.html)

---

## 🛠️ Tools and Tips

- **Code Editors**:
  - [Visual Studio Code](https://code.visualstudio.com/)
  - [Sublime Text](https://www.sublimetext.com/)
- **Browser Developer Tools**:
  - Use the console to debug and test JavaScript code.
- **Practice**:
  - Write small programs that use file input and array methods.
  - Try modifying the code to add new features or handle diffesrent data formats.

---

Happy Coding! 🎉
