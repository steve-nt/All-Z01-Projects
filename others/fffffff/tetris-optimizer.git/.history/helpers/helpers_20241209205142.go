package tetris

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// NormalizeTetrominos trims the empty space around each tetromino
func NormalizeTetrominos(tetrominos map[rune][][]int) map[rune][][]int {
	normalized := make(map[rune][][]int)
	for name, tetromino := range tetrominos {
		minRow, minCol, maxRow, maxCol := findBounds(tetromino)

		// Extract the submatrix that represents the normalized tetromino
		normalizedTetromino := make([][]int, maxRow-minRow+1)
		for i := range normalizedTetromino {
			normalizedTetromino[i] = make([]int, maxCol-minCol+1)
			for j := range normalizedTetromino[i] {
				normalizedTetromino[i][j] = tetromino[minRow+i][minCol+j]
			}
		}
		normalized[name] = normalizedTetromino
	}
	return normalized
}

// findBounds finds the minimum and maximum bounds of the tetromino
func findBounds(tetromino [][]int) (int, int, int, int) {
	minRow, minCol := len(tetromino), len(tetromino[0])
	maxRow, maxCol := 0, 0
	for i := range tetromino {
		for j := range tetromino[i] {
			if tetromino[i][j] == 1 {
				if i < minRow {
					minRow = i
				}
				if i > maxRow {
					maxRow = i
				}
				if j < minCol {
					minCol = j
				}
				if j > maxCol {
					maxCol = j
				}
			}
		}
	}
	return minRow, minCol, maxRow, maxCol
}

// SortTetrominosByArea sorts tetrominos by their area in descending order
func SortTetrominosByArea(tetrominos map[rune][][]int) []rune {
	type Tetromino struct {
		Name rune
		Area int
	}
	var sortedTetrominos []Tetromino

	// Calculate areas for each tetromino
	for name, tetromino := range tetrominos {
		area := countBlocks(tetromino)
		sortedTetrominos = append(sortedTetrominos, Tetromino{Name: name, Area: area})
	}

	// Sort by area in descending order
	sort.Slice(sortedTetrominos, func(i, j int) bool {
		return sortedTetrominos[i].Area < sortedTetrominos[j].Area
	})

	// Extract the names in sorted order
	var sortedNames []rune
	for _, tet := range sortedTetrominos {
		sortedNames = append(sortedNames, tet.Name)
	}
	return sortedNames
}

// countBlocks returns the total number of blocks in a tetromino
func countBlocks(tetromino [][]int) int {
	count := 0
	for _, row := range tetromino {
		for _, cell := range row {
			if cell == 1 {
				count++
			}
		}
	}
	return count
}

// SolveTetris uses backtracking to find the minimal grid solution
func SolveTetris(tetrominos map[rune][][]int, sortedTetrominos []rune) [][]int {
	// Calculate the total area of all tetrominos
	minArea := 0
	for _, name := range sortedTetrominos {
		minArea += countBlocks(tetrominos[name])
	}

	// Start with an estimated minimum grid size
	gridSize := int(math.Ceil(math.Sqrt(float64(minArea))))
	for {
		// Create an empty grid of the current size
		grid := make([][]int, gridSize)
		for i := range grid {
			grid[i] = make([]int, gridSize)
		}

		// Attempt to solve the puzzle with the current grid size
		if backtrackSolve(grid, tetrominos, sortedTetrominos, 0) {
			return grid
		}

		// Increment the grid size if no solution is found
		gridSize++
	}
}

// backtrackSolve tries to place each tetromino on the grid recursively
func backtrackSolve(grid [][]int, tetrominos map[rune][][]int, sortedTetrominos []rune, currentIndex int) bool {
	if currentIndex == len(sortedTetrominos) {
		return true // All tetrominos have been placed successfully
	}

	// Get the current tetromino to place
	currentTetrominoName := sortedTetrominos[currentIndex]
	currentTetromino := tetrominos[currentTetrominoName]

	// Try placing the current tetromino at each possible position
	for row := 0; row <= len(grid)-len(currentTetromino); row++ {
		for col := 0; col <= len(grid[0])-len(currentTetromino[0]); col++ {
			if canPlaceTetromino(grid, currentTetromino, row, col) {
				// Place the tetromino
				placeTetromino(grid, currentTetromino, currentTetrominoName, row, col)

				// Recurse to the next tetromino
				if backtrackSolve(grid, tetrominos, sortedTetrominos, currentIndex+1) {
					return true
				}

				// If placement fails, remove the tetromino and backtrack
				removeTetromino(grid, currentTetromino, row, col)
			}
		}
	}

	return false
}

