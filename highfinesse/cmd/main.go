package main

import (
	"log"
	"syscall"
)

func main() {
	wlmData, err := syscall.LoadDLL("wlmData.dll")
	if err != nil {
		log.Fatalf("Error loading wlmData.dll: %v", err)
	}
	defer syscall.FreeLibrary(wlmData.Handle)
}
