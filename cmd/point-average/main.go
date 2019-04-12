package main

import (
	"fmt"
	_ "github.com/alistair-english/DRC2019/internal/pkg/config"
	"gocv.io/x/gocv"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"os"
)

func main() {
	var img = gocv.NewMat()
	defer img.Close()

	img, _ = openImage("../../testImages/testLineLeft.png")
	findPoints(img)
}

func toRGB8(img image.Image) (gocv.Mat, error) {
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

// openImage
func openImage(path string) (gocv.Mat, error) {
	fimag, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fimag.Close()

	img, _, _ := image.Decode(fimag)

	return toRGB8(img)
}

func findPoints(mat gocv.Mat) {

	// The number of vertical points taken for calculations in an image
	const pointsPerImage = 20

	//Get Image Resolution
	rows := mat.Rows()
	cols := mat.Cols()

	fmt.Printf("Image is %vx%vpx \n", rows, cols)

	//Get colums per read
	rowsPerRead := rows / pointsPerImage
	averages := make([][2]int, rowsPerRead)

	var (
		unsortedTotal = 0
		unsortedCount = 1
	)
	//Iterate over the image to find the X,Y Average for each of the chosen rows
	for i := 0; i < rowsPerRead; i++ {
		//Iterate over the columns of the image.
		for j := 0; j < cols; j++ {
			//Check if bitmask is 0 or -1
			if mat.GetIntAt(j, i) == -1 {
				unsortedTotal += j
				unsortedCount++
			}
		}
		// Calculate the X average and use Y of i
		averages[i][0] = unsortedTotal / unsortedCount
		averages[i][1] = i
	}

	fmt.Printf("%v\n", averages)

	// Create Mask with points
	// mask := gocv.NewMatWithSize(rows, cols, 0)
	// defer mask.Close()
	// for i := 0; i < rowsPerRead; i++ {
	// 	point := image.Point{averages[i][0], averages[i][1] * pointsPerImage}
	// 	var value = uint8(255 / rowsPerRead * i)
	// 	gocv.Circle(&mask, point, 10, color.RGBA{value, value, value, 1}, 1)
	// }

	// displayWindow := gocv.NewWindow("Window")
	// defer displayWindow.Close()
	// for {
	// 	displayWindow.IMShow(mask)
	// 	displayWindow.WaitKey(1)
	// }

}
