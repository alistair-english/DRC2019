package main

import (
	"github.com/alistair-english/DRC2019/pkg/logginghelpers"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
)

func main() {
	// Setup a stream for file logging
	f, err := os.OpenFile("../../logs/rpi-serial.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Set the log to the file stream and write a file divider
	log.SetOutput(f)
	log.SetFlags(log.Lmicroseconds) // Set the logging mode to show microseconds
	log.Printf("-------NEWLOG-------")

	// Setup serial options
	log.Printf("Starting rpi-serial")
	options := serial.Config{
		Name:        "/dev/tty.usbserial",  // Just a testing serial port
		Baud:        115200,                // Probs the wrong baud rate
		ReadTimeout: time.Millisecond * 10, // Filler Value
		Size:        8,
		StopBits:    1,
		Parity:      'N',
	}
	// Init the Serial port
	port, err := serial.OpenPort(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	go func() {
		for {
			buf := make([]byte, 255)

			n, _ := port.Read(buf)
			if n == 0 {
				log.Printf("Failed to read from port")
			}
		}
	}()

	for {
		b := []byte{0x00, 0x01, 0x02, 0x03}
		_, err := port.Write(b)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Looping")
		time.Sleep(1000 * time.Millisecond)
	}
}
