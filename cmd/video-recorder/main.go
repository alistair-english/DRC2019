package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	cam, _ := gocv.VideoCaptureDevice(2)
	defer cam.Close()

	// Camera settings
	// cam.Set(gocv.VideoCaptureBrightness, 50)

	window := gocv.NewWindow("Window")
	defer window.Close()

	var (
		frame = gocv.NewMat()
	)

	defer frame.Close()

	ret := false
	for !ret {
		ret = cam.Read(&frame)
	}

	fmt.Printf("%v\n", frame.Size())

	fname := fmt.Sprintf("%v_output.avi", time.Now().Format("01-02-2006_5:04:05"))

	out, _ := gocv.VideoWriterFile(fname, "XVID", 20.0, frame.Size()[1], frame.Size()[0], true)

	for cam.IsOpened() {
		ret := cam.Read(&frame)
		if ret == true {

			window.IMShow(frame)
			out.Write(frame)
		}

		if window.WaitKey(1) == 32 {
			break
		}

	}
}
