package serialservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

// Serial Info Config
const (
	SerialHeaderSize = 4   // Yeet
	SerialSync1      = 255 // Filler
	SerialSync2      = 255 // Filler
)

// Header is the header information for the serial comms
type Header struct {
	Sync1 uint8
	Sync2 uint8
	Type  uint8
	Size  uint8
}

//Implementation is the interface for serial
type Implementation interface {
	RunSerialTx(writeChan <-chan []byte) error
	RunSerialRx(readChan chan<- []byte) error
}

// FakeSerial is a fake testing serial implementation
type FakeSerial struct {
}

// RunSerialTx prints to console
func (s FakeSerial) RunSerialTx(writeChan <-chan []byte) error {
	// Iterate over the channel looking for new stuff to shoot out over serial
	for v := range writeChan {
		fmt.Println("Serial Write: ", v)
	}
	return nil
}

// RunSerialRx is empty (no need to send anything)
func (s FakeSerial) RunSerialRx(readChan chan<- []byte) error {
	return nil
}

// PiSerial is the serial implementation
type PiSerial struct {
	Port       string
	Baud       int
	Timeout    time.Duration
	connection *serial.Port
}

// RunSerialTx starts serial tx
func (s PiSerial) RunSerialTx(writeChan <-chan []byte) error {
	// Iterate over the channel looking for new stuff to shoot out over serial
	for v := range writeChan {
		_, err := s.connection.Write(v)
		if err != nil {
			return err
		}
		// Debug Delay yeet delet this when actually doing stuff
		//t.Sleep(1000 * t.Millisecond)
	}
	return nil
}

// RunSerialRx starts serial rx
func (s PiSerial) RunSerialRx(readChan chan<- []byte) error {
	for {
		buf := make([]byte, 255)
		n, _ := s.connection.Read(buf)
		if n == 0 {
			// Failed to read from port in 10ms
		} else if n >= SerialHeaderSize {
			// We potentially have enought data for serial
			// Now check for serial sync characters
			for i := n; i >= SerialHeaderSize; i-- {
				if uint8(buf[i]) == SerialSync1 && uint8(buf[i+1]) == SerialSync2 {
					// We are synced up and have an entire packet
					dataPacket := make([]byte, SerialHeaderSize)
					// Copy the data without sync to the dataPacket buffer
					n := copy(dataPacket, buf[i:i+SerialHeaderSize]) //If we are truncating data, this will be the issue
					if n != SerialHeaderSize {
						return errors.New("Somehow we lost count of our buffer")
					}
					// We have a serial header, decode it and read in the memory
					headerPacket := decodeHeaderPacket(dataPacket)
					if headerPacket.Size == 0 {
						return errors.New("No data in packet")
					}
					packetBuf := make([]byte, headerPacket.Size)
					nn, _ := s.connection.Read(packetBuf)
					if uint8(nn) < headerPacket.Size {
						return errors.New("We somehow expected more data than we got")
					}
					// Place the raw data on a channel to be handled
					readChan <- packetBuf
				}
			}
		}
	}
}

func decodeHeaderPacket(arr []byte) Header {
	var packet Header

	// Decode a packet header into the SerialHeader datatype
	packet.Sync1 = arr[0]
	packet.Sync2 = arr[1]
	packet.Type = arr[2]
	packet.Size = arr[3]

	return packet
}
