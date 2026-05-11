
// Returns true if every element is a string with at least 5 characters
const longWords = (arr) => arr.every((str) => typeof str === 'string' && str.length >= 5);

// Returns true if at least one element is a string with 10 or more characters
const oneLongWord = (arr) => arr.some((str) => typeof str === 'string' && str.length >= 10);

// Returns true if there are no strings with at least 7 characters
const noLongWords = (arr) => arr.every((str) => typeof str !== 'string' || str.length < 7);
