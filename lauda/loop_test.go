package lauda_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/rydyb/godevices/lauda"
)

var loop *lauda.Loop

func TestMain(m *testing.M) {
	addr := os.Getenv("ADDR")
	if addr == "" {
		fmt.Println("ADDR not set, skipping tests")
		os.Exit(0)
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial telnet: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	loop = lauda.NewLoop(conn)
	os.Exit(m.Run())
}

func TestLoop_ReadSetTemp(t *testing.T) {
	val, err := loop.ReadSetTemp()
	if err != nil {
		t.Fatalf("ReadSetTemp failed: %v", err)
	}
	t.Logf("SetTemp: %f", val)
}

func TestLoop_ReadOutflowTemp(t *testing.T) {
	val, err := loop.ReadOutflowTemp()
	if err != nil {
		t.Fatalf("ReadOutflowTemp failed: %v", err)
	}
	t.Logf("OutflowTemp: %f", val)
}

func TestLoop_ReadOutflowTempLimitHigh(t *testing.T) {
	val, err := loop.ReadOutflowTempLimitHigh()
	if err != nil {
		t.Fatalf("ReadOutflowTempLimitHigh failed: %v", err)
	}
	t.Logf("OutflowTempLimitHigh: %f", val)
}

func TestLoop_ReadOutflowTempLimitLow(t *testing.T) {
	val, err := loop.ReadOutflowTempLimitLow()
	if err != nil {
		t.Fatalf("ReadOutflowTempLimitLow failed: %v", err)
	}
	t.Logf("OutflowTempLimitLow: %f", val)
}

func TestLoop_ReadControlParameterXp(t *testing.T) {
	val, err := loop.ReadControlParameterXp()
	if err != nil {
		t.Fatalf("ReadControlParameterXp failed: %v", err)
	}
	t.Logf("ControlParameterXp: %f", val)
}

func TestLoop_ReadControlParameterTn(t *testing.T) {
	val, err := loop.ReadControlParameterTn()
	if err != nil {
		t.Fatalf("ReadControlParameterTn failed: %v", err)
	}
	t.Logf("ControlParameterTn: %f", val)
}

func TestLoop_ReadControlParameterTv(t *testing.T) {
	val, err := loop.ReadControlParameterTv()
	if err != nil {
		t.Fatalf("ReadControlParameterTv failed: %v", err)
	}
	t.Logf("ControlParameterTv: %f", val)
}

func TestLoop_ReadControlParameterTd(t *testing.T) {
	val, err := loop.ReadControlParameterTd()
	if err != nil {
		t.Fatalf("ReadControlParameterTd failed: %v", err)
	}
	t.Logf("ControlParameterTd: %f", val)
}

func TestLoop_ReadDeviceType(t *testing.T) {
	s, err := loop.ReadDeviceType()
	if err != nil {
		t.Fatalf("ReadDeviceType failed: %v", err)
	}
	if s == "" {
		t.Fatal("ReadDeviceType returned empty")
	}
	t.Logf("DeviceType: %s", s)
}

func TestLoop_ReadSoftwareVersion(t *testing.T) {
	s, err := loop.ReadSoftwareVersion()
	if err != nil {
		t.Fatalf("ReadSoftwareVersion failed: %v", err)
	}
	if s == "" {
		t.Fatal("ReadSoftwareVersion returned empty")
	}
	t.Logf("SoftwareVersion: %s", s)
}
