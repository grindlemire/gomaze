package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/nfnt/resize"

	"github.com/grindlemire/gomaze/pkg/board"
	"github.com/grindlemire/gomaze/pkg/encode"
	"github.com/grindlemire/gomaze/pkg/traverse/dfs"
	"github.com/grindlemire/log"
)

// Opts ...
type Opts struct {
	Height    int    `          long:"height" default:"10"        description:"height of the maze"`
	Width     int    `          long:"width"  default:"20"        description:"width of the maze"`
	InputFile string `short:"i" long:"input"  default:"input.png" description:"the file containing the input biases"`
	File      string `short:"f" long:"file"   default:"out"   description:"the name of the file to write out"`
}

var opts Opts
var parser = flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)

func main() {
	log.Init(log.Default)

	_, err := parser.Parse()
	if flags.WroteHelp(err) {
		parser.WriteHelp(os.Stderr) // This writes the help when we want help. This is silenced because we are not writing any errors
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Width: %d | Height: %d", opts.Width, opts.Height)

	t := time.Now()
	inputImg, err := resizeInput(opts.InputFile, opts.Width*encode.CellWidth, opts.Height*encode.CellWidth)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Time to resize input bias: %v", time.Since(t))

	b := board.New(opts.Width, opts.Height, encode.CellWidth, inputImg)

	t = time.Now()
	b = dfs.CreateEdges(b)
	log.Infof("Time to generate maze: %v", time.Since(t))

	t = time.Now()
	path := dfs.Solve(b)
	log.Infof("Time to solve maze: %v", time.Since(t))

	p, err := encode.NewPNG(b)
	if err != nil {
		log.Fatal(err)
	}

	t = time.Now()
	err = p.EncodeBoard()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Time to encode maze: %v", time.Since(t))

	err = p.Save(fmt.Sprintf("%s.png", opts.File))
	if err != nil {
		log.Fatal(err)
	}

	t = time.Now()
	err = p.EncodeSolution(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Time to encode solution: %v", time.Since(t))

	err = p.Save(fmt.Sprintf("%s_solved.png", opts.File))
	if err != nil {
		log.Fatal(err)
	}
}

func resizeInput(filename string, mazeWidth, mazeHeight int) (img image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return img, err
	}
	defer f.Close()

	rawImg, err := png.Decode(f)
	if err != nil {
		return img, err
	}

	img = resize.Resize(uint(mazeWidth)+100, uint(mazeHeight)+100, rawImg, resize.Lanczos3)

	outf, err := os.OpenFile("resized_input.png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer outf.Close()

	png.Encode(outf, img)
	return img, nil
}
