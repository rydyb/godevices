package highfinesse

/*
#cgo windows LDFLAGS: -L. -lwlmData

#include "wlmData.h"
*/
import "C"
import "fmt"

var (
	ErrNoValue    = fmt.Errorf("no value")
	ErrNoSignal   = fmt.Errorf("no signal")
	ErrBadSignal  = fmt.Errorf("bad signal")
	ErrLowSignal  = fmt.Errorf("low signal")
	ErrBigSignal  = fmt.Errorf("big signal")
	ErrNoPulse    = fmt.Errorf("no pulse")
	ErrWlmMissing = fmt.Errorf("wlm missing")

	ErrNotAvailable = fmt.Errorf("feature not supported")
	ErrNotMeasured  = fmt.Errorf("feature not measured yet")
	ErrWLMMissing   = fmt.Errorf("wavelength meter server missing")
)

// Channels returns the number of channels the wavemeter has.
func Channels() uint8 {
	return uint8(C.GetChannelsCount(C.long(0)))
}

// Wavlenegth returns the wavelength of the specified channel in nm.
func Wavelength(channel uint32) (float64, error) {
	switch retval := C.GetWavelengthNum(C.long(channel), 0); retval {
	case 0:
		return 0, ErrNoValue
	case -1:
		return 0, ErrNoSignal
	case -2:
		return 0, ErrBadSignal
	case -3:
		return 0, ErrLowSignal
	case -4:
		return 0, ErrBigSignal
	case -5:
		return 0, ErrNoPulse
	case -6:
		return 0, ErrWlmMissing
	case -7:
		return 0, ErrNotAvailable
	default:
		return float64(retval), nil
	}
}

// Frequency returns the frequency of the specified channel in THz.
func Frequency(channel uint32) (float64, error) {
	switch retval := C.GetFrequencyNum(C.long(channel), 0); retval {
	case 0:
		return 0, ErrNoValue
	case -1:
		return 0, ErrNoSignal
	case -2:
		return 0, ErrBadSignal
	case -3:
		return 0, ErrLowSignal
	case -4:
		return 0, ErrBigSignal
	case -5:
		return 0, ErrNoPulse
	case -6:
		return 0, ErrWlmMissing
	case -7:
		return 0, ErrNotAvailable
	default:
		return float64(retval), nil
	}
}

// Temperature returns the temperature of the wavemeter in °C.
func Temperature() (float64, error) {
	switch retval := C.GetTemperature(0); retval {
	case -1:
		return 0, ErrNotMeasured
	case -2:
		return 0, ErrNotAvailable
	case -3:
		return 0, ErrWLMMissing
	default:
		return float64(retval), nil
	}
}

// Pressure returns the pressure inside the wavemeter's optical unit in mbar.
func Pressure() (float64, error) {
	switch retval := C.GetPressure(0); retval {
	case -1:
		return 0, ErrNotMeasured
	case -2:
		return 0, ErrNotAvailable
	case -3:
		return 0, ErrWLMMissing
	default:
		return -float64(retval), nil
	}
}

// WavlengthMeterVersion encodes the version details of the wavemeter.
type WavelengthMeterVersion struct {
	TypeID          uint32
	HardwareVersion uint32
	SoftwareVersion uint32
	SoftwareBuild   uint32
}

func (v WavelengthMeterVersion) String() string {
	return fmt.Sprintf("TypeID: %d, HardwareVersion: %d, SoftwareVersion: %d, SoftwareBuild: %d",
		v.TypeID, v.HardwareVersion, v.SoftwareVersion, v.SoftwareBuild)
}

// WavelengthMeterVersionInfo returns the version information of the wavemeter.
func WavelengthMeterVersionInfo() WavelengthMeterVersion {
	return WavelengthMeterVersion{
		TypeID:          uint32(C.GetWLMVersion(0)),
		HardwareVersion: uint32(C.GetWLMVersion(1)),
		SoftwareVersion: uint32(C.GetWLMVersion(2)),
		SoftwareBuild:   uint32(C.GetWLMVersion(3)),
	}
}
