package kuehlmobil

import (
	"fmt"
	"net"
	"log"
	"os"
	"testing"
)

var address = os.Getenv("ADDRESS")

func TestLoopRead(t *testing.T) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("failed to connect to controller: %s", err)
	}
	defer conn.Close()

	km := NewKuehlmobil(conn)
	temp, err := km.Read()
	if err != nil {
		t.Fatalf("failed to read loop: %s", err)
	}

	fmt.Println("Outflow Temp: ", temp)
}
