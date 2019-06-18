package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	_ "github.com/alistair-english/DRC2019/internal/pkg/config"
	"gocv.io/x/gocv"
)

func main() {
	img := gocv.NewMat()
	mask := gocv.NewMat()

	ub := gocv.NewMatFromScalar(gocv.Scalar{Val1: 255.0, Val2: 255.0, Val3: 255.0, Val4: 1.0}, gocv.MatTypeCV8UC3)
	lb := gocv.NewMatFromScalar(gocv.Scalar{Val1: 254.0, Val2: 254.0, Val3: 254.0, Val4: 0.9}, gocv.MatTypeCV8UC3)

	defer img.Close()
	defer mask.Close()

	defer ub.Close()
	defer lb.Close()

	displayWindow := gocv.NewWindow("Window")
	defer displayWindow.Close()

	// displayWindow2 := gocv.NewWindow("Window 2")
	// defer displayWindow2.Close()

	img = gocv.IMRead("../../testImages/testLineLeft.png", gocv.IMReadColor)
	fmt.Printf("Image Type: %v \n", img.Type())
	gocv.InRange(img, lb, ub, &mask)
	// mask = convertToBitmask(img)

	const pointsPerImage = 20

	y := mask.Rows()
	x := mask.Cols()

	fmt.Printf("Image is %vx%vpx \n", x, y)

	rowsPerRead := y / pointsPerImage
	averages := make([][2]int, rowsPerRead)

	var (
		unsortedTotal = 0
		unsortedCount = 1
	)

	for i := 0; i < rowsPerRead; i++ {
		//Iterate over the columns of the image.
		unsortedTotal = 0
		unsortedCount = 1

		for j := 0; j < x; j++ {
			//Check if bitmask is 0 or -1
			if mask.GetIntAt(i*pointsPerImage, j) != 0 {
				fmt.Println(mask.GetIntAt(i*pointsPerImage, j))
				unsortedTotal += j
				unsortedCount++
			}
		}
		fmt.Println()
		// Calculate the X average and use Y of i
		averages[i][0] = unsortedTotal / unsortedCount
		averages[i][1] = i
	}

	fmt.Printf("%v\n", averages)

	for i := 0; i < rowsPerRead; i++ {
		point := image.Point{averages[i][0], averages[i][1] * pointsPerImage}
		gocv.Circle(&mask, point, 10, color.RGBA{255, 255, 255, 1}, 1)
	}

	for {
		displayWindow.IMShow(mask)
		displayWindow.WaitKey(1)
	}
}

func convertToBitmask(mat gocv.Mat) gocv.Mat {
	mask := gocv.NewMat()
	gocv.InRangeWithScalar(mat, gocv.Scalar{Val1: 200.0, Val2: 200.0, Val3: 200.0, Val4: 0.9}, gocv.Scalar{Val1: 255.0, Val2: 255.0, Val3: 255.0, Val4: 1.0}, &mask)
	return mask
}

func openImage(path string) (gocv.Mat, error) {
	fimag, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fimag.Close()

	img, _, _ := image.Decode(fimag)

	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	//don't get surprised of reversed order everywhere below
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}
