// The main function that takes the puzzle and words, validates them, and attempts to solve the crossword puzzle using backtracking.
function crosswordSolver(puzzle, words) {
  const parsedPuzzle = parsePuzzle(puzzle);
  if (parsedPuzzle === "Error") return console.log("Error");

  const cleanWords = parseWords(words);
  if (cleanWords === "Error") return console.log("Error");

  const { grid, wordSlots } = parsedPuzzle;

  if (cleanWords.length !== wordSlots.length) return console.log("Error");

  const solutions = [];
  const charGrid = grid.map((row) =>
    row.split("").map((char) => (char === "." ? "." : "")),
  );

  //Initialize slot status
  wordSlots.forEach((s) => (s.filled = false));

  //Start solving
  solve(0, cleanWords, wordSlots, charGrid, solutions);

  //Final Output Validation
  if (solutions.length !== 1) {
    console.log("Error");
  } else {
    console.log(solutions[0]);
  }
}

// parsePuzzle checks that the puzzle is a non-empty string, that it only contains legal characters, and that the rows are of equal length. It then analyzes the grid to find word slots and validates the number of starts.
function parsePuzzle(puzzle) {
  if (typeof puzzle !== "string" || puzzle === "") {
    return "Error";
  }
  if (detectIllegalChars(puzzle)) {
    return "Error";
  }

  const grid = splitString(puzzle);
  if (grid === "Error") {
    return "Error";
  }

  const wordSlots = analyzeGrid(grid);
  if (wordSlots === "Error") {
    return "Error";
  }
  return { grid, wordSlots };
}

// parseWords checks that words is an array of non-empty strings and that there are no duplicates
function parseWords(words) {
  // Check that words is actually an array
  if (!Array.isArray(words)) return "Error";
  // Check all words are strings and none are empty
  for (const word of words) {
    if (typeof word !== "string" || word === "") return "Error";
  }
  // This works since sets only include unique elements
  const uniqueWords = new Set(words);
  const hasDuplicates = uniqueWords.size !== words.length;
  if (hasDuplicates) return "Error";

  return words;
}

// analyzeGrid, takes the parsed grid, finds all word slots
// records them for their length and position and validates
// whether the ammount of starts is actually true to the number in the start cell
// so for cell === 2 there have to be 2 starts etc.
function analyzeGrid(grid) {
  const slots = [];

  for (let row = 0; row < grid.length; row++) {
    for (let col = 0; col < grid[0].length; col++) {
      if (!isOpenCell(grid, row, col)) continue;

      const horizontalStart = isHorizontalStart(grid, row, col);
      const verticalStart = isVerticalStart(grid, row, col);

      const actualStarts = Number(horizontalStart) + Number(verticalStart);
      const expectedStarts = Number(grid[row][col]);
      //Compare actualstarts with expected starts to check
      // if amount of starts is actually true to the number of the start cell
      if (actualStarts !== expectedStarts) {
        return "Error";
      }

      if (horizontalStart) {
        slots.push(walkHorizontal(grid, row, col));
      }

      if (verticalStart) {
        slots.push(walkVertical(grid, row, col));
      }
    }
  }
  return slots;
}

// Uses regex to check for any illegal chars and allows only . 0 1 2 \n
function detectIllegalChars(puzzle) {
  let illegalChars = /[^.\n012]/;
  return illegalChars.test(puzzle);
}

// Validates row structure
function splitString(puzzle) {
  let grid = puzzle.split("\n");
  for (let row of grid) {
    if (row.length !== grid[0].length || row === "") {
      return "Error";
    }
  }
  return grid;
}

// isOpenCell checks if the cell is open (not a black square) by checking if it's not a dot
function isOpenCell(grid, row, col) {
  return grid[row][col] !== ".";
}

