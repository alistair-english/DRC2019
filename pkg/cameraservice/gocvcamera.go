package cameraservice

import "gocv.io/x/gocv"

// GoCVCamera is a camera implementation that reads from GoCV device 0
type GoCVCamera struct{}

// RunCameraConnection from camera Implementation
func (g *GoCVCamera) RunCameraConnection(imgRequests <-chan GetImageActionReq) {
	cam, _ := gocv.VideoCaptureDevice(0)
	for req := range imgRequests {
		cam.Read(req.img)
		req.responseChannel <- true
	}
}
