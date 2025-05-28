//go:build windows
// +build windows

package main

import (
	"os/exec"
	"syscall"
	"strings"
)

func ListSerialPorts () ([]string, error) {
	var ports []string

	// On Windows, serial ports are named COM1, COM2, etc.
	// Use mode or query the registry for serial ports
	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", "Get-WmiObject Win32_SerialPort | Select-Object -ExpandProperty DeviceID")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Suppress window

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}	

	// Parse and print the output
	devices := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, device := range devices {
		ports = append(ports, device)
	}

	return ports, nil
}
