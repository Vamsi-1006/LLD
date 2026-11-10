package models

import (
	"fmt"
)

type Board struct {
	Size  int
	jumps map[int]Jump // start cell -> Jump (snake or ladder)
}

func NewBoard(size int) *Board {
	return &Board{
		Size:  size,
		jumps: make(map[int]Jump),
	}
}

func (b *Board) AddSnake(head, tail int) error {
	if tail >= head {
		return fmt.Errorf("snake head %d must be above it's tail %d", head, tail)
	}

	return b.addJump(head, tail)
}

func (b *Board) AddLadder(bottom, top int) error {
	if top <= bottom {
		return fmt.Errorf("ladder top %d must be above it's bottom %d", top, bottom)
	}

	return b.addJump(bottom, top)
}

// addJump holds the checks that snakes and ladders share.
func (b *Board) addJump(start, end int) error {
	if start <= 1 || start >= b.Size {
		return fmt.Errorf("a jump cannot start at cell %d", start)
	}
	if end < 1 || end > b.Size {
		return fmt.Errorf("a jump cannot end at cell %d", end)
	}
	if _, exists := b.jumps[start]; exists {
		return fmt.Errorf("cell %d already has a snake or ladder", start)
	}
	b.jumps[start] = Jump{Start: start, End: end}
	return nil
}

// GetJump tells you if a cell has a snake or ladder on it.
func (b *Board) GetJump(cell int) (Jump, bool) {
	jump, ok := b.jumps[cell]
	return jump, ok
}
