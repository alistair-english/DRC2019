package main

import (
	"fmt"
	"time"

	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/serial"
)

func main() {
	// Serial Setup
	serialConfig := config.GetSerialConfig()
	ser, err := serial.NewPiSerial(
		serialConfig.Port,
		serialConfig.Baud,
		time.Duration(serialConfig.TimeoutMs)*time.Millisecond,
	)
	if err != nil {
		fmt.Println(err)
	}

	serConn, err := serial.NewConnection(ser)
	if err != nil {
		fmt.Println(err)
	}

	for {
		serConn.ControlChan <- serial.Control{
			Dir: 0,
			Spd: 50,
		}

		time.Sleep(1000 * time.Millisecond)

		serConn.ControlChan <- serial.Control{
			Dir: -90,
			Spd: 50,
		}

		time.Sleep(1000 * time.Millisecond)

		serConn.ControlChan <- serial.Control{
			Dir: 90,
			Spd: 50,
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
