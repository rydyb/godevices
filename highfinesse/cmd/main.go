package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Measure Measure `cmd:"" help:"Measures a quantity."`
	Version Version `cmd:"" help:"Prints the version of the wavemeter."`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
