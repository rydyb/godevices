package metzconnect

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/goburrow/modbus"
)

type EWIO2 struct {
	modbus.Client
}

// AnalogInput returns the voltage of the analog input channel.
func (d *EWIO2) AnalogInput(channel int) (float64, error) {
	address, err := AnalogInputAddress(channel)
	if err != nil {
		return 0.0, fmt.Errorf("failed to convert channel %d to modbus address: %s", channel, err)
	}

	registers, err := d.Client.ReadInputRegisters(address, 2)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read modbus input registers: %s", err)
	}

	return registersToFloat(registers)
}

// AnalogInputAddress returns the modbus address of the analog input channel.
func AnalogInputAddress(channel int) (uint16, error) {
	if channel > 0 && channel < 4 {
		return uint16(40 + 2*(channel-1)), nil
	}

	if channel > 3 && channel <= 3+6*8 {
		return uint16(100*((channel-4)/8+1) + 40 + 2*((channel-4)%8)), nil
	}

	return 0, fmt.Errorf("channel %d out of range", channel)
}

func registersToFloat(b []byte) (float64, error) {
	buffer := bytes.NewBuffer(b)

	var value float32
	if err := binary.Read(buffer, binary.BigEndian, &value); err != nil {
		return 0.0, err
	}

	return float64(value), nil
}
