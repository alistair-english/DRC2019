package main

import (
	"fmt"
	"os"
	"os/exec"

	"gocv.io/x/gocv"
)

var jpgStart = []byte{0xFF, 0xD8, 0xFF}

func main() {
	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	width := 1920
	height := 1080

	cmd := exec.Command("raspividyuv", "-t", "0", "-o", "-") //, "-w", strconv.Itoa(width), "-h", strconv.Itoa(height))

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

	imgBuff := make([]byte, width*height) // read in stream

	go func() {
		for {
			n, err := stdOut.Read(imgBuff)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
		}
	}()

	for {
		if len(imgBuff) > 0 {
			img, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, imgBuff)
			if err != nil {
				fmt.Println(err)
			}

			displayWindow.IMShow(img)
			// displayWindow.WaitKey(1)
			// fmt.Println(len(currImg))
			// fmt.Println(img.Size())
			img.Close()
		}
	}
}
