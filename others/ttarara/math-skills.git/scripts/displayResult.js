// scripts/displayResult.js

// Function to display results in the HTML element with id "result"
function displayResult(average, median, variance, stdDev) {
    const resultDiv = document.getElementById('result');
    resultDiv.innerText = `Average: ${Math.round(average)}\n` +
                          `Median: ${Math.round(median)}\n` +
                          `Variance: ${Math.round(variance)}\n` +
                          `Standard Deviation: ${Math.round(stdDev)}`;
}

// Function to display an error message
function displayError(message) {
    const resultDiv = document.getElementById('result');
    resultDiv.innerText = message;
}
