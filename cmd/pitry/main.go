package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func show(img gocv.Mat, window gocv.Window) {
	window.IMShow(img)
	window.WaitKey(25)
}

func main() {
	displayWindow := gocv.NewWindow("Ur Mom")
	defer displayWindow.Close()

	for i := 0; i <= 20; i++ {
		cam, e := gocv.VideoCaptureDevice(i)
		if e == nil {
			fmt.Println(i)
			cam.Close()
		}
	}

	// img := gocv.NewMat()

	// for {
	// 	cam.Read(&img)
	// 	displayWindow.IMShow(img)
	// 	displayWindow.WaitKey(1)
	// }
}
