package lauda

import (
	"fmt"
	"net"
	"log"
	"testing"
)

var address = "10.163.103.226:4196"

func TestLoopRead(t *testing.T) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("failed to connect to controller: %s", err)
	}
	defer conn.Close()

	loop := NewLoop(conn)
	temp, err := loop.Read()
	if err != nil {
		t.Fatalf("failed to read loop: %s", err)
	}

	fmt.Println("Outflow Temp: ", temp)
}
