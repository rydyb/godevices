// Copyright 2024

// leybold is a package that provides support for Leybold Combivac CM52 controllers.
//
// Usage
//
// You can either read a pressure value or a gas correction factor from a Combivac CM52 controller.
// To read a pressure value from a Combivac CM52 controller, you need to create a new CombivacCM52 and call the ReadFloat method with the Pressure command.
// Use TCP to connect to the controller and pass the connection to the CombivacCM52 constructor.
// The channel parameter specifies the channel address of the controller (1, 2 or 3).
package leybold

import (
	"fmt"
	"io"
	"bufio"
	"strconv"
	"strings"
)

type CMD string

const (
	Pressure  		CMD = "RPV"
	GasCorrection 	CMD = "RGC"
)

type CombivacCM52 struct {
	rw      io.ReadWriter
	channel uint8
}

type Status string

const (
	Ok 			Status = "0"
	TooSmall 	Status = "1"
	TooLarge 	Status = "2"
	ErrLow 		Status = "3"
	ErrHigh 	Status = "4"
	SOff 		Status = "5"
	HVOn 		Status = "6"
	SensorErr	Status = "7"
	NoSensor 	Status = "9"
	NotriG 		Status = "10"
	ErrPir 		Status = "11"
	OKDegas 	Status = "12"
	Er 			Status = "?"
)

func NewCombivacCM52(rw io.ReadWriter, channel uint8) *CombivacCM52 {
	return &CombivacCM52{rw: rw, channel: channel}
}

// Read sends one of two possible commands to a cm52 controller and returns the response.
// Sending and recieving happens according to the Leybold Combivac CM52 protocol.
//
// Read pressure value:
// Wx: RPV[channel][CR]
// channel: 1 (TM1), 2(TM2) or 3 (IONIVAC)
// Rx: s[,][TAB]x.xxxxE+-xx
// s: status byte
//
// Gas correction factor:
// Wx: RGC[channel][CR]
// channel = 3 (IONIVAC)
// Rx: gf[CR]
// gf: gas correction factor
// 0.20 < gf < 8.00, format: X.XX
func (cm52 *CombivacCM52) Read(c CMD) (string, error) {
	cmd := fmt.Sprintf("%s %d\r", c, cm52.channel)
	_, err := fmt.Fprint(cm52.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to write command to cm52: %s", err)
	}
	
	scanner := bufio.NewScanner(cm52.rw)
	scanner.Split(split)

	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read to end of command response: %s", scanner.Err())
	}

	response := scanner.Text()

	status := Status(response[0:1])

	if status == Er {
		return "", fmt.Errorf("failed to execute command (response code: %s)", status)
	}

	return response, nil
}

// ReadFloat calls the Read method with the given command and
// passes the parsed response from the Read method to the caller as a float64.
// RS232 implementation, the address does not need to be specified as oppossed to RS485.
func (cm52 *CombivacCM52) ReadFloat(c CMD) (float64, error) {
	response, err := cm52.Read(c)
	if err != nil {
		return 0, err
	}

	switch c {
	case Pressure:
		splitResponse := strings.Split(response, ",")
		if len(splitResponse) < 2 {
			return 0, fmt.Errorf("invalid response format: %s", response)
		}
		status := Status(splitResponse[0])
		data := splitResponse[1]

		if status != Ok {
			return 0, fmt.Errorf("cm52 responded with error code %v", string(status))
		}
		// Remove all whitespace and commas from the data
		cleanedData := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(data, ",", ""), "\t", ""), " ", "")

		value, err := strconv.ParseFloat(cleanedData, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse pressure value: %s", err)
		}
		return value, nil
	case GasCorrection:
		value, err := strconv.ParseFloat(response, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse gas correction factor: %s", err)
		}

		if value < 0.20 || value > 8.00 {
			return 0, fmt.Errorf("gas correction factor out of range: %f", value)
		}
		return value, nil
	}

	return 0, fmt.Errorf("unknown command: %s", c)
}

// split is a bufio.SplitFunc that splits on '\r'.
func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}
