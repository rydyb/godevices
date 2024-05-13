package spce

import (
	"fmt"
	"net"
	"os"
	"testing"
)

var address = os.Getenv("ADDRESS")

func TestOutputs(t *testing.T) {
    conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

    s := NewSPCE(conn, "05")

    pressure, err := s.Pressure()
    if err != nil {
        t.Fatalf("failed to get pressure: %v", err)
    }

    current, err := s.Current()
    if err != nil {
        t.Fatalf("failed to get current: %v", err)
    }
    
    voltage, err := s.Voltage()
    if err != nil {
        t.Fatalf("failed to get voltage: %v", err)
    }

    fmt.Printf("Pressure: %v\n", pressure)
    fmt.Printf("Current: %v\n", current)
    fmt.Printf("Voltage: %v\n", voltage)   
}