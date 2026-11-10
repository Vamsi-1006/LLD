package strategies

import "tictactoe/models"

// WinningStrategy checks whether the move just played at (row, col) won the game.
// We only check lines that pass through the last move, which is faster than scanning everything.
// Here we use strategy design pattern where we have method in an interface and we implement that interface with different structs.
type WinningStrategy interface {
	HasWon(board *models.Board, row, col int, symbol models.Symbol) bool
}
