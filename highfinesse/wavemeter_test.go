package highfinesse

import (
	"fmt"
	"testing"
)

func TestChannels(t *testing.T) {
	n := Channels()
	if n < 0 {
		t.Fatalf("failed measuring channels")
	}
	t.Logf("measured channels: %d", n)
}

func TestWavelength(t *testing.T) {
	λ, err := Wavelength(1)
	if err != nil {
		t.Fatalf("failed measuring wavelength: %s", err)
	}
	t.Logf("measured wavelength: %f nm", λ)
}

func TestFrequency(t *testing.T) {
	f, err := Frequency(1)
	if err != nil {
		fmt.Printf("failed measuring frequency: %s", err)
	}
	t.Logf("measured frequency: %f THz", f)
}

func TestTemperature(t *testing.T) {
	T, err := Temperature()
	if err != nil {
		fmt.Printf("failed measuring temperature: %s", err)
	}
	t.Logf("measured temperature: %f °C", T)
}

func TestPressure(t *testing.T) {
	p, err := Pressure()
	if err != nil {
		fmt.Printf("failed measuring pressure: %s", err)
	}
	t.Logf("measured pressure: %f mbar", p)
}

func ExampleWavelengthMeterVersionInfo() {
	fmt.Printf("%+v", WavelengthMeterVersionInfo())
}
