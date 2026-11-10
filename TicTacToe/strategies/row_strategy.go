package strategies

import "tictactoe/models"

type RowWinningStrategy struct{}

func (RowWinningStrategy) HasWon(board *models.Board, row, col int, symbol models.Symbol) bool {
	for c := 0; c < board.Size; c++ {
		if board.Grid[row][c] != symbol {
			return false
		}
	}
	return true
}
