package strategies

import "math/rand"

type NormalDice struct {
	faces int
}

func NewNormalDice(faces int) *NormalDice {
	return &NormalDice{
		faces: faces,
	}
}

func (d *NormalDice) Roll() int {
	return rand.Intn(d.faces) + 1
}
