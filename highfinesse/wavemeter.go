package highfinesse

/*
#cgo windows LDFLAGS: -L. -lwlmData

#include "wlmData.h"
*/
import "C"

func Wavelength(channel uint32) float64 {
	return float64(C.GetWavelengthNum(C.long(channel), 0))
}

func Frequency(channel uint32) float64 {
	return float64(C.GetFrequencyNum(C.long(channel), 0))
}
