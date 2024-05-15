/*
Send command to controller according to Leybold Combivac CM52 protocol.

Read pressure value:
Wx: RPV[channel][CR]
channel: 1 (TM1), 2(TM2) or 3 (IONIVAC)
Rx: s[,][TAB]x.xxxxE+-xx
s: status byte

Gas correction factor:
Wx: RGC[channel][CR]
channel = 3 (IONIVAC)
Rx: gf[CR]
gf: gas correction factor
0.20 < gf < 8.00, format: X.XX
*/

package leybold

import (
	"fmt"
	"io"
	"bufio"
	"strconv"
	"strings"
)

type Cmd string

const (
	Pressure  Cmd = "RPV"
	GasCorrection Cmd = "RGC"
)

type CombivacCM52 struct {
	rw      io.ReadWriter
	channel uint8
}

type Status string

const (
	Ok Status = "0"
	TooSmall Status = "1"
	TooLarge Status = "2"
	ErrLow Status = "3"
	ErrHigh Status = "4"
	SOff Status = "5"
	HVOn Status = "6"
	SensorErr Status = "7"
	NoSensor Status = "9"
	NotriG Status = "10"
	ErrPir Status = "11"
	OKDegas Status = "12"
	Er Status = "?"
)

func NewCombivacCM52(rw io.ReadWriter, channel uint8) *CombivacCM52 {
	return &CombivacCM52{rw: rw, channel: channel}
}


func (cm52 *CombivacCM52) Read(c Cmd) (string, error) {
	cmd := fmt.Sprintf("%s %d\r", c, cm52.channel)
	_, err := fmt.Fprint(cm52.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("failed writing command to cm52: %s", err)
	}
	
	scanner := bufio.NewScanner(cm52.rw)
	scanner.Split(split)

	if !scanner.Scan() {
		return "", fmt.Errorf("failed reading to end of command response: %s", scanner.Err())
	}

	response := scanner.Text()

	status := Status(response[0:1])

	if status == Er {
		return "", fmt.Errorf("failed to execute command successful (response code: %s)", status)
	}

	return response, nil
}

// RS232 implementation, address does not need to be specified as oppossed to RS485
func (cm52 *CombivacCM52) ReadFloat(c Cmd) (float64, error) {
	response, err := cm52.Read(c)
	if err != nil {
		return 0, err
	}

	switch c {
	case Pressure:
		status := Status(strings.Split(response, ",")[0])
		data := strings.Split(response, ",")[1]

		if status != Ok {
			return 0, fmt.Errorf("cm52 responded with error code %v", string(status))
		}
		cleanedData := strings.ReplaceAll(data, ",", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\t", "")
		cleanedData = strings.ReplaceAll(cleanedData, " ", "")

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
