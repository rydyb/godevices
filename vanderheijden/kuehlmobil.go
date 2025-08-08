package vanderheijden

import (
	"fmt"
	"io"
	"strconv"

	"github.com/rydyb/godevices/internal/telnet"
)

type Kuehlmobil struct {
	rw io.ReadWriter
}

func NewKuehlmobil(rw io.ReadWriter) *Kuehlmobil {
	return &Kuehlmobil{rw: rw}
}

func (kuehlmobil *Kuehlmobil) read(cmd string) (string, error) {
	resp, err := telnet.Exec(kuehlmobil.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get output: %s", err)
	}
	return resp, nil
}

func (kuehlmobil *Kuehlmobil) readFloat(cmd string) (float64, error) {
	resp, err := kuehlmobil.read(cmd)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(resp, 64)
}

func (kuehlmobil *Kuehlmobil) readBool(cmd string) (bool, error) {
	resp, err := kuehlmobil.read(cmd)
	if err != nil {
		return false, err
	}

	switch resp {
	case "ON":
		return true, nil
	case "OFF":
		return false, nil
	}

	return false, fmt.Errorf("unexpected response for enabled status: %s", resp)
}

func (kuehlmobil *Kuehlmobil) ReadCurrentTemperature() (float64, error) {
	return kuehlmobil.readFloat("IST")
}

func (kuehlmobil *Kuehlmobil) ReadSetTemperature() (float64, error) {
	return kuehlmobil.readFloat("SET")
}

func (kuehlmobil *Kuehlmobil) ReadSystemEnabled() (bool, error) {
	return kuehlmobil.readBool("EIN")
}

func (kuehlmobil *Kuehlmobil) ReadPump1Enabled() (bool, error) {
	return kuehlmobil.readBool("Pump1")
}

func (kuehlmobil *Kuehlmobil) ReadRegulationEnabled() (bool, error) {
	return kuehlmobil.readBool("REGEL")
}
