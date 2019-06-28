package cameraservice

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/alistair-english/DRC2019/pkg/config"
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
	cfg := config.GetPiCameraConfig()
	cmd := exec.Command("raspivid", "-t", "0", "-o", "-")
	var args []string
	if cfg.AWB != "" {
		args = append(args, "-awb", cfg.AWB)
	}
	if cfg.Bitrate != 0 {
		args = append(args, "-b", strconv.Itoa(cfg.Bitrate))
	}
	if cfg.Brightness != 0 {
		args = append(args, "-br", strconv.Itoa(cfg.Brightness))
	}
	if cfg.Contrast != 0 {
		args = append(args, "-co", strconv.Itoa(cfg.Contrast))
	}
	if cfg.Exposure != "" {
		args = append(args, "-ex", cfg.Exposure)
	}
	if cfg.FPS != 0 {
		args = append(args, "-fps", strconv.Itoa(cfg.FPS))
	}
	if cfg.Height != 0 {
		args = append(args, "-h", strconv.Itoa(cfg.Height))
	}
	if cfg.Mode != 0 {
		args = append(args, "-md", strconv.Itoa(cfg.Mode))
	}
	if cfg.ROI != "" {
		args = append(args, "-roi", cfg.ROI)
	}
	if cfg.Saturation != 0 {
		args = append(args, "-sa", strconv.Itoa(cfg.Saturation))
	}
	if cfg.Sharpness != 0 {
		args = append(args, "-sh", strconv.Itoa(cfg.Sharpness))
	}
	if cfg.Width != 0 {
		args = append(args, "-w", strconv.Itoa(cfg.Width))
	}

	cmd.Args = append(cmd.Args, args...)
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
