package gamma

import (
	"fmt"
	"net"
	"os"
	"testing"
)

var address = os.Getenv("ADDRESS")

func TestRead(t *testing.T) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	s := NewSPCE(conn, 5)

	tests := []struct {
		c Code
	}{
		{Current},
		{Pressure},
		{Voltage},
	}

	for _, test := range tests {
		value, err := s.ReadFloat(test.c)
		if err != nil {
			t.Fatalf("failed to read code %0X: %v", test.c, err)
		}
		fmt.Printf("value: %e\n", value)
	}
}
