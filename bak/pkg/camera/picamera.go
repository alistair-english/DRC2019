package camera

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

// PiCamera is a camera implementation that connects to the Pi Camera
type PiCamera struct {
	syncChan chan bool
	currImg  []byte
	stdOut   io.ReadCloser
}

// NewPiCamera creates a new Pi Camera object
func NewPiCamera() (*PiCamera, error) {
	cmd := exec.Command("raspivid", "-cd", "MJPEG", "-t", "0", "-o", "-")
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

	return &piCam, nil
}

// Init from the camera Implementation
func (cam *PiCamera) Init() {
	go cam.CameraConnectionTask()
}

// CameraConnectionTask is the task that talks to the rpi camera std out
func (cam *PiCamera) CameraConnectionTask() {
	readBuff := make([]byte, 4096) // read in stream 4 kilobyte chunks
	imgBuff := new(bytes.Buffer)

	for {
		n, err := cam.stdOut.Read(readBuff)
		if err != nil {
			// Should probably log here but for now just retry reading
			continue
		}

		// fmt.Println(n)
		// fmt.Println(readBuff[:30])

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

		// fmt.Println(foundStart)

		if !foundStart {
			imgBuff.Write(readBuff)
		}
	}
}

// RunImagePoller from the camera Implementation
func (cam *PiCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat) {
	for range imageRequest {
		<-cam.syncChan

		img, err := gocv.IMDecode(cam.currImg, gocv.IMReadUnchanged)
		if err == nil {
			img.CopyTo(outputImg)
			img.Close()
		} else {
			fmt.Println(len(cam.currImg))
			fmt.Println(err)
		}

		imageResult <- true
	}
}
