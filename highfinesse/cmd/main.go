package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Print   Print   `cmd:"" help:"Prints a measured quantity."`
	Serve   Serve   `cmd:"" help:"Starts a http server with REST API."`
	Version Version `cmd:"" help:"Prints the version of the wavemeter."`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
