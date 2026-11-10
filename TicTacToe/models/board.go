package models

import "fmt"

type Board struct {
	Size int
	Grid [][]Symbol
}

func NewBoard(size int) *Board {
	grid := make([][]Symbol, size)

	for i := range grid {
		grid[i] = make([]Symbol, size)
	}

	return &Board{
		Size: size,
		Grid: grid,
	}
}

func (b *Board) IsValidMove(row, col int) bool {
	if row < 0 || row >= b.Size || col < 0 || col >= b.Size || b.Grid[row][col] != Empty {
		return false
	}

	return true
}

func (b *Board) Place(row, col int, symbol Symbol) {
	b.Grid[row][col] = symbol
}

func (b *Board) IsFull() bool {
	for _, row := range b.Grid {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}

	return true
}

func (b *Board) Print() {
	for _, row := range b.Grid {
		for c, cell := range row {
			if c > 0 {
				fmt.Print(" ")
			}
			fmt.Print(cell)
		}
		fmt.Println()
	}
}
