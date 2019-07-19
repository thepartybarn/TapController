package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
	"github.com/tarm/serial"
)

var (
	_buildDate        string
	_buildVersion     string
	_homepageTemplate *template.Template
	log               = logrus.New()
	mqttClient        mqtt.Client

	_taps     map[int]*tapStruct
	_database Database

	_scanTimer *time.Timer

	_lastCardID    string
	_lastPerson    *Person
	_authenticated bool
)

type usbDataMessage struct {
	Key   string `json:"key"`
	Value string `json:"Value"`
}

func main() {
	var err error
	log.SetLevel(logrus.TraceLevel)
	log.Printf("---------- Program Started %v (%v) ----------", _buildVersion, _buildDate)
	/*
		err = loadPageTemplates()
		if err != nil {
			log.Panic(err)
		}
		//eth0IP
		_, err = getIPAddress("eth0")
		if err != nil {
			log.Error(err)
		}
		//wlan0IP
		_, err = getIPAddress("wlan0")
		if err != nil {
			log.Error(err)
		}*/

	err, _taps = setupTapRelays()
	if err != nil {
		log.Panic(err)
	}
	defer rpio.Close()

	/*
		go udpServer([]byte("Trailer Server"))
		go httpServer()
	*/
	go handleUSBDevice("/dev/ttyUSB0")
	go handleUSBDevice("/dev/ttyUSB1")

	//connectToMQTT()

	_database, err = CreateDatabase()
	if err != nil {
		log.Panic(err)
	}
	_database.AddAdmin("04ed19ea9c6180")

	select {
	case <-_scanTimer.C:
		log.Trace("Timer Expired")
		_authenticated = false
		_lastCardID = ""
		_lastPerson = nil
	}
}

func handleUSBDevice(device string) {
	var serialBuffer bytes.Buffer
	var serialPort *serial.Port

	serialConfig := &serial.Config{Name: device, Baud: 9600}

	var err error
	//Serial Port Open Loop
	for {
		serialPort, err = serial.OpenPort(serialConfig)
		if err != nil {
			log.Error("USB Serial Port Not Found")
			time.Sleep(time.Second)
		}
		defer serialPort.Close()

		temp := make([]byte, 256)
		var MessagePayload []byte
		var nRead int

		/*TODO handle usb ports better
		defer func() {
			if r := recover(); r != nil {
				log.Info("Recovered in f", r)
			}
		}()
		*/
		//Read Loop
		for {
			nRead, err = serialPort.Read(temp)
			if err != nil {
				log.Warn("USB Read Error:", err)
				break
			}
			if nRead > 0 {
				serialBuffer.Write(temp[:nRead])
			}
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
			//Process Message
			if len(MessagePayload) > 3 {
				log.Trace("USB Message Received: ", string(MessagePayload))
				var Message usbDataMessage
				err := json.Unmarshal(MessagePayload, &Message)
				MessagePayload = nil
				if err != nil {
					log.Warn("USB Marshal Error:", err)
				} else {
					go handleUSBMessage(Message)
				}
			}
		}
	}
	log.Warn("USB Serial Loop Exited")
	serialPort.Close()
}
func handleUSBMessage(Message usbDataMessage) {
	switch Message.Key {
	case "cardID":
		//Convert negative numbers to positive
		cardScan(Message.Value)
	case "buttonpress":
		cardButtonPress()
	case "tap1":
		state, err := strconv.ParseBool(Message.Value)
		if err != nil {
			log.Warn(err)
		}
		tapButtonPress(1, state)
	case "tap2":
		state, err := strconv.ParseBool(Message.Value)
		if err != nil {
			log.Warn(err)
		}
		tapButtonPress(2, state)
	case "tap3":
		state, err := strconv.ParseBool(Message.Value)
		if err != nil {
			log.Warn(err)
		}
		tapButtonPress(3, state)
	case "tap4":
		state, err := strconv.ParseBool(Message.Value)
		if err != nil {
			log.Warn(err)
		}
		tapButtonPress(4, state)
	}
}
func cardScan(cardID string) {
	log.Tracef("Card Scanned %v", cardID)
	_lastCardID = cardID
	if Person, ok := _database.hasUID(cardID); ok {
		_lastPerson = &Person
		_authenticated = true
		_scanTimer = time.NewTimer(4 * time.Second)
	}
}
func tapButtonPress(index int, state bool) {
	log.Tracef("Tap Button %v Pressed %v", index, state)
	//if state == true && _authenticated == true {
	if state == true {
		_taps[index].Open()
	}
	if state == false {
		_taps[index].Close()
	}
}
func cardButtonPress() {
	log.Tracef("Card Button Pressed")
	log.Tracef("Authenticated:%v CardID:%v (%v)", _authenticated, _lastCardID, _lastPerson)
	if _authenticated && _lastPerson.canAdd && _lastCardID != "" {
		_database.AddFriend(_lastCardID)
	}
}

