package strategies

import "tictactoe/models"

type DiagonalWinningStrategy struct{}

func (DiagonalWinningStrategy) HasWon(board *models.Board, row, col int, symbol models.Symbol) bool {
	n := board.Size

	// Main diagonal: (0,0), (1,1), (2,2) ... only if the move is on it
	if row == col {
		won := true
		for i := 0; i < n; i++ {
			if board.Grid[i][i] != symbol {
				won = false
				break
			}
		}
		if won {
			return true
		}
	}

	// Anti-diagonal: (0,2), (1,1), (2,0) ... only if the move is on it
	if row+col == n-1 {
		won := true
		for i := 0; i < n; i++ {
			if board.Grid[i][n-1-i] != symbol {
				won = false
				break
			}
		}
		if won {
			return true
		}
	}

	return false
}
