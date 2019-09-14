package dfs

import (
	"image"
	"math/rand"
	"sort"

	"github.com/grindlemire/gomaze/pkg/board"
)

// CreateEdges creates edges according to DFS
func CreateEdges(b board.Board) board.Board {
	placeEdges(b.Entrance, b, map[string]struct{}{})
	return b
}

func placeEdges(current *board.Cell, b board.Board, visited map[string]struct{}) {
	// visit ourselves
	visited[current.ID] = struct{}{}

	// if we have reached the end return
	if current.ID == b.Exit.ID {
		return
	}

	// get the unvisited neighbor
	directions, unvisitedNeighbors := getUnvisitedNeighbors(current, b, visited)
	// if no unvisited neighbors then return
	if len(unvisitedNeighbors) == 0 {
		return
	}

	for i, next := range unvisitedNeighbors {
		direction := directions[i]
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
		placeEdges(next, b, visited)
	}
}

func getUnvisitedNeighbors(current *board.Cell, b board.Board, visited map[string]struct{}) (directions []board.Direction, unvisitedNeighbors []*board.Cell) {
	if current.X > 0 {
		left := b.Cells[current.Y][current.X-1]
		_, leftVisited := visited[left.ID]
		if !leftVisited {
			unvisitedNeighbors = append(unvisitedNeighbors, left)
			directions = append(directions, board.Left)
		}
	}

	if current.X < len(b.Cells[0])-1 {
		right := b.Cells[current.Y][current.X+1]
		_, rightVisited := visited[right.ID]
		if !rightVisited {
			unvisitedNeighbors = append(unvisitedNeighbors, right)
			directions = append(directions, board.Right)
		}
	}

	if current.Y > 0 {
		up := b.Cells[current.Y-1][current.X]
		_, upVisited := visited[up.ID]
		if !upVisited {
			unvisitedNeighbors = append(unvisitedNeighbors, up)
			directions = append(directions, board.Up)
		}
	}

	if current.Y < len(b.Cells)-1 {
		down := b.Cells[current.Y+1][current.X]
		_, downVisited := visited[down.ID]
		if !downVisited {
			unvisitedNeighbors = append(unvisitedNeighbors, down)
			directions = append(directions, board.Down)
		}
	}

	sort.Sort(byBias{b.Bias, unvisitedNeighbors, directions, b.CellWidth})
	return directions, unvisitedNeighbors
}

type byBias struct {
	bias       image.Image
	cells      []*board.Cell
	directions []board.Direction
	cellWidth  int
}

func (b byBias) Len() int { return len(b.cells) }
func (b byBias) Swap(i, j int) {
	b.cells[i], b.cells[j] = b.cells[j], b.cells[i]
	b.directions[i], b.directions[j] = b.directions[j], b.directions[i]
}
func (b byBias) Less(i, j int) bool {
	first := b.cells[i]
	favg := 0
	for i := -b.cellWidth / 2; i <= b.cellWidth/2; i++ {
		for j := -b.cellWidth / 2; j <= b.cellWidth/2; j++ {
			fcbias := b.bias.At((first.X+1)*b.cellWidth+i, (first.Y+1)*b.cellWidth+j)
			fr, fb, fc, _ := fcbias.RGBA()
			favg += int(uint8(fr)) + int(uint8(fb)) + int(uint8(fc))
		}
	}

	second := b.cells[j]
	savg := 0
	for i := -b.cellWidth / 2; i <= b.cellWidth/2; i++ {
		for j := -b.cellWidth / 2; j <= b.cellWidth/2; j++ {
			scbias := b.bias.At((second.X+1)*b.cellWidth+i, (second.Y+1)*b.cellWidth+j)
			sr, sb, sc, _ := scbias.RGBA()
			savg += int(uint8(sr)) + int(uint8(sb)) + int(uint8(sc))
		}
	}

	if savg == favg {
		return rand.Intn(10) > 5
	}
	return favg < savg
}
