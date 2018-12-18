package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"math"
)

func main() {
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println("Unable to Open The Capture Device")
	}
	window := gocv.NewWindow("Test Capture")
	img := gocv.NewMat()

	var rho float32 = 1.0
	var theta float32 = math.Pi / 180
	threshold := 15
	lineImage := gocv.NewMat()

	var lowThreshold float32 = 0
	var highThreshold float32 = 150
	edges := gocv.NewMat()

	for {
		webcam.Read(&img)
		gocv.GaussianBlur(img, &img, image.Point{0, 0}, 2, 2, 1)
		gocv.CvtColor(img, &img, gocv.ColorBGRAToBGR)
		gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
		// gocv.InRange(img, lowerW, upperW, &whiteMask)
		gocv.Canny(img, &edges, lowThreshold, highThreshold)
		gocv.HoughLinesP(edges, &lineImage, rho, theta, threshold)
		window.IMShow(edges)
		window.WaitKey(1)
	}
}
