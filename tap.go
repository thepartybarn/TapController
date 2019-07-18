package main

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

var ()

type tapStruct struct {
	OpenRelay  rpio.Pin
	CloseRelay rpio.Pin
	IsOpen     bool
}

func (tap *tapStruct) Open() {
	log.Tracef("Open Tap")
	tap.IsOpen = true
	tap.OpenRelay.Low()
	tap.CloseRelay.High()
	time.Sleep(700 * time.Millisecond)
	tap.OpenRelay.High()
}

func (tap *tapStruct) Close() {
	log.Tracef("Close Tap")
	tap.IsOpen = false
	tap.OpenRelay.High()
	tap.CloseRelay.Low()
	time.Sleep(700 * time.Millisecond)
	tap.CloseRelay.High()
}

func setupTapRelays() (err error, taps map[int]*tapStruct) {
	log.Trace("Opening GPIO")
	err = rpio.Open()
	if err != nil {
		return
	}
	log.Trace("Opened GPIO Successfully")

	log.Trace("Setting Up GPIO")

	taps = make(map[int]*tapStruct)

	taps[1] = &tapStruct{OpenRelay: rpio.Pin(5), CloseRelay: rpio.Pin(6), IsOpen: false}
	taps[2] = &tapStruct{OpenRelay: rpio.Pin(13), CloseRelay: rpio.Pin(16), IsOpen: false}
	taps[3] = &tapStruct{OpenRelay: rpio.Pin(19), CloseRelay: rpio.Pin(20), IsOpen: false}
	taps[4] = &tapStruct{OpenRelay: rpio.Pin(21), CloseRelay: rpio.Pin(26), IsOpen: false}

	for _, tap := range taps {
		tap.OpenRelay.Output()
		tap.OpenRelay.High()
		tap.CloseRelay.Output()
		tap.CloseRelay.High()

		tap.Close()
	}
	log.Trace("Setup GPIO Successfully")
	return
}
