package kuehlmobil

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

func (kuehlmobil *Kuehlmobil) Read() (float64, error) {
	resp, err := telnet.Exec(kuehlmobil.rw, "get")
	if err != nil {
		return 0, fmt.Errorf("failed to get output: %s", err)
	}
	return strconv.ParseFloat(resp, 64)
}
