package serial

// Control contains dir and spd
type Control struct {
	Dir int8 // Direction = -90 -> 90
	Spd int8 // Speed = -100 -> 100
}

// PowerRequest is the structure for ESP -> Pi power requests
type PowerRequest struct {
	ReqType uint8
}

// PowerReqResponse is the structure for the Pi's reponse to power requests
type PowerReqResponse struct {
	Accept bool
}

// LogMessage contains a log message for ESP to log to webserver
type LogMessage struct {
	Msg string
}

// Serial info configs, we can put this in a pkg later to make neat
const (
	SerialHeaderSize = 4   // Yeet
	SerialSync1      = 255 // Filler
	SerialSync2      = 255 // Filler
)

// MsgType contains all serial message types
type MsgType int

// Regular Msg
const (
	Data       MsgType = 0
	PowerReq   MsgType = 1
	PowerConf  MsgType = 2
	PowerDeny  MsgType = 3
	ForceReset MsgType = 4
	ForceStop  MsgType = 5
	LogMsg     MsgType = 6
)

// Header is the header information for the serial comms
type Header struct {
	Sync1 uint8
	Sync2 uint8
	Type  uint8
	Size  uint8
}
