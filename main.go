package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var (
	_buildDate        string
	_buildVersion     string
	_homepageTemplate *template.Template
	log               = logrus.New()
	mqttClient        mqtt.Client
)

func main() {
	var err error
	log.SetLevel(logrus.TraceLevel)
	log.Printf("---------- Program Started %v (%v) ----------", _buildVersion, _buildDate)

	err = loadPageTemplates()
	if err != nil {
		log.Panic(err)
	}

	item, err := getIPAddress("eth0")
	if err != nil {
		log.Panic(err)
	}
	log.Tracef("Interface:%+v", item)

	go udpServer([]byte("Trailer Server"))
	go httpServer()

	connectToMQTT()

	select {}
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

func httpServer() {
	var err error
	http.HandleFunc("/", homepage)
	http.HandleFunc("/api/setRelay", setRelayHandler)

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
			ipNet = item
		}
	}
	return
}

func setRelayHandler(w http.ResponseWriter, r *http.Request) {
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
		case "FALSE":
			log.Info("Set Relay ", Relay, " to ", StateString)
		}
	}
}
