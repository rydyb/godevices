package vanderheijden_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/rydyb/godevices/vanderheijden"
)

var km *vanderheijden.Kuehlmobil

func TestMain(m *testing.M) {
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		fmt.Println("ADDRESS not set, skipping tests")
		os.Exit(0)
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	km = vanderheijden.NewKuehlmobil(conn)
	os.Exit(m.Run())
}

func TestKuehlmobil_ReadCurrentTemperature(t *testing.T) {
	val, err := km.ReadCurrentTemperature()
	if err != nil {
		t.Fatalf("ReadCurrentTemperature failed: %v", err)
	}
	t.Logf("CurrentTemperature: %f", val)
}

func TestKuehlmobil_ReadSetTemperature(t *testing.T) {
	val, err := km.ReadSetTemperature()
	if err != nil {
		t.Fatalf("ReadSetTemperature failed: %v", err)
	}
	t.Logf("SetTemperature: %f", val)
}

func TestKuehlmobil_ReadSystemEnabled(t *testing.T) {
	on, err := km.ReadSystemEnabled()
	if err != nil {
		t.Fatalf("ReadSystemEnabled failed: %v", err)
	}
	t.Logf("SystemEnabled: %t", on)
}

func TestKuehlmobil_ReadPump1Enabled(t *testing.T) {
	on, err := km.ReadPump1Enabled()
	if err != nil {
		t.Fatalf("ReadPump1Enabled failed: %v", err)
	}
	t.Logf("Pump1Enabled: %t", on)
}

func TestKuehlmobil_ReadRegulationEnabled(t *testing.T) {
	on, err := km.ReadRegulationEnabled()
	if err != nil {
		t.Fatalf("ReadRegulationEnabled failed: %v", err)
	}
	t.Logf("RegulationEnabled: %t", on)
}
