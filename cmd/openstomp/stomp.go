package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

var m sync.Mutex

var ctx, cancel = context.WithCancel(context.Background())
var BPMCtx, BPMCancel = context.WithCancel(context.Background())

var controlButtons []ControlButton
var programButtons []ProgramButton
var navButtons []NavButton

var activeProgram prog
var activeProgramNumber int

var port io.ReadWriteCloser

var currentStompPage int

var timestamps []time.Time 

func InitStomp() {

	for i := 0; i < programNum; i++ {
		programButtons = append(programButtons, ProgramButton{})
	}

	for i := 0; i < controlNum; i++ {
		controlButtons = append(controlButtons, ControlButton{})
	}

	for i := 0; i < 2; i++ {
		navButtons = append(navButtons, NavButton{})
	}

	SerialInit()
}

func SerialInit() {
	options := serial.OpenOptions {
		PortName: selectedSerialPort,
		BaudRate: 115200,
		InterCharacterTimeout: 10000,
		MinimumReadSize: 1,
		StopBits: 1,
		DataBits: 8,
		ParityMode: serial.PARITY_NONE,
	}


	serialPort, err := serial.Open(options)

	port = serialPort

	if err != nil { 
		log.Fatalf("Error opening serial port: %v", err)
	}

	go Startup()

}

func Startup() {
	WriteToPort([]byte("sConnected" + "\n"));
	WriteToPort([]byte("l1,0,0,0,0,0,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,0,0,0,0,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,0,0,0,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,1,0,0,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,1,1,0,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,1,1,1,0,0" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,1,1,1,1,1" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l1,1,1,1,1,1,1,1" + "\n"));
	time.Sleep(100 * time.Millisecond)
	WriteToPort([]byte("l0,0,0,0,0,0,0,0" + "\n"));

	QueryPort()
}

func QueryPort() {
	
	var buffer []byte
	var input string

	for {
		buf := make([]byte, 128)

		n, err := port.Read(buf)
		if err != nil {
			//log.Fatalf("Serial Read Error: %v", err)
			serialSelect.ClearSelected()
			serialSelect.Enable()
			serialConfirm.Enable()
			port.Close()
			programButtons = programButtons[1:]
			controlButtons = controlButtons[1:]
			navButtons = navButtons[1:]
			return
		}

		if n > 0 {
			buffer = append(buffer, buf[:n]...)

		}

		for i := 0; i < len(buffer); i++ {
			if buffer[i] == '\n' {
				input = strings.TrimSpace(string(buffer[:i]))

				// Debugging stuff:
				//fmt.Println("buf: ", string(buf[:i]))
				//fmt.Println("buffer: ", string(buffer[:i]))
				//fmt.Println("input: ", input)

				buffer = buffer[i+1:]


				splitInput := strings.Split(input, ",")

				//fmt.Println("splitInput Length: ", len(splitInput))

				var parts []bool

				for _, v := range splitInput {

					if v == "0" {
						parts = append(parts, false)
					}
					if v == "1" {
						parts = append(parts, true)
					}
				}

				navButtons[0].State = parts[9]
				navButtons[1].State = parts[4]

				for i := 0; i < 4; i++ {
					programButtons[i].State = parts[i]
					controlButtons[i].State = parts[i + 5]
				}

				ButtonLogic()
			}
		}
	}
}

