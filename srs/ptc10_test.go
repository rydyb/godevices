package srs

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
		t.Fatalf("failed to connect to controller: %s", err)
	}
	defer conn.Close()

	ptc10 := PTC10{conn}

	outputs, err := ptc10.Outputs()
	if err != nil {
		t.Fatalf("failed to get outputs: %s", err)
	}

	fmt.Println("HW Out: ", outputs["HW Out"])
	fmt.Println("HW TC: ", outputs["HW TC"])
	fmt.Println("Oven Out: ", outputs["Oven Out"])
	fmt.Println("Oven TC: ", outputs["Oven TC"])
}
