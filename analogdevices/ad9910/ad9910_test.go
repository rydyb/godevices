package ad9910

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinearAmplitudeScaleToASF(t *testing.T) {
	tests := []struct {
		in  float64
		out uint16
	}{
		{1.0, uint16(0x4000)},
		{0.0, uint16(0x0)},
	}

	for _, test := range tests {
		got := LinearAmplitudeScaleToASF(test.in)
		want := test.out

		if got != want {
			t.Errorf("got 0x%x; want 0x%x", got, want)
		}
	}

	require.PanicsWithValue(t, "amplitude scale cannot be less than zero", func() { LinearAmplitudeScaleToASF(-0.01) })
	require.PanicsWithValue(t, "amplitude scale cannot be greater than one", func() { LinearAmplitudeScaleToASF(1.01) })
}

func TestLogarithmicAmplitudeScaleToASF(t *testing.T) {
	tests := []struct {
		in  float64
		out uint16
	}{
		{0.0, uint16(0x4000)},
		{-84.2884, uint16(0x1)},
	}

	for _, test := range tests {
		got := LogarithmicAmplitudeScaleToASF(test.in)
		want := test.out

		if got != want {
			t.Errorf("got 0x%x; want 0x%x", got, want)
		}
	}

	require.PanicsWithValue(t, "amplitude scale cannot be greater than zero", func() { LogarithmicAmplitudeScaleToASF(0.01) })
	require.PanicsWithValue(t, "amplitude scale cannot be less than -84.2884", func() { LogarithmicAmplitudeScaleToASF(-84.2885) })
}

func TestFreqOutToFTW(t *testing.T) {
	tests := []struct {
		freqOut float64
		sysClk  float64
		out     uint32
	}{
		{41e6, 122.88e6, uint32(0x556aaaab)},
	}

	for _, test := range tests {
		got := FrequencyToFTW(test.freqOut, test.sysClk)
		want := test.out

		if got != want {
			t.Errorf("got 0x%x; want 0x%x", got, want)
		}
	}

	require.PanicsWithValue(t, "output frequency cannot be less than zero", func() { FrequencyToFTW(-1.0, 1.0) })
	require.PanicsWithValue(t, "output frequency cannot be greater than half the system clock", func() { FrequencyToFTW(1.0, 0.5) })
}
