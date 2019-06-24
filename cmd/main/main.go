package main

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/camera"
	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/serial"
)

func processImg(img gocv.Mat, motionCtrlChan chan<- serial.Control) {
	// This is where the image gets processed and then we send a move control struct to the serial.
	if img.Rows() < 1 || img.Cols() < 1 {
		fmt.Println("processing 0 image")
		return
	}
	fmt.Println(img.GetVeciAt(img.Cols()/2, img.Rows()/2))
	motionCtrlChan <- serial.Control{
		Dir: 0,
		Spd: 50,
	}
}

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

	// Serial Setup
	serialConfig := config.GetSerialConfig()
	ser, _ := serial.NewPiSerial(
		serialConfig.Port,
		serialConfig.Baud,
		time.Duration(serialConfig.TimeoutMs)*time.Millisecond,
	)
	serConn, _ := serial.NewConnection(ser)

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	camConn.RequestImage()

	for {
		select {
		case <-camConn.ImageResult:
			// Img is ready to be processed
			if camImg.Rows() > 0 || camImg.Cols() > 0 {
				displayWindow.IMShow(camImg)
				displayWindow.WaitKey(1)
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
