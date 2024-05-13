package spce

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type SPCE struct {
	rw io.ReadWriter
	channel string
}

// NewSPCE creates a new SPCE instance with the given ReadWriter and channel address in the format of two hexadezimal characters
// The default channel address is "05"
func NewSPCE(rw io.ReadWriter, channel string) *SPCE {
	// If the channel address is not provided, use the default address "05"
	if channel == "" {
		channel = "05"
	}
	return &SPCE{rw: rw, channel: channel}
}

func (s *SPCE) Outputs() (map[string]float64, error) {
	current_enc := "0A"
	pressure_enc := "0B"
	voltage_enc := "0C"

	current, err := s.Get(s.channel, current_enc)
	if err != nil {
		return nil, fmt.Errorf("failed to get current value: %s", err)
	}
	pressure, err := s.Get(s.channel, pressure_enc)
	if err != nil {
		return nil, fmt.Errorf("failed to get pressure value: %s", err)
	}
	voltage, err := s.Get(s.channel, voltage_enc)
	if err != nil {
		return nil, fmt.Errorf("failed to get voltage value: %s", err)
	}

	log.Printf("Outputs: current=%f, pressure=%f, voltage=%f", current, pressure, voltage)

	return map[string]float64{
		"current":  current,
		"pressure": pressure,
		"voltage":  voltage,
	}, nil
}

func (s *SPCE) Get(channel string, metric_enc string) (float64, error) {
	// Send command to controller according to SPCe protocol
	// ~ <channel> <metric> <checksum>
	// <channel> is the channel address
	// <metric> is the metric enconding
	// <checksum> is the checksum of the command excluding the ~, but including all the spaces
	cmd := " " + channel + " " + metric_enc + " "
	cmd = "~" + cmd + s.CalculateChecksum(cmd)
	value_string, err := s.Exec(cmd)
	if err != nil {
		return 0, fmt.Errorf("failed to get execute command and read value: %s", err)
	}

	// Remove units if they exist, the response of voltage does not have a unit
	if strings.HasSuffix(value_string, "AMPS") {
		value_string = strings.TrimSuffix(value_string, "AMPS")
	} else if strings.HasSuffix(value_string, "MBAR") {
		value_string = strings.TrimSuffix(value_string, "MBAR")
	} else if strings.HasSuffix(value_string, "Torr") {
		value_string = strings.TrimSuffix(value_string, "Torr")
		log.Printf("Warning: value is in Torr, should be in mbar")
	} else if strings.HasSuffix(value_string, "PA") {
		value_string = strings.TrimSuffix(value_string, "PA")
		log.Printf("Warning: value is in Pa, should be in mbar")
	}
	value_string = strings.TrimSpace(value_string)

	value, err := strconv.ParseFloat(value_string, 64)
	if err != nil {
		// Try parsing as an integer if ParseFloat fails
		valueInt, errInt := strconv.Atoi(value_string)
		if errInt != nil {
			return 0, fmt.Errorf("failed to parse value: %s", err)
		}
		value = float64(valueInt)
	}

	return value, nil
}

func (s *SPCE) CalculateChecksum(cmd string) string {
	var sum int
	for _, c := range cmd {
		sum += int(c)
	}
	// return the checksum as a two-digit hexadecimal number
	return fmt.Sprintf("%02X", sum%256)
}

func (s *SPCE) Exec(cmd string) (string, error) {
	_, err := fmt.Fprintf(s.rw, cmd + "\r")
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(s.rw)
	scanner.Split(bufio.ScanBytes)

	var respBuild strings.Builder
	for scanner.Scan() {
		b := scanner.Bytes()
		// The response ends with a carriage return
		if b[0] == 13 {
			break
		}
		respBuild.Write(b)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	resp := respBuild.String()

	// The response is in the format ~ <address> OK <response_code> <value> <checksum>
	// address := resp[:2]
	status := resp[3:5]
	responseCode := resp[6:8]
	value_string := resp[9 : len(resp)-3]
	checksum := resp[len(resp)-2 :]

	if status != "OK" {
		return "", fmt.Errorf("failed to execute command with response code %s", responseCode)
	}
	if checksum != s.CalculateChecksum(resp[:len(resp)-2]) {
		return "", fmt.Errorf("invalid checksum: %s, should have been %s", checksum, s.CalculateChecksum(resp[:len(resp)-2]))
	}

	return value_string, nil
}
