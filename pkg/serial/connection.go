package serial

import (
	"github.com/alistair-english/DRC2019/pkg/gohelpers"
)

// Connection e
type Connection struct {
	ControlChan          chan Control
	PowerReqChan         chan PowerRequest
	LogMsgChan           chan LogMessage
	PowerReqResponseChan chan PowerReqResponse
	writeChan            chan []byte
	readChan             chan []byte
	serial               Implementation
}

// NewConnection makes a new connection from a serial implementation
func NewConnection(serial Implementation) (*Connection, error) {
	// External Channels
	controlChan := make(chan Control, 20)
	powerReqChan := make(chan PowerRequest, 20)
	logMSgChan := make(chan LogMessage, 20)
	powerReqResponse := make(chan PowerReqResponse, 20)

	// Internal Channels
	writeChan := make(chan []byte, 20)
	readChan := make(chan []byte, 20)

	c := Connection{controlChan, powerReqChan, logMSgChan, powerReqResponse, writeChan, readChan, serial}

	go c.HandleChannels()

	return &c, nil
}

// Init inits all the serial things
func (c Connection) Init() error {
	err := c.serial.Init()
	if err != nil {
		return err
	}
	go c.serial.RunSerialTx(c.writeChan)
	go c.serial.RunSerialRx(c.readChan)
	return nil
}

// HandleChannels handles channels
func (c Connection) HandleChannels() {
	// Goroutine for handling the control channel
	go func(in <-chan Control, out chan<- []byte) {
		for con := range in {
			buf := make([]byte, SerialHeaderSize+2)
			buf[0] = byte(0xB5)
			buf[1] = byte(0x62)
			buf[2] = byte(Data)
			buf[3] = byte(2)
			buf[4] = byte(con.Dir)
			buf[5] = byte(con.Spd)

			out <- buf
		}
	}(c.ControlChan, c.writeChan)

	// Goroutine for handling the power requests from ESP
	go func(in <-chan []byte, out chan<- PowerRequest) {
		for con := range in {
			var req PowerRequest
			req.ReqType = MsgType(con[SerialHeaderSize])

			out <- req
		}
	}(c.readChan, c.PowerReqChan)

	// Goroutine for handling log messages
	go func(in <-chan LogMessage, out chan<- []byte) {
		for con := range in {
			buf := make([]byte, SerialHeaderSize+len(con.Msg))
			buf[0] = byte(0xB5)
			buf[1] = byte(0x62)
			buf[2] = byte(LogMsg)
			buf[3] = byte(len(con.Msg))
			msgBytes := []byte(con.Msg)
			copy(buf[4:], msgBytes) //This might need to be 3 not 4

			out <- buf
		}
	}(c.LogMsgChan, c.writeChan)

	// Goroutine for handling power request responses
	go func(in <-chan PowerReqResponse, out chan<- []byte) {
		for con := range in {
			buf := make([]byte, SerialHeaderSize+1)
			buf[0] = byte(0xB5)
			buf[1] = byte(0x62)
			buf[2] = byte(PowerReq)
			buf[3] = byte(1)
			accept := gohelpers.B2i(con.Accept)
			buf[4] = byte(accept)

			out <- buf
		}
	}(c.PowerReqResponseChan, c.readChan)
}
