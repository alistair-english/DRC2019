package cvservice

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"gocv.io/x/gocv"
)

// RecorderService provides recording service
type RecorderService struct {
	actionRequestChannel chan<- arch.ActionRequest
	recordFolderPath     string
}

// NewRecorderService creates a new recording service with an output folder path
func NewRecorderService(recordFolderPath string) *RecorderService {
	return &RecorderService{nil, recordFolderPath}
}

// GetActionRequestType from Service interface
func (r *RecorderService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (r *RecorderService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	r.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (r *RecorderService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

// Start from Service interface - provides main functionality
func (r *RecorderService) Start() {
	go func() {
		camImg := gocv.NewMat()
		timingChannel := make(chan bool, 1)
		getImgBlocking(r.actionRequestChannel, &camImg)
		fname := fmt.Sprintf(r.recordFolderPath+"recording_%v.avi", time.Now().Format("01-02-2006_15:04:05"))
		out, _ := gocv.VideoWriterFile(fname, "XVID", 5, camImg.Size()[1], camImg.Size()[0], true)

		r.actionRequestChannel <- cameraservice.GetImageActionReq{&camImg, timingChannel}

		for range timingChannel {
			go func() {
				out.Write(camImg)
				fmt.Println("frame written")
			}()
			r.actionRequestChannel <- cameraservice.GetImageActionReq{&camImg, timingChannel}
		}
	}()
}
