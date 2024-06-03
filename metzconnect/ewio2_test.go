package metzconnect

import (
	"os"
	"testing"

	"github.com/goburrow/modbus"
)

var address = os.Getenv("ADDRESS")

func TestEWIO2_AnalogInput(t *testing.T) {
	ewio2 := EWIO2{modbus.TCPClient(address)}

	for i := 1; i <= 3; i++ {
		voltage, err := ewio2.AnalogInput(i)
		if err != nil {
			t.Fatalf("failed to read analog input: %s", err)
		}

		t.Logf("%v", voltage)
	}
}

func TestEWIO2_getExtensionTypes(t *testing.T) {
	ewio2 := EWIO2{modbus.TCPClient(address)}

	extensions, err := ewio2.GetExtensionTypes()
	if err != nil {
		t.Fatalf("failed to read extension types: %s", err)
	}

	t.Logf("extensions: %s", extensions)
}

func TestEWIO2_getUnitRange(t *testing.T) {
	ewio2 := EWIO2{modbus.TCPClient(address)}

	unitRange, err := ewio2.GetUnitRange()
	if err != nil {
		t.Fatalf("failed to read unit range: %s", err)
	}

	t.Logf("unit range: %v", unitRange)
}
