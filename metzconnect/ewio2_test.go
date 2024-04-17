package metzconnect

import (
	"os"
	"testing"

	"github.com/goburrow/modbus"
)

var address = os.Getenv("ADDRESS")

func TestEWIO2_AnalogInput(t *testing.T) {
	ewio2 := EWIO2{modbus.TCPClient(address)}

	for i := 1; i <= 11; i++ {
		voltage, err := ewio2.AnalogInput(i)
		if err != nil {
			t.Fatalf("failed to read analog input: %s", err)
		}

		if voltage < 0.0 || voltage > 10.0 {
			t.Errorf("voltage out of range: %f", voltage)
		}
	}
}

func TestChannelToAnalogInputAddress(t *testing.T) {
	tests := []struct {
		channel int
		address uint16
	}{
		{1, 40},
		{2, 42},
		{3, 44},
		{4, 140},
		{5, 142},
		{6, 144},
		{7, 146},
		{8, 148},
		{9, 150},
		{10, 152},
		{11, 154},
		{12, 240},
		{13, 242},
		{14, 244},
		{15, 246},
		{16, 248},
		{17, 250},
		{18, 252},
		{19, 254},
		{20, 340},
	}

	for _, test := range tests {
		address, err := AnalogInputAddress(test.channel)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if address != test.address {
			t.Errorf("wrong address for channel %d: %d", test.channel, address)
		}
	}
}
