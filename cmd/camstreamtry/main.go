package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

func main() {
	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	cmd := exec.Command("raspivid", "-cd", "MJPEG", "-t", "0", "-o", "-")

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Stderr = os.Stderr

	readBuff := make([]byte, 4096) // read in stream 4 kilobyte chunks
	imgBuff := new(bytes.Buffer)

	var currImg []byte

	go func() {
		for {

			n, err := stdOut.Read(readBuff)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
			fmt.Println(len(readBuff))
			fmt.Println(readBuff[0:30])

			for i := 0; i < n; i++ {
				// check if we found the start of an image
				if bytes.Compare(readBuff[i:i+len(jpgStart)], jpgStart) == 0 {
					// we found a new image start point at i

					// write the rest of the old image into the currImg buffer
					imgBuff.Write(readBuff[0:i])

					if imgBuff.Len() > 0 {
						// the was already part of an img in here - img must be done

						// Copy the completed image out of the buffer and into the current img
						cpyImg := make([]byte, imgBuff.Len())
						copy(cpyImg, imgBuff.Bytes())
						currImg = cpyImg

						// reset the buffer
						imgBuff.Reset()
					}

					imgBuff.Write(readBuff[i:])
				} else {
					imgBuff.Write(readBuff)
				}
			}
		}
	}()

	for {
		img, err := gocv.IMDecode(currImg, gocv.IMReadUnchanged)
		if err != nil {
			fmt.Println(err)
		}

		displayWindow.IMShow(img)
		displayWindow.WaitKey(1)
		fmt.Println(len(currImg))
		fmt.Println(img.Size())
		img.Close()
	}
}
