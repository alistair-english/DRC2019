package seriallogservice

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"
	"reflect"
)

// SerialLogService is the service that logs information to the
// serial port for debugging
type SerialLogService struct {
	actionRequestChannel chan<- arch.ActionRequest
}

// NewSerialLogService returns a new SerialLogService object
func NewSerialLogService() *SerialLogService {
	return &SerialLogService{nil}
}

// Start the service
func (s *SerialLogService) Start() {
	// We dont need to start anything here
}

// LogToSerial logs a message to serial
func (s *SerialLogService) LogToSerial(msg string) {
	s.actionRequestChannel <- serialservice.SerialSendActionReq{SerialStructure: serialservice.LogMessage{Msg: msg}}
}

// GetActionRequestType returns the action request type
func (s *SerialLogService) GetActionRequestType() reflect.Type {
	return nil
}

//SetActionRequestChannel sets the action request channel
func (s *SerialLogService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	s.actionRequestChannel = channel
}

// FulfullActionRequest fullfils an action request but alistair cant spell and we are too far deep to fix it
func (s *SerialLogService) FulfullActionRequest(request arch.ActionRequest) {
	// Do nothing because we arent taking in action requests, just submitting them to the router
}
