package cvservice

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"
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

		t := time.Now()

		counter := 0

		controller := newBasicDriveController()

		for { // inifinte loop

			// start := time.Now()
			// get the image
			getImgBlocking(c.actionRequestChannel, &sourceImg)
			// fmt.Println("Img: ", time.Since(start))

			// start := time.Now()
			// blur the image
			// gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

			// convert to HSV
			gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)
			// fmt.Println(time.Since(start))

			// start := time.Now()
			// Find the HSV objects in the image
			result := cvhelpers.FindHSVObjects(hsvImg, objects)
			// fmt.Println("Calc: ", time.Since(start))

			// for _, group := range result {
			// 	if len(group.Objects) > 0 {
			// 		gocv.Rectangle(&sourceImg, group.Objects[0].BoundingBox, color.RGBA{255, 0, 0, 0}, 3)
			// 	}
			// }

			// displayWindow.IMShow(sourceImg)
			// displayWindow.WaitKey(0)

			// start := time.Now()
			control := controller.update(result)

			fmt.Println(control)
			fmt.Println()

			c.actionRequestChannel <- serialservice.SerialSendActionReq{control}
			counter++
			if time.Since(t) > time.Second {
				fmt.Println(counter)
				counter = 0
				t = time.Now()
			}
			// if time.Since(t) > time.Millisecond*100 {

			// 	t = time.Now()
			// 	fmt.Println("sent.")
			// }
			// fmt.Println(time.Since(start))
		}
	}()
}
