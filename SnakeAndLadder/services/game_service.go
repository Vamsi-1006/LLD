package services

import (
	"errors"
	"fmt"
	"snakeandladder/models"
	"snakeandladder/strategies"
)

type Game struct {
	board   *models.Board
	players []*models.Player
	dice    strategies.Dice
	turn    int
	status  models.GameStatus
	winner  *models.Player
}

func NewGame(board *models.Board, players []*models.Player, dice strategies.Dice) (*Game, error) {

	if len(players) < 2 {
		return nil, errors.New("need at least 2 players")
	}
	return &Game{
		board:   board,
		players: players,
		dice:    dice,
	}, nil
}

func (g *Game) PlayTurn() {
	player := g.players[g.turn]
	roll := g.dice.Roll()
	start := player.Position
	next := start + roll

	if next > g.board.Size {
		fmt.Printf("The final position %d exceeded board size, you must land on the board size %d", next, g.board.Size)
	} else {
		message := ""
		if jump, ok := g.board.GetJump(next); ok {
			if jump.IsSnake() {
				message = fmt.Sprintf(", snake! down to %d", jump.End)
			} else {
				message = fmt.Sprintf(", ladder! up to %d", jump.End)
			}
			next = jump.End
		}
		player.Position = next
		fmt.Printf("%s ", message)
		fmt.Printf("%s new position is %d", player.Name, player.Position)
	}

	if player.Position == g.board.Size {
		g.status = models.Finished
		g.winner = player
		return
	}

	g.turn = (g.turn + 1) % len(g.players) // next player, wraps around

}

func (g *Game) Play() *models.Player {
	for g.status == models.InProgress { // Infinite loop until a player wins.
		g.PlayTurn()
	}

	return g.winner
}
