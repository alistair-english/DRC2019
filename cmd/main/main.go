package main

import (
	"github.com/alistair-english/DRC2019/internal/pkg/camera"
	"gocv.io/x/gocv"
)

func main() {
	// Setup
	cam := camera.FakeCamera{}
	camConn := camera.NewConnection(cam)

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	camConn.ImageRequest <- true

	for {
		select {
		case img := <-camConn.ImageResult:

			// Code to handle an image is ready
			// e.g. process img - this realllyyy needs to be a goroutine otherwise huge block
			displayWindow.IMShow(img)
			displayWindow.WaitKey(1)

			// This is the gnarly part - need to remeber to ALWAYS close the img or ez mem leak
			img.Close()

			// Then request another image
			camConn.ImageRequest <- true
		}
	}
}
