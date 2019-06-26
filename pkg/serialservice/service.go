package serialservice

import (
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/gohelpers"
	"github.com/tarm/serial"
)

// SerialService provides serial connection
type SerialService struct {
	txChannel            chan []byte
	rxChannel            chan []byte
	serialImplementation Implementation
}

// NewPiSerial creates a new PiSerial object
func NewPiSerial(port string, baud int, timeout time.Duration) (*SerialService, error) {
	// Setup serial options
	options := serial.Config{
		Name:        port,
		Baud:        baud,
		ReadTimeout: timeout,
		Size:        8,
		StopBits:    1,
		Parity:      'N',
	}
	// Init the Serial port
	conn, err := serial.OpenPort(&options)
	if err != nil {
		return nil, err
	}

	serialImplementation := &PiSerial{port, baud, timeout, conn}

	txChannel := make(chan []byte, 50)
	rxChannel := make(chan []byte, 50)

	return &SerialService{txChannel, rxChannel, serialImplementation}, nil
}

func NewFakeSerial() (*SerialService, error) {
	return &SerialService{make(chan []byte, 50), make(chan []byte, 50), FakeSerial{}}, nil
}

// Start from Service interface
func (s *SerialService) Start() {
	go s.serialImplementation.RunSerialTx(s.txChannel)
	go s.serialImplementation.RunSerialRx(s.rxChannel)
}

// GetActionRequestType from Service interface
func (s *SerialService) GetActionRequestType() reflect.Type {
	return reflect.TypeOf(SerialSendActionReq{})
}

// SetActionRequestChannel from Service interface
func (s *SerialService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	// Not implemented yet - only writing atm
}

// FulfullActionRequest from Service interface
func (s *SerialService) FulfullActionRequest(request arch.ActionRequest) {
	structure := request.(SerialSendActionReq).SerialStructure

	switch structure.(type) {
	case Control:
		structure := structure.(Control)
		buf := make([]byte, SerialHeaderSize+2)
		buf[0] = byte(0xB5)
		buf[1] = byte(0x62)
		buf[2] = byte(Data)
		buf[3] = byte(2)
		buf[4] = byte(structure.Dir)
		buf[5] = byte(structure.Spd)

		s.txChannel <- buf
	case LogMessage:
		structure := structure.(LogMessage)
		buf := make([]byte, SerialHeaderSize+len(structure.Msg))
		buf[0] = byte(0xB5)
		buf[1] = byte(0x62)
		buf[2] = byte(LogMsg)
		buf[3] = byte(len(structure.Msg))
		msgBytes := []byte(structure.Msg)
		copy(buf[4:], msgBytes) //This might need to be 3 not 4

		s.txChannel <- buf
	case PowerReqResponse:
		structure := structure.(PowerReqResponse)
		buf := make([]byte, SerialHeaderSize+1)
		buf[0] = byte(0xB5)
		buf[1] = byte(0x62)
		buf[2] = byte(PowerReq)
		buf[3] = byte(1)
		accept := gohelpers.B2i(structure.Accept)
		buf[4] = byte(accept)
	}
}
