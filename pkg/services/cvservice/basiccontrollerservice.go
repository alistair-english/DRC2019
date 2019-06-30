package cvservice

import (
	"fmt"
	"image"
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"gocv.io/x/gocv"
)

// BasicControllerService provides recording service
type BasicControllerService struct {
	actionRequestChannel chan<- arch.ActionRequest
}

// GetActionRequestType from Service interface
func (c *BasicControllerService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (c *BasicControllerService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	c.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (c *BasicControllerService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

// Start from Service interface - provides main functionality
func (c *BasicControllerService) Start() {
	go func() {
		// Image Mats
		var sourceImg = gocv.NewMat()
		defer sourceImg.Close()
		var hsvImg = gocv.NewMat()
		defer hsvImg.Close()

		// Create objects
		objects := getObjectsFromConfig()

		// displayWindow := gocv.NewWindow("Display")
		// defer displayWindow.Close()

		controller := newBasicDriveController()

		for { // inifinte loop

			getImgTime := time.Now()
			// get the image
			getImgBlocking(c.actionRequestChannel, &sourceImg)
			fmt.Println(time.Since(getImgTime))

			processTime := time.Now()
			// blur the image
			gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

			// convert to HSV
			gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHSV)

			// Find the HSV objects in the image
			result := cvhelpers.FindHSVObjects(hsvImg, objects)

			found := make(map[string]cvhelpers.HSVObjectGroupResult)

			for _, obj := range result {
				found[obj.Name] = obj
			}

			control := controller.update(found)

			fmt.Println(time.Since(processTime))

			fmt.Println(control)
			fmt.Println()

			c.actionRequestChannel <- serialservice.SerialSendActionReq{control}
		}
	}()
}
