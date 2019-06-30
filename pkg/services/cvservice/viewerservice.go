package cvservice

import (
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"gocv.io/x/gocv"
)

// ViewerService provides viewing service
type ViewerService struct {
	actionRequestChannel chan<- arch.ActionRequest
	recordFolderPath     string
}

// GetActionRequestType from Service interface
func (v *ViewerService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (v *ViewerService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	v.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (v *ViewerService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

// Start from Service interface - provides main functionality
func (v *ViewerService) Start() {
	go func() {
		camImg := gocv.NewMat()
		timingChannel := make(chan bool, 1)

		window := gocv.NewWindow("Viewer")
		defer window.Close()

		v.actionRequestChannel <- cameraservice.GetImageActionReq{&camImg, timingChannel}

		for range timingChannel {
			go func() {
				window.IMShow(camImg)
				window.WaitKey(1)
			}()
			v.actionRequestChannel <- cameraservice.GetImageActionReq{&camImg, timingChannel}
		}
	}()
}
