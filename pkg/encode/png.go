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

// Encode encodes cells to a png
func (p PNG) Encode(cells [][]*board.Cell) (err error) {
	img := image.NewNRGBA(image.Rect(0, 0, cellWidth*len(cells[0])+2, cellWidth*len(cells)+2))

	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[0]); x++ {
			cell := cells[y][x]

			for i := -cellWidth / 2; i <= cellWidth/2; i++ {
				for j := -cellWidth / 2; j <= cellWidth/2; j++ {
					// calculate the pixel we are looking at
					pX := x*cellWidth + i + 1
					pY := y*cellWidth + j + 1

					_, hasLeft := cell.Connections[board.Left]
					_, hasRight := cell.Connections[board.Right]
					_, hasUp := cell.Connections[board.Up]
					_, hasDown := cell.Connections[board.Down]
					var r, g, b, a uint8 = 255, 255, 255, 255

					// if we are on an edge and there is no connection, make it black
					if (i == -cellWidth/2 && !hasLeft) ||
						(i == cellWidth/2 && !hasRight) ||
						(j == -cellWidth/2 && !hasUp) ||
						(j == cellWidth/2 && !hasDown) ||
						(math.Abs(float64(i)) == cellWidth/2 && math.Abs(float64(j)) == cellWidth/2) { // This is the case you are in the corner
						r, g, b = 0, 0, 0
					}

					img.Set(pX, pY, color.NRGBA{
						R: r,
						G: g,
						B: b,
						A: a,
					})
				}
			}
		}
	}

	// img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	// img.Set(1, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255})
	// img.Set(0, 1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})

	return png.Encode(p.File, img)

}
