package main

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"time"
)

func main() {
	// Setup serial options
	log.Printf("Starting rpi-serial.go")
	options := serial.OpenOptions{
		PortName:        "/dev/yeetPort",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	// Init the Serial port
	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	for {
		b := []byte{0x00, 0x01, 0x02, 0x03}
		_, err := port.Write(b)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Holy yeet its looping")
		time.Sleep(1000 * time.Millisecond)
	}
}
