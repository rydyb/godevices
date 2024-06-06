package metzconnect

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/goburrow/modbus"
)

type AnalogInputUnit string

const (
	UnitNone        AnalogInputUnit = ""
	UnitCelsius     AnalogInputUnit = "°C"
	UnitPercent     AnalogInputUnit = "%"
	UnitVolt        AnalogInputUnit = "V"
	UnitMilliampere AnalogInputUnit = "mA"
	UnitOhm         AnalogInputUnit = "Ω"
)

func (unit AnalogInputUnit) String() string {
	return string(unit)
}

type AnalogInputMode uint8

const (
	Voltage0to10Percent AnalogInputMode = iota
	Voltage0to5PullupPercent
	Voltage0to10Volts
	Voltage0to5PullupVolts
	ResistanceOhms
	UserDefinedRange
	TemperaturePT100Celsius
	TemperaturePT500Celsius
	TemperaturePT1000Celsius
	TemperatureNI1000TC5000Celsius
	TemperatureNI1000TC6180Celsius
	TemperatureBALCO500Celsius
	TemperatureKTY81_110Celsius
	TemperatureKTY81_210Celsius
	TemperatureNTC1k8ThermokonCelsius
	TemperatureNTC5kThermokonCelsius
	TemperatureNTC10kThermokonCelsius
	TemperatureNTC20kThermokonCelsius
	TemperatureLM235ZCelsius
	TemperatureNTC10kCarelCelsius
	TemperatureNTC5kSchneiderCelsius
	TemperatureNTC30kSchneiderCelsius
	TemperatureKP250Celsius
	ResistancePoti10kPercent
	Inactive
	Current0to20MilliamperePercent
	Current0to20Milliampere
	Current4to20MilliamperePercent
	Current4to20Milliampere
	ResistanceThreeWireSensingOhms
	ResistanceFourWireSensingOhms
	Test40OhmTo14kOhm
	Test12kOhmTo4MOhm
	Test40OhmTo650Ohm
	Test500OhmTo14kOhm
	Test12kOhmTo180kOhm
	Test140kOhmTo4MOhm
)

type ChannelOutOfRangeError struct {
	Channel    uint8
	LowerLimit uint8
	UpperLimit uint8
}

func (err ChannelOutOfRangeError) Error() string {
	return fmt.Sprintf("channel %d out of range [%d, %d]", err.Channel, err.LowerLimit, err.UpperLimit)
}

type analogInput struct {
	client       modbus.Client
	LowerChannel uint8
	UpperChannel uint8
	modeOffset   uint16
	valueOffset  uint16
}

func (ai *analogInput) Mode(channel uint8) (AnalogInputMode, error) {
	if channel < ai.LowerChannel || channel > ai.UpperChannel {
		return Inactive, ChannelOutOfRangeError{channel, ai.LowerChannel, ai.UpperChannel}
	}

	response, err := ai.client.ReadInputRegisters(ai.modeOffset+uint16(channel-1), 1)
	if err != nil {
		return Inactive, fmt.Errorf("failed to read mode via modbus: %s", err)
	}
	return AnalogInputMode(response[1]), nil
}

func (ai *analogInput) Unit(channel uint8) (AnalogInputUnit, error) {
	mode, err := ai.Mode(channel)
	if err != nil {
		return UnitNone, err
	}

	switch mode {
	case Voltage0to10Percent, Voltage0to5PullupPercent, Current0to20MilliamperePercent, Current4to20MilliamperePercent:
		return UnitPercent, nil
	case Voltage0to10Volts, Voltage0to5PullupVolts:
		return UnitVolt, nil
	case ResistanceOhms, ResistancePoti10kPercent, ResistanceThreeWireSensingOhms, ResistanceFourWireSensingOhms:
		return UnitOhm, nil
	case TemperatureBALCO500Celsius, TemperatureKTY81_110Celsius, TemperatureKTY81_210Celsius, TemperatureKP250Celsius:
		return UnitCelsius, nil
	case UserDefinedRange:
		return UnitNone, nil
	}

	return UnitNone, fmt.Errorf("unsupported mode: %d", mode)
}

func (ai *analogInput) Value(channel uint8) (float32, error) {
	if channel < ai.LowerChannel || channel > ai.UpperChannel {
		return 0.0, ChannelOutOfRangeError{channel, ai.LowerChannel, ai.UpperChannel}
	}

	response, err := ai.client.ReadInputRegisters(ai.valueOffset+uint16(2*(channel-1)), 2)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read value via modbus: %s", err)
	}

	value, err := readFloat16(response)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse value as float16: %s", err)
	}

	unit, err := ai.Unit(channel)
	if err != nil {
		return 0.0, err
	}

	switch unit {
	case UnitPercent:
		return value * 100.0, nil
	case UnitVolt:
		return value, nil
	case UnitOhm:
		return value, nil
	case UnitCelsius:
		return ai.convertValueToCelsius(value, channel)
	}

	return 0.0, fmt.Errorf("unsupported unit: %s", unit)
}

func (ai *analogInput) convertValueToCelsius(value float32, channel uint8) (float32, error) {
	mode, err := ai.Mode(channel)
	if err != nil {
		return 0.0, err
	}
	switch mode {
	default:
		return value, nil
	}
}

func readFloat16(b []byte) (float32, error) {
	buffer := bytes.NewBuffer(b)

	var value float32
	if err := binary.Read(buffer, binary.BigEndian, &value); err != nil {
		return 0.0, err
	}

	return value, nil
}
