package arch

import "reflect"

type Service interface {
	Start()
	GetActionRequestType() reflect.Type
	SetActionRequestChannel(channel chan<- ActionRequest)
	FulfullActionRequest(request ActionRequest)
}
