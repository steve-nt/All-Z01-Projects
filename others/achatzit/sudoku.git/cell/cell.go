package cell

type Cell struct {
	Row     int
	Col     int
	Value   int
	Numbers [9]bool
}

func NewCell(row, col, value int) Cell {
	cell := Cell{Row: row, Col: col, Value: value}
	if value == 0 {
		for i := 0; i < 9; i++ {
			cell.Numbers[i] = true
		}
	} else {
		for i := 0; i < 9; i++ {
			cell.Numbers[i] = false
		}
		cell.Numbers[value-1] = true
	}
	return cell
}

func (cell Cell) IsEmpty() bool {
	return cell.Value == 0
}

func (cell Cell) GetPossibleValues() []int {
	var possibleValues []int
	for i, v := range cell.Numbers {
		if v {
			possibleValues = append(possibleValues, i+1)
		}
	}
	return possibleValues
}

func (cell *Cell) SetPossibleValues(values [9]bool) {
	cell.Numbers = values
}
