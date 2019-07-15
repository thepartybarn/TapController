package main

import (
	//	"bufio"
	"bytes"
	"fmt"
	"github.com/karalabe/hid"
	//	"os"
)

func main1() {
	fmt.Println("Hello World")

	go handleRfidReader()

	select {}
}

func handleRfidReader() {
	for {
		devices := hid.Enumerate(5824, 10203)
		if len(devices) > 0 {
			fmt.Printf("%+v\n", devices[0])
			device, err := devices[0].Open()
			if err != nil {
				fmt.Println(err)
				continue
			}
			tmp := make([]byte, 100)
			var cardIDBuffer bytes.Buffer
			for {
				nRead, err := device.Read(tmp)
				if err != nil {
					fmt.Println(err)
					break
				}
				if nRead > 0 {
					readValue := tmp[2]
					switch readValue {
					case 0:
					case 40:
						cardID := make([]byte, 10)
						cardIDBuffer.Read(cardID)
						fmt.Printf("%v\n", cardID)
					default:
						value := byte(tmp[2]) - 29
						if value == 10 {
							value = 0
						}
						cardIDBuffer.WriteByte(value)
						if cardIDBuffer.Len() > 10 {
							cardIDBuffer.ReadByte()
						}
					}
				}
			}
			device.Close()
		}
	}
}
