package main

import (
	"fmt"
	_ "github.com/alistair-english/DRC2019/internal/pkg/config"
	"gocv.io/x/gocv"
	_ "image"
	_ "image/color"
)

func main() {

}

func findPoints(mat gocv.Mat) {
	fmt.Println("Hello, World!")

	// The number of vertical points taken for calculations in an image
	const pointsPerImage = 20

	//Get Image Resolution
	rows := mat.Rows()
	cols := mat.Cols()

	//Get colums per read
	rowsPerRead := rows / pointsPerImage
	averages := make([][2]int, rowsPerRead)

	//Iterate over the image to find the X,Y Average for each of the chosen rows
	for i := 0; i < rowsPerRead; i++ {
		var (
			unsortedTotal = 0
			unsortedCount = 0
		)
		//Iterate over the columns of the image.
		for j := 0; j < cols; j++ {
			//Check if bitmask is 0 or -1
			if mat.GetIntAt(i, j) == -1 {
				unsortedTotal += j
			}
		}
		// Calculate the X average and use Y of i
		averages[i][0] = unsortedTotal / unsortedCount
		averages[i][1] = i
	}
}
