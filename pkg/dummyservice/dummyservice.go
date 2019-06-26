package dummyservice

import (
	"fmt"
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/cameraservice"
	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/serialservice"
)

type DummyServiceA struct {
	actionRequestChannel chan<- arch.ActionRequest
}

func NewDummyServiceA() *DummyServiceA {
	return &DummyServiceA{nil}
}

func (d *DummyServiceA) Start() {
	go func() {
		chann := make(chan bool, 1)
		img := gocv.NewMat()
		defer img.Close()

		displayWindow := gocv.NewWindow("Display")
		defer displayWindow.Close()

		for {
			d.actionRequestChannel <- cameraservice.GetImageActionReq{&img, chann}
			<-chann
			displayWindow.IMShow(img)
			displayWindow.WaitKey(1)
			d.actionRequestChannel <- serialservice.SerialSendActionReq{serialservice.Control{50, 50}}
		}
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
