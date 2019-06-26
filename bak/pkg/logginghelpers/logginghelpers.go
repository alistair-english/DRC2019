package logginghelpers

import (
	"fmt"
)

const (
	serial     = iota
	imgproc    = iota
	debug      = iota
	serialcmds = iota
)

// TestLog is a testing function
func TestLog() {
	fmt.Println(imgproc)
}
