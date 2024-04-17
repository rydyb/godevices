package srs

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type PTC10 struct {
	rw io.ReadWriter
}

func NewPTC10(rw io.ReadWriter) *PTC10 {
	return &PTC10{rw: rw}
}

func (d *PTC10) Outputs() (map[string]float64, error) {
	resp, err := exec(d.rw, "getOutputs.names")
	if err != nil {
		return nil, fmt.Errorf("failed to get output names: %s", err)
	}
	names := strings.Split(resp, ", ")

	resp, err = exec(d.rw, "getOutputs")
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
