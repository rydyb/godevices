package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/rydyb/godevices/visa"
	"golang.org/x/exp/maps"
)

type Oscilloscope struct {
	Address      string   `name:"address" required:"" help:"The address of the oscilloscope."`
	Identity     struct{} `cmd:"identity" help:"Query the identity of the oscilloscope."`
	Measurements struct{} `cmd:"measurements" help:"Query the measurements of the oscilloscope."`
}

func (cmd *Oscilloscope) Run(ctx *kong.Context) error {
	conn, err := net.Dial("tcp", cmd.Address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", cmd.Address, err)
	}
	defer conn.Close()

	osci := visa.NewOscilloscope(conn)

	switch ctx.Command() {
	case "identity":
		identity, err := osci.Identity()
		if err != nil {
			return fmt.Errorf("failed to query identity: %w", err)
		}
		fmt.Println(identity)
	case "measurements":
		out, err := osci.Measurements()
		if err != nil {
			return fmt.Errorf("failed to query measurements: %w", err)
		}
		fmt.Println(strings.Join(maps.Keys(out), ", "))
		fmt.Println(strings.Join(maps.Values(out), ", "))
	}

	return nil
}
