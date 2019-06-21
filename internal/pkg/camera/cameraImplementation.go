package camera

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type Implentation interface {
	RunImagePoller(imageRequest <-chan bool, imageResult chan<- gocv.Mat)
}

type FakeCamera struct {
}

func (cam FakeCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- gocv.Mat) {
	for range imageRequest { // run while there are image reqests coming in
		img := gocv.NewMatWithSize(1000, 1000, gocv.MatTypeCV8UC3)
		gocv.Circle(&img, image.Point{500, 500}, 200, color.RGBA{255, 0, 0, 0}, 20)
		imageResult <- img
	}
}
