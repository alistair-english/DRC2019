package camera

import (
	"gocv.io/x/gocv"
)

// Connection represents an ansyc connection to a camera implementation
type Connection struct {
	ImageRequest chan bool
	ImageResult  chan bool
	camera       Implementation
	outputImg    *gocv.Mat
}

// NewConnection generates a new camera connection from a camera implementation and a pointer to an output image
func NewConnection(camera Implementation, outputImg *gocv.Mat) *Connection {
	req := make(chan bool, 1)
	res := make(chan bool, 1)
	c := Connection{req, res, camera, outputImg}
	c.init()
	return &c
}

func (conn Connection) init() {
	// Set our image poller running with our channels
	go conn.camera.RunImagePoller(conn.ImageRequest, conn.ImageResult, conn.outputImg)
}

// RequestImage will add a new request token to the request channel of the camera if there isn't an exisiting one.
func (conn Connection) RequestImage() {
	select {
	case conn.ImageRequest <- true:
	}
	// Select will only choose non-blocking channel operations.
	// Therefore if there is already an image request in the channel writing another
	// will be blocking and the switch won't do it.
	// This is what we want as if there is already a new image requested we don't need to request another one
}
