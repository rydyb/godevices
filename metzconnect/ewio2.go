package metzconnect

import (
	"fmt"

	"github.com/goburrow/modbus"
)

type EWIO2 struct {
	modbus.Client
}

func NewEWIO2(client modbus.Client) *EWIO2 {
	return &EWIO2{Client: client}
}

func (d *EWIO2) Version() (int, error) {
	response, err := d.Client.ReadHoldingRegisters(2, 1)
	if err != nil {
		return -1, fmt.Errorf("failed to read version via modbus: %s", err)
	}
	return int(response[1]), nil
}

func (d *EWIO2) AnalogInputs() *AnalogInputs {
	return &AnalogInputs{
		Client:       d.Client,
		LowerChannel: 1,
		UpperChannel: 3,
		modeOffset:   60,
		valueOffset:  70,
	}
}

type ExtensionType uint8

const (
	ExtensionUnknown ExtensionType = iota
	ExtensionMR_DO4
	ExtensionMR_TO4
	ExtensionMR_DI4
	ExtensionMR_DI10
	ExtensionMR_SI4
	ExtensionMR_DIO4_2
	ExtensionMR_AO4
	ExtensionMR_AOP4
	ExtensionMR_AI8
	ExtensionMR_CI4
)

type Extension interface {
	ID() uint16
	Type() ExtensionType
}

func (d *EWIO2) Extensions() ([]Extension, error) {
	extensions := []Extension{}

	for i := 1; i < 8; i++ {
		response, err := d.Client.ReadHoldingRegisters(uint16(i*100), 1)
		if err != nil {
			return nil, fmt.Errorf("failed to read extension type via modbus: %s", err)
		}

		extension := ExtensionType(response[1])
		switch extension {
		case ExtensionUnknown:
			continue
		case ExtensionMR_AI8:
			extensions = append(extensions, NewMR8AI(d.Client, uint16(i)))
		default:
			return nil, fmt.Errorf("unsupported extension type: %d", extension)
		}
	}

	return extensions, nil
}

type MR8AI struct {
	id uint16
	modbus.Client
}

func NewMR8AI(client modbus.Client, id uint16) *MR8AI {
	return &MR8AI{id: id, Client: client}
}

func (d *MR8AI) ID() uint16 {
	return d.id
}

func (d *MR8AI) Type() ExtensionType {
	return ExtensionMR_AI8
}

func (d *MR8AI) AnalogInputs() *AnalogInputs {
	return &AnalogInputs{
		Client:       d.Client,
		LowerChannel: 1,
		UpperChannel: 8,
		modeOffset:   160,
		valueOffset:  140,
	}
}
