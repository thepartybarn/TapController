package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var (
	_usbDevices      = make(map[string]*USBDevice)
	_currentLEDColor = LEDColor{Blue: 255}
)

type USBDataMessage struct {
	Key   string `json:"key"`
	Value string `json:"Value"`
}
type USBDevice struct {
	SerialPort *serial.Port
	Path       string
	Connected  bool
}

func USBHeartbeat() {
	SendLEDColor(_currentLEDColor)
}

func SendButtonLEDOn(button1, button2, button3, button4 bool) {
	for _, usbDevice := range _usbDevices {
		usbDevice.sendSetButtonLED(button1, button2, button3, button4)
	}
}
func SendLEDColor(ledColor LEDColor) {
	_currentLEDColor = ledColor
	for _, usbDevice := range _usbDevices {
		usbDevice.sendSetLEDColor(ledColor)
	}
}

func OpenUSBDevice(Path string, log *logrus.Logger) (usbDevice *USBDevice) {
	usbDevice = new(USBDevice)
	usbDevice.Path = Path

	_usbDevices[Path] = usbDevice

	go usbDevice.ProcessingLoop()
	return
}

func (usbDevice *USBDevice) ProcessingLoop() {
	var err error
	log.Tracef("%v Serial Processing Loop", usbDevice.Path)
	defer log.Tracef("%v Exited Serial Processing Loop", usbDevice.Path)
	for {
		usbDevice.SerialPort, err = serial.OpenPort(&serial.Config{Name: usbDevice.Path, Baud: 115200})
		if err != nil {
			if usbDevice.Connected {
				usbDevice.Connected = false
				log.Tracef("%v Disconnected!", usbDevice.Path)

			}
		} else {
			usbDevice.Connected = true
			log.Tracef("%v Connected!", usbDevice.Path)

			usbDevice.SerialDataProcessingLoop()
			if recover := recover(); recover != nil {
				log.Errorf("%v Recovered", usbDevice.Path)
			}
			log.Warnf("%v USB Serial Loop Exited", usbDevice.Path)
			usbDevice.SerialPort.Close()
		}
		time.Sleep(5 * time.Second)
	}
}

func (usbDevice *USBDevice) SerialDataProcessingLoop() {
	log.Trace("Serial Data Processing Loop")
	var err error
	temp := make([]byte, 256)
	var MessagePayload []byte
	var nRead int
	var serialBuffer bytes.Buffer

	//Read / Write Loop
	for {
		nRead, err = usbDevice.SerialPort.Read(temp)
		if err != nil {
			return
		}
		if nRead > 0 {
			serialBuffer.Write(temp[:nRead])
		}
		if serialBuffer.Len() > 0 {
			tempData := serialBuffer.Bytes()
			endOfMessage := bytes.IndexByte(tempData, '}')
			if endOfMessage > 0 {
				beginingOfMessage := bytes.IndexByte(tempData, '{')
				if beginingOfMessage >= 0 {
					MessagePayload = tempData[beginingOfMessage : endOfMessage+1]
				}
				//Clear buffer before }
				serialBuffer.Next(endOfMessage + 1)
			} else {
				//Didn't Find End Keep storing buffer.
				//TODO We should limit this to a particular size
			}
		}
		//Process Message
		if len(MessagePayload) > 3 {
			var Message USBDataMessage
			err = json.Unmarshal(MessagePayload, &Message)
			MessagePayload = nil
			if err != nil {
				log.Warn("USB Marshal Error:", err)
			} else {
				processUSBMessage(Message)
			}
		}
	}
}

func (usbDevice *USBDevice) sendSetButtonLED(button1, button2, button3, button4 bool) error {
	data := make(map[string]interface{})

	data["command"] = "setLEDs"
	value := 0
	if button4 {
		value = value + 1
	}
	if button3 {
		value = value + 10
	}
	if button2 {
		value = value + 100
	}
	if button1 {
		value = value + 1000
	}

	data["mask"] = fmt.Sprintf("%04d", value)

	return usbDevice.sendUSBMessage(data)
}

type LEDColor struct {
	Name  string
	Red   byte
	Green byte
	Blue  byte
}

func (usbDevice *USBDevice) sendSetLEDColor(ledColor LEDColor) error {
	data := make(map[string]interface{})
	data["command"] = "setLEDs"
	data["red"] = ledColor.Red
	data["green"] = ledColor.Green
	data["blue"] = ledColor.Blue

	return usbDevice.sendUSBMessage(data)
}
func (usbDevice *USBDevice) sendUSBMessage(data map[string]interface{}) (err error) {
	dataToSend, err := json.Marshal(data)
	if err != nil {
		return
	}
	if !usbDevice.Connected {
		return
	}
	_, err = usbDevice.SerialPort.Write(dataToSend)
	log.Tracef("Writing to USB(%v) %v", usbDevice.Path, string(dataToSend))
	time.Sleep(100 * time.Millisecond)
	return
}
