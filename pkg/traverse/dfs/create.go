package dfs

import "github.com/grindlemire/gomaze/pkg/board"

// CreateEdges creates edges according to DFS
func CreateEdges(b board.Board) board.Board {
	b.Cells = placeEdges(b.Entrance, b.Exit, b.Cells, map[string]struct{}{})
	return b
}

func placeEdges(current, exit *board.Cell, cells [][]*board.Cell, visited map[string]struct{}) [][]*board.Cell {
	// visit ourselves
	visited[current.ID] = struct{}{}

	// if we have reached the end return
	if current.ID == exit.ID {
		return cells
	}

	// get the unvisited neighbor
	unvisitedNeighbors := getUnvisitedNeighbors(current, cells, visited)

	// if no unvisited neighbors then return
	if len(unvisitedNeighbors) == 0 {
		return cells
	}

	for direction, next := range unvisitedNeighbors {
		// if it has been visited in the meantime continue
		_, isVisited := visited[next.ID]
		if isVisited {
			continue
		}

		// create a connection and set it
		connection := board.Connection{
			Cells: []*board.Cell{next, current},
		}
		current.Connections[direction] = connection
		next.Connections[direction.GetOpposite()] = connection

		// recurse down
		cells = placeEdges(next, exit, cells, visited)
	}

	return cells
}

func getUnvisitedNeighbors(current *board.Cell, cells [][]*board.Cell, visited map[string]struct{}) (unvisitedNeighbors map[board.Direction]*board.Cell) {
	unvisitedNeighbors = map[board.Direction]*board.Cell{}

	if current.X > 0 {
		left := cells[current.Y][current.X-1]
		_, leftVisited := visited[left.ID]
		if !leftVisited {
			unvisitedNeighbors[board.Left] = left
		}
	}

	if current.X < len(cells[0])-1 {
		right := cells[current.Y][current.X+1]
		_, rightVisited := visited[right.ID]
		if !rightVisited {
			unvisitedNeighbors[board.Right] = right
		}
	}

	if current.Y > 0 {
		up := cells[current.Y-1][current.X]
		_, upVisited := visited[up.ID]
		if !upVisited {
			unvisitedNeighbors[board.Up] = up
		}
	}

	if current.Y < len(cells)-1 {
		down := cells[current.Y+1][current.X]
		_, downVisited := visited[down.ID]
		if !downVisited {
			unvisitedNeighbors[board.Down] = down
		}
	}

	return unvisitedNeighbors
}
