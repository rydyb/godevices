package egnite

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/goburrow/modbus"
)

type Quantity uint16

const (
	Temperature Quantity = 10
	Humidity    Quantity = 11
	Pressure    Quantity = 18
)

type Queryx struct {
	modbus.Client
}

func (d *Queryx) ReadFloat(q Quantity) (float64, error) {
	b, err := d.Client.ReadInputRegisters(uint16(q), 1)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read modbus input registers: %s", err)
	}

	value, err := readFloat8(b)
	if err != nil {
		return 0.0, err
	}
	return value / 10, nil
}

func readFloat8(b []byte) (float64, error) {
	buffer := bytes.NewBuffer(b)
	var value int16
	if err := binary.Read(buffer, binary.BigEndian, &value); err != nil {
		return 0.0, fmt.Errorf("failed to read registers as float32: %s", err)
	}

	return float64(value), nil
}
