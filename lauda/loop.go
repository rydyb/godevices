package lauda

import (
	"fmt"
	"io"
	"strconv"

	"github.com/rydyb/godevices/internal/telnet"
)

type Loop struct {
	rw io.ReadWriter
}

var (
	// IN_PV_00 is the command to read the outflow temperature 
	outflowTemp = "IN_PV_00"
	// setTemp     = "OUT_SP_00"
	// outflowTempLimitHigh = "IN_SP_04"
	// outflowTempLimitLow = "IN_SP_05"
)

func NewLoop(rw io.ReadWriter) *Loop {
	return &Loop{rw: rw}
}

func (loop *Loop) Read() (float64, error) {
	resp, err := telnet.Exec(loop.rw, outflowTemp)
	if err != nil {
		return 0, fmt.Errorf("failed to get output names: %s", err)
	}
	return strconv.ParseFloat(resp, 64)
}
