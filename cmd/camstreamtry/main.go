package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	//displayWindow := gocv.NewWindow("Display")
	//defer displayWindow.Close()

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
	readBuff := make([]byte, 4096) // 4 kilobytes of data

	n, err := stdOut.Read(readBuff)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)

	fmt.Println(len(readBuff))
	fmt.Println(readBuff[0:50])

	// for {
	// 	readBuff := make([]byte, 4096) // 4 kilobytes of data

	// 	n, err := stdOut.Read(readBuff)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(n)

	// 	fmt.Println(len(readBuff))
	// 	fmt.Println(readBuff[0:30])

	// 	// img, err := gocv.IMDecode(data, gocv.IMReadUnchanged)
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// }

	// 	//displayWindow.IMShow(img)
	// 	//displayWindow.WaitKey(1)
	// 	//fmt.Println(len(data))
	// 	// fmt.Println(img.Size())
	// 	// img.Close()
	// }
}
