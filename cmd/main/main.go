package main

import (
	"github.com/alistair-english/DRC2019/internal/pkg/camera"
	"gocv.io/x/gocv"
)

func main() {
	// Setup
	camImg := gocv.NewMat()
	defer camImg.Close()

	cam := camera.FakeCamera{}
	camConn := camera.NewConnection(cam, &camImg)

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	camConn.RequestImage()

	for {
		select {
		case <-camConn.ImageResult:

			// Code to handle an image is ready
			// e.g. process img - this realllyyy needs to be a goroutine otherwise huge block
			displayWindow.IMShow(camImg)
			displayWindow.WaitKey(1)

			// Then request another image
			camConn.RequestImage()
		}
	}
}
