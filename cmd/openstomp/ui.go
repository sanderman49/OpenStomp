package main

import (
	"fmt"
	"image/color"
	"time"

	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var selectedSerialPort string

var serialConfirm *widget.Button
var serialSelect *widget.Select

func InitUI() {
	a := app.New()
	w := a.NewWindow("OpenStomp")

	mainPage := container.New(layout.NewVBoxLayout())


	programTitle := widget.NewEntry()
	pageText := widget.NewLabel(fmt.Sprint("Page ", currentPage + 1," of 8"))

	var selectedProgram *prog

	programTitle.OnChanged = func(updatedTitle string) {
		selectedProgram.ProgramName = updatedTitle 		
		SaveToFile()
	}

	ControlPage := func(controlID int) {
		var modes []string
		var holdModes []string

		modes = append(modes, "Toggle")
		modes = append(modes, "Momentary")
		modes = append(modes, "Selection")
		modes = append(modes, "BPM Tap")
		modes = append(modes, "Disabled")

		holdModes = append(holdModes, "None")
		holdModes = append(holdModes, "Chromatic Tuner")

		var selectedControl *control = &selectedProgram.Controls[controlID]

		//title := widget.NewLabel(fmt.Sprint(selectedProgram.ProgramName, ": ", selectedControl.ControlName))
		title := canvas.NewText(fmt.Sprint(selectedProgram.ProgramName, ": ", selectedControl.ControlName), color.White)
		title.TextSize = 32

		name := widget.NewLabel("Control Name")
		nameEntry := widget.NewEntry()
		nameEntry.SetText(selectedControl.ControlName)

		mode := widget.NewLabel("Mode")
		modeSelect := widget.NewSelect(modes, func(s string) {})
		modeSelect.SetSelected(selectedControl.Mode)

		holdMode := widget.NewLabel("Hold Mode")
		holdModeSelect := widget.NewSelect(holdModes, func(s string) {})
		holdModeSelect.SetSelected(selectedControl.HoldMode)

		defaultState := widget.NewLabel("Default State")
		defaultStateCheck := widget.NewCheck("", func(b bool) {})
		defaultStateCheck.SetChecked(selectedControl.Default)

		topSection := 
			container.NewVBox(
			title,
			name,
			nameEntry,
			mode,
			modeSelect,
			holdMode,
			holdModeSelect,
			defaultState,
			defaultStateCheck,
		)

		saveButton := widget.NewButton("Save", func() {
			selectedControl.Mode = modeSelect.Selected
			selectedControl.HoldMode = holdModeSelect.Selected
			selectedControl.Default = defaultStateCheck.Checked
			selectedControl.ControlName = nameEntry.Text

			SaveToFile()

			w.SetContent(mainPage)
		})
		exitButton := widget.NewButton("Exit", func() {
			w.SetContent(mainPage)
		})

		bottomSection :=
			container.NewHBox(
			layout.NewSpacer(),
			saveButton,
			exitButton,
			
		)

		controlPage := container.New(layout.NewVBoxLayout(), topSection, layout.NewSpacer() ,bottomSection)

		w.SetContent(controlPage)
	}

	PC_Pressed := func(ProgramId int) {
		selectedProgramIndex := ProgramId + (currentPage * programNum)

		selectedProgram = &programs[selectedProgramIndex]

		programTitle.SetText(selectedProgram.ProgramName);
	}

	CC_Pressed := func(ControlID int) {
		ControlPage(ControlID)
	}
	
	updatePage := func(change int)	{
		if (change > 0 && currentPage < (8 - 1)) || (change < 0 && currentPage > 0) {

			currentPage = currentPage + change
			pageText.SetText(fmt.Sprint("Page ", currentPage + 1, " of 8"))
		}
	}

	nextButton := widget.NewButton("Next", func() {updatePage(1)})
	prevButton := widget.NewButton("Previous", func() {updatePage(-1)})


	topSection := container.NewBorder(
		container.NewHBox(
			layout.NewSpacer(),
			pageText,
			layout.NewSpacer()),
		nil,
		prevButton,
		nextButton,
		programTitle,
	)

	//var ProgramButtons []widget.Button
	//var ControlButtons []widget.Button

	PC1 := widget.NewButton("PC0", func() {PC_Pressed(0)})
	PC2 := widget.NewButton("PC1", func() {PC_Pressed(1)})
	PC3 := widget.NewButton("PC2", func() {PC_Pressed(2)})
	PC4 := widget.NewButton("PC3", func() {PC_Pressed(3)})

	CC1 := widget.NewButton("CC0", func() {CC_Pressed(0)})
	CC2 := widget.NewButton("CC1", func() {CC_Pressed(1)})
	CC3 := widget.NewButton("CC2", func() {CC_Pressed(2)})
	CC4 := widget.NewButton("CC3", func() {CC_Pressed(3)})

	controlSection := container.New(layout.NewGridLayout(4), PC1, PC2, PC3, PC4, CC1, CC2, CC3, CC4)

	ports,_ := listSerialPorts()
	serial := widget.NewLabel("Select a Serial Port:")
	serialSelect = widget.NewSelect(ports, func(s string) {})
	serialConfirm = widget.NewButton("Connect", func () {
		selectedSerialPort = serialSelect.Selected
		go InitStomp()
		serialSelect.Disable()
		serialConfirm.Disable()
	})

	serialSection := 
		container.NewVBox(
		serial,
		serialSelect,
		serialConfirm,
	)

	mainPage = container.New(layout.NewVBoxLayout(), topSection, controlSection, serialSection)

	PC_Pressed(0)

	w.SetContent(mainPage)

	refreshPorts := func ()  {
		for {
			ports,_ = listSerialPorts()
			foundPort := false

			for _,v := range ports {
				if v == serialSelect.Selected {
					foundPort = true
				}
			}

			if foundPort == false {
				serialSelect.ClearSelected()
			}

			serialSelect.SetOptions(ports)
			time.Sleep(500 * time.Millisecond)
		}
	}

	go refreshPorts()

	w.Resize(fyne.NewSize(1000,600))

	w.ShowAndRun()
	
}


// List serial ports by scanning /dev/serial/by-id (Linux), /dev/tty.* (macOS), and COM ports (Windows)
func listSerialPorts() ([]string, error) {
	var ports []string

	switch runtime.GOOS {
	case "linux":
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

	case "darwin":
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

	case "windows":
		// On Windows, serial ports are named COM1, COM2, etc.
		// Use mode or query the registry for serial ports
		cmd := exec.Command("mode")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		// Parse the output to find COM ports
		// Look for lines starting with COM1, COM2, etc.
		for _, line := range strings.Split(string(output), "\n") {
			if strings.HasPrefix(line, "COM") {
				ports = append(ports, line)
			}
		}
	}

	return ports, nil
}

