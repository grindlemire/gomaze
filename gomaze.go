package main

import (
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/grindlemire/gomaze/pkg/board"
	"github.com/grindlemire/gomaze/pkg/encode"
	"github.com/grindlemire/log"
)

// Opts ...
type Opts struct {
	Height int `          long:"height" default:"10"      description:"height of the maze"`
	Width  int `          long:"width"  default:"20"      description:"width of the maze"`
	File   int `short:"f" long:"file"   default:"out.png" description:"the name of the file to write out"`
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
	b := board.New(opts.Width, opts.Height)

	p, err := encode.NewPNG("out.png")
	if err != nil {
		log.Fatal(err)
	}

	err = p.Encode(b.Cells)
	if err != nil {
		log.Fatal(err)
	}
}
