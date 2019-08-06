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
	cells = generateEdges(entrance, exit, cells)

	return Board{
		Entrance: entrance,
		Exit:     exit,
		Cells:    cells,
	}
}

func generateEdges(current, exit *Cell, cells [][]*Cell) [][]*Cell {
	// visit ourselves
	current.visited = true

	// if we have reached the end return
	if current.ID == exit.ID {
		return cells
	}

	// get the unvisited neighbor
	unvisitedNeighbors := getUnvisitedNeighbors(current, cells)

	// if no unvisited neighbors then return
	if len(unvisitedNeighbors) == 0 {
		return cells
	}

	for direction, next := range unvisitedNeighbors {
		// if it has been visited in the meantime continue
		if next.visited {
			continue
		}
		// create a connection and set it
		connection := Connection{
			Cells: []*Cell{next, current},
		}
		current.Connections[direction] = connection
		next.Connections[direction.getOpposite()] = connection
		// recurse down
		cells = generateEdges(next, exit, cells)
	}

	return cells
}

func getUnvisitedNeighbors(current *Cell, cells [][]*Cell) (unvisitedNeighbors map[Direction]*Cell) {
	unvisitedNeighbors = map[Direction]*Cell{}
	if current.X > 0 && !cells[current.Y][current.X-1].visited {
		unvisitedNeighbors[Left] = cells[current.Y][current.X-1]
	}
	if current.X < len(cells[0])-1 && !cells[current.Y][current.X+1].visited {
		unvisitedNeighbors[Right] = cells[current.Y][current.X+1]
	}
	if current.Y > 0 && !cells[current.Y-1][current.X].visited {
		unvisitedNeighbors[Up] = cells[current.Y-1][current.X]
	}
	if current.Y < len(cells)-1 && !cells[current.Y+1][current.X].visited {
		unvisitedNeighbors[Down] = cells[current.Y+1][current.X]
	}

	return unvisitedNeighbors
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
