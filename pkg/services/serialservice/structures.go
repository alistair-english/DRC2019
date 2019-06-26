package serialservice

// MsgType contains all serial message types
type MsgType uint8

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

// Control contains dir and spd
type Control struct {
	Dir int8 // Direction = -90 -> 90
	Spd int8 // Speed = -100 -> 100
}

// PowerRequest is the structure for ESP -> Pi power requests
type PowerRequest struct {
	ReqType MsgType
}

// PowerReqResponse is the structure for the Pi's reponse to power requests
type PowerReqResponse struct {
	Accept bool
}

// LogMessage contains a log message for ESP to log to webserver
type LogMessage struct {
	Msg string
}
