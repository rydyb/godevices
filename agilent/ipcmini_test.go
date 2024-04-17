package agilent

import (
	"fmt"
	"net"
	"os"
	"testing"
)

var address = os.Getenv("ADDRESS")

func TestIPCMini(t *testing.T) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	d := NewIPCMini(conn)

	tests := []struct {
		name string
		unit string
		w    Window
	}{
		{"Estimated pressure", "mbar", EstimatedPressure},
		{"Maximum power", "W", MaximumPower},
		{"Target voltage", "V", TargetVoltage},
		{"Power section temperature", "°C", PowerSectionTemperature},
		{"Internal controller temperature", "°C", InternalControllerTemperature},
		{"Measured voltage", "V", MeasuredVoltage},
		{"Measured current", "A", MeasuredCurrent},
		{"Estimated pressure", "mbar", EstimatedPressure},
	}

	for _, test := range tests {
		out, err := d.ReadFloat(test.w)
		if err != nil {
			t.Fatalf("failed to read %s: %v", test.name, err)
		}
		fmt.Println(test.name, ":", out, test.unit)
	}
}