// isHorizontalStart checks if the cell at (row, col) is the start of a horizontal word slot
function isHorizontalStart(grid, row, col) {
  if (!isOpenCell(grid, row, col)) return false;

  const leftIsBoundary = col === 0;
  const leftIsBlocked = !leftIsBoundary && !isOpenCell(grid, row, col - 1);
  // check that col is not the last column in the grid
  const hasRightCell = col < grid[0].length - 1;
  // check if right cell exists and its open
  const rightIsOpen = hasRightCell && isOpenCell(grid, row, col + 1);

  return (leftIsBoundary || leftIsBlocked) && rightIsOpen;
}

// isVerticalStart checks if the cell at (row, col) is the start of a vertical word slot
function isVerticalStart(grid, row, col) {
  if (!isOpenCell(grid, row, col)) return false;

  const upIsBoundary = row === 0;
  const upIsBlocked = !upIsBoundary && !isOpenCell(grid, row - 1, col);

  const hasDownCell = row < grid.length - 1;

  const downIsOpen = hasDownCell && isOpenCell(grid, row + 1, col);

  return (upIsBoundary || upIsBlocked) && downIsOpen;
}

// walkHorizontal goes through a horizontal word slot and records its length and position
function walkHorizontal(grid, startRow, startCol) {
  let col = startCol;
  let cells = [];

  while (col < grid[0].length && isOpenCell(grid, startRow, col)) {
    cells.push([startRow, col]);
    col += 1;
  }

  return {
    length: cells.length,
    cells,
  };
}

// walkVertical goes through a vertical word slot and records its length and position
function walkVertical(grid, startRow, startCol) {
  let row = startRow;
  let cells = [];

  while (row < grid.length && isOpenCell(grid, row, startCol)) {
    cells.push([row, startCol]);
    row += 1;
  }

  return {
    length: cells.length,
    cells,
  };
}

// Backtracking function to try placing words in the grid
function solve(wordIndex, cleanWords, wordSlots, charGrid, solutions) {
  // If we've placed all words, we found a solution!
  if (wordIndex === cleanWords.length) {
    solutions.push(charGrid.map((row) => row.join("")).join("\n"));
    return;
  }

  // Stop early if we find more than one solution (requirement: unique solution)
  if (solutions.length > 1) return;

  const currentWord = cleanWords[wordIndex];

  for (let i = 0; i < wordSlots.length; i++) {
    const slot = wordSlots[i];

    if (!slot.filled && slot.length === currentWord.length) {
      if (canPlace(currentWord, slot, charGrid)) {
        const originalChars = placeWord(currentWord, slot, charGrid);
        slot.filled = true;

        solve(wordIndex + 1, cleanWords, wordSlots, charGrid, solutions);

        // Backtrack: Undo the placement to try other combinations
        slot.filled = false;
        removeWord(slot, originalChars, charGrid);
      }
    }
  }
}

// Check if a word can be placed in a slot without conflicting with existing characters
function canPlace(word, slot, charGrid) {
  for (let i = 0; i < word.length; i++) {
    const [r, c] = slot.cells[i];
    const existingChar = charGrid[r][c];
    if (existingChar !== "" && existingChar !== word[i]) {
      return false;
    }
  }
  return true;
}

// Backtracking function to place a word in the grid, returning the original characters for backtracking
function placeWord(word, slot, charGrid) {
  const original = [];
  for (let i = 0; i < word.length; i++) {
    const [r, c] = slot.cells[i];
    original.push(charGrid[r][c]);
    charGrid[r][c] = word[i];
  }
  return original;
}

// Backtracking function to remove a word from the grid, restoring original characters
function removeWord(slot, originalChars, charGrid) {
  for (let i = 0; i < slot.cells.length; i++) {
    const [r, c] = slot.cells[i];
    charGrid[r][c] = originalChars[i];
  }
}

// Export for use in test files, comment out if you want to run the tests  "node test.js"
//module.exports = crosswordSolver;

/*
//Example from crossword task description, comment out if you want to run the command "node crosswordSolver.js"
const emptyPuzzle = `2001
0..0
1000
0..0`;
const words = ["casa", "alan", "ciao", "anta"];

crosswordSolver(emptyPuzzle, words);
*/
