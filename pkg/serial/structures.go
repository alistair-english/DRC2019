package serial

// Control contains dir and spd
type Control struct {
	Dir int8 // Direction = -90 -> 90
	Spd int8 // Speed = -100 -> 100
}

// PowerReq is the structure for ESP -> Pi power requests
type PowerReq struct {
	reqType uint8
}

// PowerReqResponse is the structure for the Pi's reponse to power requests
type PowerReqResponse struct {
	accept bool
}

// LogMsg contains a log message for ESP to log to webserver
type LogMsg struct {
	msg string
}

// Serial info configs, we can put this in a pkg later to make neat
const (
	SerialHeaderSize = 32  // Filler value
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
