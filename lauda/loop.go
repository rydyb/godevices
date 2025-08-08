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
	cmdReadSetTemp              = "IN_SP_00"
	cmdReadOutflowTemp          = "IN_PV_00"
	cmdReadOutflowTempLimitHigh = "IN_SP_04"
	cmdReadOutflowTempLimitLow  = "IN_SP_05"
	cmdReadControlParameterXp   = "IN_PAR_00"
	cmdReadControlParameterTn   = "IN_PAR_01"
	cmdReadControlParameterTv   = "IN_PAR_02"
	cmdReadControlParameterTd   = "IN_PAR_03"
	cmdReadDeviceType           = "TYPE"
	cmdReadSoftwareVersion      = "VERSION_R"
)

func NewLoop(rw io.ReadWriter) *Loop {
	return &Loop{rw: rw}
}

func (loop *Loop) read(cmd string) (string, error) {
	resp, err := telnet.Exec(loop.rw, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to read cmd %s output: %s", cmd, err)
	}
	return resp, nil
}

func (loop *Loop) readFloat(cmd string) (float64, error) {
	resp, err := loop.read(cmd)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(resp, 64)
}

func (loop *Loop) ReadSetTemp() (float64, error) {
	return loop.readFloat(cmdReadSetTemp)
}

func (loop *Loop) ReadOutflowTemp() (float64, error) {
	return loop.readFloat(cmdReadOutflowTemp)
}

func (loop *Loop) ReadOutflowTempLimitHigh() (float64, error) {
	return loop.readFloat(cmdReadOutflowTempLimitHigh)
}

func (loop *Loop) ReadOutflowTempLimitLow() (float64, error) {
	return loop.readFloat(cmdReadOutflowTempLimitLow)
}

func (loop *Loop) ReadControlParameterXp() (float64, error) {
	return loop.readFloat(cmdReadControlParameterXp)
}

func (loop *Loop) ReadControlParameterTn() (float64, error) {
	return loop.readFloat(cmdReadControlParameterTn)
}

func (loop *Loop) ReadControlParameterTv() (float64, error) {
	return loop.readFloat(cmdReadControlParameterTv)
}

func (loop *Loop) ReadControlParameterTd() (float64, error) {
	return loop.readFloat(cmdReadControlParameterTd)
}

func (loop *Loop) ReadDeviceType() (string, error) {
	return loop.read(cmdReadDeviceType)
}

func (loop *Loop) ReadSoftwareVersion() (string, error) {
	return loop.read(cmdReadSoftwareVersion)
}
