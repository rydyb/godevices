package metzconnect

import (
	"testing"

	"github.com/goburrow/modbus"
)

func TestMR8AI_AnalogInput(t *testing.T) {
	mr8ai := MR8AI{modbus.TCPClient(address)}

	for i := 1; i <= 8; i++ {
		voltage, err := mr8ai.AnalogInput(1, i)
		if err != nil {
			t.Fatalf("failed to read analog input: %s", err)
		}

		if voltage < 0.0 || voltage > 10.0 {
			t.Errorf("voltage out of range: %f", voltage)
		}
	}
}
