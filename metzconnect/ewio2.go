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

func (d *EWIO2) AnalogInput(port int) (float64, error) {
	if port < 1 || port > 3 {
		return 0.0, fmt.Errorf("EWIO2 has only analog inputs 1, 2, 3 not %d", port)
	}

	bytes, err := d.Client.ReadInputRegisters(0x40+uint16(port), 2)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read from modbus: %s", err)
	}

	return readFloat16(bytes)
}

func readFloat16(b []byte) (float64, error) {
	buffer := bytes.NewBuffer(b)

	var value float32
	if err := binary.Read(buffer, binary.BigEndian, &value); err != nil {
		return 0.0, err
	}

	return float64(value), nil
}
