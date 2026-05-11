#!/usr/bin/env node

// Test file for crosswordSolver
// Basic test suite with 10 core tests

const crosswordSolver = require('./crosswordSolver.js');

// Capture console output
let output = '';
const originalLog = console.log;

function captureLog(fn) {
  output = '';
  console.log = function(...args) {
    output = args.join('');
  };
  fn();
  console.log = originalLog;
  return output;
}

// Test helper
function test(name, puzzle, words, expected) {
  const result = captureLog(() => crosswordSolver(puzzle, words));
  const pass = result === expected;
  console.log(`${name} ${pass ? '[Passed]' : '[Failed]'}`);
  if (!pass) {
    console.log(`  Expected: ${JSON.stringify(expected)}`);
    console.log(`  Got:      ${JSON.stringify(result)}`);
  }
  return pass;
}

console.log(`\nTESTS:`);

let passed = 0;

// Test 1: Basic valid puzzle
passed += test(
  'Test 1: Basic 4x4 puzzle',
  '2001\n0..0\n1000\n0..0',
  ['casa', 'alan', 'ciao', 'anta'],
  'casa\ni..l\nanta\no..n'
) ? 1 : 0;

// Test 2: Single vertical word
passed += test(
  'Test 2: Single vertical word',
  '1\n0',
  ['ab'],
  'a\nb'
) ? 1 : 0;

// Test 3: Horizontal word (length 2)
passed += test(
  'Test 3: Horizontal word (length 2)',
  '10',
  ['ab'],
  'ab'
) ? 1 : 0;

// Test 4: Empty puzzle (invalid)
passed += test(
  'Test 4: Empty puzzle should error',
  '',
  ['casa'],
  'Error'
) ? 1 : 0;

// Test 5: Non-string puzzle (invalid)
passed += test(
  'Test 5: Non-string puzzle should error',
  123,
  ['casa'],
  'Error'
) ? 1 : 0;

// Test 6: Non-array words (invalid)
passed += test(
  'Test 6: Non-array words should error',
  '1\n0',
  'not-array',
  'Error'
) ? 1 : 0;

// Test 7: Duplicate words (invalid)
passed += test(
  'Test 7: Duplicate words should error',
  '2001\n0..0\n1000\n0..0',
  ['casa', 'casa', 'ciao', 'anta'],
  'Error'
) ? 1 : 0;

// Test 8: Word count mismatch (invalid)
passed += test(
  'Test 8: Word count mismatch should error',
  '2001\n0..0\n1000\n0..0',
  ['casa', 'alan'],
  'Error'
) ? 1 : 0;

// Test 9: Unsolvable puzzle (no solution)
passed += test(
  'Test 9: Incompatible words should error',
  '2001\n0..0\n1000\n0..0',
  ['aaaa', 'bbbb', 'cccc', 'dddd'],
  'Error'
) ? 1 : 0;

// Test 10: Multiple solutions (ambiguous)
passed += test(
  'Test 10: Ambiguous puzzle should error',
  '2000\n0...\n0...\n0...',
  ['abba', 'assa'],
  'Error'
) ? 1 : 0;

console.log(`\nRESULTS:`);
console.log(`Passed: ${passed}/10`);
console.log(`Failed: ${10 - passed}/10`);

if (passed === 10) {
  console.log('\n All tests passed!');
  process.exit(0);
} else {
  console.log(`\n ${10 - passed} test(s) failed.`);
  process.exit(1);
}
