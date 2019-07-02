# Camera -> Compute Module stuff
- https://www.raspberrypi.org/documentation/hardware/computemodule/cmio-camera.md
- https://lb.raspberrypi.org/forums/viewtopic.php?t=174347
- https://www.raspberrypi.org/forums/viewtopic.php?t=195810

# GOCV Install
- https://github.com/hybridgroup/gocv

# Serial Setup
 - https://www.raspberrypi.org/documentation/configuration/uart.md

# QUT Wifi
 - https://blackboard.qut.edu.au/bbcswebdav/pid-7269389-dt-content-rid-10861692_1/courses/IFB102_18se1/raspberry-pi-v49.pdf
   - under `Manual configuration of WiFi`
```
country=AU
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
network={
    ssid="QUT"
    key_mgmt=WPA-EAP
    eap=PEAP
    # put your QUT username between these quotes
    identity="my qut username"
    # replace the hash:... with a hash of your password # from ./wifi-pass.sh
    password="password"
    phase1="peaplabel=0"
    phase2="auth=MSCHAPV2"
    priority=10
}

# My Home network
network={
    ssid="my home network ssid"
    key_mgmt=WPA-PSK
    psk="my home network password"
} 
```
Your Pi is now ready to connect to the QUT network. You will need to reboot your pi (sudo reboot), you can test if its connected by typing ifconfigto see the state of the network interfaces. Note it can take several minutes to connect so be patient.