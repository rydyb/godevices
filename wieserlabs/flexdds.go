package wieserlabs

import (
	"fmt"
	"net"

	"github.com/rydyb/godevices/analogdevices/ad9910"
	"github.com/rydyb/godevices/internal/telnet"
)

// FlexDDSSlot represents a slot of a FlexDDS.
type FlexDDSSlot struct {
	clock float64
	conn  net.Conn
}

// NewFlexDDSSlot creates a new FlexDDSSlot which holds an authenticated connection to a FlexDDS controller for the respective slot.
func NewFlexDDSSlot(host string, slot uint8, clock float64) (*FlexDDSSlot, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, 26000+int(slot)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}

	out, err := telnet.Exec(conn, fmt.Sprintf("75f4a4e10dd4b6b%d", slot))
	if err != nil {
		return nil, fmt.Errorf("failed to send authentication token: %w", err)
	}
	if out != "Auth OK" {
		return nil, fmt.Errorf("failed to perform authentication: %s", out)
	}

	return &FlexDDSSlot{
		clock: clock,
		conn:  conn,
	}, nil
}

// Singletone configures channel to output a single frequency with relative logarithmic amplitude scale.
func (d *FlexDDSSlot) Singletone(channel uint8, logAmplitudeScale, frequency float64) error {
	asf := ad9910.LogarithmicAmplitudeScaleToASF(logAmplitudeScale)
	ftw := ad9910.FrequencyToFTW(frequency, d.clock)

	if _, err := telnet.Exec(d.conn, fmt.Sprintf("dcp %d spi:cfr2=0x01400820", channel)); err != nil {
		return fmt.Errorf("failed to configure CFR2 register: %w", err)
	}
	if _, err := telnet.Exec(d.conn, fmt.Sprintf("dcp %d spi:stp0=0x%x0000%x", channel, asf, ftw)); err != nil {
		return fmt.Errorf("failed to configure STP0 register: %w", err)
	}
	if _, err := telnet.Exec(d.conn, fmt.Sprintf("dcp %d update:u", channel)); err != nil {
		return fmt.Errorf("failed to update dds: %w", err)
	}

	return nil
}

// Close the connection.
func (d *FlexDDSSlot) Close() error {
	return d.conn.Close()
}
