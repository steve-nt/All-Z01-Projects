// Sum all numbers in array
const adder = (arr, initial = 0) =>
  arr.reduce((sum, num) => sum + num, initial);

// Add odd numbers, multiply even numbers
const sumOrMul = (arr, initial = 0) =>
  arr.reduce((acc, num) => (num % 2 === 0 ? acc * num : acc + num), initial);

// Execute array of functions sequentially
const funcExec = (arr, initial = 0) =>
  arr.reduce((acc, func) => func(acc), initial);
