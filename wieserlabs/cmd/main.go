package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/rydyb/godevices/wieserlabs"
)

type Context struct {
	FlexDDSSlot *wieserlabs.FlexDDSSlot
}

var cli struct {
	Host       string     `name:"host" required:"" help:"The hostname or ip address of the FlexDDS controller."`
	Slot       uint8      `name:"slot" required:"" help:"The number of the FlexDDS slot from 0 to 5."`
	SysClock   float64    `name:"system-clock" default:"1e9" help:"The system' clocks frequency in Hz."`
	Singletone Singletone `cmd:"singletone" help:"Configure a singletone output."`
}

func main() {
	ctx := kong.Parse(&cli)

	if cli.Slot > 5 {
		log.Fatalf("Slot number cannot be greater than five.")
	}

	flexddsslot, err := wieserlabs.NewFlexDDSSlot(cli.Host, cli.Slot, cli.SysClock)
	if err != nil {
		log.Fatalf("failed to open connection to %s: %s", cli.Host, err)
	}
	defer flexddsslot.Close()

	if err := ctx.Run(&Context{FlexDDSSlot: flexddsslot}); err != nil {
		log.Fatalf("failed to run command: %s", err)
	}
}
