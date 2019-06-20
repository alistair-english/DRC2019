package serial

import (
	"github.com/alistair-english/DRC2019/internal/pkg/config"
	s "github.com/tarm/serial"
	"log"
	t "time"
)

// Init the serial connection from the port given in config
func Init(c chan []byte) {
	camPort := config.GetCameraConfig()
	log.Printf("Starting Serial Communication")

	// Serial options
	options := s.Config{
		Name:        camPort.Port,       // Serial port
		Baud:        115200,             // Probs the wrong baud rate
		ReadTimeout: t.Millisecond * 10, // Filler Value
		Size:        8,                  //Frame Size
		StopBits:    1,                  //Stop Bits
		Parity:      'N',                //Parity
	}

	port, err := s.OpenPort(&options)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Write function which takes data in from the channel
	// go func(channel chan []byte) {
	// 	for {
	// 		for v := range channel {
	// 			port.Write(v)
	// 			fmt.Printf("Channel Val: %v", v)
	// 		}
	// 	}
	// }(c)
}
