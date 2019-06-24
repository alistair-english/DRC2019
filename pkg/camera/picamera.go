package camera

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

// PiCamera is a camera implementation that connects to the Pi Camera
type PiCamera struct {
	rwMutex *sync.RWMutex
	currImg []byte
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

	readBuff := make([]byte, 4096) // read in stream 4 kilobyte chunks
	imgBuff := new(bytes.Buffer)

	var currImg []byte
	piCam := PiCamera{new(sync.RWMutex), currImg}

	go func() {
		for {
			n, err := stdOut.Read(readBuff)
			if err != nil {
				// Should probably log here but for now just retry reading
				continue
			}

			fmt.Println(n)

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
						piCam.currImg = cpyImg

						// reset the buffer
						imgBuff.Reset()
					}

					imgBuff.Write(readBuff[i:])
					break
				}
			}

			fmt.Println(foundStart)

			if !foundStart {
				imgBuff.Write(readBuff)
			}
		}
	}()

	return &piCam, nil
}

// RunImagePoller from the camera Implementation
func (cam PiCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat) {
	for range imageRequest {
		for len(cam.currImg) < 1 {
			fmt.Println("waiting for big img")
		}

		img, err := gocv.IMDecode(cam.currImg, gocv.IMReadUnchanged)
		if err == nil {
			img.CopyTo(outputImg)
			img.Close()
		} else {
			fmt.Println(len(cam.currImg))
			fmt.Println(err)
		}
		fmt.Println("decode done")

		imageResult <- true
		fmt.Println("token done")
	}
}
