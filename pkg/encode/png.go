package encode

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/grindlemire/gomaze/pkg/board"
)

// PNG is a png encoder
type PNG struct {
	img *image.NRGBA
	b   board.Board
}

// NewPNG creates a new PNG encoder
func NewPNG(b board.Board) (p PNG, err error) {
	width := len(b.Cells[0])
	height := len(b.Cells)
	p = PNG{
		img: image.NewNRGBA(
			image.Rect(0, 0,
				// cellWidth*width+horizontalPadding,
				// cellWidth*height+verticalPadding,
				cellWidth*(width+1),
				cellWidth*(height+1),
			),
		),
		b: b,
	}
	return p, nil
}

// Save saves the png to a filename
func (p PNG) Save(name string) (err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, p.img)
}

const cellWidth = 28
const horizontalPadding = cellWidth / 2
const verticalPadding = cellWidth / 2

// EncodeSolution encodes a solution to the png
func (p *PNG) EncodeSolution(path []*board.Cell) (err error) {
	return p.encodeSolution(board.None, path[0], path[1:])
}

func (p *PNG) drawPathTo(d board.Direction, cell *board.Cell) {
	var xStart, yStart, xEnd, yEnd int
	switch d {
	case board.Up:
		// draw upwards
		yStart = -cellWidth / 2
		yEnd = cellWidth / 10
		xStart = -cellWidth / 10
		xEnd = cellWidth / 10
	case board.Down:
		// draw down
		yStart = -cellWidth / 10
		yEnd = cellWidth / 2
		xStart = -cellWidth / 10
		xEnd = cellWidth / 10
	case board.Left:
		// draw left
		yStart = -cellWidth / 10
		yEnd = cellWidth / 10
		xStart = -cellWidth / 2
		xEnd = cellWidth / 10
	case board.Right:
		// draw right
		yStart = -cellWidth / 10
		yEnd = cellWidth / 10
		xStart = -cellWidth / 10
		xEnd = cellWidth / 2
	case board.None:
		return
	}

	for x := xStart; x <= xEnd; x++ {
		for y := yStart; y <= yEnd; y++ {
			pX := (cell.X+1)*cellWidth + x //+ horizontalPadding/2
			pY := (cell.Y+1)*cellWidth + y //+ verticalPadding/2
			p.img.Set(pX, pY, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		}
	}
}

func (p *PNG) encodeSolution(prevDirection board.Direction, curr *board.Cell, rest []*board.Cell) (err error) {
	p.drawPathTo(prevDirection, curr)
	if len(rest) == 0 {
		return nil
	}

	next := rest[0]
	nextRest := rest[1:]

	var nextDirection board.Direction

lookfornext:
	for direction, connection := range curr.Connections {
		for _, cell := range connection.Cells {
			if cell == next {
				nextDirection = direction
				break lookfornext
			}
		}
	}

	p.drawPathTo(nextDirection, curr)
	return p.encodeSolution(nextDirection.GetOpposite(), next, nextRest)
}

// EncodeBoard encodes cells to the png
func (p *PNG) EncodeBoard() (err error) {
	cells := p.b.Cells
	entrance := p.b.Entrance
	exit := p.b.Exit

	// loop over cells
	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[0]); x++ {
			cell := cells[y][x]

			fillColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			if cell.ID == entrance.ID {
				fillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
			}

			if cell.ID == exit.ID {
				fillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
			}

			// create a square for each cell
			for i := -cellWidth / 2; i <= cellWidth/2; i++ {
				for j := -cellWidth / 2; j <= cellWidth/2; j++ {
					cellFill := fillColor
					// calculate the pixel we are looking at
					pX := (x+1)*cellWidth + i //+ horizontalPadding/2
					pY := (y+1)*cellWidth + j //+ verticalPadding/2

					_, hasLeft := cell.Connections[board.Left]
					_, hasRight := cell.Connections[board.Right]
					_, hasUp := cell.Connections[board.Up]
					_, hasDown := cell.Connections[board.Down]

					// if we are on an edge and there is no connection, make it black
					if (i == -cellWidth/2 && !hasLeft) ||
						(i == cellWidth/2 && !hasRight) ||
						(j == -cellWidth/2 && !hasUp) ||
						(j == cellWidth/2 && !hasDown) ||
						(math.Abs(float64(i)) == cellWidth/2 && math.Abs(float64(j)) == cellWidth/2) && // Fill in the corners
							(cell.ID != entrance.ID) && (cell.ID != exit.ID) { // Don't make walls for the entrance or exit
						cellFill = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
					}

					p.img.Set(pX, pY, cellFill)
				}
			}
		}
	}

	return nil
}
