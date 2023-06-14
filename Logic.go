package main

import (
	"strconv"
	"time"
)

var (
	_lastUID       string
	_currentPerson *Person
	_Override      = true
)

func processUSBMessage(Message USBDataMessage) {
	log.Tracef("Message Received: %+v", Message)
	switch Message.Key {
	case "cardID":
		//Convert negative numbers to positive
		cardScan(Message.Value)
	case "ButtonPress":
		button, err := strconv.Atoi(Message.Value)
		if err != nil {
			log.Warn(err)
		}
		tapButtonPress(button)
	}
}
func tapButtonPress(button int) {
	//TODO handle beer types here
	log.Tracef("Tap Button %v Pressed", button)
	tap := _taps[button]

	if !tap.IsOpen && (_currentPerson != nil || _Override) {
		go tap.Open()
		tap.CloseTimer.Reset(500 * time.Millisecond)
	}
	if tap.IsOpen {
		tap.CloseTimer.Reset(500 * time.Millisecond)
	}
}

func cardScan(UID string) {
	log.Tracef("Card Scanned %v", UID)
	_lastUID = UID
	Person, err := _database.GetPersonData(UID)
	if err != nil {
		return
	}
	log.Tracef("Person: %+v", Person)
	_currentPerson = &Person
	SendButtonLEDOn(true, true, true, true)

	_scanTimer.Reset(5 * time.Second)
}

func timerExpired() {
	log.Trace("Scan Timer Expired")
	SendButtonLEDOn(false, false, false, false)
	_lastUID = ""
	_currentPerson = nil
}

func cardButtonPress() {
	log.Tracef("Card Button Pressed")
	log.Tracef("Person: %v LastUID:%v", _currentPerson, _lastUID)
	if _currentPerson != nil && _currentPerson.CanAdd && _lastUID != _currentPerson.UID && _lastUID != "" {
		_database.AddFriend(_currentPerson.UID, _lastUID)
	}
}
