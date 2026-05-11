// scripts/fileReader.js

// Function to read the file and extract numbers
function readFileAndExtractNumbers(callback) {
    const fileInput = document.getElementById('fileInput');
    if (!fileInput.files.length) {
        callback(null, 'Please select a file.');
        return;
    }

    const file = fileInput.files[0];
    const reader = new FileReader();

    reader.onload = function() {
        const lines = reader.result.split('\n');
        const numbers = lines.map(line => parseFloat(line.trim())).filter(num => !isNaN(num));

        if (numbers.length === 0) {
            callback(null, 'No valid numbers found in the file.');
        } else {
            callback(numbers, null);
        }
    };

    reader.readAsText(file);
}
