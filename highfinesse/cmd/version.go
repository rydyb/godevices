package main

import (
	"fmt"

	"github.com/rydyb/godevices/highfinesse"
)

type Version struct{}

func (cmd *Version) Run() error {
	fmt.Println(highfinesse.WavelengthMeterVersionInfo())
	return nil
}
