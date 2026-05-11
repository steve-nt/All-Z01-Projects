package main

type Tetromino struct {
	Letter rune       // A, B, C, etc.
	Shape  [4][4]rune // 4x4 grid of '#' and '.'
	Width  int        // actual width based on #
	Height int        // actual height based on #
	MinRow int
	MinCol int
}

func (t *Tetromino) CalculateBounds() {
	minRow, maxRow := 4, 0
	minCol, maxCol := 4, 0

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if t.Shape[row][col] == '#' {
				if row < minRow {
					minRow = row
				}
				if row > maxRow {
					maxRow = row
				}
				if col < minCol {
					minCol = col
				}
				if col > maxCol {
					maxCol = col
				}
			}
		}
	}

	t.Height = maxRow - minRow + 1
	t.Width = maxCol - minCol + 1
	t.MinRow = minRow
	t.MinCol = minCol
}
