package visa

import (
	"fmt"
	"io"
	"strings"

	"github.com/rydyb/godevices/internal/telnet"
)

type Oscilloscope struct {
	rw io.ReadWriter
}

func NewOscilloscope(rw io.ReadWriter) *Oscilloscope {
	return &Oscilloscope{rw: rw}
}

func (d *Oscilloscope) Identity() (string, error) {
	out, err := telnet.Exec(d.rw, "*idn?")
	if err != nil {
		return "", fmt.Errorf("failed to query identity: %w", err)
	}
	return out, nil
}

func (d *Oscilloscope) MeasurementList() ([]string, error) {
	out, err := telnet.Exec(d.rw, "MEASUrement:LIST?")
	if err != nil {
		return nil, fmt.Errorf("failed to query measurement list: %w", err)
	}
	return strings.Split(out, ","), nil
}

func (d *Oscilloscope) measurementValue(name string) (string, error) {
	out, err := telnet.Exec(d.rw, fmt.Sprintf("MEASUrement:%s:VALue?", name))
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
