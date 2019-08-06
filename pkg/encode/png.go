package encode

import (
	"image"
	"image/color"
	"image/png"
	"io"
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

// Encode encodes cells to a png
func (p PNG) Encode(cells [][]*board.Cell) (err error) {
	img := image.NewNRGBA(image.Rect(0, 0, 3*len(cells[0])+2, 3*len(cells)+2))

	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[0]); x++ {
			cell := cells[y][x]

			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					// calculate the pixel we are looking at
					pX := x*3 + i + 1
					pY := y*3 + j + 1

					_, hasLeft := cell.Connections[board.Left]
					_, hasRight := cell.Connections[board.Right]
					_, hasUp := cell.Connections[board.Up]
					_, hasDown := cell.Connections[board.Down]

					var r, g, b, a uint8 = 0, 0, 0, 255
					// if we are not in the center or have a connection we are going to be black
					if (i == 0 && j == 0) ||
						(i < 0 && hasLeft && j == 0) ||
						(i > 0 && hasRight && j == 0) ||
						(i == 0 && hasUp && j < 0) ||
						(i == 0 && hasDown && j > 0) {
						r, g, b = 255, 255, 255
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
