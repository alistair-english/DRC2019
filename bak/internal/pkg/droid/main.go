package droid

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
)

type droid struct {
	steering   rpio.Pin
	drive      rpio.Pin
	maxSpeed   int
	canReverse bool
}

// Constructs New Droid Object
func New(steeringPin int, drivePin int, maxSpeed int, canReverse bool) droid {
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()

	// Intialise Pins from Pins and Set Modes to PWM for Servo/ESC
	steering := rpio.Pin(steeringPin)
	steering.Mode(rpio.Pwm)
	drive := rpio.Pin(drivePin)
	drive.Mode(rpio.Pwm)

	// Begin PWM
	rpio.StartPwm()

	d := droid{steering, drive, maxSpeed, canReverse}
	fmt.Println("Droid Constructed")
	return d
}

// Moves the Droid
func (d droid) Move(speed, direction int) {
	// Do checks to make sure the speed is an acceptable value
	if speed > d.maxSpeed {
		speed = d.maxSpeed
	}
	if speed < 0 && !d.canReverse {
		speed = 0
	}
	// Convert a speed to a PWM value, this will change in the future once we have hardware
	var multiplier float32 = 2.55
	speed = int(multiplier * float32(speed))
	// Set the PWM frequencies for the drive and steering
	d.drive.Freq(speed)
	d.steering.Freq(direction)
}

// Stops the Droid on the Spot
func (d droid) Stop() {
	d.drive.Freq(0)
	d.steering.Freq(0)
}
