package sudokuboard

import (
	"fmt"
	"strconv"

	"sudoku/cell"
)

const N = 9

type Sudoku struct {
	Board    [N][N]cell.Cell
	Solved   bool
	Solution [N][N]cell.Cell
}

func (s *Sudoku) CreateBoard(rows []string) string {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if rows[i][j] == '.' {
				s.Board[i][j] = cell.NewCell(i, j, 0)
			} else {
				num, _ := strconv.Atoi(string(rows[i][j]))
				if s.PlaceNumIsValid(cell.Cell{Row: i, Col: j}, num) {
					s.Board[i][j] = cell.NewCell(i, j, num)
				} else {
					return "Error: duplicate numbers in row, col or 3x3 subgrid"
				}
			}
		}
	}
	s.UpdatePossibleValues()
	return ""
}

func (s Sudoku) PlaceNumIsValid(c cell.Cell, num int) bool {
	row, col := c.Row, c.Col

	for i := 0; i < N; i++ {
		if s.Board[row][i].Value == num || s.Board[i][col].Value == num {
			return false
		}
	}

	startRow := (row / 3) * 3
	startCol := (col / 3) * 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.Board[startRow+i][startCol+j].Value == num {
				return false
			}
		}
	}
	return true
}

func (s *Sudoku) UpdatePossibleValues() {
	for _, row := range s.Board {
		for _, cell := range row {
			if cell.IsEmpty() {
				var possible [9]bool
				for k := 0; k < 9; k++ {
					possible[k] = s.PlaceNumIsValid(cell, k+1)
				}
				// Update the possible values of the cell in the board
				s.Board[cell.Row][cell.Col].SetPossibleValues(possible)
			}
		}
	}
}

func (s Sudoku) GetCellWithLeastPossibleValues() (cell.Cell, bool) {
	leastPossible := 10
	var smallestCell cell.Cell
	found := false

	for _, row := range s.Board {
		for _, c := range row {
			if c.IsEmpty() {
				numPossibleValues := len(c.GetPossibleValues())
				if numPossibleValues < leastPossible {
					leastPossible = numPossibleValues
					smallestCell = c
					found = true
				}
			}
		}
	}
	return smallestCell, found
}

func (s *Sudoku) SolveSudoku() bool {
	c, found := s.GetCellWithLeastPossibleValues()
	if !found {
		s.Solution = s.Board
		s.Solved = true
		return true
	}
	for _, value := range c.GetPossibleValues() {
		if s.PlaceNumIsValid(c, value) {
			s.Board[c.Row][c.Col] = cell.NewCell(c.Row, c.Col, value)
			s.UpdatePossibleValues()
			if s.SolveSudoku() {
				s.Solved = true
				return true
			}
			s.Board[c.Row][c.Col] = cell.NewCell(c.Row, c.Col, 0)
			s.UpdatePossibleValues()
		}
	}
	return false
}

func (s *Sudoku) ReverseSolveSudoku() bool {
	c, found := s.GetCellWithLeastPossibleValues()
	if !found {
		s.Solution = s.Board
		s.Solved = true
		return true
	}

	values := c.GetPossibleValues()
	reverseSlice(values)

	for _, value := range values {
		if s.PlaceNumIsValid(c, value) {
			s.Board[c.Row][c.Col] = cell.NewCell(c.Row, c.Col, value)
			s.UpdatePossibleValues()
			if s.ReverseSolveSudoku() {
				s.Solved = true
				s.Solution = s.Board
				return true
			}
			s.Board[c.Row][c.Col] = cell.NewCell(c.Row, c.Col, 0)
			s.UpdatePossibleValues()
		}
	}
	return false
}

func (s Sudoku) PrintBoard() {
	for _, row := range s.Board {
		for _, cell := range row {
			fmt.Printf("%v ", cell.Value)
		}
		fmt.Print("\n")
	}
}

func reverseSlice(slice []int) []int {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
