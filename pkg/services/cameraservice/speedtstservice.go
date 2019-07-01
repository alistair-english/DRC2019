package cameraservice

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/arch"

	"gocv.io/x/gocv"
)

// CameraSpeedTestService provides camera speed test service
type CameraSpeedTestService struct {
	actionRequestChannel chan<- arch.ActionRequest
}

// GetActionRequestType from Service interface
func (c *CameraSpeedTestService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (c *CameraSpeedTestService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	c.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (c *CameraSpeedTestService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

// Start from Service interface - provides main functionality
func (c *CameraSpeedTestService) Start() {
	go func() {
		camImg := gocv.NewMat()
		timingChannel := make(chan bool, 1)

		c.actionRequestChannel <- GetImageActionReq{&camImg, timingChannel}

		t := time.Now()
		for range timingChannel {
			go func() {
				fmt.Println("frame: ", time.Since(t))
				t = time.Now()
			}()
			c.actionRequestChannel <- GetImageActionReq{&camImg, timingChannel}
		}
	}()
}