func connectToMQTT() error {
	ServerAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:10001")
	if err != nil {
		return err
	}
	LocalAddr, err := net.ResolveUDPAddr("udp", ":10002")
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", LocalAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	buf := make([]byte, 1024)

	n, err := udpConn.WriteTo([]byte{0x01}, ServerAddr)
	if err != nil {
		return err
	}
	log.Tracef("packet-written: bytes=%d to=%s\n", n, ServerAddr.String())

	n, addr, err := udpConn.ReadFromUDP(buf)
	log.Trace("Received ", string(buf[0:n]), " from ", addr)
	if err != nil {
		return err
	}

	BrokerAddr := fmt.Sprintf("%v%v", addr.IP, ":1883")
	log.Trace("Broker ", BrokerAddr)
	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.AddBroker(BrokerAddr)
	mqttClient = mqtt.NewClient(mqttClientOptions)
	token := mqttClient.Connect()
	for token.Wait() && token.Error() != nil && mqttClient.IsConnected() == false {
		time.Sleep(2 * time.Second)
		log.Error("Trying to Connect to Broker", BrokerAddr)
		token = mqttClient.Connect()
	}
	log.Info("Connected to Broker:", BrokerAddr)

	return nil
}
func udpServer(dataToSend []byte) {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		log.Error(err)
		return
	}
	udpConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Error(err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := udpConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		deadline := time.Now().Add(time.Second)
		err = udpConn.SetWriteDeadline(deadline)
		if err != nil {
			log.Error(err)
			return
		}
		// Write the packet's contents back to the client.
		n, err = udpConn.WriteTo(dataToSend, addr)
		if err != nil {
			log.Error(err)
		}
		fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
	}
}
func getIPAddress(name string) (ipNet *net.IPNet, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return
	}
	for _, a := range addrs {
		item, ok := a.(*net.IPNet)
		if ok && !item.IP.IsLoopback() && item.IP.To4() != nil {
			log.Tracef("%v: %v", name, item.String())
			ipNet = item
		}
	}

	return
}

//HTTP STUFF
//Mike is a nerd
func httpServer() {
	var err error
	http.HandleFunc("/", homepage)
	//http.HandleFunc("/api/setRelay", setRelayHandler)

	//Serve Static Files
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("web/images/"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("web/scripts/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("web/styles/"))))

	log.Trace("Opening HTTP Server")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}
func homepage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := make(map[string]interface{})
	data["BuildDate"] = "UnSet"
	data["BuildVersion"] = "UnSet"

	err := _homepageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
func loadPageTemplates() error {
	var err error
	_homepageTemplate, err = template.ParseFiles("web/index.html")
	return err
}

/*func setRelayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	RelayString := r.URL.Query().Get("Relay")
	StateString := strings.ToUpper(r.URL.Query().Get("State"))

	log.Trace("Passed In Relay:", RelayString, " State:", StateString)

	Relay, err := strconv.Atoi(RelayString)
	if err != nil {
		return
	}
	if Relay >= 1 && Relay <= 8 {
		switch StateString {
		case "TRUE":
			log.Info("Set Relay ", Relay, " to ", StateString)
			_relays[Relay-1].Low()
		case "FALSE":
			log.Info("Set Relay ", Relay, " to ", StateString)
			_relays[Relay-1].High()
		}
	}
}
*/
