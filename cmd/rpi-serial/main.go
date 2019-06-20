package main

import (
	"github.com/alistair-english/DRC2019/internal/pkg/config"
	"github.com/tarm/serial"
	"log"
	"os"
	t "time"
)

const (
	serialPacketSize = 100 // Filler value
	serialSync1      = 255 // Filler
	serialSync2      = 255 // Filler
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

	// Get the port from config
	camConfig := config.GetCameraConfig()

	// Setup serial options
	log.Printf("Starting rpi-serial")
	options := serial.Config{
		Name:        camConfig.Port,     // Just a testing serial port
		Baud:        115200,             // Probs the wrong baud rate
		ReadTimeout: t.Millisecond * 10, // Filler Value
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

	writeChan := make(chan []byte) // Not sure if we should defer close or close these?
	readChan := make(chan []byte)

	go func(writeChannel chan []byte) {
		// Iterate over the channel looking for new stuff to shoot out over serial
		for v := range writeChannel {
			_, err := port.Write(v)
			if err != nil {
				log.Fatal(err)
			}
			// Debug Delay
			t.Sleep(1000 * t.Millisecond)
		}
	}(writeChan)

	go func(readChannel chan []byte) {
		for {
			buf := make([]byte, 255)
			n, _ := port.Read(buf)
			if n == 0 {
				// Failed to read from port in 10ms
			} else if n >= serialPacketSize {
				// We potentially have enought data for serial
				// Now check for serial sync characters
				for i := n; i >= serialPacketSize; i-- {
					if buf[i] == serialSync1 && buf[i+1] == serialSync2 {
						// We are synced up and have an entire packet
						dataPacket := make([]byte, serialPacketSize-2)
						// Copy the data without sync to the dataPacket buffer
						n := copy(dataPacket, buf[i+2:i+serialPacketSize+1]) //If we are truncating data, this will be the issue
						if n != serialPacketSize {
							log.Fatal("Somehow we lost count of our buffer")
						}
						readChannel <- dataPacket
					}
				}
			}
		}
	}(readChan)
}
