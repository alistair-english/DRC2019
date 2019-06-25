package dummyservice

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/arch"
)

type DummyServiceA struct {
	actionRequestChannel chan<- arch.ActionRequest
}

func NewDummyServiceA() *DummyServiceA {
	return &DummyServiceA{make(chan arch.ActionRequest, 100)}
}

func (d *DummyServiceA) Start() {
	go func() {
		d.actionRequestChannel <- DummyActionRequest{"hello world"}
		time.Sleep(3 * time.Second)
	}()

}

func (d *DummyServiceA) GetActionRequestType() reflect.Type {
	return nil
}

func (d *DummyServiceA) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	d.actionRequestChannel = channel
}

func (d *DummyServiceA) FulfullActionRequest(request arch.ActionRequest) {
	// Not doing anything as not responding to action requests
}

type DummyServiceB struct {
	dummyActionChannel chan DummyActionRequest
}

func NewDummyServiceB() *DummyServiceB {
	return &DummyServiceB{make(chan DummyActionRequest, 100)}
}

func (d *DummyServiceB) Start() {
	go func() {
		for request := range d.dummyActionChannel {
			fmt.Println(request.Message)
		}
	}()
}

func (d *DummyServiceB) GetActionRequestType() reflect.Type {
	return reflect.TypeOf(DummyActionRequest{})
}

func (d *DummyServiceB) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	// Not doing anything as not requesting actions
}

func (d *DummyServiceB) FulfullActionRequest(request arch.ActionRequest) {
	d.dummyActionChannel <- request.(DummyActionRequest)
}
