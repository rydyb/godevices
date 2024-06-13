package metzconnect

import (
	"fmt"

	"github.com/goburrow/modbus"
)

type EWIO2 struct {
	analogInput
	client modbus.Client
}

func NewEWIO2(client modbus.Client) *EWIO2 {
	return &EWIO2{
		client: client,
		analogInput: analogInput{
			client:       client,
			LowerChannel: 1,
			UpperChannel: 3,
			modeOffset:   60,
			valueOffset:  70,
		},
	}
}

func (d *EWIO2) Version() (int, error) {
	response, err := d.client.ReadHoldingRegisters(2, 1)
	if err != nil {
		return -1, fmt.Errorf("failed to read version via modbus: %s", err)
	}
	return int(response[1]), nil
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

var extensionTypeNames = map[ExtensionType]string{
	ExtensionUnknown:  "Unknown",
	ExtensionMR_DO4:   "DO4",
	ExtensionMR_TO4:   "TO4",
	ExtensionMR_DI4:   "DI4",
	ExtensionMR_DI10:  "DI10",
	ExtensionMR_SI4:   "SI4",
	ExtensionMR_DIO4_2: "DIO4_2",
	ExtensionMR_AO4:   "AO4",
	ExtensionMR_AOP4:  "AOP4",
	ExtensionMR_AI8:   "AI8",
	ExtensionMR_CI4:   "CI4",
}	

func (t ExtensionType) String() string {
	if name, ok := extensionTypeNames[t]; ok {
		return name
	}
	return "Unknown"
}

type Extension interface {
	ID() int
	Type() ExtensionType
}

func (d *EWIO2) Extensions() ([]Extension, error) {
	extensions := []Extension{}

	for i := 1; i < 8; i++ {
		response, err := d.client.ReadHoldingRegisters(uint16(i*100), 1)
		if err != nil {
			return nil, fmt.Errorf("failed to read extension type via modbus: %s", err)
		}

		extension := ExtensionType(response[1])
		switch extension {
		case ExtensionUnknown:
			continue
		case ExtensionMR_AI8:
			extensions = append(extensions, NewMR8AI(d.client, i))
		default:
			return nil, fmt.Errorf("unsupported extension type: %d", extension)
		}
	}

	return extensions, nil
}

type MRAI8 struct {
	analogInput
	id int
}

func NewMR8AI(client modbus.Client, id int) *MRAI8 {
	return &MRAI8{
		id: id,
		analogInput: analogInput{
			client:       client,
			LowerChannel: 1,
			UpperChannel: 8,
			modeOffset:   uint16(60 + 100*id),
			valueOffset:  uint16(70 + 100*id),
		},
	}
}

func (d *MRAI8) ID() int {
	return d.id
}

func (d *MRAI8) Type() ExtensionType {
	return ExtensionMR_AI8
}

type MRCI4 struct {
	analogInput
	id int
}

func NewMRCI4(client modbus.Client, id int) *MRCI4 {
	return &MRCI4{
		id: id,
		analogInput: analogInput{
			client:       client,
			LowerChannel: 1,
			UpperChannel: 4,
			modeOffset:   uint16(60 + 100*id),
			valueOffset:  uint16(70 + 100*id),
		},
	}
}

func (d *MRCI4) ID() int {
	return d.id
}

func (d *MRCI4) Type() ExtensionType {
	return ExtensionMR_CI4
}
