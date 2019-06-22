package camera

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// Implementation represents the code that actually talks to a camera. Is designed to be able to be swapped out with another camera Implementation.
type Implementation interface {
	RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat)
}

// FakeCamera is a dummy camera that just returns the same image.
type FakeCamera struct {
}

// RunImagePoller from the camera Implementation
func (cam FakeCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat) {
	for range imageRequest { // run while there are image reqests coming in
		img := gocv.NewMatWithSize(1000, 1000, gocv.MatTypeCV8UC3)
		gocv.Circle(&img, image.Point{500, 500}, 200, color.RGBA{255, 0, 0, 0}, 20)
		img.CopyTo(outputImg)
		imageResult <- true
		img.Close()
	}
}
