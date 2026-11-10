package strategies

import "tictactoe/models"

type ColumnWinningStrategy struct{}

func (ColumnWinningStrategy) HasWon(board *models.Board, row, col int, symbol models.Symbol) bool {
	for r := 0; r < board.Size; r++ {
		if board.Grid[r][col] != symbol {
			return false
		}
	}
	return true
}
