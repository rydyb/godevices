package telnet

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Exec sends a command to the device and returns the response for telnet-like protocols.
func Exec(rw io.ReadWriter, cmd string) (string, error) {
	_, err := fmt.Fprintf(rw, cmd+"\r\n")
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(rw)
	out, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
