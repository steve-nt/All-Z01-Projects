// scripts/calculations.js

// Function to calculate the average
function calculateAverage(numbers) {
    const total = numbers.reduce((sum, num) => sum + num, 0);
    return total / numbers.length;
}

// Function to calculate the median
function calculateMedian(numbers) {
    numbers.sort((a, b) => a - b);
    const mid = Math.floor(numbers.length / 2);
    return numbers.length % 2 === 0 ? (numbers[mid - 1] + numbers[mid]) / 2 : numbers[mid];
}

// Function to calculate the variance
function calculateVariance(numbers, mean) {
    const total = numbers.reduce((sum, num) => sum + Math.pow(num - mean, 2), 0);
    return total / numbers.length;
}

// Function to calculate the standard deviation
function calculateStandardDeviation(variance) {
    return Math.sqrt(variance);
}
