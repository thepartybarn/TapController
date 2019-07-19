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
	time.Sleep(200 * time.Millisecond)
	tap.CloseRelay.High()
	time.Sleep(700 * time.Millisecond)
	tap.OpenRelay.High()
	time.Sleep(200 * time.Millisecond)
}

func (tap *tapStruct) Close() {
	log.Tracef("Close Tap")
	tap.IsOpen = false
	tap.OpenRelay.High()
	time.Sleep(200 * time.Millisecond)
	tap.CloseRelay.Low()
	time.Sleep(700 * time.Millisecond)
	tap.CloseRelay.High()
	time.Sleep(200 * time.Millisecond)
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

	taps[1] = &tapStruct{OpenRelay: rpio.Pin(6), CloseRelay: rpio.Pin(5), IsOpen: false}
	taps[2] = &tapStruct{OpenRelay: rpio.Pin(16), CloseRelay: rpio.Pin(13), IsOpen: false}
	taps[3] = &tapStruct{OpenRelay: rpio.Pin(20), CloseRelay: rpio.Pin(19), IsOpen: false}
	taps[4] = &tapStruct{OpenRelay: rpio.Pin(26), CloseRelay: rpio.Pin(21), IsOpen: false}

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
