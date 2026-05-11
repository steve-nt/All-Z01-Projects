// script.js

let chartInstance = null; // Global variable to hold the chart instance

function calculateStatistics() {
    const fileInput = document.getElementById("fileInput");
    if (fileInput.files.length === 0) {
        alert("Please select a data file.");
        return;
    }

    const file = fileInput.files[0];
    const reader = new FileReader();
    reader.onload = function (e) {
        const data = e.target.result.trim().split('\n').map(Number);
        displayStatistics(data);
        visualizeData(data);
    };
    reader.readAsText(file);
}

function displayStatistics(data) {
    const average = Math.round(data.reduce((sum, num) => sum + num, 0) / data.length);

    const sortedData = [...data].sort((a, b) => a - b);
    let median;
    const middle = Math.floor(sortedData.length / 2);
    if (sortedData.length % 2 === 0) {
        median = Math.round((sortedData[middle - 1] + sortedData[middle]) / 2);
    } else {
        median = sortedData[middle];
    }

    const variance = Math.round(data.reduce((sum, num) => sum + Math.pow(num - average, 2), 0) / data.length);
    const stdDev = Math.round(Math.sqrt(variance));

    document.getElementById("average").textContent = average;
    document.getElementById("median").textContent = median;
    document.getElementById("variance").textContent = variance;
    document.getElementById("stdDev").textContent = stdDev;
}

function visualizeData(data) {
    const ctx = document.getElementById("dataChart").getContext("2d");

    if (chartInstance) {
        chartInstance.destroy();
    }

    const histogram = data.reduce((acc, num) => {
        const bucket = Math.floor(num / 10) * 10;
        acc[bucket] = (acc[bucket] || 0) + 1;
        return acc;
    }, {});

    const labels = Object.keys(histogram).map(bucket => `${bucket}-${Number(bucket) + 9}`);
    const values = Object.values(histogram);

    chartInstance = new Chart(ctx, {
        type: "bar",
        data: {
            labels: labels,
            datasets: [{
                label: "Frequency",
                data: values,
                backgroundColor: "rgba(75, 192, 192, 0.6)",
            }],
        },
        options: {
            responsive: true,
            scales: {
                x: {
                    ticks: {
                        color: "#ffffff" // Set x-axis tick color to white
                    },
                    grid: {
                        color: "rgba(255, 255, 255, 0.2)" // Lighten grid lines
                    }
                },
                y: {
                    beginAtZero: true,
                    ticks: {
                        color: "#ffffff" // Set y-axis tick color to white
                    },
                    grid: {
                        color: "rgba(255, 255, 255, 0.2)" // Lighten grid lines
                    }
                }
            },
            plugins: {
                legend: {
                    labels: {
                        color: "#ffffff" // Set legend text color to white
                    }
                }
            }
        }
    });
}

