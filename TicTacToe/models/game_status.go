package models

type GameStatus int

// iota is bydefault value of GameStatus and it is InProgress in this case.
const (
	InProgress GameStatus = iota
	Won
	Draw
)

func (g GameStatus) String() string {
	switch g {
	case Won:
		return "WON"
	case Draw:
		return "DRAW"
	default:
		return "IN_PROGRESS"
	}
}
