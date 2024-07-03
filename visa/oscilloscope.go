package visa

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Oscilloscope struct {
	rw io.ReadWriter
}

func NewOscilloscope(rw io.ReadWriter) *Oscilloscope {
	return &Oscilloscope{rw: rw}
}

func (d *Oscilloscope) Identity() (string, error) {
	out, err := exec(d.rw, "*idn?")
	if err != nil {
		return "", fmt.Errorf("failed to query identity: %w", err)
	}
	return out, nil
}

func (d *Oscilloscope) MeasurementList() ([]string, error) {
	out, err := exec(d.rw, "MEASUrement:LIST?")
	if err != nil {
		return nil, fmt.Errorf("failed to query measurement list: %w", err)
	}
	return strings.Split(out, ","), nil
}

func (d *Oscilloscope) measurementValue(name string) (string, error) {
	out, err := exec(d.rw, fmt.Sprintf("MEASUrement:%s:VALue?", name))
	if err != nil {
		return "", fmt.Errorf("failed to query measurement value: %w", err)
	}
	return out, nil
}

func (d *Oscilloscope) Measurements() (map[string]string, error) {
	names, err := d.MeasurementList()
	if err != nil {
		return nil, err
	}

	measurements := make(map[string]string, len(names))
	for _, name := range names {
		value, err := d.measurementValue(name)
		if err != nil {
			return nil, err
		}
		measurements[name] = value
	}
	return measurements, nil
}

func exec(rw io.ReadWriter, cmd string) (string, error) {
	_, err := fmt.Fprintf(rw, cmd+"\r\n")
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(rw)
	out, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
