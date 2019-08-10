package dfs

import "github.com/grindlemire/gomaze/pkg/board"

// Solve solves the board and returns the path
func Solve(b board.Board) (path []*board.Cell) {
	path, _ = traverse(b.Entrance, b.Exit, b.Cells, map[string]struct{}{})
	return path
}

func traverse(current, exit *board.Cell, cells [][]*board.Cell, visited map[string]struct{}) (path []*board.Cell, found bool) {
	// visit ourselves
	visited[current.ID] = struct{}{}

	// if we have reached the end return
	if current.ID == exit.ID {
		return []*board.Cell{current}, true
	}

	// get the unvisited neighbor
	unvisitedNeighbors := getConnectedUnvisitedNeighbors(current, visited)

	// if no unvisited neighbors then return
	if len(unvisitedNeighbors) == 0 {
		return append(path, current), false
	}

	for _, next := range unvisitedNeighbors {
		// if it has been visited in the meantime continue
		_, isVisited := visited[next.ID]
		if isVisited {
			continue
		}

		// recurse down
		path, found = traverse(next, exit, cells, visited)
		if found {
			path = append(path, current)
			return path, found
		}

	}

	return path, false
}

func getConnectedUnvisitedNeighbors(current *board.Cell, visited map[string]struct{}) (unvisitedNeighbors map[board.Direction]*board.Cell) {
	unvisitedNeighbors = map[board.Direction]*board.Cell{}

	for direction, connection := range current.Connections {
		for _, cell := range connection.Cells {
			if cell.ID == current.ID {
				continue
			}

			_, cellVisited := visited[cell.ID]
			if !cellVisited {
				unvisitedNeighbors[direction] = cell
			}
		}
	}

	return unvisitedNeighbors
}
