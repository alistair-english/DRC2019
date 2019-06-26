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

// NewGoCVCamera creates a new camera service that connects to a GoCVCamera
func NewGoCVCamera() (*CameraService, error) {
	return &CameraService{make(chan GetImageActionReq, 50), &GoCVCamera{}}, nil
}

// NewPiCamera creates a new camera service that connects to a PiCamera
func NewPiCamera() (*CameraService, error) {
	cam, err := NewPiCameraImplementation()
	if err != nil {
		return nil, err
	}

	return &CameraService{make(chan GetImageActionReq, 50), cam}, nil
}

// NewFileReaderCamera creates a new camera service that connects to a FileReaderCamera
func NewFileReaderCamera(file string) (*CameraService, error) {
	cam, err := NewFileReaderCameraImplementation(file)
	if err != nil {
		return nil, err
	}

	return &CameraService{make(chan GetImageActionReq, 50), cam}, nil
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
