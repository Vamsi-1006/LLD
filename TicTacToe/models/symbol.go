package models

type Symbol int

// iota is by default value of Symbol, so in this case it is Empty
const (
	Empty Symbol = iota // 0
	X                   // 1
	O                   // 2
)

func (s Symbol) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	default:
		return "-"
	}
}
