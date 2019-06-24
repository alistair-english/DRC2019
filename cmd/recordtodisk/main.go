package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/camera"
)

func main() {
	// Camera Setup
	camImg := gocv.NewMat()
	defer camImg.Close()

	cam, err := camera.NewPiCamera()
	if err != nil {
		fmt.Println(err)
		return
	}
	camConn := camera.NewConnection(cam, &camImg)

	// displayWindow := gocv.NewWindow("Display")
	// defer displayWindow.Close()

	camConn.RequestImage()

	<-camConn.ImageResult
	fname := fmt.Sprintf("recording_%v.avi", time.Now().Format("01-02-2006_5:04:05"))
	out, _ := gocv.VideoWriterFile(fname, "XVID", 20.0, camImg.Size()[1], camImg.Size()[0], true)

	camConn.RequestImage()

	for {
		select {
		case <-camConn.ImageResult:
			// Img is ready to be processed
			// if camImg.Rows() > 0 || camImg.Cols() > 0 {
			// 	displayWindow.IMShow(camImg)
			// 	displayWindow.WaitKey(1)
			// }

			// Write the img to disk
			out.Write(camImg)
			fmt.Println("frame written")

			// Then request another image
			camConn.RequestImage()
		}
	}
}
