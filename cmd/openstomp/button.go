package main

type prog struct {
	ID int
	ProgramName string
	Controls []control
}


type control struct {
	ID int
	ControlName string
	Mode string
	HoldMode string
	Default bool
}

type StompButton struct {
	State bool
	Enabled bool
}

type ProgramButton struct {
	StompButton
	ProgramName string
}

type ControlButton struct {
	StompButton
	ControlName string
	Mode string
	HoldMode string
	buttonChanged bool
	OldEnableState bool
	LED bool
	IndependentLED bool
}

type NavButton struct {
	State bool
}

