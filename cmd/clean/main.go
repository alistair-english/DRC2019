package main

import (
	"fmt"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"gocv.io/x/gocv"
)

func test(rows, cols int) gocv.Mat {
	mask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			mask.SetUCharAt3(col, row, 0, 30)
			mask.SetUCharAt3(col, row, 1, 50)
			mask.SetUCharAt3(col, row, 2, 90)
		}
	}
	return mask
}

func main() {
	rows := 1080
	cols := 1920
	// H := 100
	// S := 180
	// V := 200

	// cols := 1156 // min cols

	// mask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)

	// // copy HSV values to the mask
	// for row := 0; row < rows; row++ {
	// 	for col := 0; col < cols; col++ {
	// 		mask.SetUCharAt3(col, row, 0, 30)
	// 		mask.SetUCharAt3(col, row, 1, 50)
	// 		mask.SetUCharAt3(col, row, 2, 90)
	// 	}
	// }
	// mask := test(rows, cols)

	// mask := gocv.NewMatWithSize(cols, rows, gocv.MatTypeCV8UC3)
	thing := cvhelpers.NewHSVMask(100, 180, 200, cols, rows)

	// a := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	// a.SetTo(gocv.NewScalar(100, 180, 200, 0))
	// thing := cvhelpers.HSVMask(a)

	mask := gocv.Mat(thing)

	// mask.SetTo(gocv.NewScalar(90, 100, 180, 0))

	fmt.Println(mask.GetUCharAt3(0, 0, 0))
	fmt.Println(mask.GetUCharAt3(1, 1, 1))
	fmt.Println(mask.GetUCharAt3(2, 2, 2))
	fmt.Println(mask.ToBytes()[:100])
}
