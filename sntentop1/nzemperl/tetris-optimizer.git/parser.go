package main

import (
	"errors"
	"os"
	"strings"
)

func ParseFile(filePath string) ([]Tetromino, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		return nil, errors.New("ERROR")
	}

	text := string(content)
	blocks := strings.Split(strings.TrimSpace(text), "\n\n")

	var tetrominoes []Tetromino
	currentLetter := 'A'

	for _, block := range blocks {
		lines := strings.Split(block, "\n")
		if len(lines) != 4 {
			return nil, errors.New("ERROR")
		}

		var shape [4][4]rune
		hashCount := 0

		for i, line := range lines {
			if len(line) != 4 {
				return nil, errors.New("ERROR")
			}
			for j, char := range line {
				if char != '.' && char != '#' {
					return nil, errors.New("ERROR")
				}
				if char == '#' {
					hashCount++
				}
				shape[i][j] = char
			}
		}

		if hashCount != 4 || !isValidTetromino(shape) {
			return nil, errors.New("ERROR")
		}

		tetromino := Tetromino{
			Shape:  shape,
			Letter: currentLetter,
		}
		tetromino.CalculateBounds()

		tetrominoes = append(tetrominoes, tetromino)
		currentLetter++
	}

	return tetrominoes, nil
}
func isValidTetromino(shape [4][4]rune) bool {
	connections := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if shape[i][j] == '#' {
				// Check adjacent cells
				if i > 0 && shape[i-1][j] == '#' {
					connections++
				}
				if i < 3 && shape[i+1][j] == '#' {
					connections++
				}
				if j > 0 && shape[i][j-1] == '#' {
					connections++
				}
				if j < 3 && shape[i][j+1] == '#' {
					connections++
				}
			}
		}
	}
	return connections >= 6 // A valid tetromino has at least 6 connections
}
