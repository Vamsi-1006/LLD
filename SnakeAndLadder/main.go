package main

import (
	"fmt"

	"snakeandladder/models"
	"snakeandladder/services"
	"snakeandladder/strategies"
)

func main() {
	var n int
	fmt.Println("Enter board size: ")
	fmt.Scan(&n)

	board := models.NewBoard(n)

	snakes := [][2]int{{99, 54}, {70, 55}, {52, 42}, {25, 2}, {95, 72}}
	ladders := [][2]int{{6, 31}, {11, 40}, {60, 85}, {46, 90}, {17, 69}}

	for _, s := range snakes {
		if err := board.AddSnake(s[0], s[1]); err != nil {
			fmt.Println("error:", err)
			return
		}
	}

	for _, l := range ladders {
		if err := board.AddLadder(l[0], l[1]); err != nil {
			fmt.Println("error:", err)
			return
		}
	}

	players := []*models.Player{
		models.NewPlayer("Alice"),
		models.NewPlayer("Bob"),
	}

	dice := strategies.NewNormalDice(6)

	game, err := services.NewGame(board, players, dice)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	winner := game.Play()
	fmt.Printf("Winner is %s", winner.Name)
}
