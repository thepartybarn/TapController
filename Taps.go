package main

import (
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	_taps map[int]*tapStruct
)

type tapStruct struct {
	ID         int
	OpenRelay  rpio.Pin
	CloseRelay rpio.Pin
	CloseTimer *time.Timer
	IsOpen     bool
	Mutex      sync.Mutex
}

func (tap *tapStruct) Open() {
	log.Tracef("Open Tap %v", tap.ID)
	tap.IsOpen = true
	tap.Mutex.Lock()
	tap.OpenRelay.Low()
	time.Sleep(200 * time.Millisecond)
	tap.CloseRelay.High()
	time.Sleep(700 * time.Millisecond)
	tap.OpenRelay.High()
	time.Sleep(200 * time.Millisecond)
	tap.Mutex.Unlock()
}

func (tap *tapStruct) Close() {
	log.Tracef("Close Tap %v", tap.ID)
	tap.IsOpen = false
	tap.Mutex.Lock()
	tap.OpenRelay.High()
	time.Sleep(200 * time.Millisecond)
	tap.CloseRelay.Low()
	time.Sleep(700 * time.Millisecond)
	tap.CloseRelay.High()
	time.Sleep(200 * time.Millisecond)
	tap.Mutex.Unlock()
}

func handleTapRelays() {
	var err error
	log.Trace("Opening GPIO")
	err = rpio.Open()
	if err != nil {
		log.Panicln("Could not open GPIO")
		return
	}
	defer rpio.Close()
	log.Trace("Opened GPIO Successfully")

	log.Trace("Setting Up GPIO")

	_taps = make(map[int]*tapStruct)

	_taps[1] = &tapStruct{ID: 1, OpenRelay: rpio.Pin(6), CloseRelay: rpio.Pin(5), CloseTimer: time.NewTimer(500 * time.Millisecond)}
	_taps[2] = &tapStruct{ID: 2, OpenRelay: rpio.Pin(16), CloseRelay: rpio.Pin(13), CloseTimer: time.NewTimer(500 * time.Millisecond)}
	_taps[3] = &tapStruct{ID: 3, OpenRelay: rpio.Pin(20), CloseRelay: rpio.Pin(19), CloseTimer: time.NewTimer(500 * time.Millisecond)}
	_taps[4] = &tapStruct{ID: 4, OpenRelay: rpio.Pin(26), CloseRelay: rpio.Pin(21), CloseTimer: time.NewTimer(500 * time.Millisecond)}

	for _, tap := range _taps {
		tap.OpenRelay.Output()
		tap.OpenRelay.High()
		tap.CloseRelay.Output()
		tap.CloseRelay.High()
	}

	log.Trace("Setup GPIO Successfully")

	for {
		select {
		case <-_taps[1].CloseTimer.C:
			go _taps[1].Close()
		case <-_taps[2].CloseTimer.C:
			go _taps[2].Close()
		case <-_taps[3].CloseTimer.C:
			go _taps[3].Close()
		case <-_taps[4].CloseTimer.C:
			go _taps[4].Close()
		}
	}

	return
}
