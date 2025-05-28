package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var programs []prog

var configPath string
var configFile string

var currentPage int = 0
var pageLimit int = 8

var programNum int = 4
var controlNum int = 4

func Init() {
	configPath = GetConfigPath()

	configFile = filepath.Join(configPath, "openstomp.json")

	_, err := os.Stat(configFile)

	if err != nil {
		// If the config file doesn't exist then generate a new config from scratch.
		if os.IsNotExist(err) {
			// Create all the programs
			for i:=0; i < (programNum * pageLimit); i++ {

				// Temporary controls list before creating the program.
				var controls []control

				// Create the controls required for each program.
				for a:=0; a < controlNum; a++ {
					controls = append(controls, control{
						a, // ID
						fmt.Sprint("Control ", a), // ControlName.
						"Toggle", // Press.
						"None", // Hold.
						false, // Default.
					})
				}

				// Create each program and append them to 'programs'.
				programs = append(programs, prog{
					i, // ID
					fmt.Sprint("Program ", i), // ProgramName.
					controls, // List with each control.
				})
			}

			SaveToFile()
		}
	} else {
		LoadFromFile()
	}

	InitUI()
}

func SaveToFile() {
	//fmt.Println("SaveToFile not implemented.")
	data, err := json.MarshalIndent(programs, "", "	")
	if err != nil {
		log.Fatalf("Error Serialising Data: %v", err)
	}

	_, err = os.Stat(configPath)
	if err != nil {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			fmt.Println(configPath)
			log.Fatalf("Error Creating Config Directory: %v", err)
		}
	}

	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Error Creating File: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Error Writing to File: %v", err)
	}
}

func LoadFromFile() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error Reading File: %v", err)
	}

	err = json.Unmarshal(data, &programs)
	if err != nil {
		log.Fatalf("Error Unmarshalling Json: %v", err)
	}
}

func GetConfigPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "openstomp")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Error Getting Home Directory: %v", err)
		}

		return filepath.Join(homeDir, "Library", "Application Support", "openstomp")
	case "linux":
		return filepath.Join(os.Getenv("HOME"), ".config", "openstomp")
	}
	return ""
}
