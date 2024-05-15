package spce

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Code uint8

const (
	Current  Code = 0x0a
	Pressure Code = 0x0b
	Voltage  Code = 0x0c
)

type Status string

const (
	Ok Status = "OK"
	Er Status = "ER"
)

type SPCE struct {
	rw      io.ReadWriter
	channel uint8
}

// NewSPCE returns a new SPCE instance on a given ReadWriter for the specified channel.
func NewSPCE(rw io.ReadWriter, channel uint8) *SPCE {
	return &SPCE{rw: rw, channel: channel}
}

// Read writes a read command for code c to the SPCE and returns the response data.
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

func checksum(cmd string) string {
	var sum int
	for _, c := range cmd {
		sum += int(c)
	}
	return fmt.Sprintf("%02X", sum%256)
}

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
