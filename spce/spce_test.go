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

    outputs, err := s.Outputs()

    if err != nil {
        t.Fatalf("Outputs() returned error: %v", err)
    }

    expectedKeys := []string{"current", "pressure", "voltage"}
    for _, key := range expectedKeys {
        value, ok := outputs[key]
        if !ok {
            t.Errorf("Outputs() did not return a value for key %s", key)
        } 
		fmt.Println(key, ":", value)
    }
}