package main

import (
	"fmt"

	"tictactoe/models"
	"tictactoe/services"
	"tictactoe/strategies"
)

func main() {

	// Initialize Players
	alice := models.NewPlayer("Alice", models.X)
	bob := models.NewPlayer("Bob", models.O)

	// Winning Rules
	winRules := []strategies.WinningStrategy{
		strategies.RowWinningStrategy{},
		strategies.ColumnWinningStrategy{},
		strategies.DiagonalWinningStrategy{},
	}

	var n int
	fmt.Println(" Enter the size of the board : ")
	fmt.Scan(&n)

	game := services.NewGameService(n, []*models.Player{alice, bob}, winRules)

	var row int
	var col int

	for i := 0; i < n*n; i++ { // here we can also write for game.Status() == models.InProgress and remove i-- statement.
		fmt.Printf("%s turn", game.CurrentPlayer().Name)
		fmt.Println()
		fmt.Println(" Enter the row : ")
		fmt.Scan(&row)
		fmt.Println(" Enter the col : ")
		fmt.Scan(&col)

		if err := game.MakeMove(row, col); err != nil {
			fmt.Printf(" Error : %v\n\n", err)
			i--
			continue
		}

		game.PrintBoard()

		if game.Status() != models.InProgress {
			break
		}
	}

	switch game.Status() {
	case models.Won:
		w := game.Winner()
		fmt.Printf("Game over! Status: %s, Winner: %s (%s)\n", game.Status(), w.Name, w.Symbol)
	case models.Draw:
		fmt.Println("Game over! It's a DRAW")
	default:
		fmt.Println("Game still in progress")
	}

}
