package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/alistair-english/DRC2019/pkg/gohelpers"
	"gocv.io/x/gocv"
)

func main() {
	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	width := 1920
	height := 1080
	channels := 3
	buffLen := width * height * channels

	cmd := exec.Command("raspividyuv", "-t", "0", "-o", "-", "-w", strconv.Itoa(width), "-h", strconv.Itoa(height))

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

	readBuff := make([]byte, width*height) // read in one 'img size' at a time (will take `channels` iterations to read full img)
	imgBuff := new(bytes.Buffer)
	currImg := make([]byte, buffLen)

	go func() {
		// func() {
		for {
			n, err := stdOut.Read(readBuff)
			if err != nil {
				fmt.Println(err)
			}

			numLeftToRead := buffLen - imgBuff.Len()

			numToWrite := gohelpers.IntMin(n, numLeftToRead)

			imgBuff.Write(readBuff[:numToWrite])

			fmt.Println(imgBuff.Len())

			if imgBuff.Len() >= buffLen {
				copy(currImg, imgBuff.Bytes()[:buffLen])
				imgBuff.Reset()
				fmt.Println(len(currImg))
				fmt.Println(currImg[:30])
				fmt.Println(currImg[len(currImg)-30:])
			}

			imgBuff.Write(readBuff[numToWrite:])
		}
	}()

	for {
		if len(currImg) > 0 {
			img, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, currImg)
			if err != nil {
				fmt.Println(err)
			}

			displayWindow.IMShow(img)
			displayWindow.WaitKey(1)
			// fmt.Println(len(currImg))
			fmt.Println(img.Size())
			img.Close()
		}
	}
}
