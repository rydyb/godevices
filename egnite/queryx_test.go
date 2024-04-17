package egnite

import (
	"os"
	"testing"

	"github.com/goburrow/modbus"
)

var address = os.Getenv("ADDRESS")

func TestQueryx(t *testing.T) {
	client := modbus.TCPClient(address)
	q := Queryx{client}

	tests := []struct {
		q   Quantity
		min float64
		max float64
	}{
		{Temperature, 0.0, 300.0},
		{Humidity, 0.0, 100.0},
		{Pressure, 900.0, 1100.0},
	}

	for _, test := range tests {
		v, err := q.ReadFloat(test.q)
		if err != nil {
			t.Fatal(err)
		}
		if v < test.min || v > test.max {
			t.Errorf("value out of range: %f", v)
		}
	}
}
