package services

import (
	"errors"
	"fmt"
	"tictactoe/models"
	"tictactoe/strategies"
)

type GameService struct {
	board             *models.Board // Pointer because it is a struct
	players           []*models.Player
	turn              int                          // index of player's turn
	status            models.GameStatus            // No pointer since it is an enum
	winner            *models.Player               // Pointer because it is a struct
	winningStrategies []strategies.WinningStrategy // No Pointer because it is an interface
}

func NewGameService(size int, players []*models.Player, ws []strategies.WinningStrategy) *GameService {
	return &GameService{
		board:   models.NewBoard(size), // No pointer since it is a constructor function.
		players: players,
		//status: models.InProgress, // Not needed, since bydefault the status type is IN_PROGRESS which we declared as enum.
		winningStrategies: ws,
	}
}

func (g *GameService) CurrentPlayer() *models.Player {
	return g.players[g.turn]
}

func (g *GameService) Status() models.GameStatus {
	return g.status
}

func (g *GameService) Winner() *models.Player {
	return g.winner
}

func (g *GameService) PrintBoard() {
	g.board.Print()
}

// MakeMove: validate → update model → ask strategies → update state.
func (g *GameService) MakeMove(row, col int) error {
	// 1. Validate game state
	if g.status != models.InProgress {
		return errors.New("game is already over")
	}

	// 2. Validate the move
	if !g.board.IsValidMove(row, col) {
		return fmt.Errorf("invalid move at (%d,%d): out of bounds or cell already taken", row, col)
	}

	// 3. Update the model
	player := g.CurrentPlayer()
	g.board.Place(row, col, player.Symbol)

	// 4. Ask every winning strategy
	for _, s := range g.winningStrategies {
		if s.HasWon(g.board, row, col, player.Symbol) {
			g.status = models.Won
			g.winner = player
			return nil
		}
	}

	// 5. Draw?
	if g.board.IsFull() {
		g.status = models.Draw
		return nil
	}

	// 6. Next player's turn
	g.turn = (g.turn + 1) % len(g.players)
	return nil
}
