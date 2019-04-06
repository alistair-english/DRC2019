package main

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {
	// cam, _ := gocv.VideoCaptureDevice(2)
	// defer cam.Close()

	vid, _ := gocv.VideoCaptureFile("../video-recorder/straight_line.avi")
	defer vid.Close()

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	flatWindow := gocv.NewWindow("Flat")
	defer flatWindow.Close()

	// Image Mats
	var (
		sourceImg = gocv.NewMat()
		flatImg   = gocv.NewMat()
	)

	// get the image
	vid.Read(&sourceImg)

	h := sourceImg.Size()[0]
	w := sourceImg.Size()[1]

	var (
		bl = image.Point{0, h/2 + 60}
		br = image.Point{w, h/2 + 60}
		tl = image.Point{150, 120}
		tr = image.Point{w - 150, 120}
	)

	src := []image.Point{tl, tr, bl, br}
	dst := []image.Point{
		image.Point{0, 0},
		image.Point{210, 0},
		image.Point{0, 297},
		image.Point{210, 297},
	}

	transform := gocv.GetPerspectiveTransform(src, dst)

	// fmt.Printf("%v\n", transform.Size())

	// Image closes
	defer sourceImg.Close()
	defer flatImg.Close()

	for { // inifinte loop

		// get the image
		vid.Read(&sourceImg)

		if !transform.Empty() {
			gocv.WarpPerspective(sourceImg, &flatImg, transform, image.Point{210, 297})
		}

		gocv.Circle(&sourceImg, bl, 2, color.RGBA{0, 255, 0, 1}, 4)
		gocv.Circle(&sourceImg, br, 2, color.RGBA{0, 255, 0, 1}, 4)
		gocv.Circle(&sourceImg, tl, 2, color.RGBA{0, 255, 0, 1}, 4)
		gocv.Circle(&sourceImg, tr, 2, color.RGBA{0, 255, 0, 1}, 4)

		// Display source img
		displayWindow.IMShow(sourceImg)
		flatWindow.IMShow(flatImg)

		d := displayWindow.WaitKey(100)

		if d == 32 {
			break
		}
	}

}
