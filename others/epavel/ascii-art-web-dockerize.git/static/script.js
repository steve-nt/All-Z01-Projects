
// JavaScript to adjust the textarea width and height
const asciiOutput = document.getElementById('ascii-art-output');
const outputContainer = document.getElementById('output-container');
var lines = String(asciiInput).split('\\n');
let maxLengthLine = '';
let totalLength = 0;
for (let i = 0; i < lines.length; i++) {
  if (lines[i].length > maxLengthLine.length) {
    maxLengthLine = lines[i];
  }
  console.log(lines[i].length);
  totalLength += lines[i].length;
}
console.log(maxLengthLine.length);
console.log(totalLength);
var percentageToTotal = (maxLengthLine.length/totalLength);
const newlineCount = lines.length - 1; // Count the number of newlines
//7 us the width of the character ' '
console.log(newlineCount);
const outputLength = (asciiOutput.value.length * percentageToTotal)
const viewportWidth = window.innerWidth;
var widthInVw = (outputLength / viewportWidth) * 100 + 4; //4vw is right+left padding in vw of body

if (widthInVw < 100) {
    asciiOutput.style.width = '95vw';
    asciiOutput.style.textAlign = 'center';
} else {
    asciiOutput.style.width = widthInVw + 'vw';
    asciiOutput.style.textAlign = 'left';
    outputContainer.style.border = '1px solid #595959';
    outputContainer.style.borderRadius = '10px';
    outputContainer.style.margin = '0 20px';
}

asciiOutput.style.height = 'calc(300px + ' + (newlineCount - 1) * 150 + 'px)'; //300px is the height of the first line, +150px for each subsequent line