// canPlaceTetromino checks if a tetromino can be placed at a given position
func canPlaceTetromino(grid [][]int, tetromino [][]int, row, col int) bool {
	for i := range tetromino {
		for j := range tetromino[i] {
			if tetromino[i][j] == 1 {
				if row+i >= len(grid) || col+j >= len(grid[0]) || grid[row+i][col+j] != 0 {
					return false
				}
			}
		}
	}
	return true
}

// placeTetromino places a tetromino on the grid
func placeTetromino(grid [][]int, tetromino [][]int, name rune, row, col int) {
	for i := range tetromino {
		for j := range tetromino[i] {
			if tetromino[i][j] == 1 {
				grid[row+i][col+j] = int(name) // Use ASCII value to represent the tetromino
			}
		}
	}
}

// removeTetromino removes a tetromino from the grid
func removeTetromino(grid [][]int, tetromino [][]int, row, col int) {
	for i := range tetromino {
		for j := range tetromino[i] {
			if tetromino[i][j] == 1 {
				grid[row+i][col+j] = 0
			}
		}
	}
}

// PrintGrid prints the grid to the console
func PrintGrid(grid [][]int) {
	for _, row := range grid {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print(".")
			} else {
				fmt.Print(string(cell))
			}
		}
		fmt.Println()
	}
}

// PrintTetromino prints a tetromino to the console
func PrintTetromino(tetromino [][]int) {
	for _, row := range tetromino {
		for _, cell := range row {
			if cell == 1 {
				fmt.Print("#")
			} else if cell == 0 {
				fmt.Print(".")
			} else {
				fmt.Print("ERROR")
			}
		}
		fmt.Println()
	}
}

// GetTetrominos reads tetromino definitions from a file
func GetTetrominos(file *os.File) (map[rune][][]int, error) {
	tetrominoMap := make(map[rune][][]int)
	scanner := bufio.NewScanner(file)
	var inoName rune = 'A'
	var rowCount int
	var currentTetromino [][]int

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Save the current tetromino
			if len(currentTetromino) > 0 {
				tetrominoMap[inoName] = currentTetromino
				currentTetromino = nil
				inoName++
				rowCount = 0
			}
		} else {
			if len(currentTetromino) <= rowCount {
				currentTetromino = append(currentTetromino, []int{})
			}
			for _, char := range line {
				if char == '#' {
					currentTetromino[rowCount] = append(currentTetromino[rowCount], 1)
				} else if char == '.' {
					currentTetromino[rowCount] = append(currentTetromino[rowCount], 0)
				} else {
					return nil, fmt.Errorf("ERROR: Invalid character found in the tetromino")
				}
			}
			rowCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Add the last tetromino if necessary
	if len(currentTetromino) > 0 {
		tetrominoMap[inoName] = currentTetromino
	}

	for _, t := range tetrominoMap {
		if !IsValid(t) {
			return nil, fmt.Errorf("ERROR")
		}
		count := 0
		for _, k := range t {
			for _, l := range k {
				if l == 1 {
					count++
				}

			}

		}
		if count != 4 {
			return nil, fmt.Errorf("ERROR")
		}

	}

	return tetrominoMap, nil
}

var directions = [][2]int{
	{0, 1},
	{0, -1},
	{1, 0},
	{-1, 0},
}

func dfs(grid [][]int, x, y int, visited [][]bool) int {
	if x < 0 || x >= len(grid) || y >= len(grid[0]) || y < 0 || visited[x][y] || grid[x][y] != 1 {
		return 0
	}
	visited[x][y] = true
	count := 1

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		count += dfs(grid, nx, ny, visited)

	}
	return count
}
func IsValid(grid [][]int) bool {
	visited := make([][]bool, len(grid))
	for i := range visited {
		visited[i] = make([]bool, len(grid[0]))
	}

	connected := 0
	foundStart := false

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == 1 && !visited[i][j] {
				if foundStart {
					return false
				}
				connected = dfs(grid, i, j, visited)
				foundStart = true
			}
		}
	}
	return connected == 4
}
