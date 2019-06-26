package camera

import (
	"gocv.io/x/gocv"
)

// WebCamera is an implementation of a gocv capture deivce (most likely a webcam)
type WebCamera struct {
	Device *gocv.VideoCapture
}

// Init from the camera Implementation
func (cam *WebCamera) Init() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	cam.Device = webcam
}

// RunImagePoller from the camera Implementation
func (cam *WebCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat) {
	for range imageRequest {
		cam.Device.Read(outputImg)
		imageResult <- true
	}
}
