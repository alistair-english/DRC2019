package main

import (
	"os/exec"
	"fmt"
	"gocv.io/x/gocv"
)

func main() {
	//displayWindow := gocv.NewWindow("Display")
	//defer displayWindow.Close()

	for {
		cmd := exec.Command("raspistill", "-o", "-")
		data, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(data))

		img,err := gocv.IMDecode(data, gocv.IMReadUnchanged)
	        if err != nil {
	               fmt.Println(err)
	        }

		//displayWindow.IMShow(img)
		//displayWindow.WaitKey(1)
		//fmt.Println(len(data))
		fmt.Println(img.Size())
		img.Close()
	}
}
