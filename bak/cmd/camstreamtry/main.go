package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

func main() {
	// displayWindow := gocv.NewWindow("Display")
	// defer displayWindow.Close()

	widthPtr := flag.Int("w", -1, "width")
	heightPtr := flag.Int("h", -1, "height")
	flag.Parse()

	var cmd *exec.Cmd
	if *widthPtr != -1 && *heightPtr != -1 {
		fmt.Println("w: ", *widthPtr, " h: ", *heightPtr)
		cmd = exec.Command("raspivid", "-cd", "MJPEG", "-t", "0", "-o", "-", "-w", strconv.Itoa(*widthPtr), "-h", strconv.Itoa(*heightPtr))
	} else {
		fmt.Println("no bois")
		cmd = exec.Command("raspivid", "-cd", "MJPEG", "-t", "0", "-o", "-")
	}

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

	channel := make(chan bool, 1)

	go func(channel chan<- bool) {
		for {

			n, err := stdOut.Read(readBuff)
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Println(n)
			// fmt.Println(len(readBuff))
			// fmt.Println(readBuff[0:30])

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
						currImg = cpyImg

						select {
						case channel <- true:
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
	}(channel)

	for range channel {
		if len(currImg) > 0 {
			img, err := gocv.IMDecode(currImg, gocv.IMReadUnchanged)
			if err != nil {
				fmt.Println(err)
			}

			// fmt.Println(img.GetVeciAt(img.Cols()/2, img.Rows()/2))
			fmt.Println(img.Cols(), img.Rows())

			// displayWindow.IMShow(img)
			// displayWindow.WaitKey(1)
			// fmt.Println(len(currImg))
			// fmt.Println(img.Size())
			img.Close()
		}
	}
}
