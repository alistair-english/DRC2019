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
	fmt.Println(img.GetIntAt3(500, 500, 4))
	motionCtrlChan <- serial.Control{
		Dir: 0,
		Spd: 50,
	}
}

func main() {
	// Camera Setup
	camImg := gocv.NewMat()
	defer camImg.Close()
	cam := camera.FakeCamera{}
	camConn := camera.NewConnection(cam, &camImg)

	// Serial Setup
	serialConfig := config.GetSerialConfig()
	ser := serial.NewPiSerial(
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
			displayWindow.IMShow(camImg)
			displayWindow.WaitKey(1)

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
