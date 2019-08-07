package board

import (
	"fmt"
)

// Board is the maze board
type Board struct {
	Entrance *Cell
	Exit     *Cell
	Cells    [][]*Cell
}

// New creates a new fully initialized board
func New(width, height int) (b Board) {
	cells := [][]*Cell{}
	for i := 0; i < height; i++ {
		row := []*Cell{}
		for j := 0; j < width; j++ {
			c := NewCell(j, i)
			row = append(row, c)
		}
		cells = append(cells, row)
	}
	entrance := cells[0][0]
	exit := cells[height-1][width-1]
	return Board{
		Entrance: entrance,
		Exit:     exit,
		Cells:    cells,
	}
}

func (b Board) String() string {
	s := ""
	for _, row := range b.Cells {
		for range row {
			s += fmt.Sprintf(". ")
		}
		s += fmt.Sprintf("\n")
	}
	return s
}
