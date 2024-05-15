package gamma

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Code represents the metric code.
type Code uint8

const (
	Current  Code = 0x0a
	Pressure Code = 0x0b
	Voltage  Code = 0x0c
)

// Status represents the status of the command response.
type Status string

const (
	Ok Status = "OK"
	Er Status = "ER"
)

// SPCE represents a SPCE controller.
type SPCE struct {
	rw      io.ReadWriter
	channel uint8
}

// NewSPCE creates a new SPCE controller with the given read writer and channel.
func NewSPCE(rw io.ReadWriter, channel uint8) *SPCE {
	return &SPCE{rw: rw, channel: channel}
}

// Read reads a code from the SPCE controller and returns the data as a string.
//
// The read command is expected to follow this format:
//
//	~ <channel> <metric> <checksum>
//
// - <channel> refers to the channel address.
// - <metric> refers to the metric encoding.
// - <checksum> is the checksum of the command, excluding the '~' character but including all spaces.
//
// The response from the SPCE controller is expected to be in the following format:
//
//	~ <address> <status> <response_code> (<value>) <checksum>
//
// - <address> is the address returned by the controller.
// - <status> indicates the status of the response.
// - <response_code> is the response code from the controller.
// - <value> is the data returned by the controller, enclosed in parentheses.
// - <checksum> is the checksum of the response, excluding the '~' character but including all spaces.
func (s *SPCE) Read(c Code) (string, error) {
	cmd := fmt.Sprintf("~ %02d %02X ", s.channel, c)
	cmd = cmd + checksum(cmd[1:]) + "\r"

	_, err := fmt.Fprint(s.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("failed writing command to read writer: %s", err)
	}

	scanner := bufio.NewScanner(s.rw)
	scanner.Split(split)

	if !scanner.Scan() {
		return "", fmt.Errorf("failed reading to end of command response: %s", scanner.Err())
	}

	response := scanner.Text()

	status := response[3:5]
	code := response[6:8]
	data := response[9 : len(response)-3]
	csum := response[len(response)-2:]

	if status != string(Ok) {
		return "", fmt.Errorf("failed to execute command successful (response code: %s)", code)
	}
	if csum != checksum(response[:len(response)-2]) {
		return "", fmt.Errorf("failed checksum of response")
	}

	return data, nil
}

// ReadFloat reads a code from the SPCE controller and returns the data as a float64.
//
// For the Voltage code, the float64 value is the voltage in volts.
// For the Current code, the float64 value is the current in amperes.
// For the Pressure code, the float64 value is the pressure in mbars.
func (s *SPCE) ReadFloat(c Code) (float64, error) {
	data, err := s.Read(c)
	if err != nil {
		return 0, err
	}

	switch c {
	case Voltage:
		value, err := strconv.ParseInt(data, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("failed to parse voltage as int: %s", err)
		}
		return float64(value), nil
	case Current:
		fallthrough
	case Pressure:
		parts := strings.Split(data, " ")
		if len(parts) != 2 {
			return 0, fmt.Errorf("failed splitting response string: %s", data)
		}
		value, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse data as float: %s", err)
		}
		return value, nil
	}

	return 0, fmt.Errorf("unknown code: %d", c)
}

// checksum calculates the checksum of the command.
func checksum(cmd string) string {
	var sum int
	for _, c := range cmd {
		sum += int(c)
	}
	return fmt.Sprintf("%02X", sum%256)
}

// split is a bufio.SplitFunc that splits the data at '\r'.
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
