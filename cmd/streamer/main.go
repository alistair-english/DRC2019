package main

import (
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"net/http"
	_ "net/http/pprof"
)

var (
	err    error
	cam    *gocv.VideoCapture
	stream *mjpeg.Stream
	host   string = "localhost:8080"
)

func main() {
	cam, err = gocv.OpenVideoCapture(0)
	if err != nil {
		panic(err)
	}
	defer cam.Close()

	stream = mjpeg.NewStream()
	go mjpegCapture()

	http.Handle("/", stream)
	http.ListenAndServe(host, nil)
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		cam.Read(&img)

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
