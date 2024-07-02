package highfinesse

/*
#cgo LDFLAGS: -L. -lwlmdata.dll
#include "wavemeter.h"
*/
import "C"

func Wavelength(channel uint32) float64 {
	return float64(C.GetWavelengthNum(C.long(channel), 0))
}