func ButtonLogic() {

	runProgramLogic := true

	if navButtons[0].State && currentStompPage > 0 {
		currentStompPage--;
		SendPageMessage()
	}

	if navButtons[1].State && currentStompPage < (pageLimit - 1) {
		currentStompPage++;
		SendPageMessage()
	}

	inactiveButtons := 0

	for _, button := range programButtons {
		if button.State == false {
			inactiveButtons++
		}
	}

	if inactiveButtons >= len(programButtons) {
		runProgramLogic = false
	}

	if runProgramLogic == true {
		for i,_ := range programButtons {
			var button *ProgramButton = &programButtons[i]

			button.Enabled = false

			if button.State == true {
				button.Enabled = true
				activeProgram = programs[i + (currentStompPage * len(programButtons))]
				activeProgramNumber = i + (currentStompPage * len(programButtons))
				SendPCMidi()
				SendMessage(activeProgram.ProgramName)

				for a, control := range activeProgram.Controls {
					var controlButton *ControlButton = &controlButtons[a]

					controlButton.ControlName = control.ControlName
					controlButton.Mode = control.Mode
					controlButton.HoldMode = control.HoldMode
					controlButton.Enabled = control.Default

					if control.Mode == "BPM Tap" {
						controlButton.IndependentLED = true;
					} else {
						controlButton.IndependentLED = false;
					}
				}
			}
		}
	}

	for i,_ := range controlButtons {
		var button *ControlButton = &controlButtons[i]

		if button.State == true && button.Mode == "Disabled" {
			SendMessage("Button Disabled")
		}

		if button.State == true && button.Mode == "BPM Tap" {
			SendCCMidi(button, -1)

			now := time.Now()

			if len(timestamps) > 0 {
				lastInterval := now.Sub(timestamps[len(timestamps)-1])

				if lastInterval.Seconds() > 1 {
					timestamps = []time.Time{}
				}
			}

			timestamps = append(timestamps, now)

			var intervals []float64

			var intervalSum float64

			if len(timestamps) > 3 {
				if len(timestamps) > 4 {
					timestamps = timestamps[1:]
				}

				for i := 1; i < len(timestamps); i++ {
					prev := timestamps[i-1]
					next := timestamps[i]

					interval := next.Sub(prev).Seconds()

					intervalSum = intervalSum + interval

					intervals = append(intervals, interval)
				}
			}

			if len(intervals) > 0 {
				var intervalAverage float64 = intervalSum / float64(len(intervals))
				var BPM float64 = 60.0 / intervalAverage
				var roundedBPM int = int(math.Round(BPM))

				BPMCancel()
				BPMCtx, BPMCancel = context.WithCancel(context.Background())

				go FlashLEDOnBPM(roundedBPM, button, BPMCtx)

				SendMessage(fmt.Sprintf("BPM: %v", roundedBPM))
			}


		}

		if button.State == true && (button.Mode == "Toggle" || button.Mode == "Momentary") {
			if button.Enabled == true { 
				button.Enabled = false
			} else if button.Enabled == false { 
				button.Enabled = true
			}

			SendCCMidi(button, i)
			SendMessage(fmt.Sprintf("%v %v", button.ControlName, OnOrOff(button.Enabled)))
		}
		if button.State == false && button.Mode == "Momentary" && button.OldEnableState	{
			button.Enabled = false
			SendCCMidi(button, i)
			SendMessage(fmt.Sprintf("%v %v", button.ControlName, OnOrOff(button.Enabled)))
		}
		if button.State == true && button.Mode == "Selection" {
			for i, _ := range controlButtons {
				var selectionButton *ControlButton = &controlButtons[i]

				if selectionButton.Mode == "Selection" {
					selectionButton.Enabled = false;
				}
			}
			button.Enabled = true;
			SendCCMidi(button, i)
			SendMessage(button.ControlName)
		}

		if button.Mode == "Momentary" {
			button.OldEnableState = button.Enabled
		}
	}

	UpdateLEDS()
}

func UpdateLEDS() {
	var LEDMessage string

	for i, _ := range programButtons {
		var button *ProgramButton = &programButtons[i]

		if button.Enabled == true {
			LEDMessage = LEDMessage + "1,"
		}
		if button.Enabled == false {
			LEDMessage = LEDMessage + "0,"
		}
	}

	for i, _ := range controlButtons {
		var button *ControlButton = &controlButtons[i]

		if button.IndependentLED == false {
			if button.Enabled == true {
				LEDMessage = LEDMessage + "1,"
			}
			if button.Enabled == false {
				LEDMessage = LEDMessage + "0,"
			}
		}

		if button.IndependentLED == true {
			if button.LED == true {
				LEDMessage = LEDMessage + "1,"
			}
			if button.LED == false {
				LEDMessage = LEDMessage + "0,"
			}
		}
	}

	go WriteToPort([]byte("l" + LEDMessage + "\n"))


}

func SendCCMidi (button *ControlButton, buttonIndex int) {
	var ccNumber int
	if buttonIndex == -1 {
		ccNumber = 1
	} else {
		ccNumber = 10 + buttonIndex
	}

	var value int
	if button.Enabled == true {
		value = 127
	} else {
		value = 0
	}

	go WriteToPort([]byte(fmt.Sprintf("c%v,%v\n", ccNumber, value)))
}

func SendPCMidi() {
	go WriteToPort([]byte(fmt.Sprintf("p%v\n", activeProgramNumber)))
}

func SendMessage(message string) {
	go WriteToPort([]byte("s" + message + "\n"))
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	go MessageDelay(ctx)
}

func MessageDelay(ctx context.Context) {
	select {
	case <- time.After(2 * time.Second):
		SendPageMessage()
	case <- ctx.Done():
		return
	}
}

func SendPageMessage() {
	go WriteToPort([]byte(fmt.Sprintf("sPage %v of %v\n", (currentStompPage + 1), pageLimit)))
}

func OnOrOff(boolean bool) string {
	if boolean {
		return "On"
	} else {
		return "Off"
	}
}

func WriteToPort(message []byte) {
	m.Lock()
	port.Write(message)
	m.Unlock()
}

func FlashLEDOnBPM(roundedBPM int, button *ControlButton, BPMCtx context.Context) {
	interval := 60000 / roundedBPM

	for {
		button.LED = true
		UpdateLEDS()
		time.Sleep(100 * time.Millisecond)
		button.LED = false
		UpdateLEDS()

		select {
		case <- BPMCtx.Done():
			return
		default:
			time.Sleep(time.Duration(interval - 100) * time.Millisecond)
			continue
		}
	}



}
