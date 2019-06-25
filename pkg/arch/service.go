package arch

import "go/types"

type Service interface {
	Start()
	GetActionRequestType() types.Type
	SetActionRequestChannel(chan<- ActionRequest)
}
