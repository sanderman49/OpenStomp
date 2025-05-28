//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func ListSerialPorts () ([]string, error) {
	var ports []string

	// On Linux, use /dev/serial/by-id/
	serialByIDPath := "/dev/serial/by-id"
	if _, err := os.Stat(serialByIDPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", serialByIDPath)
	}
	// Walk through the /dev/serial/by-id directory
	err := filepath.Walk(serialByIDPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process the symbolic links (files)
		if !info.IsDir() {
			// Resolve the symlink to get the actual device file (e.g., /dev/ttyUSB0)
			//devicePath, err := os.Readlink(path)
			//if err != nil {
			//	return err
			//}

			// Add the full path of the device to the list
			ports = append(ports, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ports, nil
}
