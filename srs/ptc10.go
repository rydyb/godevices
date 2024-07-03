package srs

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/rydyb/godevices/internal/telnet"
)

type PTC10 struct {
	rw io.ReadWriter
}

func NewPTC10(rw io.ReadWriter) *PTC10 {
	return &PTC10{rw: rw}
}

func (d *PTC10) Outputs() (map[string]float64, error) {
	resp, err := telnet.Exec(d.rw, "getOutputs.names")
	if err != nil {
		return nil, fmt.Errorf("failed to get output names: %s", err)
	}
	names := strings.Split(resp, ", ")

	resp, err = telnet.Exec(d.rw, "getOutputs")
	if err != nil {
		return nil, fmt.Errorf("failed to get output values: %s", err)
	}
	values := strings.Split(resp, ", ")

	if len(names) != len(values) {
		return nil, fmt.Errorf("output names and values mismatch")
	}

	outputs := make(map[string]float64)
	for i := 0; i < len(names); i++ {
		outputs[names[i]], err = strconv.ParseFloat(values[i], 64)
		if err != nil {
			continue
		}
	}
	return outputs, nil
}
