// Task 1: Declare components array

let components = ["motor", "sensor", "battery", "camera"];

// Task 2: Extract values from robotParts array

// A variable robotParts (an array) has been provided for you
let firstPart = robotParts[0];
let lastPart = robotParts[robotParts.length - 1];
let comboParts = [lastPart, firstPart];

// Task 3: Replace and swap array elements

// Two arrays have been provided for you: replaceComponents and swapComponents

// Part A - Replace the third element (index 2) of replaceComponents with 'enhanced'
replaceComponents[2] = "enhanced";

// Part B - Swap the first and second elements of swapComponents using temp
let temp = swapComponents[0];
swapComponents[0] = swapComponents[1];
swapComponents[1] = temp;
