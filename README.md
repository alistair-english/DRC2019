# DRC2019

## The Team
* Alistair English (SW)
* Thomas Fraser (SW)
* Cooper Richmond (ELEC/MECH)
* Tom Hulbert (ELEC/MECH)
* William Plummer (DOCS)

## The Challenge
The DRC challenge is a university technology challenge hosted by Queensland University of Technology where the competitors must design, build and program an autonomous droid robot to race others around an unknown track.

For more information please visit the [DRC Website](https://qutrobotics.com/2018/01/16/drc-2019/)

## The Robot
Our team's robot was made using an old RC car, a custom PCB and some lines of code.

### Hardware
The hardware of the car contains 1 RC car with 2700KV brushless motor
1 Custom PCB containing a [Raspberry Pi Compute Module](https://www.raspberrypi.org/products/compute-module-3/), a [Raspberry Pi Cam 2](https://www.raspberrypi.org/products/camera-module-v2/) and an [ESP 32](http://esp32.net).

### Software
The Raspberry Pi program was written in Golang and utilised the following libraries:
* [GOCV](https://gocv.io)
* [COLOR](https://github.com/fatih/color)
* [PIDCTRL](https://github.com/felixge/pidctrl)
* [SERIAL](https://github.com/tarm/serial)

Go was chosen for its ease of use for concurrent programming

The ESP32 program was written in C. C was chosen because the ESPIDF with FreeRTOS is designed to work with it. It also satisfies our need for low level control of both memory and pins.

### Results
Our team (wact<sup>2</sup>) managed to place 1st overall in the challenge, crossing the line ahead of our opponent in the final run while not incurring any penalties for breaches of the rules.

[Here](https://www.linkedin.com/feed/update/urn:li:ugcPost:6552799282754355201) is a video of the final run. Our robot is on the right at the start line