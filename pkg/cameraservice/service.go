package cameraservice

import (
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/arch"
)

// CameraService provides a camera connection
type CameraService struct {
	imageRequests        chan GetImageActionReq
	cameraImplementation Implementation
}

// Start from Service interface
func (c *CameraService) Start() {
	go c.cameraImplementation.RunCameraConnection(c.imageRequests)
}

// GetActionRequestType from Service interface
func (c *CameraService) GetActionRequestType() reflect.Type {
	return reflect.TypeOf(GetImageActionReq{})
}

// SetActionRequestChannel from Service interface
func (c *CameraService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	// Don't need to send action requests
}

// FulfullActionRequest from Service interface
func (c *CameraService) FulfullActionRequest(request arch.ActionRequest) {
	c.imageRequests <- request.(GetImageActionReq)
}
