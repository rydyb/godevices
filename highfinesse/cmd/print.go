package main

import (
	"fmt"

	"github.com/rydyb/godevices/highfinesse"
)

type Print struct {
	Quantity string `arg:"" type:"string" enum:"wavelength,frequency,pressure,temperature" help:"Measure wavelength or frequency."`
}

func (cmd *Print) Run() error {
	channels := highfinesse.Channels()

	switch cmd.Quantity {
	case "wavelength":
		for channel := uint32(1); channel <= uint32(channels); channel++ {
			λ, err := highfinesse.Wavelength(channel)
			if err != nil {
				fmt.Printf("channel %d: %s\n", channel, err)
			} else {
				fmt.Printf("channel %d: %f nm\n", channel, λ)
			}
		}
	case "frequency":
		for channel := uint32(1); channel <= uint32(channels); channel++ {
			f, err := highfinesse.Frequency(channel)
			if err != nil {
				fmt.Printf("channel %d: %s\n", channel, err)
			} else {
				fmt.Printf("channel %d: %f THz\n", channel, f)
			}
		}
	case "pressure":
		p, err := highfinesse.Pressure()
		if err != nil {
			return fmt.Errorf("failed measuring pressure: %s", err)
		}
		fmt.Printf("Pressure: %f mbar\n", p)
	case "temperature":
		T, err := highfinesse.Temperature()
		if err != nil {
			return fmt.Errorf("failed measuring temperature: %s", err)
		}
		fmt.Printf("Temperature: %f °C\n", T)
	default:
		return fmt.Errorf("unknown quantity: %s", cmd.Quantity)
	}

	return nil
}
