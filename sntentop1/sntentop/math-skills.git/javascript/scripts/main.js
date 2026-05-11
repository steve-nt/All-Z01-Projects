// scripts/main.js

// Main function that connects everything
function calculateStatistics() {
    readFileAndExtractNumbers((numbers, error) => {
        if (error) {
            displayError(error);
            return;
        }

        const average = calculateAverage(numbers);
        const median = calculateMedian(numbers);
        const variance = calculateVariance(numbers, average);
        const stdDev = calculateStandardDeviation(variance);

        displayResult(average, median, variance, stdDev);
    });
}
