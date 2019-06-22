package main

import (
	"fmt"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/internal/pkg/camera"
)

func processImg(img gocv.Mat, motionCtrlChan chan<- bool) {
	// This is where the image gets processed and then we send a move control struct to the serial.
	fmt.Println(img.GetIntAt3(500, 500, 4))
	motionCtrlChan <- true
}

func main() {
	// Setup
	camImg := gocv.NewMat()
	defer camImg.Close()

	cam := camera.FakeCamera{}
	camConn := camera.NewConnection(cam, &camImg)

	serialChan := make(chan bool, 100)

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	camConn.RequestImage()

	for {
		select {
		case <-camConn.ImageResult:
			// Img is ready to be processed
			displayWindow.IMShow(camImg)
			displayWindow.WaitKey(1)

			// Spawn a go routine to do the heavy processing and then talk to serial when its done
			go processImg(camImg, serialChan)

			// Then request another image
			camConn.RequestImage()
		case <-serialChan:
			// Processing done (wouldn't actually have this here)
			fmt.Println("Processing done")
		}
	}
}
