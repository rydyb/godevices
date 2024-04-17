package agilent

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
)

// control chars
const (
	stx  = 0x02
	etx  = 0x03
	addr = 0x80
	rop  = 0x30
	wop  = 0x31
)

type Window string

const (
	PressureUnit                  Window = "600"
	MaximumPower                  Window = "612"
	TargetVoltage                 Window = "613"
	PowerSectionTemperature       Window = "800"
	InternalControllerTemperature Window = "801"
	MeasuredVoltage               Window = "810"
	MeasuredCurrent               Window = "811"
	EstimatedPressure             Window = "812"
)

const maxLength = 58

type IPCMini struct {
	rw io.ReadWriter
}

func NewIPCMini(rw io.ReadWriter) *IPCMini {
	return &IPCMini{rw: rw}
}

func (d *IPCMini) ReadFloat(w Window) (float64, error) {
	out, err := d.read(w)
	if err != nil {
		return 0, fmt.Errorf("failed to read window %s: %w", w, err)
	}
	return strconv.ParseFloat(string(out), 64)
}

func (d *IPCMini) read(w Window) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(stx)
	buf.WriteByte(addr)
	buf.WriteString(string(w))
	buf.WriteByte(rop)
	buf.WriteByte(etx)
	buf.Write(checksum(buf.Bytes()))

	message := buf.Bytes()
	_, err := d.rw.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}

	answer := make([]byte, maxLength)
	if _, err = d.rw.Read(answer); err != nil {
		return nil, fmt.Errorf("failed to read answer: %w", err)
	}
	answer = bytes.TrimRight(answer, "\x00")

	if answer[0] != stx {
		return nil, fmt.Errorf("invalid start of answer: %x", answer[0])
	}
	if answer[1] != addr {
		return nil, fmt.Errorf("invalid address of answer: %x", answer[1])
	}
	if answer[len(answer)-3] != etx {
		return nil, fmt.Errorf("invalid end of answer: %x", answer[len(answer)-3])
	}

	return answer[5 : len(answer)-3], nil
}

func checksum(message []byte) []byte {
	var checksum byte
	for _, b := range message[1:] {
		checksum ^= b
	}
	return []byte(hex.EncodeToString([]byte{checksum}))
}
