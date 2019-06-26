package cameraservice

import (
	"gocv.io/x/gocv"
)

// GetImageActionReq is an action request to get an image from a camera
type GetImageActionReq struct {
	Img             *gocv.Mat
	ResponseChannel chan<- bool
}
