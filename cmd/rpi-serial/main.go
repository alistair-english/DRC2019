package main

import (
	"log"
	"os"
	t "time"

	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/tarm/serial"
)

// Serial info configs, we can put this in a pkg later to make neat
const (
	serialHeaderSize = 32  // Filler value
	serialSync1      = 255 // Filler
	serialSync2      = 255 // Filler
)

// MsgType contains all serial message types
type MsgType int

// Enum Eq for messages
const (
	Message    MsgType = 0
	PowerReq   MsgType = 1
	PowerConf  MsgType = 2
	PowerDeny  MsgType = 3
	ForceReset MsgType = 4
	ForceStop  MsgType = 5
)

// SerialHeader is the header information for the serial comms
type SerialHeader struct {
	Sync1 uint8
	Sync2 uint8
	Type  uint8
	Size  uint8
}

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

	// Write goroutine to send serial data
	go func(writeChannel <-chan []byte) {
		// Iterate over the channel looking for new stuff to shoot out over serial
		for v := range writeChannel {
			_, err := port.Write(v)
			if err != nil {
				log.Fatal(err)
			}
			// Debug Delay yeet delet this when actually doing stuff
			t.Sleep(1000 * t.Millisecond)
		}
	}(writeChan)

	// Read goroutine to recieve serial data
	go func(readChannel chan<- []byte) {
		for {
			buf := make([]byte, 255)
			n, _ := port.Read(buf)
			if n == 0 {
				// Failed to read from port in 10ms
			} else if n >= serialHeaderSize {
				// We potentially have enought data for serial
				// Now check for serial sync characters
				for i := n; i >= serialHeaderSize; i-- {
					if uint8(buf[i]) == serialSync1 && uint8(buf[i+1]) == serialSync2 {
						// We are synced up and have an entire packet
						dataPacket := make([]byte, serialHeaderSize)
						// Copy the data without sync to the dataPacket buffer
						n := copy(dataPacket, buf[i:i+serialHeaderSize]) //If we are truncating data, this will be the issue
						if n != serialHeaderSize {
							log.Fatal("Somehow we lost count of our buffer")
						}
						// We have a serial header, decode it and read in the memory
						headerPacket := decodeHeaderPacket(dataPacket)
						if headerPacket.Size == 0 {
							log.Fatal("No data in packet")
						}
						packetBuf := make([]byte, headerPacket.Size)
						nn, _ := port.Read(packetBuf)
						if uint8(nn) < headerPacket.Size {
							log.Fatal("We somehow expected more data than we got")
						}
						// Place the raw data on a channel to be handled
						readChannel <- packetBuf
					}
				}
			}
		}
	}(readChan)

	// Goroutine to parse instructions from ESP
	go func(readChannel <-chan []byte) {
		// TODO: Read in the instructions from ESP
		for v := range readChannel {
			// TODO
			log.Println(v) // Filler to get rid of errors
		}
	}(readChan)

	// Debug testing of writeChannel
	for {
		writeChan <- []byte{0x01}
	}
}
