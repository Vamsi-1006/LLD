package strategies

import "math/rand"

// For generating even numbers only
type CrookedDice struct{}

func NewCrookedDice() *CrookedDice {
	return &CrookedDice{}
}

func (d *CrookedDice) Roll() int {
	return (rand.Intn(3) + 1) * 2
}
