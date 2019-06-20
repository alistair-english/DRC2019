package main

import (
	"log"
	t "time"
)

func main() {
	c := make(chan []byte)
	cr := make(chan int)

	go func(channel chan []byte, channelReturn chan int) {
		i := 0
		for v := range channel {
			log.Printf("Val: %v", v)
			i++
			cr <- i
		}
	}(c, cr)

	go func(channelReturn chan int) {
		for v := range channelReturn {
			log.Printf("Incr: %v", v)
		}
	}(cr)

	c <- []byte{0x11}
	c <- []byte{0x00, 0x01, 0x02}
	c <- []byte{0x01}
	c <- []byte{0x01}

	t.Sleep(1000)
}
