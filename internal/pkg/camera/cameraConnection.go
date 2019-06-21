package camera

import (
	"gocv.io/x/gocv"
)

type Connection struct {
	ImageRequest chan bool
	ImageResult  chan gocv.Mat
	camera       Implentation
}

func NewConnection(camera Implentation) *Connection {
	req := make(chan bool, 1)
	res := make(chan gocv.Mat, 1)
	c := Connection{req, res, camera}
	c.init()
	return &c
}

func (conn Connection) init() {
	// Set our image poller running with our channels
	go conn.camera.RunImagePoller(conn.ImageRequest, conn.ImageResult)
}
