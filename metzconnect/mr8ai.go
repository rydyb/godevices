package metzconnect

import (
	"fmt"

	"github.com/goburrow/modbus"
)

type MR8AI struct {
	modbus.Client
}

func (d *MR8AI) AnalogInput(ext, port int) (float64, error) {
	if ext < 1 || ext > 8 {
		return 0.0, fmt.Errorf("EWIO2 can only have up to eight MR-8AI extensions %d", ext)
	}
	if port < 1 || port > 8 {
		return 0.0, fmt.Errorf("MR-AI8 has only eight analog inputs not %d", port)
	}

	bytes, err := d.Client.ReadInputRegisters(uint16(ext*100+0x40+port), 2)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read from modbus: %s", err)
	}

	return readFloat16(bytes)
}
