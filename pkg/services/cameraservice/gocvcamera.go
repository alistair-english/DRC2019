package cameraservice

import (
	"fmt"

	"gocv.io/x/gocv"
)

// GoCVCamera is a camera implementation that reads from GoCV device 0
type GoCVCamera struct{}

// RunCameraConnection from camera Implementation
func (g *GoCVCamera) RunCameraConnection(imgRequests <-chan GetImageActionReq) {
	cam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
	}

	for req := range imgRequests {
		cam.Read(req.Img)
		select {
		case req.ResponseChannel <- true:
		default:
		}
	}
}
