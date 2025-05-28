package main

import (
	"fmt"
	"image/color"
	"time"

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

	ports,_ := ListSerialPorts()
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
			ports,_ = ListSerialPorts()
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
