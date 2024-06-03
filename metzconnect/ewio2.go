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

	bytes, err := d.Client.ReadInputRegisters(0x46+uint16(port-1), 2)
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

func (d *EWIO2) GetExtensionTypes() (extensions []string, err error) {	
	extensions = []string{}
	for i := 1; i < 50; i++ {
		response, err := d.Client.ReadHoldingRegisters(uint16(i*100), 1)
		if err != nil {
			return []string{}, err
		}
		code := response[1]
		if code == 0 {
			fmt.Printf("No (further) extension module found or device not recognized\n")
			break
		}
		switch code {
		case 1:
			extensions = append(extensions, "MR-DO4")
		case 2:
			extensions = append(extensions, "MR-TO4")
		case 3:
			extensions = append(extensions, "MR-DI4")
		case 4:
			extensions = append(extensions, "MR-DI10")
		case 5:
			extensions = append(extensions, "MR-SI4")
		case 6:
			extensions = append(extensions, "MR-DIO4/2")
		case 7:
			extensions = append(extensions, "MR-AO4")
		case 8:
			extensions = append(extensions, "MR-AOP4")
		case 9:
			extensions = append(extensions, "MR-AI8")
		case 10:
			extensions = append(extensions, "MR-CI4")
		}
	}

	return extensions, nil
}


func (d *EWIO2) GetUnitRange() (unitRange []string, err error) {
	unitRangeLookup := map[int]string{
		1:  "0-10V %",
		2:  "0-5V % Pullup",
		3:  "0-10 Volt",
		4:  "0-5 Volt Pullup",
		5:  "Ohm",
		6:  "User Defined Range",
		7:  "PT100",
		8:  "PT500",
		9:  "PT1000",
		10: "NI1000-TC5000",
		11: "NI1000-TC6180",
		12: "BALCO500",
		13: "KTY81_110",
		14: "KTY81_210",
		15: "NTC1k8 Thermokon",
		16: "NTC5k Thermokon",
		17: "NTC10k Thermokon",
		18: "NTC20k Thermokon",
		19: "LM235Z",
		20: "NTC10k Carel",
		21: "NTC5k Schneider",
		22: "NTC30k Schneider",
		23: "KP250",
		24: "Poti 10k %",
		25: "Inactive",
		26: "0-20mA %",
		27: "0-20mA",
		28: "4-20mA %",
		29: "4-20mA",
		30: "3-wire sensing (Eingang E2)",
		31: "4-wire sensing (Eingang E2)",
		32: "Test 40 Ohm - 14 kOhm",
		33: "Test 12 kOhm - 4 MOhm",
		34: "Test 40 Ohm - 650 Ohm",
		35: "Test 500 Ohm - 14 kOhm",
		36: "Test 12 kOhm - 180 kOhm",
		37: "Test 140 kOhm - 4 MOhm",
	}

	response, err := d.Client.ReadHoldingRegisters(60, 3)
	if err != nil {
		return []string{}, err
	}

	unitRange = make([]string, 3)
	for i := 0; i < 3; i++ {
		code := int(response[2*i + 1])
		
		// Determine the unit and range of the analog input:
		// V (Voltage): 1-4	
		// mA (Milliampere): 26, 27, 28, 29
		// Ohm (Resistance): 5, 32-37
		// % (Percentage): 1, 2, 26, 28
		// °C (Celsius): 7-23
		if code > 0 && code < 5 {
			unitRange[i] = "V (Voltage): " + unitRangeLookup[code]
		} else if code == 26 || code == 27 || code == 28 || code == 29 {
			unitRange[i] = "mA (Milliampere): " + unitRangeLookup[code]
		} else if code == 5 || (code > 31 && code < 38) {
			unitRange[i] = "Ohm (Resistance): " + unitRangeLookup[code]
		} else if code == 1 || code == 2 || code == 26 || code == 28 {
			unitRange[i] = "% (Percentage): " + unitRangeLookup[code]
		} else if code > 6 && code < 24 {
			unitRange[i] = "°C (Celsius): " + unitRangeLookup[code]
		} else {
			unitRange[i] = "Unknown"
		}
	}
	
	return unitRange, nil
}