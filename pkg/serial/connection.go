package serial

// Connection e
type Connection struct {
	ControlChan          chan Control
	PowerReqChan         chan PowerReq
	LogMsgChan           chan LogMsg
	PowerReqResponseChan chan PowerReqResponse
	writeChan            chan []byte
	readChan             chan []byte
	serial               Implementation
}

// NewConnection makes a new connection from a serial implementation
func NewConnection(serial Implementation) (*Connection, error) {
	// External Channels
	controlChan := make(chan Control, 20)
	powerReqChan := make(chan PowerReq, 20)
	logMSgChan := make(chan LogMsg, 20)
	powerReqResponse := make(chan PowerReqResponse, 20)

	// Internal Channels
	writeChan := make(chan []byte, 20)
	readChan := make(chan []byte, 20)

	c := Connection{controlChan, powerReqChan, logMSgChan, powerReqResponse, writeChan, readChan, serial}

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

}
