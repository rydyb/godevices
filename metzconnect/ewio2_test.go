package metzconnect

import (
	"fmt"
	"os"
	"testing"

	"github.com/goburrow/modbus"
)

var address = os.Getenv("ADDRESS")

var ewio2 = NewEWIO2(modbus.TCPClient(address))

func TestEWIO2_Version(t *testing.T) {
	version, err := ewio2.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Version: %d", version)
}

func TestEWIO2_AnalogInputs_Mode(t *testing.T) {
	for i := 1; i < 4; i++ {
		t.Run(fmt.Sprintf("E%d", i), func(t *testing.T) {
			mode, err := ewio2.Mode(uint8(i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Mode: %d", mode)
		})
	}
}

func TestEWIO2_AnalogInputs_Unit(t *testing.T) {
	for i := 1; i < 4; i++ {
		t.Run(fmt.Sprintf("E%d", i), func(t *testing.T) {
			unit, err := ewio2.Unit(uint8(i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Unit: %s", unit)
		})
	}
}

func TestEWIO2_AnalogInputs_Value(t *testing.T) {
	for i := 1; i < 4; i++ {
		t.Run(fmt.Sprintf("E%d", i), func(t *testing.T) {
			value, err := ewio2.Value(uint8(i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Value: %v", value)
		})
	}
}

func TestEWIO2_Extensions(t *testing.T) {
	extensions, err := ewio2.Extensions()
	if err != nil {
		t.Fatal(err)
	}
	for _, extension := range extensions {
		t.Logf("Extension: %d %v", extension.ID(), extension.Type())

		switch extension.Type() {
		case ExtensionMR_AI8:
			t.Run("MR-AI8", func(t *testing.T) {
				mrai8, ok := extension.(*MRAI8)
				if !ok {
					t.Fatal("invalid extension type")
				}

				for i := 1; i < 9; i++ {
					t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
						mode, err := mrai8.Mode(uint8(i))
						if err != nil {
							t.Fatal(err)
						}

						value, err := mrai8.Value(uint8(i))
						if err != nil {
							t.Fatal(err)
						}

						t.Logf("Value: %v, Mode: %d", value, mode)
					})
				}

			})
		}
	}
}
