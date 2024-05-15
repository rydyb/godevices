package leybold

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

	cm52 := NewCombivacCM52(conn, 3)

	tests := []struct {
		c CMD
	}{
		{Pressure},
		{GasCorrection},
	}

	for _, test := range tests {
		value, err := cm52.ReadFloat(test.c)
		if err != nil {
			t.Fatalf("failed to read command %s: %v", test.c, err)
		}
		fmt.Printf("value: %e\n", value)
	}
}
