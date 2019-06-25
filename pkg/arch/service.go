package arch

import "reflect"

type Service interface {
	Start()
	GetActionRequestType() reflect.Type
	SetActionRequestChannel(chan<- ActionRequest)
	FulfullActionRequest(ActionRequest)
}
