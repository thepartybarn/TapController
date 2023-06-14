package main

import (
	"html/template"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var (
	_buildDate        string
	_buildVersion     string
	_homepageTemplate *template.Template
	log               = logrus.New()
	mqttClient        mqtt.Client

	_database *DatabaseConnection

	_scanTimer *time.Timer

	_heartbeatTicker *time.Ticker
)

func main() {
	var err error
	log.SetLevel(logrus.TraceLevel)
	log.Printf("---------- Program Started %v (%v) ----------", _buildVersion, _buildDate)

	go handleTapRelays()

	_database, err = SetupDatabaseConnections(log)
	if err != nil {
		log.Panic(err)
	}

	OpenUSBDevice("/dev/ttyUSB0", log)
	OpenUSBDevice("/dev/ttyUSB1", log)

	_scanTimer = time.NewTimer(5 * time.Second)
	_scanTimer.Stop()
	_heartbeatTicker = time.NewTicker(3 * time.Second)

	for {
		select {
		case <-_scanTimer.C:
			timerExpired()
		case <-_heartbeatTicker.C:
			//USBHeartbeat()
		}
	}
}

/*
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
*/
