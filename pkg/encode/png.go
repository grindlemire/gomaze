package encode

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/grindlemire/gomaze/pkg/board"
)

// PNG is a png encoder
type PNG struct {
	File io.WriteCloser
}

// NewPNG creates a new PNG encoder
func NewPNG(name string) (p PNG, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return p, err
	}

	p = PNG{
		File: f,
	}

	return p, nil
}

const cellWidth = 20
const horizontalPadding = cellWidth / 2
const verticalPadding = cellWidth / 2

// Encode encodes cells to a png
func (p PNG) Encode(entrance, exit *board.Cell, cells [][]*board.Cell, path []*board.Cell) (err error) {
	width := len(cells[0])
	height := len(cells)
	img := image.NewNRGBA(image.Rect(0, 0, cellWidth*width+horizontalPadding, cellWidth*height+verticalPadding))

	pathCache := map[string]struct{}{}
	for _, cell := range path {
		pathCache[cell.ID] = struct{}{}
	}

	// loop over cells
	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[0]); x++ {
			cell := cells[y][x]

			fillColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			_, inPath := pathCache[cell.ID]
			if inPath {
				fillColor = color.NRGBA{R: 0, G: 0, B: 255, A: 255}
			}

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
					pX := x*cellWidth + i + horizontalPadding/2
					pY := y*cellWidth + j + verticalPadding/2

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

					img.Set(pX, pY, cellFill)
				}
			}
		}
	}

	return png.Encode(p.File, img)
}
