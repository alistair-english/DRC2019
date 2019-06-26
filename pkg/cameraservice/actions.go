package cameraservice

import (
	"gocv.io/x/gocv"
)

// GetImageActionReq is an action request to get an image from a camera
type GetImageActionReq struct {
	img             *gocv.Mat
	responseChannel chan<- bool
}
