
# Data Statistics Calculator

This program is a web-based tool for calculating key statistical measures — including the average, median, variance, and standard deviation — from a dataset. Users can upload a text file containing numerical data, and the program will display the calculated statistics alongside a visual histogram, allowing for quick insights into data distribution.

### What Are Average, Median, Variance, and Standard Deviation?

-   **Average**: The average is like finding the "middle" amount of something. Imagine you have a pile of candy, and you want to share it equally with your friends. You’d count all the candies, then divide them by how many friends you have (including yourself!). This way, each friend would get about the same amount, and that’s what we call the "average."
    
-   **Median**: The median is the number right in the middle when you line up all your numbers from smallest to largest. If you have 5 pieces of candy and sort them from least to most, the median is the candy right in the middle. If you have an even number, you take the middle two and find the average of those. So, it’s the "middle value."
    
-   **Variance**: Variance is about understanding how spread out your numbers are. If everyone has almost the same amount of candy, the variance is low. But if some people have a lot and others have a little, the variance is high. It’s a way to measure how "different" each number is from the average.
    
-   **Standard Deviation**: Standard deviation is kind of like variance’s best friend. It also shows how spread out the numbers are, but in a way that’s closer to the original values. You can think of it as a “distance” from the average — if most numbers are close to the average, the standard deviation is small. If they’re far away, it’s bigger.

-------------

### Formulas

1. **Average (Mean)**:
   - **Formula**:  
     Average = (Σ x) / N
   - **Explanation**: Add up all the numbers (denoted as x) and divide by the total number of values (N) to find the average. This tells us the typical value in our data set.

2. **Median**:
   - **Formula**:
     - If N is odd: The median is the middle value in the sorted list.
     - If N is even: The median is the average of the two middle values.
   - **Explanation**: First, sort the numbers in order. If there's an odd number of values, the median is the middle one. If there's an even number, take the middle two and find their average. This gives us the "center" of the data.

3. **Variance**:
   - **Formula**:  
     Variance = (Σ (x - Average)²) / N
   - **Explanation**: For each value, calculate the difference from the average, then square that difference (so it's always positive). Sum all these squared differences, and then divide by N. This gives us a measure of how spread out the numbers are from the average.

4. **Standard Deviation**:
   - **Formula**:  
     Standard Deviation = √Variance
   - **Explanation**: Take the square root of the variance to find the standard deviation. Standard deviation also shows how spread out the numbers are, but in the original units. A smaller standard deviation means the numbers are closer to the average, while a larger one means they are more spread out.

---
### HTML Structure Overview

This HTML file creates a clean interface for the Data Statistics Calculator.

In the `<head>`, we set up essential metadata, link our CSS file for styling, and add the Chart.js library for drawing charts. Chart.js lets us easily create dynamic data visualizations.

In the `<body>`, we wrap everything in a `container` to center and style our content.

1.  **Title**: A simple `<h1>` with the title "Data Statistics Calculator" gives users a clear idea of what the app does.
    
2.  **File Input**: We use a file input (`<input type="file">`) to let users upload a `.txt` file containing the data they want to analyze. This file should have one number per line.
    
3.  **Calculate Button**: A button triggers the JavaScript function `calculateStatistics()`, which processes the uploaded data.
    
4.  **Results Display**: We have a section to display the calculated statistics (average, median, variance, standard deviation). Each result has its own `<p>` with a `<span>` to show the values dynamically updated by JavaScript.
    
5.  **Chart Canvas**: The `<canvas>` element, which works with Chart.js, provides a place to display a histogram of the data. It’s styled with a fixed width and height.
    

Finally, we link our JavaScript file (`script.js`) at the bottom of the body to handle calculations and chart rendering after the HTML content is loaded.

----
### CSS Structure Overview

In this stylesheet, we create a stylish, centered, and cohesive look for our Data Statistics Calculator app.

-   **Body Styling**: We start by giving the entire background a dynamic gradient, going from deep blue to warm red to golden yellow, creating a modern, inviting vibe. The text color is set to white for readability, and we center everything using flexbox, so it’s perfectly aligned on any screen size.
    
-   **Container**: The `.container` class wraps the main content with a soft white overlay, using slight transparency so the background gradient subtly shows through. Padding, rounded corners, and a shadow make it visually distinct and a bit "elevated" from the background.
    
-   **Title**: The `<h1>` element is styled to be clean and prominent, using a white color that stands out without being too flashy. A bit of spacing below the title keeps it separated from the rest of the content.
    
-   **Button Styling**: The button has a vibrant red color that draws attention. It has no border, rounded corners, and a slight color transition on hover for a smooth interactive feel. When hovered, the button darkens to provide visual feedback.
    
-   **Results Section**: Each `<p>` in the results section is styled with a slightly larger font to make the statistics easy to read. Simple margins give each line of stats a bit of breathing room.
    
-   **Chart Canvas**: Finally, the `<canvas>` for the chart has a margin at the top and rounded corners to match the overall style, ensuring the chart integrates smoothly into the app’s design.
- ----
### JavaScript Overview

This script handles everything: reading the data, calculating statistics, displaying results, and visualizing data in a chart.

1.  **Global Chart Instance**: We start with `chartInstance = null`, a global variable to keep track of the chart so we can destroy and update it later without duplicating.
    
2.  **File Input and Processing** (`calculateStatistics` function):
    
    -   When the "Calculate Statistics" button is clicked, this function checks if a file is selected.
    -   If a file is chosen, it reads the file using `FileReader` and splits its content by line, mapping each line to a number.
    -   Once the data is ready, it calls `displayStatistics` to calculate and display values and `visualizeData` to draw the chart.
3.  **Statistics Calculation** (`displayStatistics` function):
    
    -   **Average**: Sums all values and divides by the count.
    -   **Median**: Sorts the data; finds the middle element (or averages the two middle elements if even).
    -   **Variance**: Measures the spread by summing the squared differences from the mean, then divides by the number of values.
    -   **Standard Deviation**: The square root of the variance, giving us an idea of how spread out the values are from the mean.
    -   After calculating, each result is inserted into the corresponding `<span>` in the HTML, updating the displayed values.
4.  **Data Visualization** (`visualizeData` function):
    
    -   First, we destroy any previous chart instance to avoid duplicates.
    -   We then create a histogram using "buckets" of 10 (e.g., 0–9, 10–19) to group values and calculate the frequency of each bucket.
    -   Chart.js is used to render a bar chart with these frequencies. We customize the color of the x- and y-axis labels, grid lines, and legend to match the app’s theme.

> Written with [StackEdit](https://stackedit.io/).
