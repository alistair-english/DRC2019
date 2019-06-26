package main

import (
	"fmt"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/camera"
	"github.com/alistair-english/DRC2019/pkg/serial"
)

func main() {
	// Camera Setup
	camImg := gocv.NewMat()
	defer camImg.Close()

	// cam, err := camera.NewPiCamera()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// camConn := camera.NewConnection(cam, &camImg)

	// // Serial Setup
	// serialConfig := config.GetSerialConfig()
	// ser, _ := serial.NewPiSerial(
	// 	serialConfig.Port,
	// 	serialConfig.Baud,
	// 	time.Duration(serialConfig.TimeoutMs)*time.Millisecond,
	// )
	// serConn, _ := serial.NewConnection(ser)

	// TESTING SETUP
	cam := camera.NewFileReaderCamera("../recordtodisk/recording_06-25-2019_16:10:16.avi")
	camConn := camera.NewConnection(cam, &camImg)

	// Serial Setup
	ser := serial.FakeSerial{}
	serConn, _ := serial.NewConnection(ser)

	// CV Process
	

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	camConn.RequestImage()

	for {
		select {
		case <-camConn.ImageResult:
			// Img is ready to be processed
			if camImg.Rows() > 0 || camImg.Cols() > 0 {
				displayWindow.IMShow(camImg)
				displayWindow.WaitKey(0)
			}

			// Spawn a go routine to do the heavy processing and then talk to serial when its done
			go processImg(camImg, serConn.ControlChan)

			// Then request another image
			camConn.RequestImage()
		case pwrReq := <-serConn.PowerReqChan:
			if pwrReq.ReqType == serial.ForceStop {
				fmt.Println("shut down!")
			}
		}
	}
}
