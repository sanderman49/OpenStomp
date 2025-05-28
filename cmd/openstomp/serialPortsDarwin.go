//go:build darwin
// +build darwin

package main

import (
	"path/filepath"
	"strings"
	"os"
)

func ListSerialPorts () ([]string, error) {
	var ports []string

	// On macOS, serial ports are usually named /dev/tty.*
	// Look for serial ports in /dev/ directory matching tty.*
	err := filepath.Walk("/dev", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Match serial devices like /dev/tty.usbserial* or /dev/tty.usbmodem*
		if strings.HasPrefix(info.Name(), "tty.") && (strings.Contains(info.Name(), "usbserial") || strings.Contains(info.Name(), "usbmodem")) {
			ports = append(ports, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ports, nil
}
