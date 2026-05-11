package board

func CreateBoard() [][]int {
	board := make([][]int, 9)
	for i := range board {
		board[i] = make([]int, 9)
	}
	return board
}
