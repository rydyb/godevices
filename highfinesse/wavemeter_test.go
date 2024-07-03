package highfinesse

import (
	"testing"
)

func TestWavelength(t *testing.T) {
	if Wavelength(1) >= 0 {
		t.Errorf("wavelength cannot be smaller than 0 but was %f", Wavelength(1))
	}
}

func TestFrequency(t *testing.T) {
	if Frequency(1) >= 0 {
		t.Errorf("frequency cannot be smaller than 0 but was %f", Wavelength(1))
	}
}
