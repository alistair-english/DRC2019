package cameraservice

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/logging"
	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

// PiCamera is a camera implementation that reads from the Pi Camera
type PiCamera struct {
	syncChan chan bool
	currImg  []byte
	stdOut   io.ReadCloser
}

func getPiCameraCmd() *exec.Cmd {
	piCfg := config.GetPiCameraConfig()
	cvCfg := config.GetCVConfig()

	cmd := exec.Command("raspivid", "-cd", "MJPEG", "-t", "0", "-o", "-")

	var args []string

	if cvCfg.ImgHeight != 0 {
		args = append(args, "-h", strconv.Itoa(cvCfg.ImgHeight))
	}
	if cvCfg.ImgWidth != 0 {
		args = append(args, "-w", strconv.Itoa(cvCfg.ImgWidth))
	}

	if piCfg.AWB != "" {
		args = append(args, "-awb", piCfg.AWB)
	}
	if piCfg.AWBGains != "" {
		args = append(args, "-awbg", piCfg.AWBGains)
	}
	if piCfg.Bitrate != 0 {
		args = append(args, "-b", strconv.Itoa(piCfg.Bitrate))
	}
	if piCfg.Brightness != 0 {
		args = append(args, "-br", strconv.Itoa(piCfg.Brightness))
	}
	if piCfg.Contrast != 0 {
		args = append(args, "-co", strconv.Itoa(piCfg.Contrast))
	}
	if piCfg.Exposure != "" {
		args = append(args, "-ex", piCfg.Exposure)
	}
	if piCfg.FPS != 0 {
		args = append(args, "-fps", strconv.Itoa(piCfg.FPS))
	}
	if piCfg.Mode != 0 {
		args = append(args, "-md", strconv.Itoa(piCfg.Mode))
	}
	if piCfg.ROI != "" {
		args = append(args, "-roi", piCfg.ROI)
	}
	if piCfg.Saturation != 0 {
		args = append(args, "-sa", strconv.Itoa(piCfg.Saturation))
	}
	if piCfg.Sharpness != 0 {
		args = append(args, "-sh", strconv.Itoa(piCfg.Sharpness))
	}

	cmd.Args = append(cmd.Args, args...)

	logging.L().Logln(TAG, logging.All, "raspivid args: %v", cmd.Args)

	return cmd

}

// NewPiCameraImplementation creates a new PiCamera Implementation
func newPiCameraImplementation() (*PiCamera, error) {
	cmd := getPiCameraCmd()
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	syncChan := make(chan bool, 1)

	var currImg []byte
	piCam := PiCamera{syncChan, currImg, stdOut}

	go piCam.CameraConnectionTask()

	return &piCam, nil
}

// CameraConnectionTask is the task that talks to the rpi camera std out
func (cam *PiCamera) CameraConnectionTask() {
	readBuff := make([]byte, 4096) // read in stream 4 kilobyte chunks
	imgBuff := new(bytes.Buffer)

	for {
		n, err := cam.stdOut.Read(readBuff)
		if err != nil {
			// Should probably log here but for now just retry reading
			fmt.Println(err)
			continue
		}

		foundStart := false

		// check if we found the start of an image
		for i := 0; i <= (n - len(jpgStart)); i++ {
			if bytes.Compare(readBuff[i:i+len(jpgStart)], jpgStart) == 0 {
				// we found a new image start point at i
				foundStart = true

				// write the rest of the old image into the currImg buffer
				imgBuff.Write(readBuff[0:i])

				if imgBuff.Len() > 0 {
					// the was already part of an img in here - img must be done

					// Copy the completed image out of the buffer and into the current img
					cpyImg := make([]byte, imgBuff.Len())
					copy(cpyImg, imgBuff.Bytes())
					cam.currImg = cpyImg

					select {
					case cam.syncChan <- true:
					default:
					}

					// reset the buffer
					imgBuff.Reset()
				}

				imgBuff.Write(readBuff[i:])
				break
			}
		}

		if !foundStart {
			imgBuff.Write(readBuff)
		}
	}
}

// RunCameraConnection from camera Implementation
func (cam *PiCamera) RunCameraConnection(imgRequests <-chan GetImageActionReq) {
	for req := range imgRequests {
		<-cam.syncChan

		img, err := gocv.IMDecode(cam.currImg, gocv.IMReadUnchanged)
		if err == nil {
			img.CopyTo(req.Img)
			img.Close()
		} else {
			fmt.Println(len(cam.currImg))
			fmt.Println(err)
		}

		select {
		case req.ResponseChannel <- true:
		default:
		}
	}
}
