package board

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Direction a cell can have a connection to
type Direction int

// the directions a cell can have connections to
const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	default:
		panic("invalid direction in print")
	}
}

// GetOpposite gets the opposite direction
func (d Direction) GetOpposite() (other Direction) {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	default:
		panic("invalid direction somehow")
	}
}

// Connection is a connection between arbitrary cells
type Connection struct {
	Cells []*Cell
}

func (c Connection) String() string {
	strs := []string{}

	for _, cell := range c.Cells {
		strs = append(strs, fmt.Sprintf("(%d,%d)", cell.X, cell.Y))
	}

	return strings.Join(strs, " | ")
}

// Cell is a cell in the board
type Cell struct {
	ID          string
	X           int
	Y           int
	Connections map[Direction]Connection

	Visited bool
}

// NewCell creates a new cell
func NewCell(x, y int) (c *Cell) {
	return &Cell{
		ID:          uuid.New().String(),
		X:           x,
		Y:           y,
		Connections: map[Direction]Connection{},
	}
}
